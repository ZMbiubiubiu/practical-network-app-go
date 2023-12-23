package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type httpConfig struct {
	url  string
	verb string
}

func HandleHttp(w io.Writer, args []string) error {
	var v string
	var fs = flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	flag.StringVar(&v, "verb", "GET", "HTTP method")
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

	switch v {
	case "GET", "POST", "HEAD":
	default:
		return errors.New("invalid HTTP method")
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	c := httpConfig{verb: v}
	c.url = fs.Arg(0)
	fmt.Fprintln(w, "Executing http command")
	return nil
}
