package store

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type StoreImpl struct {
	rdb *redis.Client
	lgr *zap.SugaredLogger
}

func NewStore(lgr *zap.SugaredLogger, rdb *redis.Client) *StoreImpl {
	return &StoreImpl{rdb, lgr}
}
