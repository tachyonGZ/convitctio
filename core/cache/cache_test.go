package cache

import (
	"bytes"
	"conviction/config"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.SetConfigType("yaml")
	var yamlExample = []byte(`
cache:
  address: localhost:6379
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	config.InitCacheConfig(viper.Sub("cache"))

	Init()
}

func TestValkeyInstance(t *testing.T) {

	token, _ := uuid.NewV7()
	uuid, _ := uuid.NewV7()

	SetUserUUID(token.String(), uuid.String())
	user_uuid_get, _ := GetUserUUID(token.String())

	assert.Equal(t, uuid.String(), user_uuid_get)
}
