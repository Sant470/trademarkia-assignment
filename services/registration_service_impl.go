package services

import (
	"github.com/sant470/trademark/common/errors"
	"github.com/sant470/trademark/dtos"
	"github.com/sant470/trademark/lib/jwt"
	"github.com/sant470/trademark/store"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationSvcImpl struct {
	lgr   *zap.SugaredLogger
	store store.Store
}

func NewRegistrationSvc(lgr *zap.SugaredLogger, store store.Store) *RegistrationSvcImpl {
	return &RegistrationSvcImpl{lgr, store}
}

func (impl *RegistrationSvcImpl) RegisterSvc(reg *dtos.RegisterRequest) (*dtos.RegisterResponse, *errors.AppError) {
	ok, err := impl.store.CheckUser(reg.UserName)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if !ok {
		return nil, errors.BadRequest("user already exist")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	reg.Password = string(hashedBytes)
	err = impl.store.AddUser(reg)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	token, err := jwt.GenerateJWT(reg.UserName, reg.Role)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	return &dtos.RegisterResponse{APIKey: token, Message: "User registered successfully"}, nil
}
