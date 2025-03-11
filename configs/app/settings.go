package app

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AllowHost     []string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	DBSSLMode     string
	JWTSecret     string
	JWTSetCookie  bool
	JWTAuthCookie string

	RedisHost     string
	RedisPort     string
	RedisPassword string
}

var (
	globalConfig *Config
	once         sync.Once
)

func LoadConfig() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("❌ Lỗi khi tải file .env")
		}
		allowHost, exists := os.LookupEnv("ALLOW_HOST")
		if !exists {
			allowHost = "*"
		}
		globalConfig = &Config{
			AllowHost:     strings.Split(allowHost, ","),
			DBHost:        os.Getenv("DB_HOST"),
			DBUser:        os.Getenv("DB_USER"),
			DBPassword:    os.Getenv("DB_PASSWORD"),
			DBName:        os.Getenv("DB_NAME"),
			DBPort:        os.Getenv("DB_PORT"),
			DBSSLMode:     os.Getenv("DB_SSLMODE"),
			JWTSecret:     os.Getenv("JWT_SECRET"),
			JWTSetCookie:  string(os.Getenv("JWT_SET_COOKIE")) == "true",
			JWTAuthCookie: os.Getenv("JWT_AUTH_COOKIE"),

			RedisHost:     os.Getenv("REDIS_HOST"),
			RedisPort:     os.Getenv("REDIS_PORT"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		}
		log.Println("✅ Config loaded successfully!")
	})
}

func GetConfig() *Config {
	if globalConfig == nil {
		LoadConfig()
	}
	return globalConfig
}
