package apis

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	v1 "github.com/sant470/trademark/apis/v1"
	"github.com/sant470/trademark/common"
	"github.com/sant470/trademark/config"
)

func InitRegistrationHlr(r *chi.Mux, rh *v1.RegistrationHlr) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/register", common.Handler(rh.RegisterHlr))
		r.Method(http.MethodPost, "/login", common.Handler(rh.LoginHlr))
		r.With(config.JWTAuthMiddleware).Method(http.MethodGet, "/admin/data", common.Handler(rh.AdminDataHlr))
		r.With(config.JWTAuthMiddleware).Method(http.MethodGet, "/user/data", common.Handler(rh.UserDataHlr))
	})
}
