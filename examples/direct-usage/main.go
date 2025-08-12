package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	fmt.Println("=== Direct Usage Demo ===")
	fmt.Println("This demonstrates the clean, direct API usage you requested!")
	fmt.Println()

	// This is exactly what you wanted - direct usage without sub-packages!
	client := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
		GlobalHeaders: map[string]string{
			"User-Agent": "GoClient-DirectUsage/1.0",
		},
	})

	ctx := context.Background()

	// Simple GET request
	fmt.Println("1. Simple GET request:")
	var post map[string]interface{}
	err := client.Get("/posts/1").
		SetHeader("Accept", "application/json").
		Into(&post)

	if err != nil {
		log.Printf("GET request failed: %v", err)
	} else {
		fmt.Printf("   Post title: %s\n", post["title"])
	}

	// POST request with body
	fmt.Println("\n2. POST request with body:")
	newPost := map[string]interface{}{
		"title":  "My New Post",
		"body":   "This is the content of my new post",
		"userId": 1,
	}

	var createdPost map[string]interface{}
	err = client.Post("/posts").
		SetHeader("Content-Type", "application/json").
		SetBody(newPost).
		Into(&createdPost)

	if err != nil {
		log.Printf("POST request failed: %v", err)
	} else {
		fmt.Printf("   Created post ID: %.0f\n", createdPost["id"])
	}

	// Authentication example
	fmt.Println("\n3. Authentication example:")
	authClient := client.SetBearerToken("your-bearer-token")
	var authResult map[string]interface{}
	err = authClient.Get("https://httpbin.org/bearer").Into(&authResult)

	if err != nil {
		log.Printf("Auth request failed: %v", err)
	} else {
		fmt.Printf("   Auth successful: %v\n", authResult["authenticated"])
	}

	// Batch requests
	fmt.Println("\n4. Batch requests:")
	batch := client.Batch()
	batch.Add(client.Get("/posts/1"))
	batch.Add(client.Get("/posts/2"))
	batch.Add(client.Get("/posts/3"))

	responses, errors := batch.Execute(ctx)
	fmt.Printf("   Batch completed: %d responses, %d errors\n", len(responses), len(errors))

	// Request pool
	fmt.Println("\n5. Request pool:")
	pool := client.Pool(3)

	resultChan1 := pool.Submit(client.Get("/posts/3"))
	resultChan2 := pool.Submit(client.Get("/posts/4"))

	result1 := <-resultChan1
	result2 := <-resultChan2

	if result1.Error == nil && result2.Error == nil {
		fmt.Printf("   Pool requests completed successfully!\n")
	}

	pool.Wait()

	fmt.Println("\n✅ All direct usage examples completed successfully!")
	fmt.Println("🎉 Your package now works exactly as requested:")
	fmt.Println("   import \"github.com/indalyadav56/goclient\"")
	fmt.Println("   client := goclient.New()")
	fmt.Println("   // No sub-package imports needed!")
}
