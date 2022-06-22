package groups

import "github.com/IBM/satcon-client-go/client/actions"

const (
	QueryUnGroupClusters       = "unGroupClusters"
	UnGroupClustersVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"uuid":{{json .UUID}},"clusters":[{{range $i,$e := .Clusters}}{{if gt $i 0}},{{end}}{{json $e}}{{end}}]{{end}}`
)

// GroupClustersVariables are the variables specific to grouping clusters.
// These include the organization ID, group UUID, and list of cluster IDs.  Rather than
// instantiating this directly, use NewGroupClustersVariables().
type UnGroupClustersVariables struct {
	actions.GraphQLQuery
	OrgID    string
	UUID     string
	Clusters []string
}

// NewGroupClustersVariables creates a correctly formed instance of GroupClustersVariables.
func NewUnGroupClustersVariables(orgID, uuid string, clusters []string) UnGroupClustersVariables {
	vars := UnGroupClustersVariables{
		OrgID:    orgID,
		UUID:     uuid,
		Clusters: clusters,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryUnGroupClusters
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
type UnGroupClustersResponse struct {
	Data *UnGroupClustersResponseData `json:"data,omitempty"`
}

type UnGroupClustersResponseData struct {
	Details *UnGroupClustersResponseDataDetails `json:"unGroupClusters,omitempty"`
}

type UnGroupClustersResponseDataDetails struct {
	Modified int `json:"modified,omitempty"`
}

func (c *Client) UnGroupClusters(orgID, uuid string, clusters []string) (*UnGroupClustersResponseDataDetails, error) {
	var response UnGroupClustersResponse

	vars := NewUnGroupClustersVariables(orgID, uuid, clusters)

	err := c.DoQuery(UnGroupClustersVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
