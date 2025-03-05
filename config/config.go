package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	REDIS     Redis
	JWT       string
	RATELIMIT RateLimit
}

type Redis struct {
	HOST     string
	PORT     string
	PASSWORD string
}

type RateLimit struct {
}

var (
	appConfig *Config
)

func GetAppConfig(filename, path string) *Config {
	if appConfig != nil {
		return appConfig
	}
	conf := loadConfig(filename, path)
	appConfig = conf
	return appConfig
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

// Middleware to verify JWT token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No Bearer prefix
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		// Validate and parse the JWT
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Ensure correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(appConfig.JWT), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Attach claims to request context
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
