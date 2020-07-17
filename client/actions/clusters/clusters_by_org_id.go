package clusters

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"strings"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryClustersByOrgID       = "clustersByOrgId"
	ClustersByOrgIDVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}"{{end}}`
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
		"orgId": "String!",
	}
	vars.Returns = []string{
		"id",
		"orgId",
		"clusterId",
		"metadata",
	}

	return vars
}

type ClustersByOrgIDResponse struct {
	Data *ClustersByOrgIDResponseData `json:"data"`
}

type ClustersByOrgIDResponseData struct {
	Clusters ClusterList `json:"clustersByOrgId"`
}

type ClusterList []Cluster

type Cluster struct {
	ID        string `json:"id,omitempty"`
	OrgID     string `json:"orgId,omitempty"`
	ClusterID string `json:"clusterId,omitempty"`
	// Metadata          []byte         `json:"metadata,omitempty"`
	Metadata          interface{}    `json:"metadata,omitempty"`
	Comments          []Comment      `json:"comments,omitempty"`
	Registration      Registration   `json:"registration,omitempty"`
	RegistrationState string         `json:"regState,omitempty"`
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

	err := c.DoQuery(ClustersByOrgIDVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Clusters, nil
	}

	return nil, err
}
