package auth

import "net/http"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AuthClient
type AuthClient interface {
	Authenticate(request *http.Request) error
}
