package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"my-project/configs/api"
	"my-project/configs/app"
	"my-project/configs/cache"
	"my-project/configs/database"
	_ "my-project/docs"
	"my-project/modules/auth"
	"my-project/pkg/logger"
)

func main() {

	logger.InitLogger()
	cache.InitRedis()
	database.InitDB()
	config := app.GetConfig()

	auth.Migrate(database.DB)

	e := echo.New()
	e.Use(logger.SetupEchoLogger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowHost,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))
	api.SetupSwagger(e)
	auth.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
