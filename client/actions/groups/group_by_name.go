package groups

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryGroupByName       = "groupByName"
	GroupByNameVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","name":"{{js .Name}}"{{end}}`
)

type GroupByNameVariables struct {
	actions.GraphQLQuery
	OrgID string
	Name  string
}

func NewGroupByNameVariables(orgID string, name string) GroupByNameVariables {
	vars := GroupByNameVariables{
		OrgID: orgID,
		Name:  name,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryGroupByName
	vars.Args = map[string]string{
		"orgId": "String!",
		"name":  "String!",
	}
	vars.Returns = []string{
		"uuid",
		"orgId",
		"name",
		"created",
		"clusters{id,orgId,clusterId,name,metadata}",
	}

	return vars
}

type GroupByNameResponse struct {
	Data *GroupByNameResponseData `json:"data,omitempty"`
}

type GroupByNameResponseData struct {
	Group *types.Group `json:"groupByName,omitempty"`
}

func (c *Client) GroupByName(orgID string, name string) (*types.Group, error) {
	var response GroupByNameResponse

	vars := NewGroupByNameVariables(orgID, name)

	err := c.DoQuery(GroupByNameVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Group, err
	}

	return nil, err
}
