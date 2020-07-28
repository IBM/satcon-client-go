package subscriptions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/subscriptions"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("AddSubscription", func() {

	var (
		orgID       string
		name        string
		groups      []string
		channelUuid string
		versionUuid string
	)

	BeforeEach(func() {
		orgID = "ganymede"
		name = "callisto"
		groups = []string{
			"europa",
			"io",
			"thebe",
		}
		channelUuid = "valetudo"
		versionUuid = "carme"
	})

	Describe("NewSubscriptionsVariables", func() {
		It("Returns a correctly populated instance of SubscriptionsVariables", func() {
			vars := NewAddSubscriptionVariables(orgID, name, channelUuid, versionUuid, groups)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryAddSubscription))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":       "String!",
				"name":        "String!",
				"groups":      "[String!]!",
				"channelUuid": "String!",
				"versionUuid": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
			))
		})
	})

	Describe("AddSubcription", func() {

		var (
			token                   string
			addSubscriptionResponse AddSubscriptionResponse
			c                       SubscriptionService
			httpClient              *webfakes.FakeHTTPClient
			response                *http.Response
		)

		BeforeEach(func() {
			token = "monolith"

			addSubscriptionResponse = AddSubscriptionResponse{
				Data: &AddSubscriptionResponseData{
					Details: &AddSubscriptionResponseDataDetails{
						UUID: "cassini",
					},
				},
			}

			respBodyBytes, err := json.Marshal(addSubscriptionResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			httpClient = &webfakes.FakeHTTPClient{}
			Expect(httpClient.DoCallCount()).To(Equal(0))
			httpClient.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", httpClient)

		})

		It("Makes a valid http request", func() {
			_, err := c.AddSubscription(orgID, name, channelUuid, versionUuid, groups, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the uuid from the AddChannelReply", func() {
			uuid, _ := c.AddSubscription(orgID, name, channelUuid, versionUuid, groups, token)
			expectedUuid := addSubscriptionResponse.Data.Details
			Expect(uuid).To(Equal(expectedUuid))
			Expect(uuid.UUID).To(MatchRegexp(expectedUuid.UUID))
		})

		Context("When mutation execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("Mutation Failure"))
			})

			It("Bubbles up the error", func() {
				_, err := c.AddSubscription(orgID, name, channelUuid, versionUuid, groups, token)
				Expect(err).To(MatchError(MatchRegexp("Mutation Failure")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(SubscriptionsResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				uuid, err := c.AddSubscription(orgID, name, channelUuid, versionUuid, groups, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(uuid).To(BeNil())
			})
		})

	})

})
