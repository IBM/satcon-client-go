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

	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	. "github.com/IBM/satcon-client-go/client/web"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
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
		type QueryResponse struct {
			Name string `json:"name"`
		}

		var (
			s              *SatConClient
			h              *webfakes.FakeHTTPClient
			endpoint       string
			fakeAuthClient authfakes.FakeAuthClient
		)

		BeforeEach(func() {
			endpoint = "https://foo.bar"
			h = &webfakes.FakeHTTPClient{}
			s = &SatConClient{
				Endpoint:   endpoint,
				HTTPClient: h,
				AuthClient: &fakeAuthClient,
			}
		})

		Describe("DoQuery", func() {
			type QueryVars struct {
				actions.GraphQLQuery
				Name string
			}

			var (
				name            string
				requestTemplate string
				vars            QueryVars
				response        *http.Response
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
				err := s.DoQuery(requestTemplate, vars, nil, &result)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Deserializes the response body into the result", func() {
				s.DoQuery(requestTemplate, vars, nil, &result)
				Expect(result.Name).To(Equal(name))
			})

			Context("When generating the request payload errors", func() {
				BeforeEach(func() {
					// This is a malformed template, see the tests for action.BuildRequestBody
					// for more granular test cases.  We just need any reason for it to error.
					requestTemplate = `{{define "vars"}}`
				})

				It("Bubbles up the RequestBodyError", func() {
					err := s.DoQuery(requestTemplate, vars, nil, &result)
					_, ok := err.(actions.RequestBodyError)
					Expect(ok).To(BeTrue())
				})
			})

			Context("When the http call errors", func() {
				BeforeEach(func() {
					h.DoReturns(nil, &url.Error{})
				})

				It("Bubbles up the error returned by the http client's .Do()", func() {
					err := s.DoQuery(requestTemplate, vars, nil, &result)
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
					err := s.DoQuery(requestTemplate, vars, nil, &result)
					Expect(err).To(MatchError("BAD BODY"))
				})
			})

			Context("When unmarshalling errors", func() {
				It("Bubbles up the unmarshal error", func() {
					err := s.DoQuery(requestTemplate, vars, nil, nil)
					_, ok := err.(*json.InvalidUnmarshalError)
					Expect(ok).To(BeTrue())
				})
			})

		})

		Describe("CheckResponseForErrors", func() {
			var errorResponse *types.RequestError

			BeforeEach(func() {
				errorResponse = &types.RequestError{
					Errors: []types.RequestErrorDetails{
						{
							Message: "ah darn",
						},
					},
				}
			})

			It("Returns an error when error messages are detected", func() {
				body, _ := json.Marshal(errorResponse)
				err := CheckResponseForErrors(body)
				Expect(err).To(HaveOccurred())
			})

			It("Does not error when error messages are not detected", func() {
				respBodyBytes, _ := json.Marshal(QueryResponse{
					Name: "joe",
				})

				err := CheckResponseForErrors(respBodyBytes)
				Expect(err).NotTo(HaveOccurred())
			})

		})
	})
})
