package pkgquery

import (
	"encoding/json"
	"io"
	"net/http"
)

type pkgData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchPackageData(url string) ([]pkgData, error) {
	var packages []pkgData

	resp, err := http.Get(url)
	if err != nil {
		return packages, err
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/json" {
		return packages, nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return packages, err
	}
	err = json.Unmarshal(data, &packages)
	return packages, err
}
