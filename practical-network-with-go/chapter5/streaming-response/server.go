package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func longRunningProcessHandler(w http.ResponseWriter, r *http.Request) {
	var done = make(chan struct{})
	pipeReader, pipeWriter := io.Pipe()
	go longRunningProcess(pipeWriter)
	go progressStreamer(pipeReader, w, done)
	<-done
}

func longRunningProcess(writer *io.PipeWriter) {
	for i := 0; i <= 20; i++ {
		fmt.Fprintf(writer, "longRunningProcess :%d", i)
		fmt.Fprintln(writer)
		time.Sleep(time.Second)
	}
	writer.Close()
}

func progressStreamer(reader *io.PipeReader, w http.ResponseWriter, done chan struct{}) {
	var buf = make([]byte, 500)

	f, isSupportFlush := w.(http.Flusher)

	defer reader.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		w.Write(buf[:n])
		if isSupportFlush {
			f.Flush()
		}
	}

	done <- struct{}{}
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/job", longRunningProcessHandler)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}
