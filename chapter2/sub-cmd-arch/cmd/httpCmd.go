package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	var fs = flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&v, "verb", "GET", "HTTP method")
	fs.Usage = func() {
		var usage = `
http: A HTTP client.

http: <options> server`
		fmt.Fprintf(w, usage)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	c := httpConfig{verb: v}
	c.url = fs.Arg(0)

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	fmt.Fprintln(w, "Executing http command")

	switch v {
	case "GET":
		data, err := fetchRemoteResource(c.url)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Got \n: %s\n", string(data))
	case "POST", "HEAD":
	default:
		return errors.New("invalid HTTP method")
	}
	return nil
}

func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}
