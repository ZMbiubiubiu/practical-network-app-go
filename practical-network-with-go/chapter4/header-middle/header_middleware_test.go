package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			writer.Header().Set(k, v[0])
		}
		fmt.Fprintf(writer, "I am the Request Header echoing program")
	}))
}

func TestAddHeaderMiddleware_RoundTrip(t *testing.T) {
	testHeaders := map[string]string{
		"X-Client-Id": "test-client",
		"X-Auth-Hash": "random$string",
	}

	ts := startHTTPServer()
	defer ts.Close()

	client := createClient(testHeaders)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected nil, but got %v\n", err)
	}

	for k := range testHeaders {
		if resp.Header.Get(k) != testHeaders[k] {
			t.Fatalf("testHeaders %s expected get %s, but got %s\n", k, testHeaders[k], resp.Header.Get(k))
		}
	}
}
