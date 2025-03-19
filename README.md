# Golang Redis Cache Module

A lightweight, configurable Redis caching module for Go applications with connection pooling support.

## Features

- Simple Redis cache operations (get/set) with type support for strings and integers
- YAML-based configuration with default path handling
- Connection pooling with configurable parameters
- Automatic connection cleanup
- Context-aware operations with timeout support
- Docker-based test environment

## Installation

```bash
go get github.com/DeepankarAcharyya/Golang-RedisCache-Module
```

## Quick Start

1. Create a configuration file at `configs/rediscache_config.yaml` (default path):

```yaml
cache:
  usage_cache_db:
    host: "localhost"
    port: "6379"
    password: "mystrongpassword"
    database: 0
    ssl_mode: "disable"
    pool_max_connections: 6
    pool_min_connections: 1
```

2. Use the module in your application:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    redis_cache "github.com/DeepankarAcharyya/Golang-RedisCache-Module/cache"
)

func main() {
    // Initialize Redis connection with default config path
    client, err := redis_cache.InitializeCacheConnection()
    if err != nil {
        log.Fatalf("failed to create Redis client: %v", err)
    }
    defer redis_cache.Close(client)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Set string value with 1 hour expiry
    err = redis_cache.SetStringDataToCache(ctx, client, "string-key", "test-value", 3600)
    if err != nil {
        log.Fatalf("could not set string key: %v", err)
    }

    // Get string value
    stringValue, err := redis_cache.GetStringDataFromCache(ctx, client, "string-key")
    if err != nil {
        log.Fatalf("could not get string key: %v", err)
    }
    fmt.Println("String Value:", stringValue)

    // Set integer value with 30 minutes expiry
    err = redis_cache.SetIntDataToCache(ctx, client, "int-key", 42, 1800)
    if err != nil {
        log.Fatalf("could not set int key: %v", err)
    }

    // Get integer value
    intValue, err := redis_cache.GetIntDataFromCache(ctx, client, "int-key")
    if err != nil {
        log.Fatalf("could not get int key: %v", err)
    }
    fmt.Println("Integer Value:", intValue)
}
```

## Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| host | Redis server hostname | localhost |
| port | Redis server port | 6379 |
| password | Redis authentication password | "" |
| database | Redis database number | 0 |
| ssl_mode | SSL mode (disable, verify-ca) | disable |
| pool_max_connections | Maximum number of connections | 6 |
| pool_min_connections | Minimum number of connections | 1 |

## Key Features

### Type-Safe Operations
- Dedicated functions for string and integer data types
- Automatic type conversion and validation
- Built-in error handling

### Connection Management
- Automatic connection pooling
- Connection cleanup with `defer Close()`
- Context-aware operations with timeout support

### Configuration
- Default configuration path support
- Custom configuration path option
- YAML-based configuration format

## Testing

The project includes a Docker-based test environment:

```bash
cd test-setup
./start_test.sh
```

This will:
- Build a custom Redis Docker image
- Start a Redis container on port 6379
- Configure the test instance with the specified password
- Verify the Redis instance is running correctly

## Dependencies

- Go 1.23 or later
- Redis server (for production)
- Docker (for testing)
- Required Go packages:
  - github.com/redis/rueidis
  - gopkg.in/yaml.v2

## License

[MIT License](LICENSE)

Copyright (c) 2025 Deepankar Acharyya
