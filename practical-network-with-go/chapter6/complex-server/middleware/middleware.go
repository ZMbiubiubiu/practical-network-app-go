package middleware

import (
	"fmt"
	"github.com/bingo/complex-server/config"
	"net/http"
	"time"
)

func loggingMiddleware(h http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			h.ServeHTTP(w, r)
			requestDuration := time.Now().Sub(t1).Seconds()
			c.Logger.Printf("protocol=%s path=%s method=%s duration =%f", r.Proto, r.URL.Path, r.Method, requestDuration)
		})
}

func panicHandler(handler http.Handler, c config.AppConfig) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					c.Logger.Println("panic detected", err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Unexpected server error")
				}
			}()
			handler.ServeHTTP(w, r)
		})
}
