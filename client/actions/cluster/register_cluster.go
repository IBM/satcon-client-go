package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryRegisterCluster       = "registerCluster"
	RegisterClusterVarTemplate = `{{define "vars"}}"org_id":"{{.OrgID}}","registration":{{printf "%s" .Registration}}{{end}}`
)

// Registration is the encapsulation of the JSON registration body, which at this
// point is used primarily to specify the name of the cluster to register.
type Registration struct {
	Name string `json:"name"`
}

// RegisterClusterVariables are the variables specific to cluster registration.
// These include the organization ID and the serialized registration.  Rather than
// instantiating this directly, use NewRegisterClusterVariables().
type RegisterClusterVariables struct {
	actions.GraphQLQuery
	OrgID        string
	Registration []byte
}

// RegisterClusterResponse is the response body we get upon a successful cluster
// registration.
type RegisterClusterResponse struct {
	Data *RegisterClusterResponseData `json:"data,omitempty"`
}

type RegisterClusterResponseData struct {
	Details RegisterClusterResponseDataDetails `json:"registerCluster"`
}

type RegisterClusterResponseDataDetails struct {
	URL          string `json:"url"`
	OrgID        string `json:"org_id"`
	OrgKey       string `json:"orgKey,omitempty"`
	ClusterID    string `json:"clusterId"`
	RegState     string `json:"regState"`
	Registration Registration
}

func (d RegisterClusterResponseDataDetails) String() string {
	return fmt.Sprintf("URL: %s\nOrg ID: %s\nOrg Key: %s\nCluster ID: %s\nRegistration State: %s\nRegistration: %+v\n",
		d.URL, d.OrgID, d.OrgKey, d.ClusterID, d.RegState, d.Registration)
}

// NewRegisterClusterVariables creates a correctly formed instance of RegisterClusterVariables.
func NewRegisterClusterVariables(orgID string, registration Registration) RegisterClusterVariables {
	regBytes, _ := json.Marshal(registration)

	vars := RegisterClusterVariables{
		OrgID:        orgID,
		Registration: regBytes,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRegisterCluster
	vars.Args = map[string]string{
		"org_id":       "String!",
		"registration": "JSON!",
	}
	vars.Returns = []string{
		"url",
		"org_id",
		"orgKey",
		"clusterId",
		"regState",
		"registration",
	}

	return vars
}

func (c *Client) RegisterCluster(orgID string, registration Registration, token string) (RegisterClusterResponseDataDetails, error) {
	var response RegisterClusterResponse

	vars := NewRegisterClusterVariables(orgID, registration)

	payload, err := actions.BuildRequestBody(RegisterClusterVarTemplate, vars, nil)

	req, _ := http.NewRequest(http.MethodPost, c.Endpoint, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return RegisterClusterResponseDataDetails{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &response)
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return RegisterClusterResponseDataDetails{}, err
}
