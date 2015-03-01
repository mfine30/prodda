package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfine30/prodda/api/middleware"
	"github.com/mfine30/prodda/api/v0"
	"github.com/mfine30/prodda/registry"
	"github.com/pivotal-golang/lager"
	"gopkg.in/robfig/cron.v2"
)

var HomeHandleFunc = homeHandleFunc

func NewHandler(
	logger lager.Logger,
	username, password string,
	prodRegistry registry.ProdRegistry,
	c *cron.Cron) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandleFunc)
	api := r.PathPrefix("/api").Subrouter()
	v0.NewSubrouter(api, prodRegistry, c, logger)

	return middleware.Chain{
		middleware.NewPanicRecovery(logger),
		middleware.NewLogger(logger),
		middleware.NewBasicAuth(username, password),
	}.Wrap(r)
}

func homeHandleFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Prodda")
}
