package redis_cache

import (
	"os"

	"gopkg.in/yaml.v2"
)

type CacheConnectionConfig struct {
	Cache struct {
		Usage_Cache_DB struct {
			Host                 string `yaml:"host"`
			Port                 string `yaml:"port"`
			Password             string `yaml:"password"`
			Database             string `yaml:"database"`
			SSL_Mode             string `yaml:"ssl_mode"`
			Pool_Max_Connections int    `yaml:"pool_max_connections"`
		} `yaml:"usage_cache_db"`
	} `yaml:"cache"`
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

	return &cache_pool_config, nil
}
