package handlers

import (
	"github.com/bingo/complex-server/config"
	"net/http"
)

func Register(mux *http.ServeMux, conf config.AppConfig) {
	mux.Handle(
		"/health",
		&app{conf: conf, handler: healthCheckHandler},
	)
	mux.Handle(
		"/api",
		&app{conf: conf, handler: apiHandler},
	)
	mux.Handle(
		"/panic",
		&app{conf: conf, handler: panicHandler},
	)
}
