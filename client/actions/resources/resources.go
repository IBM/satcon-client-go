package resources

import (
	"strings"

	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const QueryResources = "resources"

// Use strings.Replace so this template is easier to read in code
var ResourcesVarTemplate = strings.Replace(
	`{{define "vars"}}
"orgId":{{json .OrgID}},
"filter":{{json .Filter}},
"mongoQuery":{{json .MongoQuery}},
"fromDate":{{json .FromDate}},
"todate":{{json .ToDate}},
"limit":{{json .Limit}},
"kinds":{{json .Kinds}},
"sort":{{json .Sort}},
"subscriptionslimit":{{json .SubscriptionsLimit}}
{{end}}`,
	"\n", "", -1)

type ResourcesVariables struct {
	actions.GraphQLQuery
	types.ResourcesParams
}

// NewResourcesVariables returns necessary variables for query
func NewResourcesVariables(params types.ResourcesParams) ResourcesVariables {
	vars := ResourcesVariables{
		ResourcesParams: params,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResources
	vars.Args = map[string]string{
		"orgId":      "String!",
		"filter":     "String",
		"mongoQuery": "JSON",
		"fromDate":   "Date",
		"toDate":     "Date",
		"limit":      "Int",
		// TODO do we need this? skip: Int
		"kinds":              "[String!]",
		"sort":               "[SortObj!]",
		"subscriptionsLimit": "Int",
	}
	vars.Returns = []string{
		"count",
		"resources{id, orgId, clusterId, selfLink, searchableData, created, deleted, subscription{uuid, orgId, name, groups, channel{uuid, orgId, name, created}, version}}",
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
func (c *Client) Resources(params types.ResourcesParams) (*types.ResourceList, error) {
	var response ResourcesResponse

	vars := NewResourcesVariables(params)

	err := c.DoQuery(ResourcesVarTemplate, vars, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}
