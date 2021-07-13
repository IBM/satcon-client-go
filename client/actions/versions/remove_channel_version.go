package versions

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryRemoveChannelVersion       = "removeChannelVersion"
	RemoveChannelVersionVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}","uuid":"{{json .UUID}}"{{end}}`
)

// RemoveChannelVersionVariables are the variables specific to adding a group.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewRemoveChannelVersionVariables().
type RemoveChannelVersionVariables struct {
	actions.GraphQLQuery
	OrgID string
	UUID  string
}

// NewRemoveChannelVersionVariables creates a correctly formed instance of RemoveChannelVersionVariables.
func NewRemoveChannelVersionVariables(orgID, uuid string) RemoveChannelVersionVariables {
	vars := RemoveChannelVersionVariables{
		OrgID: orgID,
		UUID:  uuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRemoveChannelVersion
	vars.Args = map[string]string{
		"orgId": "String!",
		"uuid":  "String!",
	}
	vars.Returns = []string{
		"uuid",
		"success",
	}

	return vars
}

// RemoveChannelVersionResponse is the response body we get upon a successful cluster
// registration.
type RemoveChannelVersionResponse struct {
	Data *RemoveChannelVersionResponseData `json:"data,omitempty"`
}

type RemoveChannelVersionResponseData struct {
	Details *RemoveChannelVersionResponseDataDetails `json:"removeChannelVersion,omitempty"`
}

type RemoveChannelVersionResponseDataDetails struct {
	UUID    string `json:"uuid,omitempty"`
	Success bool   `json:"success,omitempty"`
}

func (c *Client) RemoveChannelVersion(orgID, uuid string) (*RemoveChannelVersionResponseDataDetails, error) {
	var response RemoveChannelVersionResponse

	vars := NewRemoveChannelVersionVariables(orgID, uuid)

	err := c.DoQuery(RemoveChannelVersionVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
