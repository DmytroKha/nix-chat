package config

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

func CreateRedisClient() {
	opt, err := redis.ParseURL("redis://redis:6379/0")
	if err != nil {
		panic(err)
	}

	rds := redis.NewClient(opt)
	Redis = rds
}
