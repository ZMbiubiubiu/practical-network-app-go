package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startBadTestHTTPServerV1() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			time.Sleep(10 * time.Second)
			fmt.Fprintf(writer, "Hello World")
		}))
}

func TestFetchBadRemoteResource(t *testing.T) {
	ts := startBadTestHTTPServerV1()
	defer ts.Close()

	client := &http.Client{Timeout: 1 * time.Second}

	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Fatalf("Expected non-nil error， but got:%v", err)
	}
}

func startBadTestHTTPServerV2(signal chan struct{}) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			<-signal
			fmt.Fprintf(writer, "Hello World")
		}))
}

func TestFetchBadRemoteResourceV2(t *testing.T) {
	var signal = make(chan struct{})
	ts := startBadTestHTTPServerV2(signal)
	defer ts.Close()
	defer func() {
		signal <- struct{}{}
	}()

	client := &http.Client{Timeout: 1 * time.Second}

	_, err := fetchRemoteResource(client, ts.URL)
	if err == nil {
		t.Fatalf("Expected non-nil error， but got:%v", err)
	}
}
