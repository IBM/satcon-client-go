package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.ibm.com/coligo/satcon-client/cli"
	"github.ibm.com/coligo/satcon-client/client"
)

const (
	DefaultEndpoint = "https://config.satellite.test.cloud.ibm.com/graphql"
)

func main() {
	if cmd, ok := cli.Cmds[os.Args[1]]; !ok {
		red := color.New(color.FgRed, color.Bold).PrintfFunc()
		red("%s is not a valid resource\n", os.Args[1])
		os.Exit(2)
	} else {
		cmd.Parse(os.Args[2:])
	}

	c, _ := client.New(cli.Endpoint(), nil)

	result, err := cli.Execute(os.Args[1], &c)

	if err != nil {
		fmt.Fprintln(os.Stderr, "KABOOM", err)
	} else {
		cli.Print(result)
	}
}
