package config

import "github.com/spf13/viper"

var SystemConfig struct {
	Debug bool
}

func InitSystemConfig(v *viper.Viper) {
	SystemConfig.Debug = v.GetBool("debug")
}
