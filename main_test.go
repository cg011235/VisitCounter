package main

import (
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestRedisConnection(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "test", "value", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value in Redis: %v", err)
	}

	val, err := rdb.Get(ctx, "test").Result()
	if err != nil {
		t.Fatalf("Failed to get value from Redis: %v", err)
	}

	if val != "value" {
		t.Fatalf("Expected 'value', got '%s'", val)
	}
}
