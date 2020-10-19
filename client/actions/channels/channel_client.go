package channels

import (
	"errors"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// ChannelService is the interface used to perform all channel-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ChannelService
type ChannelService interface {
	AddChannel(orgId, name, token string) (*AddChannelResponseDataDetails, error)
	Channel(orgId, uuid, token string) (*types.Channel, error)
	ChannelByName(orgID, channelName, token string) (*types.Channel, error)
	Channels(orgId, token string) (types.ChannelList, error)
	RemoveChannel(orgId, uuid, token string) (*RemoveChannelResponseDataDetails, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient) (ChannelService, error) {
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
