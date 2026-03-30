package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Test script for Beta Feedback System

const BASE_URL = "http://localhost:8080"

func main() {
	fmt.Println("🧪 Testing Beta Feedback System")
	fmt.Println("=" + string(make([]byte, 40)))

	// Test 1: Submit public feedback
	fmt.Println("\n1. Testing feedback submission (public endpoint)...")
	testSubmitFeedback()

	// Test 2: List feedback (requires auth - will fail without token)
	fmt.Println("\n2. Testing feedback list (protected endpoint)...")
	testListFeedback()

	// Test 3: Get feedback stats (requires auth - will fail without token)
	fmt.Println("\n3. Testing feedback stats (protected endpoint)...")
	testFeedbackStats()
}

func testSubmitFeedback() {
	feedback := map[string]interface{}{
		"feedback_type": "bug",
		"title":         "Test Bug Report",
		"description":   "This is a test bug report from the automated test script",
		"tester_name":   "Test User",
		"tester_email":  "test@example.com",
		"priority":      "high",
		"rating":        3,
		"page":          "https://todo.formatho.com/test",
	}

	jsonData, _ := json.Marshal(feedback)

	resp, err := http.Post(BASE_URL+"/feedback", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("✅ Success: Feedback submitted (Status: %d)\n", resp.StatusCode)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("   ID: %s\n", result["id"])
		fmt.Printf("   Status: %s\n", result["status"])
	} else {
		fmt.Printf("❌ Failed: Status %d\n%s\n", resp.StatusCode, string(body))
	}
}

func testListFeedback() {
	// This will fail without auth token, which is expected
	resp, err := http.Get(BASE_URL + "/feedback")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("✅ Expected: Authentication required (Status: %d)\n", resp.StatusCode)
	} else if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Success: Feedback list retrieved (Status: %d)\n", resp.StatusCode)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("   Total: %v\n", result["total"])
		fmt.Printf("   Count: %d\n", len(result["feedback"].([]interface{})))
	} else {
		fmt.Printf("❌ Unexpected status: %d\n%s\n", resp.StatusCode, string(body))
	}
}

func testFeedbackStats() {
	// This will fail without auth token, which is expected
	resp, err := http.Get(BASE_URL + "/feedback/stats")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		fmt.Printf("✅ Expected: Authentication required (Status: %d)\n", resp.StatusCode)
	} else if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Success: Stats retrieved (Status: %d)\n", resp.StatusCode)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		fmt.Printf("   Total: %v\n", result["total"])
		fmt.Printf("   Recent: %v\n", result["recent_count"])
		fmt.Printf("   Avg Rating: %.2f\n", result["average_rating"])
	} else {
		fmt.Printf("❌ Unexpected status: %d\n%s\n", resp.StatusCode, string(body))
	}
}

func printSeparator() {
	fmt.Println("\n" + string(make([]byte, 60)))
	time.Sleep(500 * time.Millisecond)
}
