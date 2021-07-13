package resources

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryResources       = "resources"
	ResourcesVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}"{{end}}`
)

// ResourcesVariables variable to query resources for specified cluster
type ResourcesVariables struct {
	actions.GraphQLQuery
	OrgID string
}

// NewResourcesVariables returns necessary variables for query
func NewResourcesVariables(orgID string) ResourcesVariables {
	vars := ResourcesVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResources
	vars.Args = map[string]string{
		"orgId": "String!",
	}
	vars.Returns = []string{
		"count",
		"resources{id, orgId, clusterId, cluster{clusterId, name}, selfLink}",
	}

	return vars
}

// ResourcesResponse query data
type ResourcesResponse struct {
	Data *ResourcesResponseData `json:"data,omitempty"`
}

// ResourcesResponseData encapsulates ResourceList response
type ResourcesResponseData struct {
	ResourceList *types.ResourceList `json:"resources,omitempty"`
}

// Resources queries specified cluster for list of resources, i.e. Pod, Deployment, Service, etc.
func (c *Client) Resources(orgID string) (*types.ResourceList, error) {
	var response ResourcesResponse

	vars := NewResourcesVariables(orgID)

	err := c.DoQuery(ResourcesVarTemplate, vars, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}
