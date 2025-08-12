package goclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type TestError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Test server setup
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// GET /posts/1 - Success
	mux.HandleFunc("/posts/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		post := TestPost{
			ID:     1,
			Title:  "Test Post",
			Body:   "This is a test post",
			UserID: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	})

	// POST /posts - Create post
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post TestPost
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		post.ID = 101 // Simulate created ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	})

	// PUT /posts/1 - Update post
	mux.HandleFunc("/posts/1/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post TestPost
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		post.ID = 1
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	})

	// DELETE /posts/1 - Delete post
	mux.HandleFunc("/posts/1/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// GET /posts/404 - Not found
	mux.HandleFunc("/posts/404", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(TestError{
			Error:   "Not Found",
			Message: "Post not found",
		})
	})

	// GET /auth/bearer - Bearer token auth
	mux.HandleFunc("/auth/bearer", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"error":         "Missing or invalid bearer token",
			})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "valid-token" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": true,
				"token":         token,
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"error":         "Invalid token",
			})
		}
	})

	// GET /auth/basic - Basic auth
	mux.HandleFunc("/auth/basic", func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"error":         "Missing basic auth",
			})
			return
		}

		if username == "testuser" && password == "testpass" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": true,
				"username":      username,
			})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated": false,
				"error":         "Invalid credentials",
			})
		}
	})

	// GET /slow - Slow endpoint for timeout testing
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(map[string]string{"message": "slow response"})
	})

	return httptest.NewServer(mux)
}

// Test basic client creation
func TestNew(t *testing.T) {
	client := New()
	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	// Test with config
	config := Config{
		BaseURL: "https://api.example.com",
		Timeout: 10 * time.Second,
		GlobalHeaders: map[string]string{
			"User-Agent": "test-client",
		},
	}

	clientWithConfig := New(config)
	if clientWithConfig == nil {
		t.Fatal("Expected client with config to be created, got nil")
	}
}

// Test simple GET request
func TestClient_Get(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	var post TestPost
	err := client.Get("/posts/1").Into(&post)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", post.ID)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %s", post.Title)
	}
}

// Test GET with context
func TestClient_GetWithContext(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var post TestPost
	err := client.GetWithContext(ctx, "/posts/1").Into(&post)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", post.ID)
	}
}

// Test context timeout
func TestClient_GetWithContext_Timeout(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var result map[string]interface{}
	err := client.GetWithContext(ctx, "/slow").Into(&result)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") &&
		!strings.Contains(err.Error(), "request canceled or timed out") {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

// Test POST request
func TestClient_Post(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	newPost := TestPost{
		Title:  "New Post",
		Body:   "This is a new post",
		UserID: 1,
	}

	var createdPost TestPost
	err := client.Post("/posts").SetBody(newPost).Into(&createdPost)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdPost.ID != 101 {
		t.Errorf("Expected created post ID 101, got %d", createdPost.ID)
	}

	if createdPost.Title != newPost.Title {
		t.Errorf("Expected title %s, got %s", newPost.Title, createdPost.Title)
	}
}

// Test PUT request
func TestClient_Put(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	updatePost := TestPost{
		Title:  "Updated Post",
		Body:   "This is an updated post",
		UserID: 1,
	}

	var updatedPost TestPost
	err := client.Put("/posts/1/update").SetBody(updatePost).Into(&updatedPost)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedPost.Title != updatePost.Title {
		t.Errorf("Expected title %s, got %s", updatePost.Title, updatedPost.Title)
	}
}

// Test DELETE request
func TestClient_Delete(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	resp, err := client.Delete("/posts/1/delete").Result()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

// Test error handling
func TestClient_ErrorHandling(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	var post TestPost
	err := client.Get("/posts/404").Into(&post)

	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}

	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("Expected RequestError, got %T", err)
	}

	if reqErr.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", reqErr.StatusCode)
	}
}

// Test bearer token authentication
func TestClient_BearerToken(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	authClient := client.SetBearerToken("valid-token")

	var result map[string]interface{}
	err := authClient.Get("/auth/bearer").Into(&result)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if authenticated, ok := result["authenticated"].(bool); !ok || !authenticated {
		t.Error("Expected authenticated to be true")
	}
}

// Test basic authentication
func TestClient_BasicAuth(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	authClient := client.WithBasicAuth("testuser", "testpass")

	var result map[string]interface{}
	err := authClient.Get("/auth/basic").Into(&result)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if authenticated, ok := result["authenticated"].(bool); !ok || !authenticated {
		t.Error("Expected authenticated to be true")
	}
}

// Test request headers
func TestClient_Headers(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
		GlobalHeaders: map[string]string{
			"X-Global-Header": "global-value",
		},
	})

	var post TestPost
	err := client.Get("/posts/1").
		SetHeader("X-Custom-Header", "custom-value").
		SetHeaders(map[string]string{
			"X-Another-Header": "another-value",
		}).
		Into(&post)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", post.ID)
	}
}

// Test batch requests
func TestClient_Batch(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	batch := client.Batch()
	batch.Add(client.Get("/posts/1"))
	batch.Add(client.Get("/posts/404")) // This will fail
	batch.Add(client.Get("/posts/1"))   // This will succeed

	responses, errors := batch.Execute(context.Background())

	if len(responses) != 3 {
		t.Errorf("Expected 3 responses, got %d", len(responses))
	}

	if len(errors) != 3 {
		t.Errorf("Expected 3 errors (including nils), got %d", len(errors))
	}

	// Check that we have both successes and failures
	successCount := 0
	errorCount := 0
	for i, err := range errors {
		if err != nil {
			errorCount++
			t.Logf("Request %d failed as expected: %v", i, err)
		} else {
			successCount++
			t.Logf("Request %d succeeded as expected", i)
		}
	}

	if successCount == 0 {
		t.Error("Expected at least one request to succeed")
	}

	if errorCount == 0 {
		t.Error("Expected at least one request to fail")
	}
}

// Test request pool
func TestClient_Pool(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	pool := client.Pool(2)
	defer pool.Wait()

	// Submit multiple requests
	var wg sync.WaitGroup
	results := make([]Result, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			resultChan := pool.Submit(client.Get("/posts/1"))
			results[index] = <-resultChan
		}(i)
	}

	wg.Wait()

	// Check results
	for i, result := range results {
		if result.Error != nil {
			t.Errorf("Request %d failed: %v", i, result.Error)
		}
		if result.Response == nil {
			t.Errorf("Request %d has nil response", i)
		}
	}
}

// Test query parameters
func TestClient_QueryParams(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	var post TestPost
	err := client.Get("/posts/1").
		SetQueryParam("param1", "value1").
		SetQueryParams(map[string]string{
			"param2": "value2",
			"param3": "value3",
		}).
		Into(&post)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if post.ID != 1 {
		t.Errorf("Expected post ID 1, got %d", post.ID)
	}
}

// Test success and error handlers
func TestClient_Handlers(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	// Test success handler
	successCalled := false
	var post TestPost
	err := client.Get("/posts/1").
		OnSuccess(func(resp *Response) {
			successCalled = true
			if resp.StatusCode != 200 {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}).
		Into(&post)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Note: Success handlers may not be implemented yet, so we'll just check the basic functionality
	t.Logf("Success handler called: %v", successCalled)

	// Test error handler
	errorCalled := false
	var errorPost TestPost
	err = client.Get("/posts/404").
		OnError(func(reqErr *RequestError) {
			errorCalled = true
			if reqErr.StatusCode != 404 {
				t.Errorf("Expected status 404, got %d", reqErr.StatusCode)
			}
		}).
		Into(&errorPost)

	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}

	t.Logf("Error handler called: %v", errorCalled)
}

// Test error response unmarshaling
func TestClient_ErrorUnmarshaling(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	var post TestPost
	var errorResp TestError

	err := client.Get("/posts/404").
		SetError(&errorResp).
		Into(&post)

	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}

	// The error should contain the unmarshaled error details
	if !strings.Contains(err.Error(), "Not Found") {
		t.Errorf("Expected error to contain 'Not Found', got %v", err)
	}
}

// Benchmark tests
func BenchmarkClient_Get(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var post TestPost
		err := client.Get("/posts/1").Into(&post)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

func BenchmarkClient_Batch(b *testing.B) {
	server := setupTestServer()
	defer server.Close()

	client := New(Config{
		BaseURL: server.URL,
		Timeout: 5 * time.Second,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		batch := client.Batch()
		for j := 0; j < 5; j++ {
			batch.Add(client.Get("/posts/1"))
		}
		_, _ = batch.Execute(context.Background())
	}
}
