package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	REDIS Redis
}

type Redis struct {
	HOST     string
	PORT     string
	PASSWORD string
}

func GetAppConfig(filename, path string) *Config {
	return loadConfig(filename, path)
}

func InitRouters() *chi.Mux {
	r := chi.NewRouter()
	// setup cors here ...
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
	return r
}

func loadConfig(filename, path string) *Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("fatal: error reading config file")
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("fatal: error reading config variable")
	}
	return &conf
}

// GetDBConn ...
func GetDBConn(log *zap.SugaredLogger, conf Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.HOST, conf.PORT),
		Password: conf.PASSWORD,
		DB:       0,
	})
	return rdb
}

// GetConsoleLogger ...
func GetConsoleLogger() *zap.SugaredLogger {
	encoder := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
