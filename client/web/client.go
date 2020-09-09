package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	"github.ibm.com/coligo/satcon-client/client/actions"
	"github.ibm.com/coligo/satcon-client/client/types"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HTTPClient
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SatConClient struct {
	Endpoint   string
	HTTPClient HTTPClient
}

// DoQuery makes the graphql query request and returns the result
func (s *SatConClient) DoQuery(requestTemplate string, vars interface{}, funcs template.FuncMap, result interface{}, token string) error {
	payload, err := actions.BuildRequestBody(requestTemplate, vars, funcs)
	if err != nil {
		return err
	}

	req := actions.BuildRequest(payload, s.Endpoint, token)

	response, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if response.Body != nil {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if err = CheckResponseForErrors(body); err != nil {
			return err
		}

		err = json.Unmarshal(body, result)

		if err != nil {
			return err
		}
	}

	return nil
}

/*
 * CheckResponseForErrors takes the request and determines if an "errors" field is present. This is
 * done because as long as the graphql request receives a properly formed request, it will return a
 * 200OK, even if the request is "bad". For example the, attempting to access resources using an
 * orgID that is not accessible via the token will still return a 200OK, but will contain an error
 * message in the body. This function will parse that error message and return to user to provide
 * better information about the request and better error handling.
 */
func CheckResponseForErrors(body []byte) error {
	if strings.Contains(string(body), "errors") {
		var errorDetails *types.RequestError
		err := json.Unmarshal(body, &errorDetails)
		if err != nil {
			return err
		}

		if errorDetails != nil {
			var errorMessage string
			for i := range errorDetails.Errors {
				errorMessage += fmt.Sprintf("%s ", errorDetails.Errors[i].Message)
			}
			return errors.New(errorMessage)
		}
	}
	return nil
}
