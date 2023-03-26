package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var ctx = context.Background()
var client *redis.Client

func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: "000000",                        // 密码
		DB:       viper.GetInt("redis.db"),        // 数据库
		PoolSize: viper.GetInt("redis.pool_size"), // 连接池大小
	})
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("error happens in redis by %w", err)
	}
	return
}

func Close() {
	_ = client.Close()
}
