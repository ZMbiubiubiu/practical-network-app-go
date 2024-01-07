package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func createContextWithTimout(duration time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	return ctx, cancel
}

func executeCommand(ctx context.Context, cmd, arg string) error {
	return exec.CommandContext(ctx, cmd, arg).Run()
}

func setupSignalHandler(w io.Writer, cancel context.CancelFunc) {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		fmt.Fprintf(w, "Got signal: %v\n", s.String())
		cancel()
	}()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stdout, "Usage: %s <command> <argument>\n", os.Args[0])
		os.Exit(1)
	}
	command, arg := os.Args[1], os.Args[2]

	ctx, cancel := createContextWithTimout(5 * time.Second)
	defer cancel()

	setupSignalHandler(os.Stdout, cancel)

	err := executeCommand(ctx, command, arg)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
