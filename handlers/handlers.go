package handlers

import (
	"github.com/gorilla/mux"
	"github.com/simplyianm/httputil"
	"github.com/simplyianm/inject"
)

// Router returns a router built from the injector
func Router(injector inject.Injector) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/spider", httputil.MakeHandler(newSpiderRequest(injector)))
	return r
}
