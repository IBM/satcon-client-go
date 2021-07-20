package resources

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryResourcesByOrgID       = "resources"
	ResourcesByOrgIDVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}}{{end}}`
)

// ResourcesByOrgIDVariables variable to query resources for specified cluster
type ResourcesByOrgIDVariables struct {
	actions.GraphQLQuery
	OrgID string
}

// NewResourcesByOrgIDVariables returns necessary variables for query
func NewResourcesByOrgIDVariables(orgID string) ResourcesByOrgIDVariables {
	vars := ResourcesByOrgIDVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResourcesByOrgID
	vars.Args = map[string]string{
		"orgId": "String!",
	}
	vars.Returns = []string{
		"count",
		"resources{id, orgId, clusterId, cluster{clusterId, name}, selfLink}",
	}

	return vars
}

// ResourcesByOrgIDResponse query data
type ResourcesByOrgIDResponse struct {
	Data *ResourcesByOrgIDResponseData `json:"data,omitempty"`
}

// ResourcesByOrgIDResponseData encapsulates ResourceList response
type ResourcesByOrgIDResponseData struct {
	ResourceList *types.ResourceList `json:"resources,omitempty"`
}

// ResourcesByOrgID queries specified cluster for list of resources, i.e. Pod, Deployment, Service, etc.
func (c *Client) ResourcesByOrgID(orgID string) (*types.ResourceList, error) {
	var response ResourcesByOrgIDResponse

	vars := NewResourcesByOrgIDVariables(orgID)

	err := c.DoQuery(ResourcesByOrgIDVarTemplate, vars, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}
