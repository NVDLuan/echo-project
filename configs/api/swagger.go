package api

import (
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	_ "my-project/docs"
)

func SetupSwagger(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
