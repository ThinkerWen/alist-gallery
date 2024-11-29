package config

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/redis/go-redis/v9"
)

func initRedis() {
	if !CONFIG.Redis.Enable {
		return
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", CONFIG.Redis.Host, CONFIG.Redis.Port),
		Password: CONFIG.Redis.Password,
		DB:       CONFIG.Redis.Database,
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}
	log.Info("Redis loaded")
}
