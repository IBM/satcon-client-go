package groups

import "github.com/IBM/satcon-client-go/client/actions"

const (
	QueryRemoveGroupByName       = "removeGroupByName"
	RemoveGroupByNameVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"name":{{json .Name}}{{end}}`
)

// RemoveGroupByNameVariables are the variables specific to removing a group by name.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewRemoveGroupByNameVariables().
type RemoveGroupByNameVariables struct {
	actions.GraphQLQuery
	OrgID string
	Name  string
}

// NewRemoveGroupByNameVariables creates a correctly formed instance of RemoveGroupByNameVariables.
func NewRemoveGroupByNameVariables(orgID, name string) RemoveGroupByNameVariables {
	vars := RemoveGroupByNameVariables{
		OrgID: orgID,
		Name:  name,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRemoveGroupByName
	vars.Args = map[string]string{
		"orgId": "String!",
		"name":  "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

type RemoveGroupByNameResponse struct {
	Data *RemoveGroupByNameResponseData `json:"data,omitempty"`
}

type RemoveGroupByNameResponseData struct {
	Details *RemoveGroupByNameResponseDataDetails `json:"removeGroupByName,omitempty"`
}

type RemoveGroupByNameResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

func (c *Client) RemoveGroupByName(orgID, name string) (*RemoveGroupByNameResponseDataDetails, error) {
	var response RemoveGroupByNameResponse

	vars := NewRemoveGroupByNameVariables(orgID, name)

	err := c.DoQuery(RemoveGroupByNameVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
