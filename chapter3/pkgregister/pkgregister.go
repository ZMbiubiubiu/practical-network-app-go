package pkgregister

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type pkgRegisterResult struct {
	Id string `json:"id"`
}

func registerPackageData(url string, pkg pkgData) (pkgRegisterResult, error) {
	var p pkgRegisterResult

	b, err := json.Marshal(pkg)
	if err != nil {
		return p, err
	}
	body := bytes.NewReader(b)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return p, err
	}
	if resp.StatusCode != http.StatusOK {
		return p, errors.New(string(respData))
	}

	err = json.Unmarshal(respData, &p)

	return p, err
}
