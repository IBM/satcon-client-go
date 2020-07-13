package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryRegisterCluster       = "registerCluster"
	RegisterClusterVarTemplate = `{{define "vars"}}"org_id":"{{.OrgID}}","registration":{{printf "%s" .Registration}}{{end}}`
)

type Registration struct {
	Name string `json:"name"`
}

type RegisterClusterVariables struct {
	actions.GraphQLQuery
	OrgID        string
	Registration []byte
}

type RegisterClusterResponse struct {
	URL          string `json:"url"`
	OrgID        string `json:"org_id"`
	OrgKey       string `json:"orgKey,omitempty"`
	ClusterID    string `json:"clusterId"`
	RegState     string `json:"regState"`
	Registration Registration
}

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

func (c *Client) RegisterCluster(orgID string, registration Registration, token string) error {
	vars := NewRegisterClusterVariables(orgID, registration)

	payload, err := actions.BuildRequestBody(RegisterClusterVarTemplate, vars, nil)

	req, _ := http.NewRequest(http.MethodPost, c.Endpoint, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)

	if res.Body != nil {

		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		fmt.Fprintln(os.Stderr, res)
		fmt.Fprintln(os.Stderr, string(body))
	}

	return err
}
