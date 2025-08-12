package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	// Create client
	client := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
	})

	ctx := context.Background()


	// Example 1: Basic Batch Request
	fmt.Println("=== Basic Batch Request ===")
	batch := client.Batch()

	// Add multiple requests to the batch
	batch.Add(client.Get("/posts/1"))
	batch.Add(client.Get("/posts/2"))
	batch.Add(client.Get("/posts/3"))
	batch.Add(client.Get("/users/1"))
	batch.Add(client.Get("/users/2"))

	// Execute all requests concurrently
	start := time.Now()
	responses, errors := batch.Execute(ctx)
	duration := time.Since(start)

	fmt.Printf("Executed %d requests in %v\n", len(responses), duration)

	for i, resp := range responses {
		if errors[i] != nil {
			log.Printf("Request %d failed: %v", i+1, errors[i])
			continue
		}
		fmt.Printf("Request %d: Status %d, Body length: %d bytes\n", 
			i+1, resp.StatusCode, len(resp.Body))
	}

	// Example 2: Batch with Different HTTP Methods
	fmt.Println("\n=== Mixed HTTP Methods Batch ===")
	mixedBatch := client.Batch()

	// Add different types of requests
	mixedBatch.Add(client.Get("/posts/1"))
	mixedBatch.Add(client.Post("/posts").SetBody(map[string]interface{}{
		"title":  "New Post",
		"body":   "This is a new post",
		"userId": 1,
	}))
	mixedBatch.Add(client.Put("/posts/1").SetBody(map[string]interface{}{
		"id":     1,
		"title":  "Updated Post",
		"body":   "This is an updated post",
		"userId": 1,
	}))
	mixedBatch.Add(client.Delete("/posts/1"))

	responses, errors = mixedBatch.Execute(ctx)

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	for i, resp := range responses {
		if errors[i] != nil {
			log.Printf("%s request failed: %v", methods[i], errors[i])
			continue
		}
		fmt.Printf("%s request: Status %d\n", methods[i], resp.StatusCode)
	}

	// Example 3: Batch with Error Handling
	fmt.Println("\n=== Batch with Error Handling ===")
	errorBatch := client.Batch()

	// Add some requests that will succeed and some that will fail
	errorBatch.Add(client.Get("/posts/1"))           // Should succeed
	errorBatch.Add(client.Get("/posts/999999"))      // Should fail (404)
	errorBatch.Add(client.Get("/invalid-endpoint"))  // Should fail (404)
	errorBatch.Add(client.Get("/users/1"))           // Should succeed

	responses, errors = errorBatch.Execute(ctx)

	successCount := 0
	errorCount := 0

	for i, resp := range responses {
		if errors[i] != nil {
			errorCount++
			if reqErr, ok := errors[i].(*goclient.RequestError); ok {
				fmt.Printf("Request %d failed with status %d: %v\n", 
					i+1, reqErr.StatusCode, reqErr)
			} else {
				fmt.Printf("Request %d failed: %v\n", i+1, errors[i])
			}
			continue
		}
		successCount++
		fmt.Printf("Request %d succeeded: Status %d\n", i+1, resp.StatusCode)
	}

	fmt.Printf("Summary: %d successful, %d failed\n", successCount, errorCount)

	// Example 4: Large Batch Request
	fmt.Println("\n=== Large Batch Request ===")
	largeBatch := client.Batch()

	// Add many requests
	for i := 1; i <= 10; i++ {
		largeBatch.Add(client.Get(fmt.Sprintf("/posts/%d", i)))
	}

	start = time.Now()
	responses, errors = largeBatch.Execute(ctx)
	duration = time.Since(start)

	fmt.Printf("Executed %d requests in %v (avg: %v per request)\n", 
		len(responses), duration, duration/time.Duration(len(responses)))

	// Count successes and failures
	successCount = 0
	errorCount = 0
	for _, err := range errors {
		if err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	fmt.Printf("Results: %d successful, %d failed\n", successCount, errorCount)
}
