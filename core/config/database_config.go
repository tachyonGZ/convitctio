package config

import "github.com/spf13/viper"

var DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     uint16
}

func InitDatabaseConfig(v *viper.Viper) {
	DatabaseConfig.Host = v.GetString("host")
	DatabaseConfig.User = v.GetString("user")
	DatabaseConfig.Password = v.GetString("password")
	DatabaseConfig.Name = v.GetString("name")
	DatabaseConfig.Port = v.GetUint16("port")
}
