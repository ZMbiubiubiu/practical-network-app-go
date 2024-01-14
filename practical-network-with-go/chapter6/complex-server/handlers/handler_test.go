package handlers

import (
	"bytes"
	"github.com/bingo/complex-server/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api", nil)
	conf := config.InitConfig(bytes.NewBuffer(nil))

	apiHandler(w, r, conf)
	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, string(body), "Hello, world!")
}
