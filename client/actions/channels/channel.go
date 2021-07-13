package channels

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryChannel       = "channel"
	ChannelVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}","uuid":"{{json .UUID}}"{{end}}`
)

// ChannelVariables to query channel
type ChannelVariables struct {
	actions.GraphQLQuery
	OrgID string
	UUID  string
}

// NewChannelVariables returns required query variables
func NewChannelVariables(orgID, uuid string) ChannelVariables {
	vars := ChannelVariables{
		OrgID: orgID,
		UUID:  uuid,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryChannel
	vars.Args = map[string]string{
		"orgId": "String!",
		"uuid":  "String!",
	}
	vars.Returns = []string{
		"uuid",
		"orgId",
		"name",
		"created",
		"versions{uuid, name, location}",
		"subscriptions{uuid, orgId, name, groups}",
	}

	return vars
}

// ChannelResponse channel data
type ChannelResponse struct {
	Data *ChannelResponseData `json:"data,omitempty"`
}

// ChannelResponseData channel details
type ChannelResponseData struct {
	Channel *types.Channel `json:"channel,omitempty"`
}

// Channel returns channel specified by channeUuid
func (c *Client) Channel(orgID, uuid string) (*types.Channel, error) {
	var response ChannelResponse

	vars := NewChannelVariables(orgID, uuid)

	err := c.DoQuery(ChannelVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Channel, err
	}

	return nil, err
}
