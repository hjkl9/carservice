package svc

import (
	"carservice/internal/config"
	"context"
	"testing"
)

func TestNewRedis(t *testing.T) {
	c := config.Config{
		RedisConf: config.RedisConf{
			Addr:        "127.0.0.1:6379",
			Requirepass: "",
			DB:          0,
		},
	}
	rdb := NewRedis(c)
	cmd := rdb.Ping(context.Background())
	result, err := cmd.Result()
	if err != nil {
		t.Errorf("failed to ping the redis client, err: %s\n", err.Error())
		return
	}

	t.Logf("%s\n", result)
}
