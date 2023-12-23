package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		config
		args   []string
		output string
		err    error
	}{
		{
			args:   []string{"-h"},
			config: config{numTimes: 0},
			output: `
A greeter application which prints the name you entered a specified number of times.

Usage of greeter: <options> [name]`,
			err: errors.New("flag: help requested"),
		},
		{
			args:   []string{"-n", "10"},
			config: config{numTimes: 10},
			output: "",
			err:    nil,
		},
		{
			config: config{numTimes: 0},
			args:   []string{"-n", "abc"},
			err: errors.New("invalid value \"abc\" for flag -n: " +
				"parse error"),
		},
		{
			config: config{numTimes: 10, name: "bingo"},
			output: "",
			args:   []string{"-n", "10", "bingo"},
			err:    nil,
		},
		{
			config: config{numTimes: 1, name: ""},
			args:   []string{"-n", "1", "foo", "ok"},
			err:    errors.New("more than one positional argument specified"),
		},
	}
	byteBuf := new(bytes.Buffer)
	for _, tt := range tests {
		c, err := parseArgs(byteBuf, tt.args)
		if err != nil && tt.err == nil {
			t.Fatalf("expected nil, but got: %v\n", err)
		}
		if err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("expected %v, but got: %v\n", tt.err, err)
		}
		if c.numTimes != tt.numTimes {
			t.Fatalf("expected %d, but got %d\n", tt.numTimes, c.numTimes)
		}
		if c.name != tt.name {
			t.Fatalf("expected %s, but got %s\n", tt.name, c.name)
		}
		if tt.output != "" && !strings.Contains(byteBuf.String(), tt.output) {
			t.Fatalf("expected %s, but got %s\n", tt.output, byteBuf.String())
		}
		byteBuf.Reset()
	}
}

func TestRunCmd(t *testing.T) {
	var tests = []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:    errors.New("you didn't enter your name"),
		},
		{
			c:     config{numTimes: 5},
			input: "bingo",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1) +
				strings.Repeat("Nice to meet you bingo!\n", 5),
			err: nil,
		},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.input)
		w := new(bytes.Buffer)
		err := runCmd(r, w, tt.c)
		if err != nil && tt.err == nil {
			t.Fatalf("Expected nil error, but got: %v\n", err)
		}
		if err != nil && err.Error() != tt.err.Error() {
			t.Fatalf("Expected error %q, but got %q\n", tt.err, err)
		}
		if w.String() != tt.output {
			t.Errorf("runCmd(%q) output = %q, want %q", tt.input, w.String(), tt.output)
		}
	}
}
