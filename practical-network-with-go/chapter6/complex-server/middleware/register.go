package middleware

import (
	"github.com/bingo/complex-server/config"
	"net/http"
)

func RegisterMiddleware(mux *http.ServeMux, c config.AppConfig) http.Handler {
	return loggingMiddleware(panicHandler(mux, c), c)
}
