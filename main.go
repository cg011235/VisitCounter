package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST or REDIS_PORT environment variables are not set")
	}

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // no password set if empty
		DB:       0,             // use default DB
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter, err := rdb.Incr(ctx, "counter").Result()
		if err != nil {
			log.Fatalf("Failed to increment counter: %v", err)
		}
		message := fmt.Sprintf("Hello, Docker Compose with Redis! You are visitor number %d.", counter)
		fmt.Fprintf(w, "Message: %s\n", message)
	})

	wrappedMux := loggingMiddleware(mux)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("RemoteAddr: %s, Method: %s, URL: %s, Protocol: %s, Duration: %dms",
			r.RemoteAddr, r.Method, r.URL, r.Proto, duration.Milliseconds())
	})
}
