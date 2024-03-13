package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func startHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			writer.Header().Set(k, v[0])
		}
		fmt.Fprintf(writer, "I am the Request Header echoing program")
	}))
}

func TestAddHeaderMiddleware(t *testing.T) {
	testHeaders := map[string]string{
		"X-Client-Id": "test-client",
		"X-Auth-Hash": "random$string",
	}

	ts := startHTTPServer()
	defer ts.Close()

	client := createClient(testHeaders)
	resp, err := client.Get(ts.URL)
	assert.Nil(t, err)

	for k := range testHeaders {
		assert.Equal(t, testHeaders[k], resp.Header.Get(k))
	}
}
