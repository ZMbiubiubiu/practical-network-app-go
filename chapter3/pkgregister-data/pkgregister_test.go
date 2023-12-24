package pkgregister

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func handlePackageRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//if r.Header.Get("Content-Type") != "application/multi"
	d := pkgRegisterResult{}
	err := r.ParseMultipartForm(5000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mForm := r.MultipartForm
	f := mForm.File["filedata"][0]
	d.Filename = f.Filename
	d.Size = f.Size
	d.Id = fmt.Sprintf("%s-%s", mForm.Value["name"][0], mForm.Value["version"][0])

	data, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(data))
}

func startTestPackageServer() *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(handlePackageRegister))
	return ts
}

func TestRegisterPackageData(t *testing.T) {
	ts := startTestPackageServer()
	defer ts.Close()

	input := pkgData{
		Name:     "golang",
		Version:  "1.8",
		Filename: "golang-explain",
		Bytes:    strings.NewReader("data"),
	}

	client := createHTTPClientWithTimeout(10 * time.Second)

	got, err := registerPackageData(client, ts.URL, input)
	if err != nil {
		t.Fatalf("expected nil, but got %v\n", err)
	}

	if got.Filename != input.Filename {
		t.Fatalf("expected Filename:%s, but got %s\n", input.Filename, got.Filename)
	}

	if got.Size != 4 {
		t.Fatalf("expected Size:%d, but got %d\n", 4, got.Size)

	}

	expectedId := fmt.Sprintf("%s-%s", input.Name, input.Version)
	if got.Id != expectedId {
		t.Fatalf("expected id:%s, but got %s\n", expectedId, got.Id)
	}
}
