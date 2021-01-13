package local_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/IBM/satcon-client-go/client/auth/local"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Signing in a User", func() {
	var (
		endpoint, login, password string
		h                         *webfakes.FakeHTTPClient
		response                  *http.Response
	)

	BeforeEach(func() {
		endpoint = "https://foo.bar"
		login = "foo@bar.ibm.com"
		password = "supersecretpassword"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewSignInVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := local.NewSignInVariables(login, password)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(local.MutationSignIn))
			Expect(vars.Login).To(Equal(login))
			Expect(vars.Password).To(Equal(password))
			Expect(vars.Args).To(Equal(map[string]string{
				"login":    "String!",
				"password": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"token",
			))
		})
	})

	Describe("SignIn", func() {
		var (
			signInResponse local.SignInResponse
		)

		BeforeEach(func() {
			signInResponse = local.SignInResponse{
				Data: &local.SignInResponseData{
					Details: &local.SignInResponseDataDetails{
						Token: "ey123.mytoken",
					},
				},
			}

			respBodyBytes, err := json.Marshal(signInResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := local.SignIn(h, endpoint, login, password)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the sign in token", func() {
			details, _ := local.SignIn(h, endpoint, login, password)
			Expect(details).NotTo(BeNil())

			expected := signInResponse.Data.Details.Token
			Expect(*details).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := local.SignIn(h, endpoint, login, password)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(local.SignInResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := local.SignIn(h, endpoint, login, password)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
