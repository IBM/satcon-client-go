package local

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web"
	"github.com/golang-jwt/jwt"
)

const TokenValidityDuration = 40 * time.Minute
const MinimumTimeTokenStillValid = 5 * time.Minute
const AuthorizationHeaderKey = "Authorization"

type LocalRazeeClient struct {
	HTTPClient web.HTTPClient
	url        string
	login      string
	password   string

	// cached jwt token from previous request
	token types.Token
	// timestamp of the time when the previous
	// cached token was retrieved
	tokenTimestamp time.Time
	// timestamp that was set in the 'exp' claim of
	// the previous retrieved jwt token.
	expireTimestamp time.Time
}

func (l *LocalRazeeClient) Authenticate(request *http.Request) error {
	invalidExpiredTimestamp := time.Until(l.expireTimestamp) < MinimumTimeTokenStillValid
	invalidTokenTimestamp := time.Since(l.tokenTimestamp) >= TokenValidityDuration
	if l.token == "" || invalidExpiredTimestamp || invalidTokenTimestamp {
		token, err := SignIn(l.HTTPClient, l.url, l.login, l.password)
		if err != nil {
			return err
		}
		if token == nil {
			return fmt.Errorf("Could not get a token by signing in %v to %v", l.login, l.url)
		}
		l.token = *token
		parsedToken, err := jwt.Parse(string(*token), nil)
		if parsedToken != nil {
			claims, _ := parsedToken.Claims.(jwt.MapClaims)
			if expiredTimestamp, ok := claims["exp"]; ok {
				if d, ok := expiredTimestamp.(float64); ok {
					l.expireTimestamp = time.Unix(int64(d), 0)
				}
			}
		}

		// Token timestamp is a backup mechanism if the token does not contain the 'exp' field
		// or the value of that field could not be parsed
		l.tokenTimestamp = time.Now()
	}
	request.Header.Add(AuthorizationHeaderKey, "Bearer "+string(l.token))
	return nil
}

func NewClient(url string, login string, password string) (*LocalRazeeClient, error) {
	return NewClientWithHttpClient(http.DefaultClient, url, login, password)
}

func NewClientWithHttpClient(httpClient web.HTTPClient, url string, login string, password string) (*LocalRazeeClient, error) {
	if url == "" {
		return nil, fmt.Errorf("Field 'url' cannot be empty!")
	}
	if login == "" {
		return nil, fmt.Errorf("Field 'login' cannot be empty!")
	}
	if password == "" {
		return nil, fmt.Errorf("Field 'password' cannot be empty!")
	}
	client := &LocalRazeeClient{
		HTTPClient: httpClient,
		url:        url,
		login:      login,
		password:   password,
	}
	return client, nil
}
