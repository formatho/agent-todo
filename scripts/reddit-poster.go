package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Reddit API structures
type RedditPost struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Subreddit string   `json:"subreddit"`
	FlairID   string   `json:"flair_id,omitempty"`
	FlairText string   `json:"flair_text,omitempty"`
}

type RedditAuth struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RedditError struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type RedditResponse struct {
	Data struct {
		URL     string `json:"url"`
		Name    string `json:"name"`
		Created int64  `json:"created_utc"`
	} `json:"data"`
}

// Configuration
type RedditConfig struct {
	ClientID     string
	ClientSecret string
	UserAgent    string
	Username     string
	Password     string
}

// Get Reddit access token
func getRedditAccessToken(config RedditConfig) (*RedditAuth, error) {
	url := "https://www.reddit.com/api/v1/access_token"
	
	data := strings.NewReader("grant_type=password&username=" + config.Username + "&password=" + config.Password)
	
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.SetBasicAuth(config.ClientID, config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", config.UserAgent)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var auth RedditAuth
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %v", err)
	}
	
	return &auth, nil
}

// Post to Reddit
func postToReddit(auth *RedditAuth, post RedditPost) (*RedditResponse, error) {
	url := fmt.Sprintf("https://oauth.reddit.com/r/%s/api/submit", post.Subreddit)
	
	payload := map[string]interface{}{
		"api_type": "json",
		"sr":       post.Subreddit,
		"title":    post.Title,
		"kind":     "self",
		"text":     post.Content,
	}
	
	if post.FlairID != "" {
		payload["flair_id"] = post.FlairID
	}
	
	if post.FlairText != "" {
		payload["flair_text"] = post.FlairText
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal post data: %v", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create post request: %v", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", auth.TokenType, auth.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Agent-Todo/1.0")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to post to Reddit: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		var redditError RedditError
		if err := json.Unmarshal(body, &redditError); err == nil {
			return nil, fmt.Errorf("Reddit API error: %s (%s)", redditError.Message, redditError.Reason)
		}
		return nil, fmt.Errorf("Reddit API error with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response RedditResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	
	return &response, nil
}

// Create a Reddit post for blog content
func createRedditPost(blogPath string) (*RedditResponse, error) {
	// Read blog content
	content, err := os.ReadFile(blogPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read blog content: %v", err)
	}
	
	// Extract title and content (simplified)
	lines := strings.Split(string(content), "\n")
	title := lines[0] // First line is usually the title
	if strings.HasPrefix(title, "# ") {
		title = strings.TrimPrefix(title, "# ")
	}
	
	// Create content excerpt
	var contentLines []string
	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "# ") {
			// Skip other headers
			continue
		}
		if strings.TrimSpace(line) != "" {
			contentLines = append(contentLines, line)
			if len(contentLines) >= 10 { // Limit excerpt length
				break
			}
		}
	}
	
	contentText := strings.Join(contentLines, "\n")
	if len(contentText) > 500 {
		contentText = contentText[:500] + "... [Read more](https://todo.formatho.com/blog/" + blogPath + ")"
	}
	
	// Create Reddit post
	redditPost := RedditPost{
		Title:     fmt.Sprintf("🤖 AI Task Management: %s", title),
		Content:   contentText + "\n\n\n---\n\n🔗 Read the full article: https://todo.formatho.com/blog/" + blogPath,
		Subreddit: "artificial", // Main AI subreddit
		FlairText: "AI Tool",
	}
	
	// Load configuration from environment
	config := RedditConfig{
		ClientID:     os.Getenv("REDDIT_CLIENT_ID"),
		ClientSecret: os.Getenv("REDDIT_CLIENT_SECRET"),
		UserAgent:    "Agent-Todo/1.0",
		Username:     os.Getenv("REDDIT_USERNAME"),
		Password:     os.Getenv("REDDIT_PASSWORD"),
	}
	
	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("Reddit API credentials not found in environment variables")
	}
	
	// Get access token
	fmt.Println("Authenticating with Reddit...")
	auth, err := getRedditAccessToken(config)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Reddit: %v", err)
	}
	
	// Post to Reddit
	fmt.Printf("Posting to r/%s: %s\n", redditPost.Subreddit, redditPost.Title)
	response, err := postToReddit(auth, redditPost)
	if err != nil {
		return nil, fmt.Errorf("failed to post to Reddit: %v", err)
	}
	
	return response, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: reddit-poster <blog-path>")
		os.Exit(1)
	}
	
	blogPath := os.Args[1]
	fmt.Printf("Processing blog post: %s\n", blogPath)
	
	response, err := createRedditPost(blogPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("✅ Successfully posted to Reddit!\n")
	fmt.Printf("   Post ID: %s\n", response.Data.Name)
	fmt.Printf("   URL: https://reddit.com/%s\n", response.Data.Name)
	fmt.Printf("   Created: %s\n", time.Unix(response.Data.Created, 0).Format(time.RFC1123))
}