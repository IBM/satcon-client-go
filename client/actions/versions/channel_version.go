package versions

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	//QueryChannelVersion specifies the query
	QueryChannelVersion = "channelVersion"
	// ChannelVersionVarTemplate is the template used to create the graphql query
	ChannelVersionVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"channelUuid":{{json .ChannelUUID}},"versionUuid":{{json .VersionUUID}}{{end}}`
)

// ChannelVersionVariables are the variables used for the subscription query
type ChannelVersionVariables struct {
	actions.GraphQLQuery
	OrgID       string
	ChannelUUID string
	VersionUUID string
}

// NewChannelVersionVariables returns variables required for channelVersion query
func NewChannelVersionVariables(orgID, channelUuid, versionUuid string) ChannelVersionVariables {
	vars := ChannelVersionVariables{
		OrgID:       orgID,
		ChannelUUID: channelUuid,
		VersionUUID: versionUuid,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryChannelVersion
	vars.Args = map[string]string{
		"orgId":       "String!",
		"channelUuid": "String!",
		"versionUuid": "String!",
	}
	vars.Returns = []string{
		"orgId",
		"uuid",
		"channelId",
		"channelName",
		"name",
		"type",
		"description",
		"content",
		"created",
	}

	return vars
}

// ChannelVersionResponse top level response struct
type ChannelVersionResponse struct {
	Data *ChannelVersionResponseData `json:"data,omitempty"`
}

// ChannelVersionResponseData provides response details
type ChannelVersionResponseData struct {
	Details *types.DeployableVersion `json:"channelVersion,omitempty"`
}

// ChannelVersion queries a channel version given orgID, channelUuid, and versionUuid
func (c *Client) ChannelVersion(orgID, channelUuid, versionUuid string) (*types.DeployableVersion, error) {
	var response ChannelVersionResponse

	vars := NewChannelVersionVariables(orgID, channelUuid, versionUuid)

	err := c.DoQuery(ChannelVersionVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
