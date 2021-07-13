package channels

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryChannelByName       = "channelByName"
	ChannelByNameVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}","name":"{{json .Name}}"{{end}}`
)

// ChannelByNameVariables to query channel by name
type ChannelByNameVariables struct {
	actions.GraphQLQuery
	OrgID string
	Name  string
}

// NewChannelByNameVariables returns appropriate variables to query channel
func NewChannelByNameVariables(orgID, channelName string) ChannelByNameVariables {
	vars := ChannelByNameVariables{
		OrgID: orgID,
		Name:  channelName,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryChannelByName
	vars.Args = map[string]string{
		"orgId": "String!",
		"name":  "String!",
	}
	vars.Returns = []string{
		"uuid",
		"orgId",
		"name",
		"created",
		"versions{uuid, name, description, location, created}",
		"subscriptions{uuid, orgId, name, groups, channelUuid, channelName, version, versionUuid, created, updated}",
	}

	return vars
}

// ChannelVersionByNameResponse top level response struct
type ChannelByNameResponse struct {
	Data *ChannelByNameResponseData `json:"data,omitempty"`
}

// ChannelVersionByNameResponseData provides response details
type ChannelByNameResponseData struct {
	Details *types.Channel `json:"channelByName,omitempty"`
}

// ChannelVersionByName queries a channel version given orgID, channelName, and versionName
func (c *Client) ChannelByName(orgID, channelName string) (*types.Channel, error) {
	var response ChannelByNameResponse

	vars := NewChannelByNameVariables(orgID, channelName)

	err := c.DoQuery(ChannelByNameVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
