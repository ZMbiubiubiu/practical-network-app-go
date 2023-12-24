package pkgregister

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlePkgRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method specified", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid Body type", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var pkg pkgData
	err = json.Unmarshal(data, &pkg)
	if err != nil || pkg.Version == "" || pkg.Name == "" {
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	resp := pkgRegisterResult{Id: pkg.Name + "-" + pkg.Version}

	data, err = json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
	return
}

func TestRegisterPackageData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(handlePkgRegister))
	defer ts.Close()

	pkg := pkgData{
		Name:    "goland",
		Version: "1.18",
	}
	expectedId := fmt.Sprintf("%s-%s", pkg.Name, pkg.Version)
	pkgResult, err := registerPackageData(ts.URL, pkg)
	if err != nil {
		t.Fatalf("expected nil, but got %v\n", err)
	}

	if pkgResult.Id != expectedId {
		t.Fatalf("expected %s, but got %s\n", expectedId, pkgResult.Id)
	}
}
