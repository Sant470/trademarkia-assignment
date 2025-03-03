package store

import (
	"github.com/go-redis/redis"
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
	key := userKey(username)
	val, err := impl.rdb.Exists(key).Result()
	return val == 0, err
}

func (impl *StoreImpl) AddUser(user *dtos.RegisterRequest) error {
	key := userKey(user.UserName)
	_, err := impl.rdb.HMSet(key, common.StructToMap(user)).Result()
	return err
}

func (impl *StoreImpl) GetUserDetails(username string) (map[string]string, error) {
	key := userKey(username)
	return impl.rdb.HGetAll(key).Result()
}
