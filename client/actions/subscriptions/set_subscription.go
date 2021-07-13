package subscriptions

import (
	"github.com/IBM/satcon-client-go/client/actions"
)

const (
	QuerySetSubscription       = "setSubscription"
	SetSubscriptionVarTemplate = `{{define "vars"}}"orgId":"{{json .OrgID}}","uuid":"{{json .UUID}}","versionUuid":"{{json .VersionUUID}}"{{end}}`
)

type SetSubscriptionVariables struct {
	actions.GraphQLQuery
	OrgID       string
	UUID        string
	VersionUUID string
}

func NewSetSubscriptionVariables(orgID string, subscriptionUuid string, versionUuid string) SetSubscriptionVariables {
	vars := SetSubscriptionVariables{
		OrgID:       orgID,
		UUID:        subscriptionUuid,
		VersionUUID: versionUuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QuerySetSubscription
	vars.Args = map[string]string{
		"orgId":       "String!",
		"uuid":        "String!",
		"versionUuid": "String!",
	}
	vars.Returns = []string{
		"uuid",
	}

	return vars
}

// SetSubscriptionResponse for unmarshalling the response data
type SetSubscriptionResponse struct {
	Data *SetSubscriptionResponseData `json:"data,omitempty"`
}

// SetSubscriptionResponseData for unmarshalling response details
type SetSubscriptionResponseData struct {
	Details *SetSubscriptionResponseDataDetails `json:"setSubscription,omitempty"`
}

// SetSubscriptionResponseDataDetails for unmarshalling response uuid
type SetSubscriptionResponseDataDetails struct {
	UUID string `json:"uuid,omitempty"`
}

// SetSubscription changes a subscription to a new version
func (c *Client) SetSubscription(orgID string, subscriptionUUID string, versionUUID string) (*SetSubscriptionResponseDataDetails, error) {
	var response SetSubscriptionResponse

	vars := NewSetSubscriptionVariables(orgID, subscriptionUUID, versionUUID)

	err := c.DoQuery(SetSubscriptionVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, nil
	}

	return nil, err
}
