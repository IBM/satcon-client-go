package client

import (
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HTTPClient
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
