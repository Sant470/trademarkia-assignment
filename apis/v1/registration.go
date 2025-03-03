package v1

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sant470/trademark/common"
	"github.com/sant470/trademark/common/errors"
	"github.com/sant470/trademark/common/respond"
	"github.com/sant470/trademark/dtos"
	"github.com/sant470/trademark/services"
	"go.uber.org/zap"
)

type RegistrationHlr struct {
	lgr *zap.SugaredLogger
	svc services.RegistrationSvc
}

func NewRegistrationHlr(lgr *zap.SugaredLogger, svc services.RegistrationSvc) *RegistrationHlr {
	return &RegistrationHlr{lgr, svc}
}

func (rh *RegistrationHlr) RegisterHlr(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	var reg dtos.RegisterRequest
	if err := common.Decode(r, &reg); err != nil {
		return errors.BadRequest("invalid params")
	}
	res, err := rh.svc.RegisterSvc(&reg)
	if err != nil {
		return err
	}
	return respond.OK(rw, res)
}

func (rh *RegistrationHlr) LoginHlr(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	var cred dtos.LoginRequest
	if err := common.Decode(r, &cred); err != nil {
		return errors.BadRequest("invalid params")
	}
	res, err := rh.svc.LoginSvc(&cred)
	if err != nil {
		return err
	}
	return respond.OK(rw, res)
}

func (rh *RegistrationHlr) AdminDataHlr(rw http.ResponseWriter, r *http.Request) *errors.AppError {
	claims, _ := r.Context().Value("userClaims").(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "admin" {
		return errors.Unauthorized("Unauthorized Access")
	}
	return respond.OK(rw, map[string]interface{}{"data": "Sensitive admin information"})
}
