package resources

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryResourceContent       = "resourceContent"
	ResourceContentVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}", "clusterId":"{{js .ClusterID}}", "resourceSelfLink":"{{js .ResourceSelfLink}}"{{end}}`
)

// ResourceContentVariables variable to query resources for specified cluster
type ResourceContentVariables struct {
	actions.GraphQLQuery
	OrgID            string
	ClusterID        string
	ResourceSelfLink string
}

// NewResourceContentVariables returns necessary variables for query
func NewResourceContentVariables(orgID, clusterID, resourceSelfLink string) ResourceContentVariables {
	vars := ResourceContentVariables{
		OrgID:            orgID,
		ClusterID:        clusterID,
		ResourceSelfLink: resourceSelfLink,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResourceContent
	vars.Args = map[string]string{
		"orgId":            "String!",
		"clusterId":        "String!",
		"resourceSelfLink": "String!",
	}
	vars.Returns = []string{
		"id",
		"histId",
		"content",
		"updated",
	}

	return vars
}

type ResourceContentResponse struct {
	Data *ResourceContentResponseData `json:"data,omitempty"`
}

type ResourceContentResponseData struct {
	ResourceContent *types.ResourceContentObj `json:"resourceContent,omitempty"`
}

//ResourceContent retrieves resource content
func (c *Client) ResourceContent(orgID, clusterID, resourceSelfLink string) (*types.ResourceContentObj, error) {
	var response ResourceContentResponse

	vars := NewResourceContentVariables(orgID, clusterID, resourceSelfLink)

	err := c.DoQuery(ResourceContentVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceContent, err
	}

	return nil, err
}
