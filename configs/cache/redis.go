package cache

import (
	"context"
	"fmt"
	"log"
	"my-project/configs/app"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	rdb *redis.Client
)

func InitRedis() {
	config := app.GetConfig()

	// Danh sách Sentinel hosts (thay đổi nếu cần)
	sentinelAddrs := []string{
		config.RedisSentinel1, // Ví dụ: "redis-sentinel-1:26379"
		config.RedisSentinel2, // Ví dụ: "redis-sentinel-2:26379"
		config.RedisSentinel3, // Ví dụ: "redis-sentinel-3:26379"
	}

	// Tên master trong Redis Sentinel
	masterName := "mymaster"

	rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       masterName,
		SentinelAddrs:    sentinelAddrs,
		SentinelPassword: config.RedisSentinelPassword, // Nếu có mật khẩu Sentinel
		Password:         config.RedisPassword,         // Mật khẩu Redis (nếu có)
		DB:               0,
	})

	// Kiểm tra kết nối
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis Sentinel: %v", err)
	}
	fmt.Println("🔗 Kết nối Redis Sentinel thành công!")
}

func GetRedisClient() *redis.Client {
	return rdb
}
