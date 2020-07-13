package main

import (
	"flag"
	"fmt"
	"github.ibm.com/coligo/satcon-client/client/actions/cluster"
	"os"
)

var (
	action, endpoint, clusterName, orgID, clusterID, token string
)

const (
	DefaultEndpoint = "https://config.satellite.test.cloud.ibm.com/graphql"
)

func init() {
	flag.StringVar(&action, "a", "ClustersByOrgID", "-a <action>")
	flag.StringVar(&endpoint, "e", DefaultEndpoint, "-e <satcon endpoint URL>")
	flag.StringVar(&clusterName, "c", "", "-c <cluster name>")
	flag.StringVar(&orgID, "o", "d4653c3af7a142fe957a3602f467f183", "-o <organization ID>")
	flag.StringVar(&clusterID, "clusterid", "", "-clusterid <cluster ID>")
	flag.StringVar(&token, "token", "", "-token <IAM token>")
}

func main() {
	flag.Parse()

	c, _ := cluster.NewClient(endpoint, nil)

	var (
		result interface{}
		err    error
	)

	switch action {
	case "ClustersByOrgID":
		result, err = c.ClustersByOrgID(orgID, token)
	case "RegisterCluster":
		result, err = c.RegisterCluster(orgID, cluster.Registration{Name: clusterName}, token)
	case "DeleteClusterByClusterID":
		result, err = c.DeleteClusterByClusterID(orgID, clusterID, token)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "KABOOM", err)
	} else {
		fmt.Fprintf(os.Stdout, "%s\n", result)
	}
}
