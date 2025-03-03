package v1

import (
	"github.com/sant470/trademark/services"
	"go.uber.org/zap"
)

type RegistrationHlr struct {
	lgr *zap.SugaredLogger
	svc *services.RegistrationSvc
}

func NewRegistrationHlr(lgr *zap.SugaredLogger, svc *services.RegistrationSvc) *RegistrationHlr {
	return &RegistrationHlr{lgr, svc}
}

// func (sh *SearchHandler) Search(rw http.ResponseWriter, r *http.Request) *errors.AppError {
// 	var req apptypes.SearchReq
// 	if err := common.Decode(r, &req); err != nil {
// 		return err
// 	}
// 	if req.SearchKeyword == "" || req.From > req.To {
// 		return errors.BadRequest("invalid params")
// 	}
// 	result, err := sh.svc.Search(&req)
// 	if err != nil {
// 		return err
// 	}
// 	return respond.OK(rw, result)
// }
