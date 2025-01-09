package redisdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func HandleConnections() (*redis.Client, error) {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-19096.c252.ap-southeast-1-1.ec2.redns.redis-cloud.com:19096",
		Password: "LTBJl47CNZaEHaDCthIMcVottHUL9Hyg",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()

	return rdb, err
}
