package subscriptions

import (
	"fmt"
	"strings"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

const (
	QuerySubscriptions       = "subscriptions"
	SubscriptionsVarTemplate = `{{define "vars"}}"orgId":"{{js .OrgID}}"{{end}}`
)

type SubscriptionsVariables struct {
	actions.GraphQLQuery
	OrgID string
}

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
		"id",
		"orgId",
		"name",
		"groups",
	}

	return vars
}

type SubscriptionsResponse struct {
	Data *SubscriptionsResponseData `json:"data"`
}

type SubscriptionsResponseData struct {
	Subscriptions SubscriptionList `json:"subscriptions"`
}

type SubscriptionList []Subscription

func (l SubscriptionList) String() string {
	if len(l) == 0 {
		return "[]"
	}

	subscriptions := make([]string, 1)

	for _, s := range l {
		subscriptions = append(subscriptions, fmt.Sprintf("UUID: %s\nOrgID: %s\nName: %s\n", s.UUID, s.OrgID, s.Name))
	}
	return strings.Join(subscriptions, "\n==\n")
}

type Subscription struct {
	UUID        string              `json:"uuid,omitempty"`
	OrgID       string              `json:"orgId,omitempty"`
	Name        string              `json:"name,omitempty"`
	Groups      []SubscriptionGroup `json:"groups,omitempty"`
	ChannelUUID string              `json:"channelUuid,omitempty"`
	ChannelName string              `json:"channelName,omitempty"`
	Channel     Channel             `json:"channel,omitempty"`
	Version     string              `json:"version,omitempty"`
	VersionUUID string              `json:"versionUuid,omitempty"`
	Owner       BasicUser           `json:"owner,omitempty"`
	Created     string              `json:"created,omitempty"`
	Updated     string              `json:"updated,omitempty"`
}

type SubscriptionGroup struct {
	Scalar string `json:"scalar"`
}

type Channel struct {
	UUID     string           `json:"uuid"`
	OrgID    string           `json:"orgId"`
	Name     string           `json:"name"`
	Created  string           `json:"created"`
	Versions []ChannelVersion `json:"versions"`
}

type ChannelVersion struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Created     string `json:"created"`
}

type BasicChannelSubscriptions struct {
	UUID        string              `json:"uuid"`
	OrgID       string              `json:"orgId"`
	Name        string              `json:"name"`
	Groups      []SubscriptionGroup `json:"groups"`
	ChannelUUID string              `json:"channelUuid"`
	ChannelName string              `json:"channelName"`
	Version     string              `json:"version"`
	Created     string              `json:"created"`
	Updated     string              `json:"updated"`
}

type BasicUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) Subscriptions(orgID, token string) (SubscriptionList, error) {
	var response SubscriptionsResponse

	vars := NewSubscriptionsVariables(orgID)

	err := c.DoQuery(SubscriptionsVarTemplate, vars, nil, &response, token)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.Subscriptions, nil
	}

	return nil, err
}
