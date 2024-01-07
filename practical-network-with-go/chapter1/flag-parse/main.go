// feature:
// 1.自定义flag的Usage，flag解析出错、help查询会调用Usage方法
// 2.去掉重复的错误打印
// 3.支持可选参数确定greet的人
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type config struct {
	numTimes int
	name     string
}

func parseArgs(w io.Writer, args []string) (c config, err error) {
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		var usage = `
A greeter application which prints the name you entered a specified number of times.

Usage of %s: <options> [name] `
		fmt.Fprintf(w, usage, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	if err = fs.Parse(args); err != nil {
		return c, err
	}

	if fs.NArg() > 1 {
		return c, errPosArgSpecified
	}
	if fs.NArg() == 1 {
		c.name = fs.Args()[0]
	}
	return c, nil
}

var errPosArgSpecified = errors.New("more than one positional argument specified")

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

func runCmd(r io.Reader, w io.Writer, c config) (err error) {
	if c.name == "" {
		c.name, err = getName(r, w)
		if err != nil {
			return err
		}
	}

	greetUser(c, w)
	return nil
}

func greetUser(c config, w io.Writer) error {
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprintf(w, "Nice to meet you %s!\n", c.name)
	}

	return nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		return
	}

	if err = validateArgs(c); err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s\n", err)
		os.Exit(1)
	}

	if err = runCmd(os.Stdin, os.Stdout, c); err != nil {
		fmt.Fprintf(os.Stdout, "Error: %s\n", err)
		os.Exit(1)
	}
}
