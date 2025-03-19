package redis_cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCacheConfigFromFile(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
		validate    func(*CacheConnectionConfig) error
	}{
		{
			name: "valid configuration",
			yamlContent: `
cache:
  usage_cache_db:
    host: "test-host"
    port: "6379"
    password: "test-pass"
    database: "0"
    ssl_mode: "disable"
    pool_max_connections: 10
`,
			wantErr: false,
			validate: func(cfg *CacheConnectionConfig) error {
				if cfg.Cache.Usage_Cache_DB.Host != "test-host" {
					t.Errorf("expected host to be 'test-host', got %s", cfg.Cache.Usage_Cache_DB.Host)
				}
				if cfg.Cache.Usage_Cache_DB.Port != "6379" {
					t.Errorf("expected port to be '6379', got %s", cfg.Cache.Usage_Cache_DB.Port)
				}
				if cfg.Cache.Usage_Cache_DB.Password != "test-pass" {
					t.Errorf("expected password to be 'test-pass', got %s", cfg.Cache.Usage_Cache_DB.Password)
				}
				if cfg.Cache.Usage_Cache_DB.Database != "0" {
					t.Errorf("expected database to be '0', got %s", cfg.Cache.Usage_Cache_DB.Database)
				}
				if cfg.Cache.Usage_Cache_DB.Pool_Max_Connections != 10 {
					t.Errorf("expected pool_max_connections to be 10, got %d", cfg.Cache.Usage_Cache_DB.Pool_Max_Connections)
				}
				return nil
			},
		},
		{
			name: "invalid yaml",
			yamlContent: `
invalid:
  - yaml:
    content
`,
			wantErr: true,
			validate: func(cfg *CacheConnectionConfig) error {
				return nil
			},
		},
		{
			name: "missing required fields",
			yamlContent: `
cache:
  usage_cache_db:
    host: "test-host"
`,
			wantErr: false,
			validate: func(cfg *CacheConnectionConfig) error {
				if cfg.Cache.Usage_Cache_DB.Host != "test-host" {
					t.Errorf("expected host to be 'test-host', got %s", cfg.Cache.Usage_Cache_DB.Host)
				}
				if cfg.Cache.Usage_Cache_DB.Port != "" {
					t.Errorf("expected port to be empty, got %s", cfg.Cache.Usage_Cache_DB.Port)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for the test
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "config.yaml")

			// Write the test content to the file
			err := os.WriteFile(tmpFile, []byte(tt.yamlContent), 0644)
			if err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Test the LoadCacheConfigFromFile function
			got, err := LoadCacheConfigFromFile(tmpFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadCacheConfigFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we don't expect an error, validate the configuration
			if !tt.wantErr {
				if err := tt.validate(got); err != nil {
					t.Errorf("validation failed: %v", err)
				}
			}
		})
	}
}

func TestLoadCacheConfigFromNonExistentFile(t *testing.T) {
	_, err := LoadCacheConfigFromFile("non_existent_file.yaml")
	if err == nil {
		t.Error("LoadCacheConfigFromFile() expected error for non-existent file, got nil")
	}
}

func TestLoadCacheConfigWithInvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	// Write some content
	err := os.WriteFile(tmpFile, []byte("test"), 0000)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Test loading the file
	_, err = LoadCacheConfigFromFile(tmpFile)
	if err == nil {
		t.Error("LoadCacheConfigFromFile() expected error for unreadable file, got nil")
	}
}
