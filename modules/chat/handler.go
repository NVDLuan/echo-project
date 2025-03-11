package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"my-project/configs/app"
	"net/http"
)

func verifyToken(tokenString string) (string, error) {
	config := app.GetConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	}
	return "", fmt.Errorf("invalid token")
}

func (manager *WebSocketManager) HandleWebSocket(c echo.Context) error {
	group := c.QueryParam("group")
	token := c.QueryParam("token")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
	}

	username, err := verifyToken(token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &WebSocketConnection{Conn: conn, User: username, Group: group}

	manager.mu.Lock()
	manager.Connections[client] = true
	manager.mu.Unlock()

	// Start Redis listener if not already running
	manager.mu.Lock()
	if _, exists := manager.RedisChannels[group]; !exists {
		ctx, cancel := context.WithCancel(context.Background())
		manager.RedisChannels[group] = cancel
		go manager.listenRedis(ctx, group)
	}
	manager.mu.Unlock()

	// Handle incoming messages
	go manager.receiveMessages(client)

	return nil
}

func (manager *WebSocketManager) receiveMessages(client *WebSocketConnection) {
	defer func() {
		manager.removeConnection(client)
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Invalid message format")
			continue
		}

		recipients, _ := message["to"].([]interface{})
		textMessage, _ := message["message"].(string)

		if len(recipients) > 0 {
			// Send private message
			for _, recipient := range recipients {
				manager.sendToUser(fmt.Sprintf("%v", recipient), fmt.Sprintf("📩 %s says: %s", client.User, textMessage))
			}
		} else {
			// Broadcast via Redis
			manager.RedisClient.Publish(ctx, client.Group, fmt.Sprintf("%s: %s", client.User, textMessage))
		}
	}
}

// Listen for Redis messages
func (manager *WebSocketManager) listenRedis(ctx context.Context, group string) {
	pubsub := manager.RedisClient.Subscribe(ctx, group)
	defer pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err == nil {
				manager.broadcast(msg.Payload, group)
			}
		}
	}
}

// Broadcast message to a group
func (manager *WebSocketManager) broadcast(message, group string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for conn := range manager.Connections {
		if conn.Group == group {
			conn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

// Send private message to a user
func (manager *WebSocketManager) sendToUser(user string, message string) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for conn := range manager.Connections {
		if conn.User == user {
			conn.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

// Remove disconnected WebSocket connection
func (manager *WebSocketManager) removeConnection(client *WebSocketConnection) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	delete(manager.Connections, client)
	if len(manager.Connections) == 0 {
		if cancel, exists := manager.RedisChannels[client.Group]; exists {
			cancel()
			delete(manager.RedisChannels, client.Group)
		}
	}
}
