package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryClustersByOrgID       = "clustersByOrgID"
	ClustersByOrgIDVarTemplate = `{{define "vars"}}"org_id":"{{.OrgID}}"{{end}}`
)

type ClustersByOrgIDVariables struct {
	actions.GraphQLQuery
	OrgID string
}

func NewClustersByOrgIDVariables(orgID string) ClustersByOrgIDVariables {
	vars := ClustersByOrgIDVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryClustersByOrgID
	vars.Args = map[string]string{
		"org_id": "String!",
	}
	vars.Returns = []string{
		"_id",
		"org_id",
		"cluster_id",
		"metadata",
	}

	return vars
}

type ClustersByOrgIDResponse struct {
	Data *ClustersByOrgIDResponseData `json:"data"`
}

type ClustersByOrgIDResponseData struct {
	Clusters ClusterList `json:"clustersByOrgID"`
}

type ClusterList []Cluster

type Cluster struct {
	ID        string `json:"_id,omitempty"`
	OrgID     string `json:"org_id,omitempty"`
	ClusterID string `json:"cluster_id,omitempty"`
	// Metadata          []byte         `json:"metadata,omitempty"`
	Metadata          interface{}    `json:"metadata,omitempty"`
	Comments          []Comment      `json:"comments,omitempty"`
	Registration      Registration   `json:"registration,omitempty"`
	RegistrationState string         `json:"reg_state,omitempty"`
	Groups            []ClusterGroup `json:"groups,omitempty"`
	Created           string         `json:"created,omitempty"`
	Updated           string         `json:"updated,omitempty"`
	Dirty             bool           `json:"dirty,omitempty"`
}

type Comment struct {
	UserId  string `json:"user_id"`
	Content string `json:"content"`
	Created string `json:"created"`
}

type ClusterGroupList []ClusterGroup

func (l ClusterGroupList) String() string {
	if len(l) == 0 {
		return "[]"
	}

	groups := make([]string, 1)

	for _, g := range l {
		groups = append(groups, fmt.Sprintf("{UUID: %s\tName: %s}", g.UUID, g.Name))
	}

	return strings.Join(groups, ", ")
}

type ClusterGroup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func (l ClusterList) String() string {
	if len(l) == 0 {
		return "No clusters found"
	}

	clusters := make([]string, 1)

	for _, c := range l {
		clusters = append(clusters, fmt.Sprintf("ID: %s\nOrg ID: %s\nCluster ID: %s\nMetadata: %+v\n", c.ID, c.OrgID, c.ClusterID, c.Metadata))
	}

	return strings.Join(clusters, "\n==\n")
}

func (c *Client) ClustersByOrgID(orgID, token string) (ClusterList, error) {
	var response ClustersByOrgIDResponse

	vars := NewClustersByOrgIDVariables(orgID)

	payload, err := actions.BuildRequestBody(ClustersByOrgIDVarTemplate, vars, nil)

	req := actions.BuildRequest(payload, c.Endpoint, token)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	// TODO: do more than simply dump the body to output.
	if res.Body != nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(body, &response)
	}

	if response.Data != nil {
		return response.Data.Clusters, err
	}

	return nil, err
}
