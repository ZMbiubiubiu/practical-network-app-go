package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type appConfig struct {
	logger *log.Logger
}

type app struct {
	config  appConfig
	handler func(w http.ResponseWriter, r *http.Request, config appConfig)
}

func (a app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r, a.config)
}

func apiHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	config.logger.Println("Handing API request")
	fmt.Fprintf(w, "API handler, hello world")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config.logger.Println("Handling healthcheck request")
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
