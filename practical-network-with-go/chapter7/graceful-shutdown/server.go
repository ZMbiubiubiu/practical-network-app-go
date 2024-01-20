package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "pong")
}

func shutdown(ctx context.Context, server *http.Server, waitForTerminateCh chan struct{}) {
	var sigCh = make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Got signal: %v. Server shutting down", sig)

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	waitForTerminateCh <- struct{}{}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       100 * time.Second,
	}

	var waitForTerminateCh = make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go shutdown(ctx, server, waitForTerminateCh)

	// we don't call this function inside a call to log.Fatal like before
	// because when the Shutdown() is called, the ListenAndServe will immediately return
	// and hence the server will exit without waiting for the Shutdown()
	err := server.ListenAndServe()
	log.Print("Waiting for shutdown to complete...")

	<-waitForTerminateCh
	log.Fatal(err)
}
