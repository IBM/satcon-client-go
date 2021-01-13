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
	. "github.com/IBM/satcon-client-go/client/actions/users"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Adding a User", func() {
	var (
		endpoint, username, email, password, orgName, role string
		c                                                  UserService
		h                                                  *webfakes.FakeHTTPClient
		response                                           *http.Response
		fakeAuthClient                                     authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		endpoint = "https://foo.bar"
		username = "foo"
		email = "foo@bar.ibm.com"
		password = "supersecretpassword"
		orgName = "myorg"
		role = "admin"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		c, _ = NewClient(endpoint, h, &fakeAuthClient)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewSignUpVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := local.NewSignUpVariables(username, email, password, orgName, role)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(local.MutationSignUp))
			Expect(vars.Username).To(Equal(username))
			Expect(vars.Email).To(Equal(email))
			Expect(vars.Password).To(Equal(password))
			Expect(vars.OrgName).To(Equal(orgName))
			Expect(vars.Role).To(Equal(role))
			Expect(vars.Args).To(Equal(map[string]string{
				"username": "String!",
				"email":    "String!",
				"password": "String!",
				"orgName":  "String!",
				"role":     "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"token",
			))
		})
	})

	Describe("SignUp", func() {
		var (
			signUpResponse local.SignUpResponse
		)

		BeforeEach(func() {
			signUpResponse = local.SignUpResponse{
				Data: &local.SignUpResponseData{
					Details: &local.SignUpResponseDataDetails{
						Token: "ey123.mytoken",
					},
				},
			}

			respBodyBytes, err := json.Marshal(signUpResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := local.SignUp(h, endpoint, username, email, password, orgName, role)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the token of the new user", func() {
			details, _ := local.SignUp(h, endpoint, username, email, password, orgName, role)
			Expect(details).NotTo(BeNil())

			expected := signUpResponse.Data.Details.Token
			Expect(*details).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := local.SignUp(h, endpoint, username, email, password, orgName, role)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(local.SignUpResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := local.SignUp(h, endpoint, username, email, password, orgName, role)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
