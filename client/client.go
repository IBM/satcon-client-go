package client

import (
	"github.com/IBM/satcon-client-go/client/actions/channels"
	"github.com/IBM/satcon-client-go/client/actions/clusters"
	"github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/actions/subscriptions"
	"github.com/IBM/satcon-client-go/client/actions/versions"
	"github.com/IBM/satcon-client-go/client/web"
)

type SatCon struct {
	Channels      channels.ChannelService
	Clusters      clusters.ClusterService
	Groups        groups.GroupService
	Resources     resources.ResourceService
	Subscriptions subscriptions.SubscriptionService
	Versions      versions.VersionService
}

func New(endpointURL string, httpClient web.HTTPClient) (SatCon, error) {
	var (
		err error
		s   SatCon
	)
	s.Channels, err = channels.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Clusters, err = clusters.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Groups, err = groups.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Resources, err = resources.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Subscriptions, err = subscriptions.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}
	s.Versions, err = versions.NewClient(endpointURL, httpClient)
	if err != nil {
		return SatCon{}, err
	}

	return s, nil
}
