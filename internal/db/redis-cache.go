package db

import (
	"alist-gallery/config"
	"context"
	"time"
)

func RedisGet(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	val, err := config.RDB.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

func RedisSet(key, value string, expiration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result := config.RDB.SetEx(ctx, key, value, expiration)
	return result.Err() == nil
}
