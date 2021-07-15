package groups

import "github.com/IBM/satcon-client-go/client/actions"

const (
	QueryRemoveGroup       = "removeGroup"
	RemoveGroupVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"uuid":{{json .UUID}}{{end}}`
)

// RemoveGroupVariables are the variables specific to removing a group by name.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewRemoveGroupVariables().
type RemoveGroupVariables struct {
	actions.GraphQLQuery
	OrgID string
	UUID  string
}

// NewRemoveGroupVariables creates a correctly formed instance of RemoveGroupVariables.
func NewRemoveGroupVariables(orgID, uuid string) RemoveGroupVariables {
	vars := RemoveGroupVariables{
		OrgID: orgID,
		UUID:  uuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRemoveGroup
	vars.Args = map[string]string{
		"orgId": "String!",
		"uuid":  "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

type RemoveGroupResponse struct {
	Data *RemoveGroupResponseData `json:"data,omitempty"`
}

type RemoveGroupResponseData struct {
	Details *RemoveGroupResponseDataDetails `json:"removeGroup,omitempty"`
}

type RemoveGroupResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

func (c *Client) RemoveGroup(orgID, uuid string) (*RemoveGroupResponseDataDetails, error) {
	var response RemoveGroupResponse

	vars := NewRemoveGroupVariables(orgID, uuid)

	err := c.DoQuery(RemoveGroupVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
