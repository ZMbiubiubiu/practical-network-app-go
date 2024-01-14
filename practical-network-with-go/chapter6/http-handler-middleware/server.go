package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	a.handler(w, r)
	a.config.logger.Printf("path=%s method=%s duration=%f", r.URL.Path, r.Method, time.Since(startTime).Seconds())
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API handler, hello world")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "ok")
}

func setupHandlers(mux *http.ServeMux, config appConfig) {
	mux.Handle("/api", &app{
		config:  config,
		handler: apiHandler,
	})
	mux.Handle("/health", &app{
		config:  config,
		handler: healthCheckHandler,
	})
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	config := appConfig{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)}
	mux := http.NewServeMux()
	setupHandlers(mux, config)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
