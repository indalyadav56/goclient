package main

import (
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	// Create a new client with basic configuration
	client := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
		GlobalHeaders: map[string]string{
			"User-Agent": "GoClient-Example/1.0",
		},
	})

	// Example 1: Simple GET request
	fmt.Println("=== Simple GET Request ===")
	var post map[string]interface{}
	err := client.Get("/posts/1").
		SetHeader("Accept", "application/json").
		Into(&post)

	if err != nil {
		log.Printf("GET request failed: %v", err)
	} else {
		fmt.Printf("Post: %+v\n", post)
	}

	// Example 2: POST request with body
	fmt.Println("\n=== POST Request ===")
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
		fmt.Printf("Created Post: %+v\n", createdPost)
	}

	// Example 3: GET with query parameters
	fmt.Println("\n=== GET with Query Parameters ===")
	var posts []map[string]interface{}
	err = client.Get("/posts").
		SetQueryParam("userId", "1").
		SetQueryParam("_limit", "3").
		Into(&posts)

	if err != nil {
		log.Printf("GET with params failed: %v", err)
	} else {
		fmt.Printf("Found %d posts for user 1\n", len(posts))
		for i, post := range posts {
			fmt.Printf("  Post %d: %s\n", i+1, post["title"])
		}
	}

	// Example 4: Using Result() method for more control
	fmt.Println("\n=== Using Result() Method ===")
	resp, err := client.Get("/posts/1").
		SetHeader("Accept", "application/json").
		Result()

	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		fmt.Printf("Status Code: %d\n", resp.StatusCode)
		fmt.Printf("Content-Type: %s\n", resp.Headers.Get("Content-Type"))
		fmt.Printf("Response Body: %s\n", string(resp.Body))
	}

	// Example 5: Error handling with custom error type
	fmt.Println("\n=== Error Handling ===")
	var errorResponse struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}

	var result map[string]interface{}
	err = client.Get("/posts/999999"). // Non-existent post
						SetError(&errorResponse).
						Into(&result)

	if err != nil {
		if reqErr, ok := err.(*goclient.RequestError); ok {
			fmt.Printf("Request failed with status %d\n", reqErr.StatusCode)
			fmt.Printf("Error response: %+v\n", errorResponse)
		} else {
			fmt.Printf("Other error: %v\n", err)
		}
	}
}
