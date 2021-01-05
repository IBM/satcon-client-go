package users

import (
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/types"
)

const (
	MutationSignIn    = "signIn"
	SignInVarTemplate = `{{define "vars"}}"login":"{{js .Login}}","password":"{{js .Password}}"{{end}}`
)

// AddSignInVariables are the variables specific to log in a user.
// These include the login and password.  Rather than
// instantiating this directly, use NewSignInVariables().
type AddSignInVariables struct {
	actions.GraphQLQuery
	Login    string
	Password string
}

// NewSignInVariables creates a correctly formed instance of AddSignInVariables.
func NewSignInVariables(login string, password string) AddSignInVariables {
	vars := AddSignInVariables{
		Login:    login,
		Password: password,
	}

	vars.Type = actions.QueryTypeMutation
	vars.QueryName = MutationSignIn
	vars.Args = map[string]string{
		"login":    "String!",
		"password": "String!",
	}
	vars.Returns = []string{
		"token",
	}

	return vars
}

// SignInResponse is the response body we get upon a successful user
// creation.
type SignInResponse struct {
	Data *SignInResponseData `json:"data,omitempty"`
}

type SignInResponseData struct {
	Details *SignInResponseDataDetails `json:"signIn,omitempty"`
}

type SignInResponseDataDetails struct {
	Token types.Token `json:"token,omitempty"`
}

func (c *Client) SignIn(login string, password string) (*types.Token, error) {
	var response SignInResponse

	vars := NewSignInVariables(login, password)

	err := c.DoQueryNoAuth(SignInVarTemplate, vars, nil, &response)

	if err != nil {
		return nil, err
	}

	if response.Data != nil {
		return &response.Data.Details.Token, err
	}

	return nil, err
}
