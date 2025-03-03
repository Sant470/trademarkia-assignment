package main

import (
	"net/http"

	"github.com/sant470/trademark/apis"
	handlers "github.com/sant470/trademark/apis/v1"
	"github.com/sant470/trademark/config"
	"github.com/sant470/trademark/services"
	"github.com/sant470/trademark/store"
)

func main() {
	conf := config.GetAppConfig("config.yaml", "./")
	lgr := config.GetConsoleLogger()
	rdb := config.GetDBConn(lgr, conf.REDIS)
	store := store.NewStore(lgr, rdb)
	registrationSvc := services.NewRegistrationSvc(lgr, store)
	registrationHlr := handlers.NewRegistrationHlr(lgr, registrationSvc)
	router := config.InitRouters()
	apis.InitRegistrationHlr(router, registrationHlr)
	http.ListenAndServe("localhost:8000", router)
}
