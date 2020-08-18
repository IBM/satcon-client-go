package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.ibm.com/coligo/satcon-client/client"
	"github.ibm.com/coligo/satcon-client/client/types"
)

type SubCommand struct {
	Resource string
	Action   string
	Name     string
	Id       string
	OrgID    string
	Endpoint string
	Token    string
}

type ChannelMetadata struct {
	ChannelUUID string
}

type GroupMetadata struct {
	Clusters Clusters
}

type SubscriptionMetadata struct {
	ChannelUUID string
	VersionUUID string
	Groups      Groups
}

type VersionMetadata struct {
	ChannelUUID string
	ChannelName string
	Filename    string
	Content     []byte
	Description string
}

type Clusters []string

func (c *Clusters) String() string {
	return strings.Join(*c, ", ")
}

func (c *Clusters) Set(v string) error {
	*c = append(*c, v)
	return nil
}

type Groups []string

func (g *Groups) String() string {
	return strings.Join(*g, ", ")
}

func (g *Groups) Set(v string) error {
	*g = append(*g, v)
	return nil
}

var (
	subCmd               SubCommand
	channelMetadata      ChannelMetadata
	groupMetadata        GroupMetadata
	subscriptionMetadata SubscriptionMetadata
	versionMetadata      VersionMetadata
)

func TargetResource(resourceType string) {
	subCmd.Resource = resourceType
}

func Endpoint() string {
	return subCmd.Endpoint
}

func Execute(resource string, s *client.SatCon) (interface{}, error) {
	subCmd.Resource = resource
	return subCmd.execute(s)
}

func (cmd *SubCommand) execute(s *client.SatCon) (interface{}, error) {
	var (
		result interface{}
		err    error
	)
	switch cmd.Resource {

	case TypeChannel:
		switch cmd.Action {
		case ActionAdd:
			result, err = s.Channels.AddChannel(cmd.OrgID, cmd.Name, cmd.Token)
		case ActionGet:
			result, err = s.Channels.ChannelByName(cmd.OrgID, cmd.Name, cmd.Token)
		case ActionList:
			result, err = s.Channels.Channels(cmd.OrgID, cmd.Token)
		case ActionDelete:
			fallthrough
		case ActionRemove:
			result, err = s.Channels.RemoveChannel(cmd.OrgID, cmd.Id, cmd.Token)
		}

	case TypeCluster:
		switch cmd.Action {
		case ActionList:
			result, err = s.Clusters.ClustersByOrgID(cmd.OrgID, cmd.Token)
		case ActionRemove:
			fallthrough
		case ActionDelete:
			result, err = s.Clusters.DeleteClusterByClusterID(cmd.OrgID, cmd.Id, cmd.Token)
		case ActionRegister:
			result, err = s.Clusters.RegisterCluster(cmd.OrgID, types.Registration{Name: cmd.Name}, cmd.Token)
		default:
			err = fmt.Errorf("%s is not a recognized action for resource type %s", cmd.Action, cmd.Resource)
		}

	case TypeGroup:
		switch cmd.Action {
		case ActionAdd:
			result, err = s.Groups.AddGroup(cmd.OrgID, cmd.Name, cmd.Token)
		case ActionList:
			result, err = s.Groups.Groups(cmd.OrgID, cmd.Token)
		case ActionClusters:
			result, err = s.Groups.GroupClusters(cmd.OrgID, cmd.Id, groupMetadata.Clusters, cmd.Token)
		default:
			err = fmt.Errorf("%s is not a recognized action for resource type %s", cmd.Action, cmd.Resource)
		}

	case TypeSubscription:
		switch cmd.Action {
		case ActionList:
			result, err = s.Subscriptions.Subscriptions(cmd.OrgID, cmd.Token)
		case ActionAdd:
			result, err = s.Subscriptions.AddSubscription(cmd.OrgID, cmd.Name, subscriptionMetadata.ChannelUUID,
				subscriptionMetadata.VersionUUID, subscriptionMetadata.Groups, cmd.Token)
		case ActionDelete:
			fallthrough
		case ActionRemove:
			result, err = s.Subscriptions.RemoveSubscription(cmd.OrgID, cmd.Id, cmd.Token)
		default:
			err = fmt.Errorf("%s is not a recognized action for resource type %s", cmd.Action, cmd.Resource)
		}

	case TypeVersion:
		switch cmd.Action {
		case ActionAdd:
			if versionMetadata.Filename == "" {
				err = fmt.Errorf("Must specify content file with -%s flag", FLAG_FILENAME)
				break
			}

			// versionMetadata.Content, err = MarshalYAMLFromFile(versionMetadata.Filename)
			versionMetadata.Content, err = ioutil.ReadFile(versionMetadata.Filename)
			if err != nil {
				err = fmt.Errorf("Unable to read content file %s: %s", versionMetadata.Filename, err)
				break
			}

			result, err = s.Versions.AddChannelVersion(cmd.OrgID, versionMetadata.ChannelUUID, cmd.Name, versionMetadata.Content, versionMetadata.Description, cmd.Token)
		case ActionDelete:
			fallthrough
		case ActionRemove:
			result, err = s.Versions.RemoveChannelVersion(cmd.OrgID, cmd.Id, cmd.Token)
		case ActionGet:
			result, err = s.Versions.ChannelVersionByName(cmd.OrgID, versionMetadata.ChannelName, cmd.Name, cmd.Token)
		default:
			err = fmt.Errorf("%s is not a recognized action for resource type %s", cmd.Action, cmd.Resource)
		}

	default:
		err = fmt.Errorf("%s is not a valid resource type", cmd.Resource)
	}

	return result, err
}
