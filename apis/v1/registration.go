package v1

import (
	"net/http"

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
