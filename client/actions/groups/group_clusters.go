package groups

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryGroupClusters       = "groupClusters"
	GroupClustersVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"uuid":{{json .UUID}},"clusters":[{{range $i,$e := .Clusters}}{{if gt $i 0}},{{end}}{{json $e}}{{end}}]{{end}}`
)

// GroupClustersVariables are the variables specific to grouping clusters.
// These include the organization ID, group UUID, and list of cluster IDs.  Rather than
// instantiating this directly, use NewGroupClustersVariables().
type GroupClustersVariables struct {
	actions.GraphQLQuery
	OrgID    string
	UUID     string
	Clusters []string
}

// NewGroupClustersVariables creates a correctly formed instance of GroupClustersVariables.
func NewGroupClustersVariables(orgID, uuid string, clusters []string) GroupClustersVariables {
	vars := GroupClustersVariables{
		OrgID:    orgID,
		UUID:     uuid,
		Clusters: clusters,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryGroupClusters
	vars.Args = map[string]string{
		"orgId":    "String!",
		"uuid":     "String!",
		"clusters": "[String!]!",
	}
	vars.Returns = []string{
		"modified",
	}

	return vars
}

// GroupClustersResponse is the response body we get upon a successful cluster
// registration.
type GroupClustersResponse struct {
	Data *GroupClustersResponseData `json:"data,omitempty"`
}

type GroupClustersResponseData struct {
	Details *GroupClustersResponseDataDetails `json:"groupClusters,omitempty"`
}

type GroupClustersResponseDataDetails struct {
	Modified int `json:"modified,omitempty"`
}

func (c *Client) GroupClusters(orgID, uuid string, clusters []string) (*GroupClustersResponseDataDetails, error) {
	var response GroupClustersResponse

	vars := NewGroupClustersVariables(orgID, uuid, clusters)

	err := c.DoQuery(GroupClustersVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
