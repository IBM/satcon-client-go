package subscriptions

import (
	"errors"
	"net/http"

	"github.com/IBM/satcon-client-go/client/auth"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// SubscriptionService is the interface used to perform all cluster-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SubscriptionService
type SubscriptionService interface {
	AddSubscription(orgID, name, channelUuid, versionUuid string, groups []string) (*AddSubscriptionResponseDataDetails, error)
	SetSubscription(orgID string, subscriptionUuid string, versionUuid string) (*SetSubscriptionResponseDataDetails, error)
	RemoveSubscription(orgID, uuid string) (*RemoveSubscriptionResponseDataDetails, error)
	Subscriptions(orgID string) (types.SubscriptionList, error)
	SubscriptionIdsForCluster(orgID string, clusterID string) ([]string, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient, authClient auth.AuthClient) (SubscriptionService, error) {
	if endpointURL == "" {
		return nil, errors.New("Must supply a valid endpoint URL")
	}

	s := web.SatConClient{
		Endpoint:   endpointURL,
		HTTPClient: http.DefaultClient,
		AuthClient: authClient,
	}

	if httpClient != nil {
		s.HTTPClient = httpClient
	}

	return &Client{
		s,
	}, nil
}
