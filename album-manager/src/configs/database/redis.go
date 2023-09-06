package database

import (
	"album-manager/src/configs"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type RDConfig struct {
	DB  *redis.Client
	ctx context.Context
}

func (c *RDConfig) Ping() error {
	err := c.DB.Ping(c.ctx).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RDConfig) Close() error {
	err := c.DB.Conn().Close()
	if err != nil {
		return err
	}

	return nil
}

func InitRedis(ctx context.Context) RDConfig {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", configs.Env.Redis.Host, configs.Env.Redis.Port),
		Password: configs.Env.Redis.Password, // no password set
		DB:       configs.Env.Redis.Database, // use default DB
	})

	fmt.Println(configs.Env.Redis)

	RedisClient = rdb

	return RDConfig{rdb, ctx}
}
