package redis_cache

import "fmt"

// This file contains the following methods :
// 1. Method to establish a connection pool with a redis cache
// 2. Method to set data to cache - with the option to add expiry
// 3. Method to get data from cache

const DefaultConfigPath = "configs/rediscache_config.yaml"

// Initialize creates a new Redis connection pool.
// It can be called without arguments to use the default config path,
// or with a custom config file path.
func Initialize(configFilePath ...string) (*RedisConnectionPool, error) {
	// Use default path if none provided
	path := DefaultConfigPath
	if len(configFilePath) > 0 && configFilePath[0] != "" {
		path = configFilePath[0]
	}

	// Load the configuration from file
	config, err := LoadCacheConfigFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load cache config: %v", err)
	}

	// NOTE :
	// 	While auto pipelining maximizes throughput, it relies on additional goroutines to process requests
	// and responses and may add some latencies due to goroutine scheduling and head of line blocking.
	// You can avoid this by setting DisableAutoPipelining to true, then it will switch to connection pooling
	// approach and serve each request with dedicated connection on the same goroutine.
	// Ref : https://pkg.go.dev/github.com/redis/rueidis#section-readme

	// Create the Redis connection pool
	// Options to pass :
	// ClientOption{
	// 	InitAddress: []string{"127.0.0.1:6379"}, // Redis server address
	// 	Password:    "mystrongpassword",         // Redis password
	// 	SelectDB:    0,                          // Redis database number
	// 	BlockingPoolCleanup						// BlockingPoolCleanup is the duration for cleaning up idle connections.
	// 	BlockingPoolMinSize
	// 	BlockingPoolMaxSize
	// DisableCache
	// DisableAutoPipelining
	//}

}

type RedisConnectionPool struct {
	// Define the fields for your Redis connection pool here

}

func SetDataToCache() {}

func GetDataFromCache() {}
