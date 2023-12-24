package pkgregister

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type pkgRegisterResult struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func registerPackageData(client *http.Client, url string, data pkgData) (pkgRegisterResult, error) {
	var res pkgRegisterResult

	bs, contentType, err := createMultiPartMessage(data)
	if err != nil {
		return res, err
	}

	body := bytes.NewReader(bs)
	resp, err := client.Post(url, contentType, body)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(respBody, &res)
	return res, err
}

func createHTTPClientWithTimeout(d time.Duration) *http.Client {
	client := http.Client{Timeout: d}
	return &client
}
