package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type httpConfig struct {
	url    string
	output string
	verb   string
}

func HandleHttp(w io.Writer, args []string) error {
	c := httpConfig{}
	var fs = flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&c.output, "output", "", "result save to file-path")
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

	c.url = fs.Arg(0)

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	fmt.Fprintln(w, "Executing http command")

	switch c.verb {
	case "GET":
		data, err := fetchRemoteResource(c.url)
		if err != nil {
			return err
		}
		if c.output != "" {
			fmt.Fprintf(w, "should write data to file, path is %s\n", c.output)
		} else {
			fmt.Fprintf(w, "Got \n: %s\n", string(data))
		}
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
