package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Whey do we need such a type at all?
// It enables us to write a function that wraps around any other http.Handler value
// and returns another http.Handler
func logTimeHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			handler.ServeHTTP(w, r)
			fmt.Printf("path=%s method=%s duration=%f\n", r.URL.Path, r.Method, time.Since(startTime).Seconds())
		})
}

func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("handle recover, ", err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "Unexpected server error")
				}
			}()
			handler.ServeHTTP(w, r)
		})
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API handler, hello world")
}

func createPanicHandler(w http.ResponseWriter, r *http.Request) {
	panic("divide by zero")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "ok")
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/panic", createPanicHandler)

	handler := logTimeHandler(panicHandler(mux))
	log.Fatal(http.ListenAndServe(listenAddr, handler))
}
