package resources

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
	"github.ibm.com/coligo/satcon-client/client/types"
)

const (
	QueryResources       = "resources"
	ResourcesVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}", "filter":"{{js .Filter}}", "fromDate":"{{js .FromDate}}", "toDate": "{{js .ToDate}}","limit":{{js .Limit}}, "kinds":[{{range $i,$e := .Kinds}}{{if gt $i 0}},{{end}}"{{js $e}}"{{end}}], "sort":[{{range $i, $e := .Sort}}{{if gt $i 0}},{{end}}{"field": "{{js $e.Field}}", "desc": {{js $e.Desc}}}{{end}}], "subscriptionsLimit": {{js .SubscriptionsLimit}}{{end}}`
)

// ResourcesVariables variable to query resources for specified cluster
type ResourcesVariables struct {
	actions.GraphQLQuery
	OrgID              string
	Filter             string
	FromDate           string
	ToDate             string
	Limit              int
	Kinds              []string
	Sort               []SortObj
	SubscriptionsLimit int
}

// SortObj will sort on input field. Options are:
// ['_id', 'cluster_id', 'selfLink', 'created', 'updated', 'lastModified', 'deleted', 'hash']
type SortObj struct {
	Field string
	Desc  bool
}

// NewResourcesVariables returns necessary variables for query
func NewResourcesVariables(orgID, filter, fromDate, toDate string, limit int, kinds []string, sort []SortObj, subscriptionsLimit int) ResourcesVariables {
	vars := ResourcesVariables{
		OrgID:              orgID,
		Filter:             filter,
		FromDate:           fromDate,
		ToDate:             toDate,
		Limit:              limit,
		Kinds:              kinds,
		Sort:               sort,
		SubscriptionsLimit: subscriptionsLimit,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryResources
	vars.Args = map[string]string{
		"orgId":              "String!",
		"filter":             "String",
		"fromDate":           "Date",
		"toDate":             "Date",
		"limit":              "Int",
		"kinds":              "[String!]",
		"sort":               "[SortObj!]",
		"subscriptionsLimit": "Int",
	}
	vars.Returns = []string{
		"count",
		"resources{id, orgId, clusterId, cluster{clusterId, name}, selfLink, hash, data, deleted, created, updated, lastModified, searchableData, searchableDataHash, subscription{uuid, orgId, name, groups, channelUuid, channelName, version, versionUuid, owner{id, name}, created, updated}}",
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
func (c *Client) Resources(orgID, filter, fromDate, toDate string, limit int, kinds []string,
	sort []SortObj, subscriptionsLimit int, token string) (*types.ResourceList, error) {
	var response ResourcesResponse

	vars := NewResourcesVariables(orgID, filter, fromDate, toDate, limit, kinds, sort, subscriptionsLimit)

	err := c.DoQuery(ResourcesVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.ResourceList, err
	}

	return nil, err
}
