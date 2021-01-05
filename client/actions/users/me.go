package users

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	QueryMe       = "me"
	MeVarTemplate = ``
)

// MeVariables to query user
type MeVariables struct {
	actions.GraphQLQuery
}

// NewMeVariables returns required query variables
func NewMeVariables() MeVariables {
	vars := MeVariables{}

	vars.Type = actions.QueryTypeQuery
	vars.QueryName = QueryMe
	vars.Args = map[string]string{}
	vars.Returns = []string{
		"id",
		"type",
		"orgId",
		"identifier",
		"email",
		"role",
	}

	return vars
}

// MeResponse user data
type MeResponse struct {
	Data *MeResponseData `json:"data,omitempty"`
}

// MeResponseData user details
type MeResponseData struct {
	User *types.User `json:"me,omitempty"`
}

// Channel returns channel specified by channeUuid
func (c *Client) Me() (*types.User, error) {
	var response MeResponse

	vars := NewMeVariables()

	err := c.DoQuery(MeVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return response.Data.User, err
	}

	return nil, err
}
