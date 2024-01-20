package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api", nil)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	r.WithContext(ctx)

	apiHandler(w, r)
}
