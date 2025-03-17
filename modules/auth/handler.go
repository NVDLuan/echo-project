package auth

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"my-project/configs/app"
	"my-project/configs/cache"
	"my-project/pkg/logger"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var luaScript = redis.NewScript(`
    local key = KEYS[1]
    local data = redis.call("GET", key)
    if data then
        return data
    else
        return nil
    end
`)
var cacheMutex = sync.Mutex{}

// GetUsersHandler lấy danh sách user
// @Summary Lấy danh sách người dùng
// @Description API này trả về danh sách tất cả người dùng
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsersHandler(c echo.Context) error {
	// Kiểm tra dữ liệu trong Redis bằng Lua

	redisClient := cache.GetRedisClient()
	log := logger.GetLogger()
	result, err := luaScript.Run(cache.Ctx, redisClient, []string{"hello"}).Result()
	if err == redis.Nil || result == nil {
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		result, err = luaScript.Run(cache.Ctx, redisClient, []string{"hello"}).Result()
		if err == redis.Nil || result == nil {
			log.Warn("🔍 Không tìm thấy dữ liệu trong Redis, lấy từ database...")
			users, err := GetUsers()
			if err != nil {
				log.Error("❌ Lỗi khi lấy dữ liệu từ database:", err)
				return c.JSON(http.StatusInternalServerError, "Database error")
			}

			usersJSONBytes, err := json.Marshal(users)
			if err != nil {
				log.Error("❌ Lỗi khi encode JSON:", err)
				return c.JSON(http.StatusInternalServerError, "JSON encoding error")
			}

			err = redisClient.Set(cache.Ctx, "hello", string(usersJSONBytes), 10*time.Minute).Err()
			if err != nil {
				log.Error("❌ Lỗi khi lưu dữ liệu vào Redis:", err)
				return c.JSON(http.StatusInternalServerError, "Redis store error")
			}
			log.Info("✅ Lưu dữ liệu mới vào Redis thành công!")
			return c.JSON(http.StatusOK, users)
		}
	} else if err != nil {
		log.Error("❌ Lỗi khi chạy Lua script:", err)
		return c.JSON(http.StatusInternalServerError, "Redis Lua error")
	}

	// Nếu Redis có dữ liệu
	var users []User
	if err := json.Unmarshal([]byte(result.(string)), &users); err != nil {
		log.Error("❌ Lỗi khi parse JSON từ Redis:", err)
		return c.JSON(http.StatusInternalServerError, "JSON parsing error")
	}
	log.Info("✅ Lấy dữ liệu từ Redis thành công!")
	return c.JSON(http.StatusOK, users)
}

// GetUserHandler lấy thông tin user theo ID
// @Summary Lấy user theo ID
// @Description API này trả về thông tin một user dựa trên ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := GetUserByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUserHandler xóa user
// @Summary Xóa user theo ID
// @Description API này xóa một user từ database
// @Tags users
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUserHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
}

// RegisterHandler
// @Summary Đăng ký user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Thông tin user"
// @Success 201 {object} map[string]string
// @Router /auth/register [post]
func RegisterHandler(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	req.Password = hashedPassword
	if err := CreateUser(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered"})
}

// LoginHandler
// @Summary Đăng nhập user
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Email và Password"
// @Success 200 {object} map[string]string
// @Router /auth/login [post]
func LoginHandler(c echo.Context) error {
	var req LoginRequest
	log := logger.GetLogger()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		log.Error("user not found")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	if !CheckPassword(user.PasswordHash, req.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	log.Info("Login success")
	config := app.GetConfig()
	if config.JWTSetCookie {
		return SetCookie(config.JWTAuthCookie, token, c)
	}
	return c.JSON(http.StatusOK, map[string]string{"access_token": token})
}
