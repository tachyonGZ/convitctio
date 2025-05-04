package config

import "github.com/spf13/viper"

type CacheConfig struct {
	Address string
}

var config CacheConfig

func InitCacheConfig(v *viper.Viper) {
	config.Address = v.GetString("address")
}

func GetCacheConfig() *CacheConfig {
	return &config
}
