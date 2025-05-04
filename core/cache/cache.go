package cache

import (
	"context"
	"conviction/config"
	"conviction/serializer"
	"encoding/json"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type ValkeyInstance struct {
	client valkey.Client
}

var instance ValkeyInstance

func Init() {

	cfg := config.GetCacheConfig()

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{cfg.Address},
	})
	if err != nil {
		panic(err)
	}

	instance.client = client
}

func SetUserUUID(token string, user_uuid string) (this_err error) {
	this_err = nil

	ctx := context.Background()
	res := instance.client.Do(ctx,
		instance.client.B().
			Set().Key(token).Value(user_uuid).
			Build())

	if err := res.Error(); err != nil {
		this_err = fmt.Errorf("fatal set: %w", err)
		return
	}

	return
}

func GetUserUUID(token string) (user_uuid string, this_err error) {
	this_err = nil

	ctx := context.Background()
	res := instance.client.Do(ctx,
		instance.client.B().
			Get().Key(token).
			Build())
	if err := res.Error(); err != nil {
		this_err = fmt.Errorf("fatal get: %w", err)
		return
	}

	user_uuid, _ = res.ToString()
	return
}

func SetRawWithExpire(key string, raw string, expire int64) (this_err error) {
	this_err = nil

	client := instance.client
	ctx := context.Background()
	resp := instance.client.DoMulti(ctx,
		client.B().Multi().Build(),
		client.B().Set().Key(key).Value(raw).Build(),
		client.B().Expire().Key(key).Seconds(expire).Build(),
		client.B().Exec().Build())
	for _, res := range resp {
		if err := res.Error(); err != nil {
			this_err = fmt.Errorf("fatal set: %w", err)
			return
		}
	}
	return
}

func GetRaw(key string) (raw *string, this_err error) {
	this_err = nil

	ctx := context.Background()
	res := instance.client.Do(ctx,
		instance.client.B().Get().Key(key).Build())
	if err := res.Error(); err != nil {
		this_err = fmt.Errorf("fatal get: %w", err)
		return
	}

	_raw, _ := res.ToString()
	raw = &_raw
	return
}

func SetDownloadSession(key string, session *serializer.DownloadSession, expire int64) (this_err error) {
	sessionRaw, _ := json.Marshal(session)
	return SetRawWithExpire("download_"+key, string(sessionRaw), expire)
}

func GetDownloadSession(key string) (ds *serializer.DownloadSession, this_err error) {
	sessionRaw, err := GetRaw("download_" + key)
	if err != nil {
		this_err = err
		return
	}
	ds = &serializer.DownloadSession{}
	json.Unmarshal([]byte(*sessionRaw), ds)
	return
}

func SetUploadSession(key string, session *serializer.UploadSession, expire int64) error {
	sessionRaw, _ := json.Marshal(session)
	return SetRawWithExpire("callback_"+key, string(sessionRaw), expire)
}

func GetUploadSession(key string) (us *serializer.UploadSession, this_err error) {
	sessionRaw, err := GetRaw("callback_" + key)
	if err != nil {
		this_err = err
		return
	}
	us = &serializer.UploadSession{}
	json.Unmarshal([]byte(*sessionRaw), us)
	return
}

func DeleteUploadSession(key string) (this_err error) {
	this_err = nil

	key = "callback_" + key

	ctx := context.Background()
	res := instance.client.Do(ctx,
		instance.client.B().Del().Key(key).Build())
	if err := res.Error(); err != nil {
		this_err = fmt.Errorf("fatal get: %w", err)
		return
	}
	return
}
