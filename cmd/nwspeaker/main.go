package main

import (
	"fmt"
	"os"

	"github.com/mas9612/nwspeaker/pkg/command"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("nwspeaker", "0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"arp": func() (cli.Command, error) {
			return &command.ArpResolverCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(exitStatus)
}
