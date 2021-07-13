package channels

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QueryAddChannel       = "addChannel"
	AddChannelVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}","name":"{{json .Name}}"{{end}}`
)

// AddChannelVariables are the variables specific to adding a group.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewAddChannelVariables().
type AddChannelVariables struct {
	actions.GraphQLQuery
	OrgID string
	Name  string
}

// NewAddChannelVariables creates a correctly formed instance of AddChannelVariables.
func NewAddChannelVariables(orgID, name string) AddChannelVariables {
	vars := AddChannelVariables{
		OrgID: orgID,
		Name:  name,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryAddChannel
	vars.Args = map[string]string{
		"orgId": "String!",
		"name":  "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

// AddChannelResponse is the response body we get upon a successful cluster
// registration.
type AddChannelResponse struct {
	Data *AddChannelResponseData `json:"data,omitempty"`
}

type AddChannelResponseData struct {
	Details *AddChannelResponseDataDetails `json:"addChannel,omitempty"`
}

type AddChannelResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

func (c *Client) AddChannel(orgID, name string) (*AddChannelResponseDataDetails, error) {
	var response AddChannelResponse

	vars := NewAddChannelVariables(orgID, name)

	err := c.DoQuery(AddChannelVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
