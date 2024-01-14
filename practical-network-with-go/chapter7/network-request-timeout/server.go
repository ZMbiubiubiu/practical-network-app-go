package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("begin to process apiHandler")
	doSomeWorkd()

	ctx := r.Context()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:8080/ping", nil)
	if err != nil {
		log.Printf("Error building request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Printf("DNSStart info:%v\n", info)
		},
		DNSDone: nil,
		ConnectStart: func(network, addr string) {
			log.Printf("ConnectStart network:%s addr:%s\n", network, addr)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			log.Printf("WroteRequest info:%v\n", info)
		},
	}
	traceCtx := httptrace.WithClientTrace(r.Context(), trace)
	req = req.WithContext(traceCtx)

	log.Println("Outgoing HTTP request")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error doing request: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Fprintf(w, string(body))
	log.Println("I finished processing the request")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	fmt.Fprintf(w, "pong from ther server")
}

func doSomeWorkd() {
	time.Sleep(2 * time.Second)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.Handle("/api", http.TimeoutHandler(http.HandlerFunc(apiHandler), 3*time.Second, "ran out of time"))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
