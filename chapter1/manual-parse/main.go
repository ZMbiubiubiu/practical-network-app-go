package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type config struct {
	numTimes   int
	printUsage bool
}

var usage = fmt.Sprintf(`Usage: %s <integer> [-h|--help]\n

A Greeter application which prints the name you entered  <integer> number of time.
`, os.Args[0])

func printUsage(w io.Writer) {
	fmt.Fprintf(w, usage)
}

func validateArgs(c config) error {
	if c.numTimes <= 0 {
		return fmt.Errorf("number of times must be greater than 0")
	}
	return nil
}

func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprint(w, msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	name := scanner.Text()
	if len(name) == 0 {
		return "", fmt.Errorf("you didn't enter your name")
	}

	return name, nil
}

func parseArgs(args []string) (config, error) {
	c := config{}
	if len(args) != 1 {
		return c, fmt.Errorf("unexpected argument")
	}

	if args[0] == "-h" || args[0] == "--help" {
		c.printUsage = true
		return c, nil
	}

	numTimes, err := strconv.Atoi(args[0])
	if err != nil {
		return c, fmt.Errorf("invalid number of times")
	}
	c.numTimes = numTimes

	return c, nil
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	if c.printUsage {
		printUsage(w)
		return nil
	}

	name, err := getName(r, w)
	if err != nil {
		return err
	}

	greetUser(c, name, w)

	return nil
}

func greetUser(c config, name string, w io.Writer) error {
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, "Nice to meet you %s!\n", name)
	}

	return nil

}

func main() {
	c, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		printUsage(os.Stdout)
		os.Exit(1)
	}

	if err = validateArgs(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		printUsage(os.Stdout)
		os.Exit(1)
	}

	if err = runCmd(os.Stdin, os.Stdout, c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
