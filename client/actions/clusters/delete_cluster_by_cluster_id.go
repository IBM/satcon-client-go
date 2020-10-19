package clusters

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryDeleteClusterByClusterID       = "deleteClusterByClusterId"
	DeleteClusterByClusterIDVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","clusterId":"{{js .ClusterID}}"{{end}}`
)

type DeleteClusterByClusterIDVariables struct {
	actions.GraphQLQuery
	OrgID     string
	ClusterID string
}

func NewDeleteClusterByClusterIDVariables(orgID, clusterID string) DeleteClusterByClusterIDVariables {
	vars := DeleteClusterByClusterIDVariables{
		OrgID:     orgID,
		ClusterID: clusterID,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryDeleteClusterByClusterID
	vars.Args = map[string]string{
		"orgId":     "String!",
		"clusterId": "String!",
	}
	vars.Returns = []string{
		"deletedClusterCount",
		"deletedResourceCount",
	}

	return vars
}

type DeleteClustersResponse struct {
	Data *DeleteClustersResponseData `json:"data,omitempty"`
}

type DeleteClustersResponseData struct {
	Details *DeleteClustersResponseDataDetails `json:"deleteClusterByClusterId,omitempty"`
}

type DeleteClustersResponseDataDetails struct {
	DeletedClusterCount  int `json:"deletedClusterCount,omitempty"`
	DeletedResourceCount int `json:"deletedResourceCount,omitempty"`
}

func (c *Client) DeleteClusterByClusterID(orgID, clusterID, token string) (*DeleteClustersResponseDataDetails, error) {
	var response DeleteClustersResponse

	vars := NewDeleteClusterByClusterIDVariables(orgID, clusterID)

	err := c.DoQuery(DeleteClusterByClusterIDVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
