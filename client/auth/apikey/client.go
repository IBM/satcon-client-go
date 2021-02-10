package apikey

import "net/http"

const APIKeyHeader = "x-api-key"

// This client is used for authentication using OAuth e.g. GitHub/GitHub Enterprise
type RazeeApiKeyAuthClient struct {
	apiKey string
}

func (r *RazeeApiKeyAuthClient) Authenticate(request *http.Request) error {
	request.Header.Add(APIKeyHeader, r.apiKey)
	return nil
}

func NewClient(apiKey string) (*RazeeApiKeyAuthClient, error) {
	client := &RazeeApiKeyAuthClient{
		apiKey: apiKey,
	}

	return client, nil
}
