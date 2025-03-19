package redis_cache

import (
	"os"

	"gopkg.in/yaml.v2"
)

type CacheConnectionConfig struct {
	Cache struct {
		Usage_Cache_DB struct {
			Host                   string `yaml:"host"`
			Port                   string `yaml:"port"`
			Password               string `yaml:"password"`
			Database               string `yaml:"database"`
			Pool_Max_Connections   int    `yaml:"pool_max_connections"`
			Pool_Min_Connections   int    `yaml:"pool_min_connections"`
			Auto_Pipelining_Mode   bool   `yaml:"auto_pipelining_mode"`
			DisableClientSideCache bool   `yaml:"disable_cache"`
			Pool_Max_Idle_Time     string `yaml:"pool_max_idle_time"`
		} `yaml:"usage_cache_db"`
	} `yaml:"cache"`
}

// setDefaults sets default values for fields that are empty
func (c *CacheConnectionConfig) setDefaults() {
	if c.Cache.Usage_Cache_DB.Port == "" {
		c.Cache.Usage_Cache_DB.Port = "6379"
	}
	if c.Cache.Usage_Cache_DB.Database == "" {
		c.Cache.Usage_Cache_DB.Database = "0"
	}
	if c.Cache.Usage_Cache_DB.SSL_Mode == "" {
		c.Cache.Usage_Cache_DB.SSL_Mode = "disable"
	}
	if !c.Cache.Usage_Cache_DB.Auto_Pipelining_Mode {
		c.Cache.Usage_Cache_DB.Auto_Pipelining_Mode = false
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
