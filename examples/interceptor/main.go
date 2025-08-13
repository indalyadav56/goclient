package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/indalyadav56/goclient"
)

// CustomInterceptor demonstrates how to create a custom interceptor
type CustomInterceptor struct {
	Next http.RoundTripper
}

func (c *CustomInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("🚀 Custom Interceptor: Making %s request to %s\n", req.Method, req.URL.String())

	req.Header.Set("X-Custom-Client", "GoClient-Example")
	req.Header.Set("X-Request-Time", time.Now().Format(time.RFC3339))

	// Measure request duration
	start := time.Now()

	// Execute the request
	resp, err := c.Next.RoundTrip(req)

	duration := time.Since(start)

	// Add custom logic after request
	if err != nil {
		fmt.Printf("❌ Custom Interceptor: Request failed after %v: %v\n", duration, err)
	} else {
		fmt.Printf("✅ Custom Interceptor: Request completed in %v with status %d\n",
			duration, resp.StatusCode)
	}

	return resp, err
}

// AuthInterceptor demonstrates authentication interceptor
type AuthInterceptor struct {
	Next  http.RoundTripper
	Token string
}

func (a *AuthInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Printf("🔐 Auth Interceptor: Adding authentication to %s %s\n",
		req.Method, req.URL.Path)

	// Add authentication header
	req.Header.Set("Authorization", a.Token)

	return a.Next.RoundTrip(req)
}

// RetryInterceptor demonstrates retry logic
type RetryInterceptor struct {
	Next       http.RoundTripper
	MaxRetries int
	RetryDelay time.Duration
}

func (r *RetryInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		if attempt > 0 {
			fmt.Printf("🔄 Retry attempt %d/%d after %v\n", attempt, r.MaxRetries, r.RetryDelay)
			time.Sleep(r.RetryDelay)
		}

		resp, err := r.Next.RoundTrip(req)

		// If successful or non-retryable error, return immediately
		if err == nil || !r.shouldRetry(resp, err) {
			return resp, err
		}

		lastErr = err
		fmt.Printf("⚠️  Request failed (attempt %d): %v\n", attempt+1, err)
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", r.MaxRetries, lastErr)
}

func (r *RetryInterceptor) shouldRetry(resp *http.Response, err error) bool {
	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on 5xx status codes
	if resp != nil && resp.StatusCode >= 500 {
		return true
	}

	return false
}

func main() {

	// Example 1: Custom Interceptor
	fmt.Println("=== Custom Interceptor ===")
	customClient := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
		Interceptor: &CustomInterceptor{
			Next: http.DefaultTransport,
		},
	})

	var post map[string]interface{}
	err := customClient.Get("/posts/1").
		SetHeader("Accept", "application/json").
		Into(&post)

	if err != nil {
		log.Printf("Custom interceptor request failed: %v", err)
	} else {
		fmt.Printf("Post title: %s\n\n", post["title"])
	}

	// Example 2: Retry Interceptor
	fmt.Println("=== Retry Interceptor ===")
	retryClient := goclient.New(goclient.Config{
		BaseURL: "https://httpbin.org",
		Timeout: 10 * time.Second,
		Interceptor: &RetryInterceptor{
			Next:       http.DefaultTransport,
			MaxRetries: 3,
			RetryDelay: 1 * time.Second,
		},
	})

	// This should succeed immediately
	var successResponse map[string]interface{}
	err = retryClient.Get("/get").
		SetQueryParam("test", "retry-success").
		Into(&successResponse)

	if err != nil {
		log.Printf("Retry success request failed: %v", err)
	} else {
		fmt.Printf("Success response received\n")
	}

	// This should trigger retries (500 status)
	fmt.Println("\n--- Testing retry on 500 status ---")
	var errorResponse map[string]interface{}
	err = retryClient.Get("/status/500").
		Into(&errorResponse)

	if err != nil {
		fmt.Printf("Expected retry failure: %v\n\n", err)
	}

	// Example 3: Chained Interceptors
	fmt.Println("=== Chained Interceptors ===")
	chainedClient := goclient.New(goclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
		Interceptor: &CustomInterceptor{
			Next: &RetryInterceptor{
				Next:       http.DefaultTransport,
				MaxRetries: 2,
				RetryDelay: 500 * time.Millisecond,
			},
		},
	})

	var chainedResponse map[string]interface{}
	err = chainedClient.Get("/posts/2").
		Into(&chainedResponse)

	if err != nil {
		log.Printf("Chained interceptor request failed: %v", err)
	} else {
		fmt.Printf("Chained request successful: %s\n\n", chainedResponse["title"])
	}

	// Example 4: Authentication Interceptor
	fmt.Println("=== Authentication Interceptor ===")

	authInterceptor := &AuthInterceptor{
		Next:  http.DefaultTransport,
		Token: "Bearer my-secret-token",
	}

	authClient := goclient.New(goclient.Config{
		BaseURL:     "https://httpbin.org",
		Timeout:     30 * time.Second,
		Interceptor: authInterceptor,
	})

	var authResponse map[string]interface{}
	err = authClient.Get("/bearer").
		Into(&authResponse)

	if err != nil {
		log.Printf("Auth interceptor request failed: %v", err)
	} else {
		fmt.Printf("Auth response: %+v\n\n", authResponse)
	}

	// // Example 5: Combined with Logging
	// fmt.Println("=== Custom Interceptor + Logging ===")
	// logger := goclient.NewStandardLogger()

	// // Create a logging interceptor first
	// loggingInterceptor := goclient.NewWithOptions(
	// 	goclient.WithBaseURL("https://jsonplaceholder.typicode.com"),
	// 	goclient.WithTimeout(30*time.Second),
	// 	goclient.WithLoggingInterceptor(logger, goclient.LoggingOptions{
	// 		LogRequestBody: false,
	// 		LogHeaders:     true,
	// 		MaxBodySize:    500,
	// 	}),
	// )

	// // Use the logging client directly for this example
	// combinedClient := loggingInterceptor

	// var combinedResponse map[string]interface{}
	// err = combinedClient.Get("/posts/3").
	// 	SetHeader("X-Example", "combined-interceptors").
	// 	Into(&combinedResponse)

	// if err != nil {
	// 	log.Printf("Combined interceptor request failed: %v", err)
	// } else {
	// 	fmt.Printf("Combined request successful\n")
	// }

	// fmt.Println("All interceptor examples completed!")
}
