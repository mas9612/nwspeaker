package main

import (
	"log"
	"os"

	"github.com/mas9612/nwspeaker/pkg/command"
	"github.com/mitchellh/cli"
)

func main() {
	c := cli.NewCLI("craftpkt", "0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"arp": func() (cli.Command, error) {
			return &command.ArpCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitStatus)
}
