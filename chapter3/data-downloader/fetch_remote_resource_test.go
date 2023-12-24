package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestHTTPServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello World")
			}))

	return ts
}

func TestDownloadRemoteResource(t *testing.T) {
	expected := "Hello World"

	ts := setupTestHTTPServer()

	data, err := fetchRemoteResource(ts.URL)
	if err != nil {
		t.Fatalf("expected nil, but got %v\n", err)
	}
	if expected != string(data) {
		t.Fatalf("expected %s, but got %v\n", expected, string(data))
	}
}
