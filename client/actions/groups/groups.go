package groups

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryGroups       = "groups"
	GroupsVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}"{{end}}`
)

type GroupsVariables struct {
	actions.GraphQLQuery
	OrgID string
}

func NewGroupsVariables(orgID string) GroupsVariables {
	vars := GroupsVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryGroups
	vars.Args = map[string]string{
		"orgId": "String!",
	}
	vars.Returns = []string{
		"uuid",
		"orgId",
		"name",
		"created",
	}

	return vars
}

type GroupsResponse struct {
	Data *GroupsResponseData `json:"data,omitempty"`
}

type GroupsResponseData struct {
	Groups types.GroupList `json:"groups,omitempty"`
}

func (c *Client) Groups(orgID, token string) (types.GroupList, error) {
	var response GroupsResponse

	vars := NewGroupsVariables(orgID)

	err := c.DoQuery(GroupsVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Groups, err
	}

	return nil, err
}
