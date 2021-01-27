package channels

import (
	"errors"
	"github.com/IBM/satcon-client-go/client/auth"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// ChannelService is the interface used to perform all channel-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ChannelService
type ChannelService interface {
	AddChannel(orgId, name string) (*AddChannelResponseDataDetails, error)
	Channel(orgId, uuid string) (*types.Channel, error)
	ChannelByName(orgID, channelName string) (*types.Channel, error)
	Channels(orgId string) (types.ChannelList, error)
	RemoveChannel(orgId, uuid string) (*RemoveChannelResponseDataDetails, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient, authClient auth.AuthClient) (ChannelService, error) {
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
