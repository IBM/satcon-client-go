package main

import (
	"fmt"
	"github.ibm.com/coligo/satcon-client/client/actions/cluster"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Must supply cluster name")
		os.Exit(1)
	}

	clusterName := os.Args[1]

	url := "https://config.satellite.test.cloud.ibm.com/graphql"

	c := cluster.NewClient(url, nil)

	err := c.RegisterCluster("d4653c3af7a142fe957a3602f467f183", cluster.Registration{Name: clusterName}, token)

	if err != nil {
		fmt.Fprintln(os.Stderr, "KABOOM", err)
	}
}
