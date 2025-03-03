package services

import (
	"github.com/sant470/trademark/common/errors"
	"github.com/sant470/trademark/dtos"
)

type RegistrationSvc interface {
	RegisterSvc(reg *dtos.RegisterRequest) (*dtos.RegisterResponse, *errors.AppError)
	LoginSvc(cred *dtos.LoginRequest) (*dtos.LoginResponse, *errors.AppError)
}
