package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"my-project/configs/app"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitRedis() {
	config := app.GetConfig()
	redisUrl := config.RedisHost + ":" + config.RedisPort
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: config.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}
	fmt.Println("🔗 Kết nối Redis thành công!")
}
