package main

import (
	"fmt"
	"log"
	"time"

	"github.com/indalyadav56/goclient"
)

func main() {

	// Example 1: Bearer Token Authentication
	fmt.Println("=== Bearer Token Authentication ===")
	client := goclient.New(goclient.Config{
		BaseURL: "https://api.github.com",
		Timeout: 30 * time.Second,
	}).SetBearerToken("your-github-token-here")

	var user map[string]interface{}
	err := client.Get("/user").
		SetHeader("Accept", "application/vnd.github.v3+json").
		Into(&user)

	if err != nil {
		log.Printf("GitHub API request failed: %v", err)
	} else {
		fmt.Printf("GitHub User: %s\n", user["login"])
	}

	// Example 2: Basic Authentication
	fmt.Println("\n=== Basic Authentication ===")
	basicClient := goclient.New(goclient.Config{
		BaseURL: "https://httpbin.org",
		Timeout: 30 * time.Second,
	}).WithBasicAuth("testuser", "testpass")

	var authResult map[string]interface{}
	err = basicClient.Get("/basic-auth/testuser/testpass").
		Into(&authResult)

	if err != nil {
		log.Printf("Basic auth request failed: %v", err)
	} else {
		fmt.Printf("Basic Auth Result: %+v\n", authResult)
	}

	// Example 3: Custom Headers for API Key Authentication
	fmt.Println("\n=== API Key Authentication ===")
	apiKeyClient := goclient.New(goclient.Config{
		BaseURL: "https://api.openweathermap.org/data/2.5",
		Timeout: 30 * time.Second,
		GlobalHeaders: map[string]string{
			"X-API-Key": "your-api-key-here",
		},
	})

	var weather map[string]interface{}
	err = apiKeyClient.Get("/weather").
		SetQueryParam("q", "London").
		SetQueryParam("appid", "your-api-key-here").
		Into(&weather)

	if err != nil {
		log.Printf("Weather API request failed: %v", err)
	} else {
		fmt.Printf("Weather in London: %+v\n", weather)
	}

	// Example 4: Dynamic Authentication (changing tokens)
	fmt.Println("\n=== Dynamic Authentication ===")
	dynamicClient := goclient.New(goclient.Config{
		BaseURL: "https://httpbin.org",
		Timeout: 30 * time.Second,
	})

	// First request with one token
	firstToken := "token-123"
	var firstResult map[string]interface{}
	err = dynamicClient.SetBearerToken(firstToken).Get("/bearer").
		Into(&firstResult)

	if err != nil {
		log.Printf("First token request failed: %v", err)
	} else {
		fmt.Printf("First token result: %+v\n", firstResult)
	}

	// Second request with different token
	secondToken := "token-456"
	var secondResult map[string]interface{}
	err = dynamicClient.SetBearerToken(secondToken).Get("/bearer").
		Into(&secondResult)

	if err != nil {
		log.Printf("Second token request failed: %v", err)
	} else {
		fmt.Printf("Second token result: %+v\n", secondResult)
	}
}
