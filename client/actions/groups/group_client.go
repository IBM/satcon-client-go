package groups

import (
	"errors"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// ClusterService is the interface used to perform all group-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . GroupService
type GroupService interface {
	Groups(orgID, token string) (types.GroupList, error)
	AddGroup(orgID, name, token string) (*AddGroupResponseDataDetails, error)
	RemoveGroup(orgID, uuid, token string) (*RemoveGroupResponseDataDetails, error)
	RemoveGroupByName(orgID, name, token string) (*RemoveGroupByNameResponseDataDetails, error)
	GroupClusters(orgID, uuid string, clusters []string, token string) (*GroupClustersResponseDataDetails, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of GroupService which can then be used
// to perform group queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient) (GroupService, error) {
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
