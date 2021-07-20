package resources

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryResourcesByCluster       = "resourcesByCluster"
	ResourcesByClusterVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"clusterId":{{json .ClusterID}}, "filter":{{json .Filter}},"limit":{{json .Limit}}{{end}}`
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
		"resources{id, orgId, clusterId, selfLink, searchableData, subscription{uuid, orgId, name, groups, channel{uuid, orgId, name, created}, version}}",
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
func (c *Client) ResourcesByCluster(orgID, clusterID, filter string, limit int, lastResource *types.Resource) (*types.ResourceList, error) {
	var response ResourcesByClusterResponse

	vars := NewResourcesByClusterVariables(orgID, clusterID, filter, limit)

	err := c.DoQuery(ResourcesByClusterVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}


/*

	fromDate, limit 		

	-For the first call, don't include fromDate
	

	fromDate: some 2020ish date 

	call the api using fromDate and the limit to get the first page

	1) the resources need to be sorted by creation time. 

	2) loop through fromDate to now
		get the *limit number of resources (the current batch/page of reources)

		- now update query using the lastResource and use the limit again 


	



*/