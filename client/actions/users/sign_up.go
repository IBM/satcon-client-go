package users

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	MutationSignUp    = "signUp"
	SignUpVarTemplate = `{{define "vars"}}"username":"{{js .Username}}","email":"{{js .Email}}","password":"{{js .Password}}","orgName":"{{js .OrgName}}","role":"{{js .Role}}"{{end}}`
)

// AddSignUpVariables are the variables specific to adding a user.
// These include the username, email, password, organization name and role.  Rather than
// instantiating this directly, use NewSignUpVariables().
type AddSignUpVariables struct {
	actions.GraphQLQuery
	Username string
	Email    string
	Password string
	OrgName  string
	Role     string
}

// NewSignUpVariables creates a correctly formed instance of AddSignUpVariables.
func NewSignUpVariables(username string, email string, password string, orgName string, role string) AddSignUpVariables {
	vars := AddSignUpVariables{
		Username: username,
		Email:    email,
		Password: password,
		OrgName:  orgName,
		Role:     role,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = MutationSignUp
	vars.Args = map[string]string{
		"username": "String!",
		"email":    "String!",
		"password": "String!",
		"orgName":  "String!",
		"role":     "String!",
	}
	vars.Returns = []string{
		"token",
	}

	return vars
}

// SignUpResponse is the response body we get upon a successful user
// creation.
type SignUpResponse struct {
	Data *SignUpResponseData `json:"data,omitempty"`
}

type SignUpResponseData struct {
	Details *SignUpResponseDataDetails `json:"signUp,omitempty"`
}

type SignUpResponseDataDetails struct {
	Token types.Token `json:"token,omitempty"`
}

func (c *Client) SignUp(username string, email string, password string, orgName string, role string) (*types.Token, error) {
	var response SignUpResponse

	vars := NewSignUpVariables(username, email, password, orgName, role)

	err := c.DoQueryNoAuth(SignUpVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return &response.Data.Details.Token, err
	}

	return nil, err
}
