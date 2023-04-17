package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Cache struct {
	Redis *redis.Client
}

var cacheObj Cache

func Initialize(ctx context.Context) {
	fmt.Println("Initializing cache")
	address := viper.GetString("REDIS_ADDR")
	dbNum := viper.GetInt("REDIS_DB")
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       dbNum,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("error in connecting to redis: %s", err.Error()))
	}
	cacheObj.Redis = client
}

func GetCacheInstance() *Cache {
	return &cacheObj
}
