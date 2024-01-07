package pkgquery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupPackageDataHTTPTestServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(writer, `[
{"name": "python", "version":"3.12"},
{"name": "goland", "version":"1.18"}
]`)
		}))
}

func TestFetchPackageData(t *testing.T) {
	ts := setupPackageDataHTTPTestServer()
	defer ts.Close()

	pkgs, err := fetchPackageData(ts.URL)
	if err != nil {
		t.Fatalf("expected nil, but got %v\n", err)
	}

	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, but got %d\n", len(pkgs))
	}
}
