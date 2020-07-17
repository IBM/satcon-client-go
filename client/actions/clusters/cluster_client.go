package clusters

import (
	"errors"
	"net/http"

	"github.ibm.com/coligo/satcon-client/client/web"
)

// ClusterService is the interface used to perform all cluster-centric actions
// in Satellite Config.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ClusterService
type ClusterService interface {
	// RegisterCluster registers a new cluster under the specified organization ID.
	RegisterCluster(string, Registration, string) (*RegisterClusterResponseDataDetails, error)
	// ClustersByOrgID lists the clusters registered under the specified organization.
	ClustersByOrgID(string, string) (ClusterList, error)
	// DeleteClusterByClusterID deletes the specified cluster from the specified org,
	// including all resources under that cluster.
	DeleteClusterByClusterID(string, string, string) (*DeleteClustersResponseDataDetails, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient) (ClusterService, error) {
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
