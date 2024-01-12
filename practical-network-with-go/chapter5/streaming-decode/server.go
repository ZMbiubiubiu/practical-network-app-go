package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type logEntry struct {
	UserIP string `json:"user_ip"`
	Event  string `json:"event"`
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var dec = json.NewDecoder(r.Body)

	var e *json.UnmarshalTypeError

	for {
		var log logEntry
		err := dec.Decode(&log)
		if err == io.EOF {
			break
		}

		if errors.As(err, &e) {
			fmt.Println(err)
			continue
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("receive log entry:%s %s\n", log.UserIP, log.Event)
	}
	fmt.Fprintf(w, "OK")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/decode", decodeHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
