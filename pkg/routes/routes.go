package routes

import (
	"github.com/araminian/argo-appset-pr-label-filter/pkg/controller"
	"github.com/gorilla/mux"
)

var RegisterPluginRoutes = func(router *mux.Router) {
	router.HandleFunc("/api/v1/getparams.execute", controller.FilterPRs).Methods("POST")
}
