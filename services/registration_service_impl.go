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

// CheckPassword verifies a password against a stored hash
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Returns true if passwords match
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

func (impl *RegistrationSvcImpl) LoginSvc(cred *dtos.LoginRequest) (*dtos.LoginResponse, *errors.AppError) {
	user, err := impl.store.GetUserDetails(cred.UserName)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	ok := checkPassword(cred.Password, user["password"])
	if !ok {
		return nil, errors.Unauthorized("invalid username or password")
	}
	token, err := jwt.GenerateJWT(user["username"], user["role"])
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	return &dtos.LoginResponse{APIKey: token}, nil
}
