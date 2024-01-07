package main

import "net/http"

type AddHeaderMiddleware struct {
	headers map[string]string
}

func (a *AddHeaderMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	copyReq := r.Clone(r.Context())
	for k, v := range a.headers {
		copyReq.Header.Add(k, v)
	}
	return http.DefaultTransport.RoundTrip(copyReq)
}

func createClient(headers map[string]string) *http.Client {
	h := &AddHeaderMiddleware{headers: headers}
	client := &http.Client{Transport: h}
	return client
}
