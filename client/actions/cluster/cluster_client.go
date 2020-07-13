package cluster

import (
	"errors"
	"net/http"

	"github.ibm.com/coligo/satcon-client/client"
)

type ClusterService interface {
	RegisterCluster(string, Registration, string) error
	ClustersByOrgID(string, string) error
}

type Client struct {
	Endpoint   string
	HTTPClient client.HTTPClient
}

func NewClient(endpointURL string, httpClient client.HTTPClient) (ClusterService, error) {
	if endpointURL == "" {
		return nil, errors.New("Must supply a valid endpoint URL")
	}

	if httpClient == nil {
		return &Client{
			Endpoint:   endpointURL,
			HTTPClient: http.DefaultClient,
		}, nil
	}

	return &Client{
		Endpoint:   endpointURL,
		HTTPClient: httpClient,
	}, nil
}
