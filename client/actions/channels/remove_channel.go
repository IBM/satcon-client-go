package channels

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryRemoveChannel       = "removeChannel"
	RemoveChannelVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"uuid":{{json .UUID}}{{end}}`
)

// RemoveChannelVariables are the variables specific to adding a group.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewRemoveChannelVariables().
type RemoveChannelVariables struct {
	actions.GraphQLQuery
	OrgID string
	UUID  string
}

// NewRemoveChannelVariables creates a correctly formed instance of RemoveChannelVariables.
func NewRemoveChannelVariables(orgID, uuid string) RemoveChannelVariables {
	vars := RemoveChannelVariables{
		OrgID: orgID,
		UUID:  uuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRemoveChannel
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

// RemoveChannelResponse is the response body we get upon a successful cluster
// registration.
type RemoveChannelResponse struct {
	Data *RemoveChannelResponseData `json:"data,omitempty"`
}

type RemoveChannelResponseData struct {
	Details *RemoveChannelResponseDataDetails `json:"removeChannel,omitempty"`
}

type RemoveChannelResponseDataDetails struct {
	UUID    string `json:"uuid,omitempty"`
	Success bool   `json:"success,omitempty"`
}

func (c *Client) RemoveChannel(orgID, uuid string) (*RemoveChannelResponseDataDetails, error) {
	var response RemoveChannelResponse

	vars := NewRemoveChannelVariables(orgID, uuid)

	err := c.DoQuery(RemoveChannelVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
