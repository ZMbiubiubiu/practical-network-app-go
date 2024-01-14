package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sleepHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("I will begin to sleep")
	time.Sleep(15 * time.Second)
	fmt.Fprintf(w, "Hello World")
	log.Println("I am done the work")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/sleep", http.TimeoutHandler(http.HandlerFunc(sleepHandler), 14*time.Second, "I ran out of time"))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
