package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"my-project/configs/app"
	"my-project/configs/database"
	"net/http"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		config := app.GetConfig()
		var tokenString string
		if config.JWTSetCookie {
			cookie, err := c.Cookie(config.JWTAuthCookie)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}
			tokenString = cookie.Value
		} else {
			tokenString = c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}

		userID := uint(claims["user_id"].(float64))

		var user User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
		}

		c.Set("user", user)

		return next(c)
	}
}

func RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role").(string)

			if role != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Access denied",
				})
			}

			return next(c)
		}
	}
}
