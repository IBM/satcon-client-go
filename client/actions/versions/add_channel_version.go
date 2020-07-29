package versions

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	ContentType                  = "application/yaml"
	QueryAddChannelVersion       = "addChannelVersion"
	AddChannelVersionVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","channelUuid":"{{js .ChannelUUID}}","name":"{{js .Name}}","type":"{{js .ContentType}}","content":"{{js .Content}}","description":"{{js .Description}}"{{end}}`
)

// AddChannelVersionVariables to create addChannelVersion graphql request
type AddChannelVersionVariables struct {
	actions.GraphQLQuery
	OrgID       string
	ChannelUUID string
	Name        string
	ContentType string
	Content     string
	File        string
	Description string
}

// NewAddChannelVersionVariables creates query variable
func NewAddChannelVersionVariables(orgID, channelUuid, name, contentType string, content string, file, description string) AddChannelVersionVariables {
	vars := AddChannelVersionVariables{
		OrgID:       orgID,
		ChannelUUID: channelUuid,
		Name:        name,
		ContentType: contentType,
		Content:     content,
		File:        file,
		Description: description,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryAddChannelVersion
	vars.Args = map[string]string{
		"orgId":       "String!",
		"channelUuid": "String!",
		"name":        "String!",
		"type":        "String!",
		"content":     "String",
		"file":        "Upload",
		"description": "String",
	}
	vars.Returns = []string{
		"versionUuid",
		"success",
	}

	return vars
}

// AddChannelVersionResponse for unmarshalling the response data
type AddChannelVersionResponse struct {
	Data *AddChannelVersionResponseData `json:"data,omitempty"`
}

// AddChannelVersionResponseData for unmarshalling response details
type AddChannelVersionResponseData struct {
	Details *AddChannelVersionResponseDataDetails `json:"addChannelVersion,omitempty"`
}

// AddChannelVersionResponseDataDetails for unmarshalling response uuid
type AddChannelVersionResponseDataDetails struct {
	VersionUUID string `json:"versionUuid,omitempty"`
	Success     bool   `json:"success,omitempty"`
}

// AddChannelVersion creates a new channelVersion for valid channel.
// contentFile is path to yaml file
func (c *Client) AddChannelVersion(orgID, channelUuid, name string, content []byte, description, token string) (*AddChannelVersionResponseDataDetails, error) {
	var response AddChannelVersionResponse

	vars := NewAddChannelVersionVariables(orgID, channelUuid, name, ContentType, string(content), "", description)

	err := c.DoQuery(AddChannelVersionVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
