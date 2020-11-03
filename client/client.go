package client

import (
	"github.com/IBM/satcon-client-go/client/actions/channels"
	"github.com/IBM/satcon-client-go/client/actions/channels/channelsfakes"
	"github.com/IBM/satcon-client-go/client/actions/clusters"
	"github.com/IBM/satcon-client-go/client/actions/clusters/clustersfakes"
	"github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/actions/groups/groupsfakes"
	"github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/actions/resources/resourcesfakes"
	"github.com/IBM/satcon-client-go/client/actions/subscriptions"
	"github.com/IBM/satcon-client-go/client/actions/subscriptions/subscriptionsfakes"
	"github.com/IBM/satcon-client-go/client/actions/versions"
	"github.com/IBM/satcon-client-go/client/actions/versions/versionsfakes"
	"github.com/IBM/satcon-client-go/client/auth"
	"github.com/IBM/satcon-client-go/client/web"
)

//SatCon struct for satellite configuration entities
type SatCon struct {
	Channels      channels.ChannelService
	Clusters      clusters.ClusterService
	Groups        groups.GroupService
	Resources     resources.ResourceService
	Subscriptions subscriptions.SubscriptionService
	Versions      versions.VersionService
}

//New creates new SatCon clients
func New(endpointURL string, httpClient web.HTTPClient, authClient auth.AuthClient) (SatCon, error) {
	var (
		err error
		s   SatCon
	)

	s.Channels, err = channels.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Clusters, err = clusters.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Groups, err = groups.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Resources, err = resources.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Subscriptions, err = subscriptions.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Versions, err = versions.NewClient(endpointURL, httpClient, authClient)
	if err != nil {
		return SatCon{}, err
	}

	return s, nil
}

// NewTesting is a convenience method which creates a client using only fakes
// for the type-specific service interfaces.  See the counterfeiter documentation
// for how to customize these fakes e.g. to provide stub implementations, etc.
func NewTesting(endpointURL string, httpClient web.HTTPClient) SatCon {
	var s SatCon
	s.Channels = &channelsfakes.FakeChannelService{}
	s.Clusters = &clustersfakes.FakeClusterService{}
	s.Groups = &groupsfakes.FakeGroupService{}
	s.Resources = &resourcesfakes.FakeResourceService{}
	s.Subscriptions = &subscriptionsfakes.FakeSubscriptionService{}
	s.Versions = &versionsfakes.FakeVersionService{}

	return s
}
