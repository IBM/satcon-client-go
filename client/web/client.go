package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.ibm.com/coligo/satcon-client/client/actions"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . HTTPClient
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SatConClient struct {
	Endpoint   string
	HTTPClient HTTPClient
}

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

		err = json.Unmarshal(body, result)
		if err != nil {
			return err
		}
	}

	return nil
}
