package cli

import (
	"flag"
	"fmt"
)

const (
	DefaultEndpoint = "https://config.satellite.test.cloud.ibm.com/graphql"

	TypeChannel      = "channel"
	TypeCluster      = "cluster"
	TypeGroup        = "group"
	TypeResource     = "resource"
	TypeSubscription = "subscription"
	TypeVersion      = "version"

	FLAG_ACTION   = "a"
	FLAG_NAME     = "n"
	FLAG_ID       = "id"
	FLAG_ORGID    = "o"
	FLAG_ENDPOINT = "e"
	FLAG_TOKEN    = "t"

	FLAG_CLUSTERS = "c"

	FLAG_CHANNEL_UUID = "channel"
	FLAG_VERSION_UUID = "version"
	FLAG_GROUPS       = "g"

	FLAG_FILENAME     = "f"
	FLAG_DESCRIPTION  = "desc"
	FLAG_CHANNEL_NAME = "channelName"

	FLAG_CLUSTER_ID = "clusterID"
	FLAG_FILTER     = "filter"
	FLAG_LIMIT      = "l"

	ActionList     = "list"
	ActionRegister = "register"
	ActionDelete   = "delete"
	ActionAdd      = "add"
	ActionClusters = "clusters"
	ActionRemove   = "remove"
	ActionGet      = "get"
)

var (
	ChannelCmd      *flag.FlagSet
	ClusterCmd      *flag.FlagSet
	GroupCmd        *flag.FlagSet
	ResourceCmd     *flag.FlagSet
	SubscriptionCmd *flag.FlagSet
	VersionCmd      *flag.FlagSet
	Cmds            map[string]*flag.FlagSet
)

func init() {
	ChannelCmd = flag.NewFlagSet(TypeChannel, flag.ExitOnError)
	ClusterCmd = flag.NewFlagSet(TypeCluster, flag.ExitOnError)
	GroupCmd = flag.NewFlagSet(TypeGroup, flag.ExitOnError)
	ResourceCmd = flag.NewFlagSet(TypeResource, flag.ExitOnError)
	SubscriptionCmd = flag.NewFlagSet(TypeSubscription, flag.ExitOnError)
	VersionCmd = flag.NewFlagSet(TypeVersion, flag.ExitOnError)

	Cmds = map[string]*flag.FlagSet{
		TypeChannel:      ChannelCmd,
		TypeCluster:      ClusterCmd,
		TypeGroup:        GroupCmd,
		TypeResource:     ResourceCmd,
		TypeSubscription: SubscriptionCmd,
		TypeVersion:      VersionCmd,
	}

	// All subcommands accept these common arguments, even where they may not actually apply
	for _, fs := range Cmds {
		fs.StringVar(&(subCmd.Action), FLAG_ACTION, "", fmt.Sprintf("-%s list|add|remove|get", FLAG_ACTION))
		fs.StringVar(&(subCmd.Name), FLAG_NAME, "", fmt.Sprintf("-%s <name>", FLAG_NAME))
		fs.StringVar(&(subCmd.Id), FLAG_ID, "", fmt.Sprintf("-%s <id>", FLAG_ID))
		fs.StringVar(&(subCmd.OrgID), FLAG_ORGID, "", fmt.Sprintf("-%s <organization_id>", FLAG_ORGID))
		fs.StringVar(&(subCmd.Endpoint), FLAG_ENDPOINT, DefaultEndpoint, fmt.Sprintf("-%s <satcon endpoint URL>", FLAG_ENDPOINT))
		fs.StringVar(&(subCmd.Token), FLAG_TOKEN, "", fmt.Sprintf("-%s <IAM token>", FLAG_TOKEN))
	}

	// Channel-specific arguments
	ChannelCmd.StringVar(&(channelMetadata.ChannelUUID), FLAG_CHANNEL_UUID, "", fmt.Sprintf("-%s <channel_uuid>", FLAG_CHANNEL_UUID))

	// Group-specific arguments
	GroupCmd.Var(&(groupMetadata.Clusters), FLAG_CLUSTERS, fmt.Sprintf("-%s <cluster_id> [ -%s <cluster_id_1> ... -%s <cluster_id_n> ]", FLAG_CLUSTERS, FLAG_CLUSTERS, FLAG_CLUSTERS))

	// Subscription-specific arguments
	SubscriptionCmd.StringVar(&(subscriptionMetadata.ChannelUUID), FLAG_CHANNEL_UUID, "", fmt.Sprintf("-%s <channel_uuid>", FLAG_CHANNEL_UUID))
	SubscriptionCmd.StringVar(&(subscriptionMetadata.VersionUUID), FLAG_VERSION_UUID, "", fmt.Sprintf("-%s <channel_uuid>", FLAG_VERSION_UUID))
	SubscriptionCmd.Var(&(subscriptionMetadata.Groups), FLAG_GROUPS, fmt.Sprintf("-%s <group_id>", FLAG_GROUPS))

	// Version-specific arguments
	VersionCmd.StringVar(&(versionMetadata.ChannelUUID), FLAG_CHANNEL_UUID, "", fmt.Sprintf("-%s <channel_uuid>", FLAG_CHANNEL_UUID))
	VersionCmd.StringVar(&(versionMetadata.Filename), FLAG_FILENAME, "", fmt.Sprintf("-%s <path_to_yaml>", FLAG_FILENAME))
	VersionCmd.StringVar(&(versionMetadata.Description), FLAG_DESCRIPTION, "", fmt.Sprintf("-%s <description>", FLAG_DESCRIPTION))
	VersionCmd.StringVar(&(versionMetadata.ChannelName), FLAG_CHANNEL_NAME, "", fmt.Sprintf("-%s <channel_name>", FLAG_CHANNEL_NAME))

	//Resource-specific commands
	ResourceCmd.StringVar(&(resourceMetadata.ClusterID), FLAG_CLUSTER_ID, "", fmt.Sprintf("-%s <cluster_id>", FLAG_CLUSTER_ID))
	ResourceCmd.StringVar(&(resourceMetadata.Filter), FLAG_FILTER, "", fmt.Sprintf("-%s <kind>", FLAG_FILTER))
	ResourceCmd.StringVar(&(resourceMetadata.Limit), FLAG_LIMIT, "50", fmt.Sprintf("-%s <limit>", FLAG_LIMIT))
}
