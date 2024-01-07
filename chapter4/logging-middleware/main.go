package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LoggingClient struct {
	log *log.Logger
}

func (l *LoggingClient) RoundTrip(req *http.Request) (*http.Response, error) {
	l.log.Printf("Sending a %s request to %s by %s\n", req.Method, req.URL, req.Proto)
	r, err := http.DefaultTransport.RoundTrip(req)
	l.log.Printf("Got back a response over %s\n", r.Proto)

	return r, err
}

func fetchRemoteResource(client *http.Client, url string) ([]byte, error) {
	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Must specify a HTTP URL to get data from")
		os.Exit(1)
	}
	myTransport := &LoggingClient{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: myTransport,
	}
	body, err := fetchRemoteResource(client, os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%#v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "Bytes in response: %d\n", len(body))
}
