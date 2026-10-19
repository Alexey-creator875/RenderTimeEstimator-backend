package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	JWT JWTConfig
	Redis RedisConfig
}

type JWTConfig struct {
	ExpiresIn time.Duration
	Token     string
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
   envRedisHost = "REDIS_HOST"
   envRedisPort = "REDIS_PORT"
   envRedisUser = "REDIS_USER"
   envRedisPass = "REDIS_PASSWORD"
)

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{} // создаем объект конфига
	err = viper.Unmarshal(cfg) // читаем информацию из файла, 
	// конвертируем и затем кладем в нашу переменную cfg
	if err != nil {
		return nil, err
	}

	cfg.JWT.Token = os.Getenv("JWT_KEY")
	if cfg.JWT.Token == "" {
		return nil, fmt.Errorf("JWT_KEY is empty")
	}

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))

	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}

	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)

	log.Info("config parsed")

	return cfg, nil
}