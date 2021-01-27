package iam

import (
	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/IBM/satcon-client-go/client/auth"
)

//IAMClient manages authorization for Satcon Client requests
type IAMClient struct {
	Client auth.AuthClient
}

//New returns a new core.IamAuthenticator struct and also returns the error
func NewIAMClient(apiKey string) (*IAMClient, error) {

	iamClient, err := core.NewIamAuthenticator(apiKey, "", "", "", false, nil)

	if err == nil {
		return &IAMClient{Client: iamClient}, nil
	}

	return nil, err

}
