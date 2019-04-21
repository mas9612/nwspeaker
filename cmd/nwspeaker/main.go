package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("nwspeaker", "0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(exitStatus)
}
