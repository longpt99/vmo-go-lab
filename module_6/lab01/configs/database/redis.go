package database

import (
	"context"
	"fmt"
	"load-balancer/commons"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis(c *commons.Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     string(fmt.Sprintf("%s:%d", c.Database.Redis.Host, c.Database.Redis.Port)),
		Password: c.Database.Redis.Password,
		DB:       c.Database.Redis.Db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	RedisClient = client
}
