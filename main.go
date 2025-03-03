package main

import (
	"net/http"

	"github.com/sant470/trademark/apis"
	handlers "github.com/sant470/trademark/apis/v1"
	"github.com/sant470/trademark/config"
	"github.com/sant470/trademark/services"
)

func main() {
	conf := config.GetAppConfig("config.yaml", "./")
	lgr := config.GetConsoleLogger()
	rdb := config.GetDBConn(lgr, conf.REDIS)
	router := config.InitRouters()
	searchSvc := services.NewSearchService(lgr)
	searchHlr := handlers.NewSearchHandler(lgr, searchSvc)
	apis.InitSerachRoutes(router, searchHlr)
	http.ListenAndServe("localhost:8000", router)
}
