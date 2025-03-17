package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func InitLogger() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true, // Tắt màu để log rõ ràng hơn
		FullTimestamp:   true, // Hiển thị thời gian đầy đủ
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})
	log.SetLevel(logrus.DebugLevel)
}

func SetupEchoLogger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom}  |  ${method}  |  ${uri}  |  ${status}  |  ${latency_human}\n",
		CustomTimeFormat: "2006-01-02 15:04:05.000", // Hiển thị đến millisecond
		Output:           os.Stdout,
	})
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return log
}
