package iam

import (
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/satcon-client-go/client/auth"
)

//Client manages authorization for Satcon Client requests
type Client struct {
	Client auth.AuthClient
}

//NewIAMClient returns a new core.IamAuthenticator struct and also returns the error
func NewIAMClient(apiKey string, url string) (*Client, error) {

	iamClient, err := core.NewIamAuthenticator(apiKey, url, "", "", false, nil)

	if err == nil {
		return &Client{Client: iamClient}, nil
	}

	return nil, err

}
