package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {
	client := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
	})

	// Example 1: Basic Request Pool
	fmt.Println("=== Basic Request Pool ===")
	pool := client.Pool(5) // 5 workers

	// Submit multiple requests
	var results []<-chan goclient.Result
	for i := 1; i <= 10; i++ {
		resultChan := pool.Submit(client.Get(fmt.Sprintf("/posts/%d", i)))
		results = append(results, resultChan)
	}

	// Collect results
	start := time.Now()
	successCount := 0
	errorCount := 0

	for i, resultChan := range results {
		result := <-resultChan
		if result.Error != nil {
			errorCount++
			log.Printf("Request %d failed: %v", i+1, result.Error)
		} else {
			successCount++
			fmt.Printf("Request %d: Status %d\n", i+1, result.Response.StatusCode)
		}
	}

	duration := time.Since(start)
	fmt.Printf("Processed %d requests in %v (%d successful, %d failed)\n",
		len(results), duration, successCount, errorCount)

	pool.Wait() // Wait for all workers to finish

	// Example 2: High-Throughput Processing
	fmt.Println("\n=== High-Throughput Processing ===")
	highThroughputPool := client.Pool(10) // 10 workers

	const numRequests = 50
	resultChannels := make([]<-chan goclient.Result, numRequests)

	// Submit many requests quickly
	start = time.Now()
	for i := 0; i < numRequests; i++ {
		postID := (i % 100) + 1 // Cycle through posts 1-100
		resultChannels[i] = highThroughputPool.Submit(
			client.Get(fmt.Sprintf("/posts/%d", postID)),
		)
	}
	submitDuration := time.Since(start)

	// Process results as they come in
	var wg sync.WaitGroup
	successCount = 0
	errorCount = 0
	var mu sync.Mutex

	for i, resultChan := range resultChannels {
		wg.Add(1)
		go func(index int, ch <-chan goclient.Result) {
			defer wg.Done()
			result := <-ch

			mu.Lock()
			if result.Error != nil {
				errorCount++
			} else {
				successCount++
			}
			mu.Unlock()

			if index%10 == 0 { // Log every 10th request
				if result.Error != nil {
					log.Printf("Request %d failed: %v", index+1, result.Error)
				} else {
					fmt.Printf("Request %d completed: Status %d\n",
						index+1, result.Response.StatusCode)
				}
			}
		}(i, resultChan)
	}

	wg.Wait()
	highThroughputPool.Wait()
	totalDuration := time.Since(start)

	fmt.Printf("High-throughput results:\n")
	fmt.Printf("  Submitted %d requests in %v\n", numRequests, submitDuration)
	fmt.Printf("  Completed all requests in %v\n", totalDuration)
	fmt.Printf("  Average: %v per request\n", totalDuration/numRequests)
	fmt.Printf("  Success: %d, Errors: %d\n", successCount, errorCount)

	// Example 3: Mixed Request Types in Pool
	fmt.Println("\n=== Mixed Request Types ===")
	mixedPool := client.Pool(3) // 3 workers

	// Submit different types of requests
	requests := []struct {
		name string
		req  goclient.RequestBuilder
	}{
		{"GET Post 1", client.Get("/posts/1")},
		{"GET Post 2", client.Get("/posts/2")},
		{"GET User 1", client.Get("/users/1")},
		{"POST New Post", client.Post("/posts").SetBody(map[string]interface{}{
			"title":  "Pool Test Post",
			"body":   "Created via request pool",
			"userId": 1,
		})},
		{"PUT Update Post", client.Put("/posts/1").SetBody(map[string]interface{}{
			"id":     1,
			"title":  "Updated via Pool",
			"body":   "Updated using request pool",
			"userId": 1,
		})},
	}

	var mixedResults []<-chan goclient.Result
	for _, req := range requests {
		resultChan := mixedPool.Submit(req.req)
		mixedResults = append(mixedResults, resultChan)
	}

	// Process mixed results
	for i, resultChan := range mixedResults {
		result := <-resultChan
		if result.Error != nil {
			log.Printf("%s failed: %v", requests[i].name, result.Error)
		} else {
			fmt.Printf("%s: Status %d, Body length: %d bytes\n",
				requests[i].name, result.Response.StatusCode, len(result.Response.Body))
		}
	}

	mixedPool.Wait()

	// Example 4: Pool with Error Handling and Retries
	fmt.Println("\n=== Pool with Error Handling ===")
	errorPool := client.Pool(2) // 2 workers

	// Submit requests that might fail
	errorRequests := []string{"/posts/1", "/posts/999999", "/invalid", "/users/1"}
	var errorResults []<-chan goclient.Result

	for _, endpoint := range errorRequests {
		resultChan := errorPool.Submit(client.Get(endpoint))
		errorResults = append(errorResults, resultChan)
	}

	// Process with detailed error handling
	for i, resultChan := range errorResults {
		result := <-resultChan
		endpoint := errorRequests[i]

		if result.Error != nil {
			if reqErr, ok := result.Error.(*goclient.RequestError); ok {
				fmt.Printf("Request to %s failed with status %d\n",
					endpoint, reqErr.StatusCode)
			} else {
				fmt.Printf("Request to %s failed: %v\n", endpoint, result.Error)
			}
		} else {
			fmt.Printf("Request to %s succeeded: Status %d\n",
				endpoint, result.Response.StatusCode)
		}
	}

	errorPool.Wait()
	fmt.Println("All pool examples completed!")
}
