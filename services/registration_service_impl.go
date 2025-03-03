package services

import (
	"go.uber.org/zap"
)

type RegistrationSvcImpl struct {
	lgr   *zap.SugaredLogger
	store *store.Store
}

func NewRegistrationSvc(lgr *zap.SugaredLogger, store *store.Store) *RegistrationSvcImpl {
	return &RegistrationSvcImpl{lgr, store}
}
