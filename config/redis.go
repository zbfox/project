package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() *redis.Client {
	r := Conf.Redis

	redisDb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
		Username: r.UserName,
	})
	_, err := redisDb.Ping(context.Background()).Result()
	if err != nil {
		panic("redis连接失败: " + err.Error())
	} else {
		fmt.Println("redis连接成功")
	}

	rdb = redisDb
	return redisDb
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return rdb
}
