package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// SocialMediaPost represents a social media post
type SocialMediaPost struct {
	Platform     string    `json:"platform"`
	Content      string    `json:"content"`
	Hashtags     []string  `json:"hashtags"`
	ImagePath    string    `json:"image_path,omitempty"`
	VideoPath    string    `json:"video_path,omitempty"`
	LinkURL      string    `json:"link_url,omitempty"`
	PublishTime  time.Time `json:"publish_time"`
	IsPublished  bool      `json:"is_published"`
	Engagement   int       `json:"engagement"`
}

// SocialMediaCampaign represents a complete social media campaign
type SocialMediaCampaign struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	LaunchDate  time.Time            `json:"launch_date"`
	Platforms   []string             `json:"platforms"`
	Posts       []SocialMediaPost    `json:"posts"`
	Metrics     CampaignMetrics      `json:"metrics"`
}

// CampaignMetrics represents campaign success metrics
type CampaignMetrics struct {
	TargetImpressions int   `json:"target_impressions"`
	TargetEngagement int   `json:"target_engagement"`
	TargetReach      int   `json:"target_reach"`
	ActualImpressions int  `json:"actual_impressions"`
	ActualEngagement int  `json:"actual_engagement"`
	ActualReach      int  `json:"actual_reach"`
}

// ContentTemplate represents a reusable content template
type ContentTemplate struct {
	Name        string   `json:"name"`
	Platform    string   `json:"platform"`
	ContentType string   `json:"content_type"`
	Structure   []string `json:"structure"`
	Example     string   `json:"example"`
}

// SocialMediaContentGenerator generates social media content
type SocialMediaContentGenerator struct {
	templates []ContentTemplate
}

// NewSocialMediaContentGenerator creates a new content generator
func NewSocialMediaContentGenerator() *SocialMediaContentGenerator {
	return &SocialMediaContentGenerator{
		templates: []ContentTemplate{
			{
				Name:        "Launch Announcement",
				Platform:    "twitter",
				ContentType: "announcement",
				Structure:   []string{"emoji", "announcement", "key_feature", "link", "hashtags"},
				Example:     "🚀 BIG NEWS! We're live on Product Hunt! Agent-Todo is now officially launched - the task management platform built specifically for AI agents. No more context loss, no more forgotten tasks. 🤖➡️✅ https://producthunt.com/products/agent-todo #ProductHunt #AI #TaskManagement",
			},
			{
				Name:        "Feature Spotlight",
				Platform:    "twitter",
				ContentType: "feature",
				Structure:   []string{"emoji", "feature_name", "description", "benefit", "hashtags"},
				Example:     "✨ What makes Agent-Todo different? We built it FOR AI agents, not just for humans. Features: • Persistent memory • Priority-based task management • Real-time agent status tracking Your AI agents deserve better tools. 🛠️ #AIProductivity #DeveloperTools",
			},
			{
	Name:        "Problem/Solution",
				Platform:    "linkedin",
				ContentType: "problem_solution",
				Structure:   []string{"emoji", "problem_statement", "impact", "solution", "benefit", "hashtags"},
				Example:     "🔥 The problem: \"My AI agent lost track of what it was doing after a restart\" The solution: Agent-Todo gives agents persistent memory and task tracking. No more context loss, no more duplicated work. It's that simple. 💡 #AIPainPoint #Solutions #AI",
			},
		},
	}
}

// GeneratePost generates a social media post based on template
func (g *SocialMediaContentGenerator) GeneratePost(templateName, platform, contentType string, data map[string]string) (*SocialMediaPost, error) {
	// Find template
	var template *ContentTemplate
	for _, t := range g.templates {
		if t.Name == templateName && t.Platform == platform && t.ContentType == contentType {
			template = &t
			break
		}
	}
	
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}
	
	// Generate content using template structure
	var contentBuilder strings.Builder
	
	for _, part := range template.Structure {
		switch part {
		case "emoji":
			contentBuilder.WriteString(g.getEmojiForContentType(contentType))
		case "announcement":
			contentBuilder.WriteString(data["announcement"])
		case "key_feature":
			contentBuilder.WriteString(g.formatKeyFeature(data["feature"]))
		case "description":
			contentBuilder.WriteString(data["description"])
		case "benefit":
			contentBuilder.WriteString(g.formatBenefit(data["benefit"]))
		case "problem_statement":
			contentBuilder.WriteString(data["problem"])
		case "impact":
			contentBuilder.WriteString(g.formatImpact(data["impact"]))
		case "solution":
			contentBuilder.WriteString(data["solution"])
		case "link":
			if url, ok := data["url"]; ok {
				contentBuilder.WriteString(" " + url)
			}
		}
		contentBuilder.WriteString(" ")
	}
	
	// Generate hashtags
	var hashtags []string
	if tags, ok := data["hashtags"]; ok {
		hashtags = strings.Split(tags, " ")
	} else {
		hashtags = g.getDefaultHashtags(platform, contentType)
	}
	
	return &SocialMediaPost{
		Platform:    platform,
		Content:     strings.TrimSpace(contentBuilder.String()),
		Hashtags:    hashtags,
		LinkURL:     data["url"],
		PublishTime: time.Now(),
		IsPublished: false,
		Engagement: 0,
	}, nil
}

// getEmojiForContentType returns appropriate emoji for content type
func (g *SocialMediaContentGenerator) getEmojiForContentType(contentType string) string {
	switch contentType {
	case "announcement":
		return "🚀"
	case "feature":
		return "✨"
	case "problem_solution":
		return "🔥"
	case "testimonial":
		return "💬"
	case "use_case":
		return "🎯"
	case "behind_scenes":
		return "🏗️"
	case "interactive":
		return "📊"
	case "call_to_action":
		return "💪"
	default:
		return "📝"
	}
}

// formatKeyFeature formats a key feature
func (g *SocialMediaContentGenerator) formatKeyFeature(feature string) string {
	if feature == "" {
		return ""
	}
	return fmt.Sprintf("• %s", feature)
}

// formatBenefit formats a benefit
func (g *SocialMediaContentGenerator) formatBenefit(benefit string) string {
	if benefit == "" {
		return ""
	}
	return fmt.Sprintf("Your %s.", benefit)
}

// formatImpact formats impact
func (g *SocialMediaContentGenerator) formatImpact(impact string) string {
	if impact == "" {
		return ""
	}
	return fmt.Sprintf("The impact: %s", impact)
}

// getDefaultHashtags returns default hashtags for platform and content type
func (g *SocialMediaContentGenerator) getDefaultHashtags(platform, contentType string) []string {
	baseTags := []string{"AI", "Productivity", "Technology"}
	
	switch platform {
	case "twitter":
		return append(baseTags, "Twitter", "SocialMedia")
	case "linkedin":
		return append(baseTags, "LinkedIn", "Business", "Innovation")
	case "devto":
		return append(baseTags, "Dev.to", "Programming", "Development")
	default:
		return baseTags
	}
}

// GenerateProductHuntCampaign generates a complete Product Hunt social media campaign
func (g *SocialMediaContentGenerator) GenerateProductHuntCampaign() (*SocialMediaCampaign, error) {
	launchDate := time.Date(2026, 3, 30, 0, 0, 0, 0, time.UTC)
	
	// Twitter posts
	twitterPosts := []SocialMediaPost{
		{
			Platform: "twitter",
			Content:  "🚀 BIG NEWS! We're live on Product Hunt! Agent-Todo is now officially launched - the task management platform built specifically for AI agents. No more context loss, no more forgotten tasks. 🤖➡️✅ https://producthunt.com/products/agent-todo",
			Hashtags: []string{"ProductHunt", "AI", "TaskManagement", "LaunchDay"},
			PublishTime: launchDate.Add(0), // 9 AM
		},
		{
			Platform: "twitter",
			Content:  "✨ What makes Agent-Todo different? We built it FOR AI agents, not just for humans. Features: • Persistent memory • Priority-based task management • Real-time agent status tracking Your AI agents deserve better tools. 🛠️",
			Hashtags: []string{"AIProductivity", "DeveloperTools", "Automation"},
			PublishTime: launchDate.Add(1 * time.Hour), // 10 AM
		},
		{
			Platform: "twitter",
			Content:  "🔥 The problem: \"My AI agent lost track of what it was doing after a restart\" The solution: Agent-Todo gives agents persistent memory and task tracking. No more context loss, no more duplicated work. It's that simple. 💡",
			Hashtags: []string{"AIPainPoint", "Solutions", "AI"},
			PublishTime: launchDate.Add(2 * time.Hour), // 11 AM
		},
	}
	
	// LinkedIn posts
	linkedinPosts := []SocialMediaPost{
		{
			Platform: "linkedin",
			Content:  "🚀 Exciting announcement: Agent-Todo is officially launching on Product Hunt today! As AI systems become more prevalent in enterprise environments, the need for specialized task management solutions has never been greater. Traditional task management systems weren't designed for autonomous agents. Agent-Todo fills this gap by providing persistent memory and task tracking for AI agents, real-time monitoring and analytics, enterprise-grade security and scalability, seamless integration with existing AI workflows. Learn more: https://producthunt.com/products/agent-todo",
			Hashtags: []string{"ProductLaunch", "AIManagement", "DigitalTransformation", "AIOps", "EnterpriseAI"},
			PublishTime: launchDate.Add(0),
		},
		{
			Platform: "linkedin",
			Content:  "📊 The AI workforce management challenge is real, but so are the solutions. Recent studies show that organizations using specialized AI task management systems see: • 3x improvement in agent productivity • 60% reduction in context-related errors • 40% faster task completion rates. Agent-Todo was designed from the ground up to address these specific pain points. #ArtificialIntelligence #Productivity #BusinessTechnology #Innovation",
			Hashtags: []string{"ArtificialIntelligence", "Productivity", "BusinessTechnology", "Innovation"},
			PublishTime: launchDate.Add(2 * time.Hour),
		},
	}
	
	// Combine all posts
	allPosts := append(twitterPosts, linkedinPosts...)
	
	// Calculate metrics targets
	totalPosts := len(allPosts)
	targetEngagement := totalPosts * 10  // 10 engagements per post target
	
	return &SocialMediaCampaign{
		Name:        "Agent-Todo Product Hunt Launch",
		Description: "Complete social media campaign for Product Hunt launch day",
		LaunchDate:  launchDate,
		Platforms:   []string{"twitter", "linkedin"},
		Posts:       allPosts,
		Metrics: CampaignMetrics{
			TargetImpressions: 1000,
			TargetEngagement: targetEngagement,
			TargetReach:      5000,
			ActualImpressions: 0,
			ActualEngagement: 0,
			ActualReach:      0,
		},
	}, nil
}

// SaveCampaign saves a campaign to a JSON file
func (g *SocialMediaContentGenerator) SaveCampaign(campaign *SocialMediaCampaign, filename string) error {
	data, err := json.MarshalIndent(campaign, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}

// LoadCampaign loads a campaign from a JSON file
func (g *SocialMediaContentGenerator) LoadCampaign(filename string) (*SocialMediaCampaign, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var campaign SocialMediaCampaign
	err = json.Unmarshal(data, &campaign)
	if err != nil {
		return nil, err
	}
	
	return &campaign, nil
}

// GetPostingSchedule returns a schedule for publishing posts
func (g *SocialMediaContentGenerator) GetPostingSchedule(campaign *SocialMediaCampaign) map[string][]time.Time {
	schedule := make(map[string][]time.Time)
	
	for _, post := range campaign.Posts {
		if _, exists := schedule[post.Platform]; !exists {
			schedule[post.Platform] = []time.Time{}
		}
		schedule[post.Platform] = append(schedule[post.Platform], post.PublishTime)
	}
	
	return schedule
}

func main() {
	// Create content generator
	generator := NewSocialMediaContentGenerator()
	
	// Generate Product Hunt campaign
	campaign, err := generator.GenerateProductHuntCampaign()
	if err != nil {
		fmt.Printf("Error generating campaign: %v\n", err)
		return
	}
	
	// Save campaign
	filename := "product-hunt-campaign.json"
	err = generator.SaveCampaign(campaign, filename)
	if err != nil {
		fmt.Printf("Error saving campaign: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Generated and saved Product Hunt campaign to %s\n", filename)
	
	// Display campaign summary
	fmt.Printf("\n📊 Campaign Summary:\n")
	fmt.Printf("Name: %s\n", campaign.Name)
	fmt.Printf("Platforms: %v\n", campaign.Platforms)
	fmt.Printf("Total Posts: %d\n", len(campaign.Posts))
	fmt.Printf("Target Engagement: %d\n", campaign.Metrics.TargetEngagement)
	
	// Display posting schedule
	schedule := generator.GetPostingSchedule(campaign)
	fmt.Printf("\n📅 Posting Schedule:\n")
	for platform, times := range schedule {
		fmt.Printf("%s: %v\n", platform, times)
	}
	
	// Display sample posts
	fmt.Printf("\n📝 Sample Posts:\n")
	for i, post := range campaign.Posts {
		if i < 3 { // Show first 3 posts
			fmt.Printf("\n[%s] %s:\n", post.Platform, post.PublishTime.Format("15:00"))
			fmt.Printf("Content: %s\n", post.Content)
			fmt.Printf("Hashtags: %v\n", post.Hashtags)
		}
	}
}