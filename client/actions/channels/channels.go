package channels

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryChannels       = "channels"
	ChannelsVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}}{{end}}`
)

type ChannelsVariables struct {
	actions.GraphQLQuery
	OrgID string
}

func NewChannelsVariables(orgID string) ChannelsVariables {
	vars := ChannelsVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryChannels
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

type ChannelsResponse struct {
	Data *ChannelsResponseData `json:"data,omitempty"`
}

type ChannelsResponseData struct {
	Channels types.ChannelList `json:"channels,omitempty"`
}

func (c *Client) Channels(orgID string) (types.ChannelList, error) {
	var response ChannelsResponse

	vars := NewChannelsVariables(orgID)

	err := c.DoQuery(ChannelsVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Channels, err
	}

	return nil, err
}
