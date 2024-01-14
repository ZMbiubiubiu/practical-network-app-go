package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSetupServer(t *testing.T) {
	b := new(bytes.Buffer)
	mux := http.NewServeMux()
	wrappedMux := setupServer(mux, b)

	ts := httptest.NewServer(wrappedMux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/panic")
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "Unexpected server error")

	logs := b.String()
	expectedLogFragments := []string{
		"path=/panic method=GET duration=",
		"panic detected",
	}

	for _, log := range expectedLogFragments {
		if !strings.Contains(logs, log) {
			t.Errorf("Expected logs to contain: %s, Got: %s", log, logs)
		}
	}
}
