package config

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	REDIS      Redis
	JWT        string
	RateLimits []RateLimitPolicy
}

type Redis struct {
	HOST     string
	PORT     string
	PASSWORD string
}

type RateLimitPolicy struct {
	Domain    string
	Value     string
	RateLimit struct {
		Unit  string
		Limit int
	}
}

var (
	appConfig *Config
	rdb       *redis.Client
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
		RateLimitMiddleware,
	)
	return r
}

func loadConfig(filename, path string) *Config {
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("fatal: error reading config file", err.Error())
	}
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("fatal: error reading config variable", err.Error())
	}
	return &conf
}

// GetDBConn ...
func GetDBConn(log *zap.SugaredLogger, conf Redis) *redis.Client {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", conf.HOST, conf.PORT),
			Password: conf.PASSWORD,
			DB:       0,
		})
	}
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

// RateLimitMiddleware uses Redis to enforce rate limits based on URL domain and IP address.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Loop through policies and check if the request path matches any.
		for _, policy := range appConfig.RateLimits {
			fmt.Println("domain: ", policy.Domain)
			if strings.HasSuffix(path, policy.Domain) {
				ip, _, _ := net.SplitHostPort(r.RemoteAddr)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				key := "ratelimit:" + policy.Domain + ":" + ip

				// Determine the time window based on the policy's Unit.
				var ttl time.Duration
				switch policy.RateLimit.Unit {
				case "minute":
					ttl = time.Minute
				case "second":
					ttl = time.Second
				case "hour":
					ttl = time.Hour
				default:
					ttl = time.Minute
				}
				fmt.Println("ttl: ", ttl)
				// Increment the counter for this key.
				count, err := rdb.Incr(ctx, key).Result()
				fmt.Println("after: ", count)
				if err != nil {
					http.Error(w, "Rate limiter error", http.StatusInternalServerError)
					return
				}
				// If this is the first request, set the TTL.
				if count == 1 {
					rdb.Expire(ctx, key, ttl)
				}

				// If the count exceeds the limit, reject the request.
				fmt.Println("limit: ", policy.RateLimit.Limit)
				if int(count) > policy.RateLimit.Limit {
					fmt.Println("we are here: ", count)
					retryAfter, _ := rdb.TTL(ctx, key).Result()
					w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
				// Only enforce one matching policy.
				break
			}
		}
		next.ServeHTTP(w, r)
	})
}
