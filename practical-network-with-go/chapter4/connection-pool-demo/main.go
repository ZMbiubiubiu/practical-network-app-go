package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"
)

func createHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	client := &http.Client{Timeout: timeout}
	return client
}

func createHTTPGetRequestWithTrace(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return req, err
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			fmt.Printf("GotConn info: %+v\n", info)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			fmt.Printf("DNSDone info: %+v\n", info)
		},
	}
	traceCtx := httptrace.WithClientTrace(ctx, trace)
	req = req.WithContext(traceCtx)
	return req, nil
}

func main() {
	d := 5 * time.Second
	ctx := context.Background()
	client := createHTTPClientWithTimeout(d)
	req, err := createHTTPGetRequestWithTrace(ctx, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for {
		client.Do(req)
		time.Sleep(1 * time.Second)
		fmt.Println("--------")
	}
}
