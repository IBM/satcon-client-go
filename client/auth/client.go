package auth


import (
	"fmt"
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/auth/local"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
	"net/http"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
)

//IAMClient manages authorization for Satcon Client requests
type IAMClient struct {
	Client actions.AuthClient
}

//New returns a new core.IamAuthenticator struct and also returns the error
func NewIAMClient(apiKey string) (*IAMClient, error) {

	iamClient, err := core.NewIamAuthenticator(apiKey, "", "", "", false, nil)

	if err == nil {
		return &IAMClient{Client: iamClient}, nil
	}

	return nil, err

}

const TokenValidityDuration = 40 * time.Minute

type LocalRazeeClient struct {
	HTTPClient web.HTTPClient
	url      string
	login    string
	password string
	tokenTimestamp time.Time
	token    types.Token
}

func (l *LocalRazeeClient) Authenticate(request *http.Request) error {
	time.Now()
	if l.token == "" || time.Since(l.tokenTimestamp) >= TokenValidityDuration {
		token, err := local.SignIn(l.HTTPClient, l.url, l.login, l.password)
		if err != nil {
			return err
		}
		if token == nil {
			return fmt.Errorf("Could not get a token by signing in %v to %v", l.login, l.url)
		}
		l.token = *token
		l.tokenTimestamp = time.Now()
	}
	request.Header.Add("Authorization", "Bearer " + string(l.token))
	return nil
}

func NewLocalRazeeClient(url string, login string, password string) (*LocalRazeeClient, error) {
	client := &LocalRazeeClient{
		url:      url,
		login:    login,
		password: password,
	}
	return client, nil
}