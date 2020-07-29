package versions

import (
	"errors"
	"net/http"

	"github.ibm.com/coligo/satcon-client/client/types"
	"github.ibm.com/coligo/satcon-client/client/web"
)

// VersionService is the interface used to perform all channel-version-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . VersionService
type VersionService interface {
	AddChannelVersion(orgId, channelUuid, name string, content []byte, description, token string) (*AddChannelVersionResponseDataDetails, error)
	RemoveChannelVersion(orgId, uuid, token string) (*RemoveChannelVersionResponseDataDetails, error)
	ChannelVersionByName(orgID, channelName, versionName, token string) (*types.DeployableVersion, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient) (VersionService, error) {
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
