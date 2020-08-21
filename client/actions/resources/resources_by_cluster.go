package resources

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
	"github.ibm.com/coligo/satcon-client/client/types"
)

const (
	QueryResourcesByCluster       = "resourcesByCluster"
	ResourcesByClusterVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","clusterId":"{{js .ClusterID}}", "filter":"{{js .Filter}}","limit":{{js .Limit}}{{end}}`
)

// ResourcesByClusterVariables variable to query resources for specified cluster
type ResourcesByClusterVariables struct {
	actions.GraphQLQuery
	OrgID     string
	ClusterID string
	Filter    string
	Limit     int
}

// NewResourcesByClusterVariables returns necessary variables for query
func NewResourcesByClusterVariables(orgID, clusterID, filter string, limit int) ResourcesByClusterVariables {
	vars := ResourcesByClusterVariables{
		OrgID:     orgID,
		ClusterID: clusterID,
		Filter:    filter,
		Limit:     limit,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResourcesByCluster
	vars.Args = map[string]string{
		"orgId":     "String!",
		"clusterId": "String!",
		"filter":    "String",
		"limit":     "Int",
	}
	vars.Returns = []string{
		"count",
		"resources{id, orgId, clusterId, cluster{clusterId, name}, selfLink, hash, data, deleted, created, updated, lastModified, searchableData, searchableDataHash, subscription{uuid, orgId, name, groups, channelUuid, channelName, version, versionUuid, owner{id, name}, created, updated}}",
	}

	return vars
}

// ResourcesByClusterResponse query data
type ResourcesByClusterResponse struct {
	Data *ResourcesByClusterResponseData `json:"data,omitempty"`
}

// ResourcesByClusterResponseData encapsulates ResourceList response
type ResourcesByClusterResponseData struct {
	ResourceList *types.ResourceList `json:"resourcesByCluster,omitempty"`
}

// ResourcesByCluster queries specified cluster for list of resources, i.e. Pod, Deployment, Service, etc.
func (c *Client) ResourcesByCluster(orgID, clusterID, filter string, limit int, token string) (*types.ResourceList, error) {
	var response ResourcesByClusterResponse

	vars := NewResourcesByClusterVariables(orgID, clusterID, filter, limit)

	err := c.DoQuery(ResourcesByClusterVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}
