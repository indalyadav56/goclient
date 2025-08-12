package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	fmt.Println("=== Package-Level Functions Demo ===")
	fmt.Println("Now you can use goclient.Get(), goclient.Post(), etc. directly!")
	fmt.Println()

	// Configure the default client (optional)
	goclient.SetDefaultClient(goclient.Config{
		Timeout: 30 * time.Second,
		GlobalHeaders: map[string]string{
			"User-Agent": "goclient-package-level-demo",
		},
	})

	// 1. DIRECT GET - No client creation needed!
	fmt.Println("1. Direct GET request:")
	var post map[string]interface{}
	err := goclient.Get("https://jsonplaceholder.typicode.com/posts/1").
		SetHeader("Accept", "application/json").
		Into(&post)

	if err != nil {
		log.Printf("GET failed: %v", err)
	} else {
		fmt.Printf("   ✅ GET: %s\n", post["title"])
	}

	// 2. DIRECT POST - No client creation needed!
	fmt.Println("\n2. Direct POST request:")
	newPost := map[string]interface{}{
		"title":  "Direct Package-Level Post",
		"body":   "Created using goclient.Post() directly!",
		"userId": 1,
	}

	var createdPost map[string]interface{}
	err = goclient.Post("https://jsonplaceholder.typicode.com/posts").
		SetBody(newPost).
		Into(&createdPost)

	if err != nil {
		log.Printf("POST failed: %v", err)
	} else {
		fmt.Printf("   ✅ POST: Created post ID %.0f\n", createdPost["id"])
	}

	// 3. DIRECT PUT - No client creation needed!
	fmt.Println("\n3. Direct PUT request:")
	updatePost := map[string]interface{}{
		"id":     1,
		"title":  "Updated via Package-Level",
		"body":   "Updated using goclient.Put() directly!",
		"userId": 1,
	}

	var updatedPost map[string]interface{}
	err = goclient.Put("https://jsonplaceholder.typicode.com/posts/1").
		SetBody(updatePost).
		Into(&updatedPost)

	if err != nil {
		log.Printf("PUT failed: %v", err)
	} else {
		fmt.Printf("   ✅ PUT: %s\n", updatedPost["title"])
	}

	// 4. DIRECT DELETE - No client creation needed!
	fmt.Println("\n4. Direct DELETE request:")
	resp, err := goclient.Delete("https://jsonplaceholder.typicode.com/posts/1").Result()

	if err != nil {
		log.Printf("DELETE failed: %v", err)
	} else {
		fmt.Printf("   ✅ DELETE: Status %d\n", resp.StatusCode)
	}

	// 5. CONTEXT-AWARE DIRECT REQUESTS
	fmt.Println("\n5. Context-aware direct requests:")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var contextPost map[string]interface{}
	err = goclient.GetWithContext(ctx, "https://jsonplaceholder.typicode.com/posts/2").
		Into(&contextPost)

	if err != nil {
		log.Printf("Context GET failed: %v", err)
	} else {
		fmt.Printf("   ✅ Context GET: %s\n", contextPost["title"])
	}

	// 6. DIRECT AUTHENTICATION
	fmt.Println("\n6. Direct authentication:")
	goclient.SetBearerToken("demo-token")

	var authResult map[string]interface{}
	err = goclient.Get("https://httpbin.org/bearer").Into(&authResult)
	if err != nil {
		log.Printf("Auth GET failed: %v", err)
	} else {
		fmt.Printf("   ✅ Auth GET: %v\n", authResult["authenticated"])
	}

	// 7. DIRECT BATCH REQUESTS
	fmt.Println("\n7. Direct batch requests:")
	batch := goclient.Batch()
	batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/1"))
	batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/2"))
	batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/3"))

	responses, _ := batch.Execute(context.Background())
	fmt.Printf("   ✅ Batch: %d responses completed\n", len(responses))

	// 8. DIRECT REQUEST POOL
	fmt.Println("\n8. Direct request pool:")
	pool := goclient.Pool(2)
	defer pool.Wait()

	resultChan1 := pool.Submit(goclient.Get("https://jsonplaceholder.typicode.com/posts/4"))
	resultChan2 := pool.Submit(goclient.Get("https://jsonplaceholder.typicode.com/posts/5"))

	result1 := <-resultChan1
	result2 := <-resultChan2

	if result1.Error == nil && result2.Error == nil {
		fmt.Printf("   ✅ Pool: Both requests completed successfully\n")
	}

	fmt.Println("\n🎉 Perfect! Package-level functions work beautifully!")
	fmt.Println("   • Simple usage: goclient.Get(url)")
	fmt.Println("   • Context usage: goclient.GetWithContext(ctx, url)")
	fmt.Println("   • All HTTP methods: Get, Post, Put, Patch, Delete")
	fmt.Println("   • Authentication: goclient.SetBearerToken(token)")
	fmt.Println("   • Batch: goclient.Batch()")
	fmt.Println("   • Pool: goclient.Pool(workers)")
	fmt.Println("   • No client creation needed - just import and use!")
}
