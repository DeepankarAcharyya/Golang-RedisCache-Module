package redis_cache

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type CacheConnectionConfig struct {
	Cache struct {
		Usage_Cache_DB struct {
			Host                   string        `yaml:"host"`
			Port                   string        `yaml:"port"`
			Password               string        `yaml:"password"`
			Database               int           `yaml:"database"`
			Pool_Max_Connections   int           `yaml:"pool_max_connections"`
			Pool_Min_Connections   int           `yaml:"pool_min_connections"`
			Auto_Pipelining_Mode   bool          `yaml:"auto_pipelining_mode"`
			DisableClientSideCache bool          `yaml:"disable_cache"`
			Pool_Max_Idle_Time     time.Duration `yaml:"pool_max_idle_time"`
		} `yaml:"usage_cache_db"`
	} `yaml:"cache"`
}

// setDefaults sets default values for fields that are empty
func (c *CacheConnectionConfig) setDefaults() {
	if c.Cache.Usage_Cache_DB.Host == "" {
		c.Cache.Usage_Cache_DB.Host = "127.0.0.1"
	}
	if c.Cache.Usage_Cache_DB.Port == "" {
		c.Cache.Usage_Cache_DB.Port = "6379"
	}
	if !c.Cache.Usage_Cache_DB.Auto_Pipelining_Mode {
		c.Cache.Usage_Cache_DB.Auto_Pipelining_Mode = false
	}
	if !c.Cache.Usage_Cache_DB.DisableClientSideCache {
		c.Cache.Usage_Cache_DB.DisableClientSideCache = true
	}
	if c.Cache.Usage_Cache_DB.Pool_Max_Idle_Time == 0 {
		c.Cache.Usage_Cache_DB.Pool_Max_Idle_Time = 90 * 60 * time.Second
	}
}

func LoadCacheConfigFromFile(yaml_config_file_path string) (*CacheConnectionConfig, error) {
	data, err := os.ReadFile(yaml_config_file_path)
	if err != nil {
		return nil, err
	}

	var cache_pool_config CacheConnectionConfig
	if err = yaml.Unmarshal(data, &cache_pool_config); err != nil {
		return nil, err
	}

	cache_pool_config.setDefaults()

	return &cache_pool_config, nil
}
