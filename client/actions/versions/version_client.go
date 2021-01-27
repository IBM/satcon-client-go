package versions

import (
	"errors"
	"github.com/IBM/satcon-client-go/client/auth"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// VersionService is the interface used to perform all channel-version-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . VersionService
type VersionService interface {
	AddChannelVersion(orgId, channelUuid, name string, content []byte, description string) (*AddChannelVersionResponseDataDetails, error)
	RemoveChannelVersion(orgId, uuid string) (*RemoveChannelVersionResponseDataDetails, error)
	ChannelVersion(orgID, channelUuid, versionUuid string) (*types.DeployableVersion, error)
	ChannelVersionByName(orgID, channelName, versionName string) (*types.DeployableVersion, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient, authClient auth.AuthClient) (VersionService, error) {
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
