package subscriptions

import (
	"errors"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// SubscriptionService is the interface used to perform all cluster-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SubscriptionService
type SubscriptionService interface {
	AddSubscription(orgID, name, channelUuid, versionUuid string, groups []string, token string) (*AddSubscriptionResponseDataDetails, error)
	RemoveSubscription(orgID, uuid, token string) (*RemoveSubscriptionResponseDataDetails, error)
	Subscriptions(orgID, token string) (types.SubscriptionList, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient) (SubscriptionService, error) {
	if endpointURL == "" {
		return nil, errors.New("Must supply a valid endpoint URL")
	}

	s := web.SatConClient{
		Endpoint:   endpointURL,
		HTTPClient: http.DefaultClient,
	}

	if httpClient != nil {
		s.HTTPClient = httpClient
	}

	return &Client{
		s,
	}, nil
}
