package subscriptions

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	//QuerySubscriptions specifies the query
	QuerySubscriptions = "subscriptions"
	//SubscriptionsVarTemplate is the template used to create the graphql query
	SubscriptionsVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}"{{end}}`
)

//SubscriptionsVariables are the variables used for the subscription query
type SubscriptionsVariables struct {
	actions.GraphQLQuery
	OrgID string
}

//NewSubscriptionsVariables generates variables used for query
func NewSubscriptionsVariables(orgID string) SubscriptionsVariables {
	vars := SubscriptionsVariables{
		OrgID: orgID,
	}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QuerySubscriptions
	vars.Args = map[string]string{
		"orgId": "String!",
	}
	vars.Returns = []string{
		"orgId",
		"name",
		"uuid",
		"groups",
		"channelName",
		"channelUuid",
		"version",
	}

	return vars
}

//SubscriptionsResponse response from query
type SubscriptionsResponse struct {
	Data *SubscriptionsResponseData `json:"data,omitempty"`
}

// SubscriptionsResponseData data response from query
type SubscriptionsResponseData struct {
	Subscriptions types.SubscriptionList `json:"subscriptions,omitempty"`
}

func (c *Client) Subscriptions(orgID string) (types.SubscriptionList, error) {
	var response SubscriptionsResponse

	vars := NewSubscriptionsVariables(orgID)

	err := c.DoQuery(SubscriptionsVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Subscriptions, nil
	}

	return nil, err
}
