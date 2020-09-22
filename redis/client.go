package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

var once sync.Once
var redisClient *redis.Client

// 单例模式
func RedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: ":6379",
		})

	})
	return redisClient
}
func init() {
	cli := RedisClient()
	err := cli.Ping().Err()
	if err != nil {
		panic(errors.New("please start redis-server \n" + err.Error()))
	}
}
