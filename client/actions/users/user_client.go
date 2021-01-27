package users

import (
	"errors"
	"github.com/IBM/satcon-client-go/client/auth"
	"net/http"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
)

// UserService is the interface used to perform all user-centric actions.
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . UserService
type UserService interface {
	Me() (*types.User, error)
}

// Client is an implementation of a satcon client.
type Client struct {
	web.SatConClient
}

// NewClient returns a configured instance of ClusterService which can then be used
// to perform cluster queries against Satellite Config.
func NewClient(endpointURL string, httpClient web.HTTPClient, authClient auth.AuthClient) (UserService, error) {
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
