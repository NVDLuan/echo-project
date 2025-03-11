package auth

import (
	_ "errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"my-project/configs/app"
	"net/http"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func GenerateJWT(userID uint) (string, error) {
	config := app.GetConfig()
	claims := jwt.MapClaims{
		"user_id": userID,
		"jit":     uuid.New().String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func SetCookie(cookieName string, token string, c echo.Context) error {
	// ✅ Set cookie
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.MaxAge = 3600 // 1 giờ
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}
