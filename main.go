// v1 : simple script to connect to redis and do a ping - healthcheck -- DONE
// v2 : create a connection pool
// v3 : create the functions to set and get values
// v4 : create the functions to set and get values with expiration
// v5 : create the serdes module with support for timestamp, dictionary, list, set, zset

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	redis_cache "github.com/DeepankarAcharyya/Golang-RedisCache-Module/cache"
)

func main() {
	// Initialize the Redis client with connection pooling using default config path
	client, err := redis_cache.InitializeCacheConnection()
	if err != nil {
		log.Fatalf("failed to create Redis client: %v", err)
	}

	// Make sure to clean up the client after use
	defer redis_cache.Close(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set a string key-value pair with 1 hour expiry
	err = redis_cache.SetStringDataToCache(ctx, client, "string-key", "test-value", 3600)
	if err != nil {
		log.Fatalf("could not set string key: %v", err)
	}

	// Retrieve the string value
	stringValue, err := redis_cache.GetStringDataFromCache(ctx, client, "string-key")
	if err != nil {
		log.Fatalf("could not get string key: %v", err)
	}
	fmt.Println("String Value:", stringValue)

	// Set an integer key-value pair with 30 minutes expiry
	err = redis_cache.SetIntDataToCache(ctx, client, "int-key", 42, 1800)
	if err != nil {
		log.Fatalf("could not set int key: %v", err)
	}

	// Retrieve the integer value
	intValue, err := redis_cache.GetIntDataFromCache(ctx, client, "int-key")
	if err != nil {
		log.Fatalf("could not get int key: %v", err)
	}
	fmt.Println("Integer Value:", intValue)
}
