package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/web"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

type BadBody struct{}

func (b *BadBody) Read(p []byte) (int, error) {
	return 0, errors.New("BAD BODY")
}

func (b *BadBody) Close() error {
	return nil
}

var _ = Describe("Client", func() {
	Describe("SatConClient", func() {
		var (
			s        *SatConClient
			h        *webfakes.FakeHTTPClient
			endpoint string
		)

		BeforeEach(func() {
			endpoint = "https://foo.bar"
			h = &webfakes.FakeHTTPClient{}
			s = &SatConClient{
				Endpoint:   endpoint,
				HTTPClient: h,
			}
		})

		Describe("DoQuery", func() {
			type QueryVars struct {
				actions.GraphQLQuery
				Name string
			}

			type QueryResponse struct {
				Name string `json:"name"`
			}

			var (
				name            string
				requestTemplate string
				vars            QueryVars
				response        *http.Response
				token           string
				result          QueryResponse
			)

			BeforeEach(func() {
				// Setup the http stuff
				name = "george"
				respBodyBytes, _ := json.Marshal(QueryResponse{
					Name: name,
				})
				response = &http.Response{Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes))}
				h.DoReturns(response, nil)
				token = "this is a token"

				// Setup the template
				requestTemplate = `{{define "vars"}}"name":"{{js .Name}}"{{end}}`
				vars = QueryVars{
					Name: "foo",
				}

				vars.Type = actions.QueryTypeQuery
				vars.QueryName = "SomeQuery"
				vars.Args = map[string]string{"name": "String!"}
				vars.Returns = []string{"name"}
			})

			It("Does not error", func() {
				err := s.DoQuery(requestTemplate, vars, nil, &result, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Deserializes the response body into the result", func() {
				s.DoQuery(requestTemplate, vars, nil, &result, token)
				Expect(result.Name).To(Equal(name))
			})

			Context("When generating the request payload errors", func() {
				BeforeEach(func() {
					// This is a malformed template, see the tests for action.BuildRequestBody
					// for more granular test cases.  We just need any reason for it to error.
					requestTemplate = `{{define "vars"}}`
				})

				It("Bubbles up the RequestBodyError", func() {
					err := s.DoQuery(requestTemplate, vars, nil, &result, token)
					_, ok := err.(actions.RequestBodyError)
					Expect(ok).To(BeTrue())
				})
			})

			Context("When the http call errors", func() {
				BeforeEach(func() {
					h.DoReturns(nil, &url.Error{})
				})

				It("Bubbles up the error returned by the http client's .Do()", func() {
					err := s.DoQuery(requestTemplate, vars, nil, &result, token)
					_, ok := err.(*url.Error)
					Expect(ok).To(BeTrue())
				})
			})

			Context("When the response body cannot be read", func() {
				BeforeEach(func() {
					response.Body = &BadBody{}
				})

				// This would work the same for an error on Close(),
				// they both flow through ioutil.ReadAll()
				It("Bubbles up the read error", func() {
					err := s.DoQuery(requestTemplate, vars, nil, &result, token)
					Expect(err).To(MatchError("BAD BODY"))
				})
			})

			Context("When unmarshalling errors", func() {
				It("Bubbles up the unmarshal error", func() {
					err := s.DoQuery(requestTemplate, vars, nil, nil, token)
					_, ok := err.(*json.InvalidUnmarshalError)
					Expect(ok).To(BeTrue())
				})
			})
		})
	})
})
