package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	var tests = []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c: config{
				printUsage: true,
			},
			input:  "",
			output: usage,
			err:    nil,
		},
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
