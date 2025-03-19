package redis_cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/redis/rueidis"
)

func setupTestConfig(t *testing.T) string {
	// Create a temporary directory and file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_config.yaml")

	// Test configuration content
	configContent := `
cache:
  usage_cache_db:
    host: "localhost"
    port: "6379"
    password: "mystrongpassword"
    database: 0
    pool_max_connections: 10
    pool_min_connections: 2
    auto_pipelining_mode: true
    disable_cache: false
    pool_max_idle_time: 3600s
`
	err := os.WriteFile(tmpFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	return tmpFile
}

func setupTestClient(t *testing.T) rueidis.Client {
	configPath := setupTestConfig(t)
	client, err := InitializeCacheConnection(configPath)
	if err != nil {
		t.Fatalf("Failed to initialize test client: %v", err)
	}
	return client
}

func TestInitializeCacheConnection(t *testing.T) {
	tests := []struct {
		name        string
		configPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid config",
			configPath: setupTestConfig(t),
			wantErr:    false,
		},
		{
			name:        "invalid config path",
			configPath:  "nonexistent/config.yaml",
			wantErr:     true,
			errContains: "failed to load cache config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := InitializeCacheConnection(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeCacheConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
			}
			if client != nil {
				defer Close(client)
			}
		})
	}
}

func TestSetAndGetStringDataToCache(t *testing.T) {
	client := setupTestClient(t)
	defer Close(client)

	ctx := context.Background()
	tests := []struct {
		name        string
		key         string
		value       string
		expiry      int64
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid set without expiry",
			key:     "test-key-1",
			value:   "test-value-1",
			expiry:  0,
			wantErr: false,
		},
		{
			name:    "valid set with expiry",
			key:     "test-key-2",
			value:   "test-value-2",
			expiry:  60,
			wantErr: false,
		},
		{
			name:        "invalid expiry",
			key:         "test-key-3",
			value:       "test-value-3",
			expiry:      -1,
			wantErr:     true,
			errContains: "expiry time must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Set
			err := SetStringDataToCache(ctx, client, tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetStringDataToCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			// Test Get if Set was successful
			if !tt.wantErr {
				got, err := GetStringDataFromCache(ctx, client, tt.key)
				if err != nil {
					t.Errorf("GetStringDataFromCache() error = %v", err)
					return
				}
				if got != tt.value {
					t.Errorf("GetStringDataFromCache() = %v, want %v", got, tt.value)
				}

				// Test expiry if set
				if tt.expiry > 0 {
					time.Sleep(time.Duration(tt.expiry+1) * time.Second)
					got, err := GetStringDataFromCache(ctx, client, tt.key)
					if err != nil {
						t.Errorf("GetStringDataFromCache() after expiry error = %v", err)
					}
					if got != "" {
						t.Errorf("Expected empty string after expiry, got %v", got)
					}
				}
			}
		})
	}
}

func TestSetAndGetIntDataToCache(t *testing.T) {
	client := setupTestClient(t)
	defer Close(client)

	ctx := context.Background()
	tests := []struct {
		name        string
		key         string
		value       int
		expiry      int64
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid set without expiry",
			key:     "int-key-1",
			value:   42,
			expiry:  0,
			wantErr: false,
		},
		{
			name:    "valid set with expiry",
			key:     "int-key-2",
			value:   100,
			expiry:  60,
			wantErr: false,
		},
		{
			name:        "invalid expiry",
			key:         "int-key-3",
			value:       200,
			expiry:      -1,
			wantErr:     true,
			errContains: "expiry time must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Set
			err := SetIntDataToCache(ctx, client, tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIntDataToCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			// Test Get if Set was successful
			if !tt.wantErr {
				got, err := GetIntDataFromCache(ctx, client, tt.key)
				if err != nil {
					t.Errorf("GetIntDataFromCache() error = %v", err)
					return
				}
				if got != tt.value {
					t.Errorf("GetIntDataFromCache() = %v, want %v", got, tt.value)
				}

				// Test expiry if set
				if tt.expiry > 0 {
					time.Sleep(time.Duration(tt.expiry+1) * time.Second)
					got, err := GetIntDataFromCache(ctx, client, tt.key)
					if err != nil {
						t.Errorf("GetIntDataFromCache() after expiry error = %v", err)
					}
					if got != -1 {
						t.Errorf("Expected -1 after expiry, got %v", got)
					}
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
}
