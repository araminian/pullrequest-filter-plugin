package main

import (
	"log"
	"net/http"

	"github.com/araminian/argo-appset-pr-label-filter/pkg/routes"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	routes.RegisterPluginRoutes(r)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":4355", r))

}
