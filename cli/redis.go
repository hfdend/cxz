package cli

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/hfdend/cxz/conf"
)

var (
	Redis    *redis.Client
	RedisNil = redis.Nil
)

func InitializeRedis() {
	redisConfig := conf.Config.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:        redisConfig.Addr,
		Password:    redisConfig.Password,
		DB:          redisConfig.DB,
		PoolSize:    redisConfig.PoolSize,
		IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
	})
	err := Redis.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}
}
