// v1 : simple script to connect to redis and do a ping - healthcheck

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/rueidis"
)

func main() {
	// Initialize the Redis client with connection pooling
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"}, // Redis server address
		Password:    "mystrongpassword",         // Redis password
		SelectDB:    0,                          // Redis database number
	})
	if err != nil {
		log.Fatalf("failed to create Redis client: %v", err)
	}
	fmt.Println("Connected to Redis!")

	// Make sure to clean up the client after use
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform a health check
	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		log.Fatalf("Redis health check failed: %v", err)
	}
	fmt.Println("Redis is healthy!")

	// Set a key-value pair
	err = client.Do(ctx, client.B().Set().Key("key").Value("value").Nx().Build()).Error()
	if err != nil {
		log.Fatalf("could not set key: %v", err)
	}

	// Retrieve the value
	resp, err := client.Do(ctx, client.B().Get().Key("key").Build()).ToString()
	if err != nil {
		log.Fatalf("could not get key: %v", err)
	}
	fmt.Println("Value:", resp)
}
