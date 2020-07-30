package subscriptions

import (
	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QueryRemoveSubscription       = "removeSubscription"
	RemoveSubscriptionVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}","uuid":"{{js .UUID}}"{{end}}`
)

// RemoveSubscriptionVariables are the variables specific to adding a group.
// These include the organization ID and the group name.  Rather than
// instantiating this directly, use NewRemoveSubscriptionVariables().
type RemoveSubscriptionVariables struct {
	actions.GraphQLQuery
	OrgID string
	UUID  string
}

// NewRemoveSubscriptionVariables creates a correctly formed instance of RemoveSubscriptionVariables.
func NewRemoveSubscriptionVariables(orgID, uuid string) RemoveSubscriptionVariables {
	vars := RemoveSubscriptionVariables{
		OrgID: orgID,
		UUID:  uuid,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = QueryRemoveSubscription
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

// RemoveSubscriptionResponse is the response body we get upon a successful cluster
// registration.
type RemoveSubscriptionResponse struct {
	Data *RemoveSubscriptionResponseData `json:"data,omitempty"`
}

type RemoveSubscriptionResponseData struct {
	Details *RemoveSubscriptionResponseDataDetails `json:"removeSubscription,omitempty"`
}

type RemoveSubscriptionResponseDataDetails struct {
	UUID    string `json:"uuid,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// RemoveSubscription deletes specified subscription
func (c *Client) RemoveSubscription(orgID, uuid, token string) (*RemoveSubscriptionResponseDataDetails, error) {
	var response RemoveSubscriptionResponse

	vars := NewRemoveSubscriptionVariables(orgID, uuid)

	err := c.DoQuery(RemoveSubscriptionVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Details, err
	}

	return nil, err
}
