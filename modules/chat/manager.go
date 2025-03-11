package chat

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"my-project/configs/app"
	"net/http"
	"sync"
)

var (
	ctx      = context.Background()
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

var (
	globalWebSocketManager *WebSocketManager
	once                   sync.Once
)

type WebSocketConnection struct {
	Conn  *websocket.Conn
	User  string
	Group string
}

type WebSocketManager struct {
	Connections   map[*WebSocketConnection]bool
	RedisClient   *redis.Client
	RedisChannels map[string]context.CancelFunc
	mu            sync.Mutex
}

func newWebSocketManager() {
	config := app.GetConfig()
	var redisUrl = config.RedisHost + ":" + config.RedisPort
	once.Do(func() {
		globalWebSocketManager = &WebSocketManager{
			Connections:   make(map[*WebSocketConnection]bool),
			RedisClient:   redis.NewClient(&redis.Options{Addr: redisUrl}),
			RedisChannels: make(map[string]context.CancelFunc),
		}
	})
}

func GetWebSocketManager() *WebSocketManager {
	if globalWebSocketManager == nil {
		newWebSocketManager()
	}
	return globalWebSocketManager
}
