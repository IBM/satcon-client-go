package client

import (
	"github.ibm.com/coligo/satcon-client/client/actions/channels"
	"github.ibm.com/coligo/satcon-client/client/actions/clusters"
	"github.ibm.com/coligo/satcon-client/client/actions/groups"
	"github.ibm.com/coligo/satcon-client/client/actions/resources"
	"github.ibm.com/coligo/satcon-client/client/actions/subscriptions"
	"github.ibm.com/coligo/satcon-client/client/web"
)

type SatCon struct {
	Channels      channels.ChannelService
	Clusters      clusters.ClusterService
	Groups        groups.GroupService
	Resources     resources.ResourceService
	Subscriptions subscriptions.SubscriptionService
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

	return s, nil
}
