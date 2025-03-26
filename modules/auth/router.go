package auth

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	//e.POST("/auth/register", RegisterHandler)
	//e.POST("/auth/login", LoginHandler)

	protected := e.Group("/users")
	protected.GET("", GetUsersHandler)
	//protected.GET("/:id", GetUserHandler)
	//protected.DELETE("/:id", DeleteUserHandler)
}
