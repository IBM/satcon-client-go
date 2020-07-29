package cli

import (
	"fmt"
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

type SubscriptionMetadata struct {
	ChannelUUID string
	VersionUUID string
	Groups      Groups
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
	subscriptionMetadata SubscriptionMetadata
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

	case TypeCluster:
		switch cmd.Action {
		case ActionList:
			result, err = s.Clusters.ClustersByOrgID(cmd.OrgID, cmd.Token)
		case ActionDelete:
			result, err = s.Clusters.DeleteClusterByClusterID(cmd.OrgID, cmd.Id, cmd.Token)
		case ActionRegister:
			result, err = s.Clusters.RegisterCluster(cmd.OrgID, types.Registration{Name: cmd.Name}, cmd.Token)
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
		default:
			err = fmt.Errorf("%s is not a recognized action for resource type %s", cmd.Action, cmd.Resource)
		}

	default:
		err = fmt.Errorf("%s is not a valid resource type", cmd.Resource)
	}

	return result, err
}
