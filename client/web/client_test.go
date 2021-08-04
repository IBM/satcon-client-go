package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
				requestTemplate = `{{define "vars"}}"name":{{json .Name}}{{end}}`
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
				err := s.DoQuery(requestTemplate, vars, nil, &result)
				Expect(err).NotTo(HaveOccurred())
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

			Context("When unmarshaling errors", func() {
				It("Bubbles up the unmarshal error", func() {
					err := s.DoQuery(requestTemplate, vars, nil, nil)
					_, ok := err.(*json.InvalidUnmarshalError)
					Expect(ok).To(BeTrue())
				})
			})

			Context("When BuildRequest returns an error", func() {
				BeforeEach(func() {
					fakeAuthClient.AuthenticateStub = func(r *http.Request) error {
						return errors.New("Failed to Authenticate!")
					}
				})

				It("Bubbles up the Authenticate error", func() {
					err := s.DoQuery(requestTemplate, vars, nil, nil)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fmt.Errorf("Failed to Authenticate!")))
				})
			})

			Context("When a request fails", func() {
				var response *http.Response

				BeforeEach(func() {
					fakeAuthClient.AuthenticateStub = func(r *http.Request) error {
						return nil
					}

					response = &http.Response{
						Body: ioutil.NopCloser(bytes.NewBufferString(`{"errors": [{"message": "Context creation failed: Your session expired. Sign in again","extensions": {"code": "UNAUTHENTICATED"}}]}`)),
					}

					h.DoReturns(response, nil)
				})

				It("Returns an error when the response body is parsed", func() {
					err := s.DoQuery(requestTemplate, vars, nil, nil)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fmt.Errorf("Context creation failed: Your session expired. Sign in again")))
				})

			})

		})

		Describe("CheckResponseForErrors", func() {
			var errorResponse *types.RequestError
			var badBytes []byte

			BeforeEach(func() {
				errorResponse = &types.RequestError{
					Errors: []types.RequestErrorDetails{
						{
							Message: "ah darn",
						},
					},
				}

				// JSON syntax error
				badBytes = []byte(`{"errors": [{"message": "Context creation failed: Your session expired. Sign in again","extensions": {"code": "UNAUTHENTICATED"}}`)
				// Expect(err).NotTo(HaveOccurred())
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

			It("Errors during Unmarshal", func() {
				err := CheckResponseForErrors(badBytes)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("unexpected end of JSON input"))
			})

			Context("When there are multiple error messages from request", func() {
				var bytes []byte
				BeforeEach(func() {
					bytes = []byte(`{"errors": [{"message": "First Error: Your session expired. Sign in again","extensions": {"code": "UNAUTHENTICATED"}}, {"message": "Second Error: Your session expired. Sign in again","extensions": {"code": "UNAUTHENTICATED"}}]}`)
				})

				It("Returns all of the error messages", func() {
					err := CheckResponseForErrors(bytes)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(fmt.Errorf("First Error: Your session expired. Sign in again, Second Error: Your session expired. Sign in again")))
				})
			})

		})
	})
})
