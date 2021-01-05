package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/users"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Signing in a User", func() {
	var (
		login, password string
		c               UserService
		h               *webfakes.FakeHTTPClient
		response        *http.Response
		fakeAuthClient  authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		login = "foo@bar.ibm.com"
		password = "supersecretpassword"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		c, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewSignInVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewSignInVariables(login, password)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(MutationSignIn))
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
			signInResponse SignInResponse
		)

		BeforeEach(func() {
			signInResponse = SignInResponse{
				Data: &SignInResponseData{
					Details: &SignInResponseDataDetails{
						Token: "ey123.mytoken",
					},
				},
			}

			respBodyBytes, err := json.Marshal(signInResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.SignIn(login, password)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the sign in token", func() {
			details, _ := c.SignIn(login, password)
			Expect(details).NotTo(BeNil())

			expected := signInResponse.Data.Details.Token
			Expect(*details).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.SignIn(login, password)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(SignInResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.SignIn(login, password)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
