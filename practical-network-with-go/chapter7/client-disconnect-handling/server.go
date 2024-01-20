package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ping: Got a request")
	time.Sleep(10 * time.Second)
	fmt.Fprintf(w, "pong")
}

func doSomeWork(data []byte) {
	time.Sleep(15 * time.Second)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("begin to process apiHandler")

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

	var done = make(chan struct{})
	go func() {
		doSomeWork(body)
		done <- struct{}{}
	}()

	select {
	case <-done:
		log.Println("Processed the outgoing response data")
	case <-r.Context().Done():
		log.Printf("Aborting request processing: %v\n", r.Context().Err())
		return
	}

	fmt.Fprintf(w, string(body))
	log.Println("I finished processing the request")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.Handle("/api", http.TimeoutHandler(http.HandlerFunc(apiHandler), 30*time.Second, "ran out of time"))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
