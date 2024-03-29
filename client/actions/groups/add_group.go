package groups

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryAddGroup       = "addGroup"
	AddGroupVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"name":{{json .Name}}{{end}}`
)

// AddGroupVariables are the variables specific to adding a group.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewAddGroupVariables().
type AddGroupVariables struct {
	actions.GraphQLQuery
	OrgID string
	Name  string
}

// NewAddGroupVariables creates a correctly formed instance of AddGroupVariables.
func NewAddGroupVariables(orgID, name string) AddGroupVariables {
	vars := AddGroupVariables{
		OrgID: orgID,
		Name:  name,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryAddGroup
	vars.Args = map[string]string{
		"orgId": "String!",
		"name":  "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

// AddGroupResponse is the response body we get upon a successful cluster
// registration.
type AddGroupResponse struct {
	Data *AddGroupResponseData `json:"data,omitempty"`
}

type AddGroupResponseData struct {
	Details *AddGroupResponseDataDetails `json:"addGroup,omitempty"`
}

type AddGroupResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

func (c *Client) AddGroup(orgID, name string) (*AddGroupResponseDataDetails, error) {
	var response AddGroupResponse

	vars := NewAddGroupVariables(orgID, name)

	err := c.DoQuery(AddGroupVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
