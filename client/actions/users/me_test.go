package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/users"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Me", func() {

	var (
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
	})

	Describe("NewMeVariables", func() {
		It("Returns a correctly populated instance of MeVariables", func() {
			vars := NewMeVariables()
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryMe))
			Expect(vars.Returns).To(ConsistOf(
				"id",
				"type",
				"orgId",
				"identifier",
				"email",
				"role",
			))
		})
	})

	Describe("Me", func() {
		var (
			c          UserService
			h          *webfakes.FakeHTTPClient
			response   *http.Response
			meResponse MeResponse
		)

		BeforeEach(func() {
			meResponse = MeResponse{
				Data: &MeResponseData{
					User: &types.User{
						Id:         "1",
						Type:       "local",
						OrgId:      "myorg",
						Identifier: "admin",
						Email:      "admin@foo.ibm.com",
						Role:       "admin",
					},
				},
			}

			respBodyBytes, err := json.Marshal(meResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			Expect(h.DoCallCount()).To(Equal(0))
			h.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		})

		It("Makes a valid http request", func() {
			_, err := c.Me()
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the logged in user", func() {
			channel, _ := c.Me()
			expected := meResponse.Data.User
			Expect(channel).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Kablooie!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.Me()
				Expect(err).To(MatchError("Kablooie!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(MeResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				channel, err := c.Me()
				Expect(err).NotTo(HaveOccurred())
				Expect(channel).To(BeNil())
			})
		})
	})

})
