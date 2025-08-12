# GoClient - Advanced HTTP Client for Go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

GoClient is a comprehensive, feature-rich HTTP client library for Go that provides a fluent API, advanced authentication, request pooling, batch operations, and sophisticated logging capabilities.

## Features

### 🚀 Core Features
- **Fluent API**: Intuitive request builder pattern
- **HTTP Methods**: GET, POST, PUT, PATCH, DELETE support
- **Authentication**: Bearer token and Basic authentication
- **Context Support**: Proper context handling for cancellation and timeouts
- **Error Handling**: Detailed error information with custom error types

### 🔧 Advanced Features
- **Batch Requests**: Execute multiple requests concurrently
- **Request Pooling**: Worker pool pattern for high-throughput scenarios
- **Interceptors**: Middleware pattern for extending functionality
- **Logging**: Structured logging with request/response details
- **Connection Management**: Configurable connection pooling
- **Object Pooling**: Efficient memory usage with sync.Pool

### 📊 Performance Features
- **Connection Pooling**: Configurable idle connections and timeouts
- **Keep-Alive Support**: Reuse TCP connections for better performance
- **Compression**: Optional request/response compression
- **Memory Efficient**: Object pooling to reduce GC pressure

## Installation

```bash
go get github.com/indalyadav56/goclient
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/indalyadav56/goclient"
)

func main() {
    // Create a new client
    client := goclient.New(goclient.Config{
        BaseURL: "https://api.example.com",
        Timeout: 30 * time.Second,
    })

    // Make a simple GET request
    var result map[string]interface{}
    err := client.Get(context.Background(), "/users/1").
        SetHeader("Accept", "application/json").
        Into(&result)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

### Using Functional Options

```go
client := goclient.NewWithOptions(
    goclient.WithBaseURL("https://api.example.com"),
    goclient.WithTimeout(30*time.Second),
    goclient.WithGlobalHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    }),
    goclient.WithMaxIdleConns(100),
)
```

## Advanced Usage

### Authentication

```go
// Bearer Token
client := goclient.New().SetBearerToken("your-token-here")

// Basic Auth
client := goclient.New().WithBasicAuth("username", "password")
```

### Request Building

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

var user User
err := client.Post(ctx, "/users").
    SetHeader("Content-Type", "application/json").
    SetBody(map[string]string{
        "name": "John Doe",
        "email": "john@example.com",
    }).
    SetQueryParam("include", "profile").
    Into(&user)
```

### Error Handling

```go
var user User
var apiError struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
}

err := client.Get(ctx, "/users/1").
    SetError(&apiError).
    Into(&user)

if err != nil {
    if reqErr, ok := err.(*goclient.RequestError); ok {
        fmt.Printf("Request failed: %d - %s\n", reqErr.StatusCode, apiError.Message)
    }
}
```

### Batch Requests

```go
batch := client.Batch()

// Add multiple requests to the batch
batch.Add(client.Get(ctx, "/users/1"))
batch.Add(client.Get(ctx, "/users/2"))
batch.Add(client.Get(ctx, "/users/3"))

// Execute all requests concurrently
responses, errors := batch.Execute(ctx)

for i, resp := range responses {
    if errors[i] != nil {
        fmt.Printf("Request %d failed: %v\n", i, errors[i])
        continue
    }
    fmt.Printf("Response %d: %d\n", i, resp.StatusCode)
}
```

### Request Pool (High Throughput)

```go
pool := client.Pool(10) // 10 workers

// Submit requests to the pool
resultChan := pool.Submit(client.Get(ctx, "/users/1"))
resultChan2 := pool.Submit(client.Get(ctx, "/users/2"))

// Process results
result := <-resultChan
if result.Error != nil {
    log.Printf("Request failed: %v", result.Error)
} else {
    fmt.Printf("Status: %d\n", result.Response.StatusCode)
}

// Wait for all workers to finish
pool.Wait()
```

### Logging Interceptor

```go
logger := goclient.NewStandardLogger()

client := goclient.NewWithOptions(
    goclient.WithBaseURL("https://api.example.com"),
    goclient.WithLoggingInterceptor(logger, goclient.LoggingOptions{
        LogRequestBody:  true,
        LogHeaders:      true,
        MaxBodySize:     1024, // Log up to 1KB of body
    }),
)
```

### Custom Interceptor

```go
type CustomInterceptor struct {
    Next http.RoundTripper
}

func (c *CustomInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
    // Add custom logic before request
    req.Header.Set("X-Custom-Header", "MyValue")
    
    // Execute request
    resp, err := c.Next.RoundTrip(req)
    
    // Add custom logic after request
    if resp != nil {
        fmt.Printf("Response received: %d\n", resp.StatusCode)
    }
    
    return resp, err
}

client := goclient.New(goclient.Config{
    Interceptor: &CustomInterceptor{Next: http.DefaultTransport},
})
```

## Configuration Options

```go
config := goclient.Config{
    BaseURL:               "https://api.example.com",
    Timeout:               30 * time.Second,
    GlobalHeaders:         map[string]string{"User-Agent": "MyApp/1.0"},
    MaxIdleConns:          100,
    MaxIdleConnsPerHost:   10,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ResponseHeaderTimeout: 10 * time.Second,
    DisableKeepAlives:     false,
    DisableCompression:    false,
}
```

## Examples

Check out the [examples](./examples) directory for more detailed usage examples:

- [Basic Usage](./examples/basic/main.go)
- [Authentication](./examples/auth/main.go)
- [Batch Requests](./examples/batch/main.go)
- [Request Pool](./examples/pool/main.go)
- [Logging](./examples/logging/main.go)
- [Custom Interceptor](./examples/interceptor/main.go)

## Performance

GoClient is designed for high performance with:

- **Object Pooling**: Reduces memory allocations and GC pressure
- **Connection Pooling**: Reuses HTTP connections for better throughput
- **Concurrent Execution**: Batch requests and worker pools for parallel processing
- **Efficient Memory Usage**: Minimal allocations in hot paths

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, please:

1. Check the [examples](./examples) directory
2. Open an issue on GitHub
3. Read the [documentation](https://pkg.go.dev/github.com/indalyadav56/goclient)

## Acknowledgments

- Inspired by popular HTTP clients in other languages
- Built with Go's excellent standard library
- Designed for real-world production use cases
