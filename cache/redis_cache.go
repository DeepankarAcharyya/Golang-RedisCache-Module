package redis_cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/rueidis"
)

// This file contains the following methods :
// 1. Method to establish a connection pool with a redis cache
// 2. Method to set data to cache - with the option to add expiry
// 3. Method to get data from cache

// Initialize creates a new Redis connection pool.
// It can be called without arguments to use the default config path,
// or with a custom config file path.
func InitializeCacheConnection(configFilePath ...string) (rueidis.Client, error) {
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
	// 	BlockingPoolCleanup						 // BlockingPoolCleanup is the duration for cleaning up idle connections.
	// 	BlockingPoolMinSize						 // BlockingPoolMinSize is the minimum number of connections in the pool.
	// 	BlockingPoolSize						// BlockingPoolMaxSize is the maximum number of connections in the pool.
	// DisableCache								// DisableCache is used to disable the client side cache.
	// DisableAutoPipelining					// DisableAutoPipelining is used to disable the auto pipelining. So that it will use classic connection pooling approach.
	//}

	_redis_cache_host := fmt.Sprintf("%s:%s", config.Cache.Usage_Cache_DB.Host, config.Cache.Usage_Cache_DB.Port)
	init_address := []string{_redis_cache_host}

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:           init_address,                         // Redis server address
		Password:              config.Cache.Usage_Cache_DB.Password, // Redis password
		SelectDB:              config.Cache.Usage_Cache_DB.Database, // Redis database number
		DisableCache:          config.Cache.Usage_Cache_DB.DisableClientSideCache,
		BlockingPoolCleanup:   config.Cache.Usage_Cache_DB.Pool_Max_Idle_Time,
		BlockingPoolMinSize:   config.Cache.Usage_Cache_DB.Pool_Min_Connections,
		BlockingPoolSize:      config.Cache.Usage_Cache_DB.Pool_Max_Connections,
		DisableAutoPipelining: config.Cache.Usage_Cache_DB.Auto_Pipelining_Mode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %v", err)
	}
	fmt.Println("Connected to Redis!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform a health check
	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		return nil, fmt.Errorf("Redis health check failed: %v", err)
	}
	fmt.Println("Redis is healthy!")

	return client, nil

}

func SetStringDataToCache(ctx context.Context, client rueidis.Client, key string, value string, expiry int64) error {
	// Set the data in cache with the expiry time
	err := client.Do(ctx, client.B().Set().Key(key).Value(value).ExSeconds(expiry).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to set data to cache: %v", err)
	}
	return nil
}

func GetStringDataFromCache(ctx context.Context, client rueidis.Client, key string) (string, error) {
	resp, err := client.Do(ctx, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return "", fmt.Errorf("failed to get data from cache: %v", err)
	}
	return resp, nil
}

func Close(client rueidis.Client) {
	client.Close()
	fmt.Println("Redis connection closed!")
}
