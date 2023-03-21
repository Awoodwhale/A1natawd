package cache

import (
	"github.com/go-redis/redis"
	"go_awd/pkg/wlog"
	"strconv"
)

var RedisClient *redis.Client

// InitDatabase
// @Description: 初始化redis
// @param redisAddr string
// @param redisName string
// @param redisPassword string
func InitDatabase(redisAddr, redisName, redisPassword string) {
	db, _ := strconv.ParseInt(redisName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       int(db),
		Password: redisPassword,
	})
	_, err := client.Ping().Result()
	if err != nil {
		wlog.Logger.Errorln("cache common redis ping, ", err)
		panic(err)
	}
	RedisClient = client
}
