package store

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/sant470/trademark/common"
	"github.com/sant470/trademark/dtos"
	"go.uber.org/zap"
)

type StoreImpl struct {
	rdb *redis.Client
	lgr *zap.SugaredLogger
}

func NewStore(lgr *zap.SugaredLogger, rdb *redis.Client) *StoreImpl {
	return &StoreImpl{rdb, lgr}
}

func (impl *StoreImpl) CheckUser(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := userKey(username)
	val, err := impl.rdb.Exists(ctx, key).Result()
	return val == 0, err
}

func (impl *StoreImpl) AddUser(user *dtos.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := userKey(user.UserName)
	_, err := impl.rdb.HMSet(ctx, key, common.StructToMap(user)).Result()
	return err
}

func (impl *StoreImpl) GetUserDetails(username string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	key := userKey(username)
	return impl.rdb.HGetAll(ctx, key).Result()
}
