package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.ibm.com/coligo/satcon-client/cli"
	"github.ibm.com/coligo/satcon-client/client/actions/clusters"
	"github.ibm.com/coligo/satcon-client/client/actions/subscriptions"
	"github.ibm.com/coligo/satcon-client/client/actions/versions"
	"github.ibm.com/coligo/satcon-client/client/types"
)

var (
	action, endpoint, clusterName, orgID, clusterID, token string
	subscriptionName, channelUuid, versionUuid, filePath   string
	versionName, description, content                      string
	groups                                                 Groups
)

const (
	DefaultEndpoint = "https://config.satellite.test.cloud.ibm.com/graphql"
)

type Groups []string

func (g *Groups) String() string {
	return strings.Join(*g, ", ")
}

func (g *Groups) Set(v string) error {
	*g = append(*g, v)
	return nil
}

func init() {
	flag.StringVar(&action, "a", "ClustersByOrgID", "-a <action>")
	flag.StringVar(&endpoint, "e", DefaultEndpoint, "-e <satcon endpoint URL>")
	flag.StringVar(&clusterName, "c", "", "-c <cluster name>")
	flag.StringVar(&orgID, "o", "d4653c3af7a142fe957a3602f467f183", "-o <organization ID>")
	flag.StringVar(&clusterID, "clusterid", "", "-clusterid <cluster ID>")
	flag.StringVar(&token, "token", "", "-token <IAM token>")
	flag.StringVar(&subscriptionName, "s", "", "-s <subscriptionName>")
	flag.StringVar(&channelUuid, "channelUuid", "", "-channelUuid <channelUuid>")
	flag.StringVar(&versionUuid, "versionUuid", "", "-versionUuid <versionUuid>")
	flag.Var(&groups, "g", "-g <group1> -g <...> -g <groupN>")
	flag.StringVar(&filePath, "f", "", "-f path/to/version.yml")
	flag.StringVar(&versionName, "v", "", "-v <versionName>")
	flag.StringVar(&description, "d", "", "-d <description>")
	flag.StringVar(&content, "content", "", "-content <yaml_as_string>")
}

func main() {
	flag.Parse()

	c, _ := clusters.NewClient(endpoint, nil)
	s, _ := subscriptions.NewClient(endpoint, nil)
	v, _ := versions.NewClient(endpoint, nil)

	var (
		result interface{}
		err    error
	)

	switch action {
	case "ClustersByOrgID":
		result, err = c.ClustersByOrgID(orgID, token)
	case "RegisterCluster":
		result, err = c.RegisterCluster(orgID, types.Registration{Name: clusterName}, token)
	case "DeleteClusterByClusterID":
		result, err = c.DeleteClusterByClusterID(orgID, clusterID, token)
	case "Subscriptions":
		result, err = s.Subscriptions(orgID, token)
	case "AddSubscription":
		result, err = s.AddSubscription(orgID, subscriptionName, channelUuid, versionUuid, groups, token)
	case "AddChannelVersion":
		var versionContent []byte
		if strings.Compare(filePath, "") != 0 {
			encodedBytes, err := cli.MarshalYAMLFromFile(filePath)
			if err != nil {
				break
			}
			versionContent, err = base64.StdEncoding.DecodeString(string(encodedBytes))
		} else {
			versionContent = []byte(content)
		}
		result, err = v.AddChannelVersion(orgID, channelUuid, versionName, versionContent, description, token)
	case "RemoveChannelVersion":
		result, err = v.RemoveChannelVersion(orgID, versionUuid, token)
	default:
		red := color.New(color.FgRed, color.Bold).PrintfFunc()
		red("%s is not a valid action\n", action)
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "KABOOM", err)
	} else {
		cli.Print(result)
	}
}
