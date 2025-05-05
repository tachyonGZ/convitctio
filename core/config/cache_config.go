package config

import "github.com/spf13/viper"

var CacheConfig struct {
	Address string
}

func InitCacheConfig(v *viper.Viper) {
	CacheConfig.Address = v.GetString("address")
}
