package clusters

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
	"github.ibm.com/coligo/satcon-client/client/types"
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
	Clusters types.ClusterList `json:"clustersByOrgId"`
}

func (c *Client) ClustersByOrgID(orgID, token string) (types.ClusterList, error) {
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
