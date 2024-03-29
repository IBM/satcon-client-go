package versions

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	//QueryChannelVersionByName specifies the query
	QueryChannelVersionByName = "channelVersionByName"
	// ChannelVersionByNameVarTemplate is the template used to create the graphql query
	ChannelVersionByNameVarTemplate = `{{define "vars"}}"orgId":{{json .OrgID}},"channelName":{{json .ChannelName}},"versionName":{{json .VersionName}}{{end}}`
)

//SubscriptionsVariables are the variables used for the subscription query
type ChannelVersionByNameVariables struct {
	actions.GraphQLQuery
	OrgID       string
	ChannelName string
	VersionName string
}

func NewChannelVersionByNameVariables(orgID, channelName, versionName string) ChannelVersionByNameVariables {
	vars := ChannelVersionByNameVariables{
		OrgID:       orgID,
		ChannelName: channelName,
		VersionName: versionName,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryChannelVersionByName
	vars.Args = map[string]string{
		"orgId":       "String!",
		"channelName": "String!",
		"versionName": "String!",
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

// ChannelVersionByNameResponse top level response struct
type ChannelVersionByNameResponse struct {
	Data *ChannelVersionByNameResponseData `json:"data,omitempty"`
}

// ChannelVersionByNameResponseData provides response details
type ChannelVersionByNameResponseData struct {
	Details *types.DeployableVersion `json:"channelVersionByName,omitempty"`
}

// ChannelVersionByName queries a channel version given orgID, channelName, and versionName
func (c *Client) ChannelVersionByName(orgID, channelName, versionName string) (*types.DeployableVersion, error) {
	var response ChannelVersionByNameResponse

	vars := NewChannelVersionByNameVariables(orgID, channelName, versionName)

	err := c.DoQuery(ChannelVersionByNameVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
