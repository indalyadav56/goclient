package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	fmt.Println("=== Context Flexibility Demo ===")
	fmt.Println("This demonstrates both simple and context-aware usage patterns!")
	fmt.Println()

	// Create client
	client := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
	})

	// 1. SIMPLE USAGE - No context needed (uses context.Background() internally)
	fmt.Println("1. Simple usage (no context required):")
	var post1 map[string]interface{}
	err := client.Get("/posts/1").
		SetHeader("Accept", "application/json").
		Into(&post1)

	if err != nil {
		log.Printf("Simple GET failed: %v", err)
	} else {
		fmt.Printf("   ✅ Simple GET: %s\n", post1["title"])
	}

	// 2. CONTEXT-AWARE USAGE - Explicit context control
	fmt.Println("\n2. Context-aware usage (explicit context control):")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var post2 map[string]interface{}
	err = client.GetWithContext(ctx, "/posts/2").
		SetHeader("Accept", "application/json").
		Into(&post2)

	if err != nil {
		log.Printf("Context GET failed: %v", err)
	} else {
		fmt.Printf("   ✅ Context GET: %s\n", post2["title"])
	}

	// 3. SIMPLE POST - No context needed
	fmt.Println("\n3. Simple POST (no context required):")
	newPost := map[string]interface{}{
		"title":  "Simple Post",
		"body":   "Created without explicit context",
		"userId": 1,
	}

	var createdPost1 map[string]interface{}
	err = client.Post("/posts").
		SetBody(newPost).
		Into(&createdPost1)

	if err != nil {
		log.Printf("Simple POST failed: %v", err)
	} else {
		fmt.Printf("   ✅ Simple POST: Created post ID %.0f\n", createdPost1["id"])
	}

	// 4. CONTEXT-AWARE POST - With timeout
	fmt.Println("\n4. Context-aware POST (with timeout):")
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shortCancel()

	newPost2 := map[string]interface{}{
		"title":  "Context Post",
		"body":   "Created with explicit context and timeout",
		"userId": 1,
	}

	var createdPost2 map[string]interface{}
	err = client.PostWithContext(shortCtx, "/posts").
		SetBody(newPost2).
		Into(&createdPost2)

	if err != nil {
		log.Printf("Context POST failed: %v", err)
	} else {
		fmt.Printf("   ✅ Context POST: Created post ID %.0f\n", createdPost2["id"])
	}

	// 5. MIXED USAGE - Batch with both patterns
	fmt.Println("\n5. Mixed usage in batch operations:")
	batch := client.Batch()

	// Add simple requests (no context)
	batch.Add(client.Get("/posts/3"))
	batch.Add(client.Get("/posts/4"))

	// Add context-aware requests
	batchCtx := context.Background()
	batch.Add(client.GetWithContext(batchCtx, "/posts/5"))

	responses, _ := batch.Execute(context.Background())
	fmt.Printf("   ✅ Batch completed: %d responses\n", len(responses))

	// 6. AUTHENTICATION - Works with both patterns
	fmt.Println("\n6. Authentication with both patterns:")
	authClient := client.SetBearerToken("demo-token")

	// Simple auth request
	var authResult1 map[string]interface{}
	err = authClient.Get("https://httpbin.org/bearer").Into(&authResult1)
	if err != nil {
		log.Printf("Simple auth failed: %v", err)
	} else {
		fmt.Printf("   ✅ Simple auth: %v\n", authResult1["authenticated"])
	}

	// Context-aware auth request
	authCtx := context.Background()
	var authResult2 map[string]interface{}
	err = authClient.GetWithContext(authCtx, "https://httpbin.org/bearer").Into(&authResult2)
	if err != nil {
		log.Printf("Context auth failed: %v", err)
	} else {
		fmt.Printf("   ✅ Context auth: %v\n", authResult2["authenticated"])
	}

	fmt.Println("\n🎉 Perfect! Your API now provides maximum flexibility:")
	fmt.Println("   • Simple usage: client.Get(\"/path\")")
	fmt.Println("   • Context control: client.GetWithContext(ctx, \"/path\")")
	fmt.Println("   • Users can choose based on their needs!")
}
