package db

import (
	"github.com/redis/go-redis/v9"
)

var DB *redis.Client

func ConnectDatabase() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	DB = rdb

	// Output: key value
	// key2 does not exist
}
