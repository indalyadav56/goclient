package main

import (
	"fmt"
	"log"

	"github.com/indalyadav56/goclient"
)

// // Custom logger example
// type CustomLogger struct {
// 	logger *log.Logger
// }

// func NewCustomLogger() *CustomLogger {
// 	return &CustomLogger{
// 		logger: log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags),
// 	}
// }

// func (l *CustomLogger) Log(level goclient.LogLevel, message string, fields map[string]interface{}) {
// 	// Custom formatting with colors
// 	var color string
// 	switch level {
// 	case goclient.LogLevelDebug:
// 		color = "\033[36m" // Cyan
// 	case goclient.LogLevelInfo:
// 		color = "\033[32m" // Green
// 	case goclient.LogLevelWarn:
// 		color = "\033[33m" // Yellow
// 	case goclient.LogLevelError:
// 		color = "\033[31m" // Red
// 	default:
// 		color = "\033[0m" // Reset
// 	}

// 	l.logger.Printf("%s[%s]\033[0m %s", color, level.String(), message)

// 	// Pretty print fields
// 	for k, v := range fields {
// 		switch k {
// 		case "headers", "response_headers":
// 			fmt.Printf("  📋 %s:\n", k)
// 			if headers, ok := v.(map[string]string); ok {
// 				for hk, hv := range headers {
// 					fmt.Printf("    %s: %s\n", hk, hv)
// 				}
// 			}
// 		case "query_params":
// 			fmt.Printf("  🔍 %s: %v\n", k, v)
// 		case "body", "response_body":
// 			fmt.Printf("  📄 %s: %v\n", k, v)
// 		case "status_code":
// 			fmt.Printf("  📊 %s: %v\n", k, v)
// 		case "duration_ms":
// 			fmt.Printf("  ⏱️  %s: %vms\n", k, v)
// 		default:
// 			fmt.Printf("  🔸 %s: %v\n", k, v)
// 		}
// 	}
// 	fmt.Println()
// }

func main() {
	goclient.EnableDebug()

	var post map[string]interface{}
	err := goclient.Get("https://jsonplaceholder.typicode.com/posts/1").
		SetHeader("User-Agent", "goclient-debug-demo").
		SetQueryParam("debug", "true").
		Into(&post)

	if err != nil {
		log.Printf("❌ GET failed: %v", err)
	} else {
		fmt.Printf("✅ GET Success: %s\n", post["title"])
	}

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 2. CUSTOM LOGGER
	// fmt.Println("2️⃣ Custom Logger with Pretty Formatting:")
	// fmt.Println("----------------------------------------")

	// customLogger := NewCustomLogger()
	// goclient.SetLogger(customLogger)

	// newPost := map[string]interface{}{
	// 	"title":  "Debug Demo Post",
	// 	"body":   "This request will be logged with custom formatting!",
	// 	"userId": 1,
	// }

	// var createdPost map[string]interface{}
	// err = goclient.Post("https://jsonplaceholder.typicode.com/posts").
	// 	SetHeader("Content-Type", "application/json").
	// 	SetHeader("X-Debug", "custom-logger").
	// 	SetBody(newPost).
	// 	Into(&createdPost)

	// if err != nil {
	// 	log.Printf("❌ POST failed: %v", err)
	// } else {
	// 	fmt.Printf("✅ POST Success: Created post ID %.0f\n", createdPost["id"])
	// }

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 3. CLIENT INSTANCE WITH DEBUG
	// fmt.Println("3️⃣ Client Instance with Debug Logging:")
	// fmt.Println("--------------------------------------")

	// client := goclient.New(goclient.Config{
	// 	BaseURL: "https://httpbin.org",
	// 	Timeout: 30 * time.Second,
	// 	GlobalHeaders: map[string]string{
	// 		"X-Client-Debug": "enabled",
	// 	},
	// }).EnableDebug()

	// var httpbinResp map[string]interface{}
	// err = client.Get("/get").
	// 	SetQueryParam("show", "headers").
	// 	SetQueryParam("show", "args").
	// 	Into(&httpbinResp)

	// if err != nil {
	// 	log.Printf("❌ Httpbin GET failed: %v", err)
	// } else {
	// 	fmt.Printf("✅ Httpbin GET Success\n")
	// }

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 4. ERROR LOGGING
	// fmt.Println("4️⃣ Error Request Logging:")
	// fmt.Println("-------------------------")

	// var errorResp map[string]interface{}
	// err = goclient.Get("https://httpbin.org/status/404").
	// 	SetHeader("X-Test", "error-logging").
	// 	Into(&errorResp)

	// if err != nil {
	// 	fmt.Printf("✅ Error logged correctly: %v\n", err)
	// }

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 5. BATCH REQUESTS WITH LOGGING
	// fmt.Println("5️⃣ Batch Requests with Debug Logging:")
	// fmt.Println("-------------------------------------")

	// batch := goclient.Batch()
	// batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/1"))
	// batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/2"))
	// batch.Add(goclient.Get("https://jsonplaceholder.typicode.com/posts/3"))

	// responses, errors := batch.Execute(context.Background())
	// fmt.Printf("✅ Batch completed: %d responses, %d errors\n", len(responses), len(errors))

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 6. AUTHENTICATION WITH LOGGING
	// fmt.Println("6️⃣ Authentication with Debug Logging:")
	// fmt.Println("-------------------------------------")

	// goclient.SetBearerToken("demo-token-12345")

	// var authResp map[string]interface{}
	// err = goclient.Get("https://httpbin.org/bearer").Into(&authResp)
	// if err != nil {
	// 	log.Printf("Auth request failed: %v", err)
	// } else {
	// 	fmt.Printf("✅ Auth request completed\n")
	// }

	// fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	// // 7. DISABLE DEBUG
	// fmt.Println("7️⃣ Disable Debug Logging:")
	// fmt.Println("-------------------------")

	// goclient.DisableDebug()

	// var silentResp map[string]interface{}
	// err = goclient.Get("https://jsonplaceholder.typicode.com/posts/4").Into(&silentResp)
	// if err != nil {
	// 	log.Printf("Silent request failed: %v", err)
	// } else {
	// 	fmt.Printf("✅ Silent request completed (no debug logs)\n")
	// }

	// fmt.Println("\n🎉 Logging & Debugging Demo Complete!")
	// fmt.Println("\n📋 What you can see with debug logging:")
	// fmt.Println("   • 🌐 Full request URL with query parameters")
	// fmt.Println("   • 📤 Request method (GET, POST, etc.)")
	// fmt.Println("   • 📋 All request headers (auth headers redacted)")
	// fmt.Println("   • 📄 Request body content (truncated if large)")
	// fmt.Println("   • 📊 Response status code and message")
	// fmt.Println("   • 📋 Response headers")
	// fmt.Println("   • 📄 Response body content (truncated if large)")
	// fmt.Println("   • ⏱️  Request duration in milliseconds")
	// fmt.Println("   • 🎨 Custom logger support with your own formatting")
	// fmt.Println("   • 🔴 Error-level logging for failed requests (4xx/5xx)")
	// fmt.Println("\n💡 Usage:")
	// fmt.Println("   goclient.EnableDebug()  // Enable for package-level functions")
	// fmt.Println("   client.EnableDebug()    // Enable for specific client")
	// fmt.Println("   goclient.SetLogger(customLogger) // Use custom logger")
}
