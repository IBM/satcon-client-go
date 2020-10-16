package clusters

import (
	"encoding/json"
	"fmt"

	// "io/ioutil"

	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryRegisterCluster       = "registerCluster"
	RegisterClusterVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","registration":{{printf "%s" .Registration}}{{end}}`
)

// RegisterClusterVariables are the variables specific to cluster registration.
// These include the organization ID and the serialized registration.  Rather than
// instantiating this directly, use NewRegisterClusterVariables().
type RegisterClusterVariables struct {
	actions.GraphQLQuery
	OrgID        string
	Registration []byte
}

// NewRegisterClusterVariables creates a correctly formed instance of RegisterClusterVariables.
func NewRegisterClusterVariables(orgID string, registration types.Registration) RegisterClusterVariables {
	regBytes, _ := json.Marshal(registration)

	vars := RegisterClusterVariables{
		OrgID:        orgID,
		Registration: regBytes,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRegisterCluster
	vars.Args = map[string]string{
		"orgId":        "String!",
		"registration": "JSON!",
	}
	vars.Returns = []string{
		"url",
		"orgId",
		"orgKey",
		"clusterId",
		"regState",
		"registration",
	}

	return vars
}

// RegisterClusterResponse is the response body we get upon a successful cluster
// registration.
type RegisterClusterResponse struct {
	Data *RegisterClusterResponseData `json:"data,omitempty"`
}

type RegisterClusterResponseData struct {
	Details *RegisterClusterResponseDataDetails `json:"registerCluster,omitempty"`
}

type RegisterClusterResponseDataDetails struct {
	URL          string `json:"url,omitempty"`
	OrgID        string `json:"orgId,omitempty"`
	OrgKey       string `json:"orgKey,omitempty"`
	ClusterID    string `json:"clusterId,omitempty"`
	RegState     string `json:"regState,omitempty"`
	Registration types.Registration
}

func (d RegisterClusterResponseDataDetails) String() string {
	return fmt.Sprintf("URL: %s\nOrg ID: %s\nOrg Key: %s\nCluster ID: %s\nRegistration State: %s\nRegistration: %+v\n",
		d.URL, d.OrgID, d.OrgKey, d.ClusterID, d.RegState, d.Registration)
}

func (c *Client) RegisterCluster(orgID string, registration types.Registration, token string) (*RegisterClusterResponseDataDetails, error) {
	var response RegisterClusterResponse

	vars := NewRegisterClusterVariables(orgID, registration)

	err := c.DoQuery(RegisterClusterVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
