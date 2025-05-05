package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigFile("./config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal config file: %s", err))
	}

	InitSystemConfig(viper.Sub("system"))
	InitCacheConfig(viper.Sub("cache"))
}
