package clusters

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryClusterByName       = "clusterByName"
	ClusterByNameVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","clusterName":"{{js .ClusterName}}"{{end}}`
)

type ClusterByNameVariables struct {
	actions.GraphQLQuery
	OrgID       string
	ClusterName string
}

func NewClusterByNameVariables(orgID string, clusterName string) ClusterByNameVariables {
	vars := ClusterByNameVariables{
		OrgID:       orgID,
		ClusterName: clusterName,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryClusterByName
	vars.Args = map[string]string{
		"orgId":       "String!",
		"clusterName": "String!",
	}
	vars.Returns = []string{
		"id",
		"orgId",
		"clusterId",
		"name",
		"metadata",
	}

	return vars
}

type ClusterByNameResponse struct {
	Data *ClusterByNameResponseData `json:"data,omitempty"`
}

type ClusterByNameResponseData struct {
	Cluster *types.Cluster `json:"clusterByName,omitempty"`
}

func (c *Client) ClusterByName(orgID string, clusterName string) (*types.Cluster, error) {
	var response ClusterByNameResponse

	vars := NewClusterByNameVariables(orgID, clusterName)

	err := c.DoQuery(ClusterByNameVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Cluster, nil
	}

	return nil, err
}
