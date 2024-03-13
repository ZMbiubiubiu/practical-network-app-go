package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestHTTPServer() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(60 * time.Second)
				fmt.Fprintf(w, "Hello World")
			}))

	return ts
}

func TestDownloadRemoteResource(t *testing.T) {
	expected := "Hello World"

	ts := setupTestHTTPServer()
	defer ts.Close()

	data, err := fetchRemoteResource(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(data))
}
