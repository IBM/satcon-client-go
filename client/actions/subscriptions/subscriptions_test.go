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
	"github.ibm.com/coligo/satcon-client/client/types"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("SubscriptionsByOrgId", func() {

	var (
		orgID string
	)

	BeforeEach(func() {
		orgID = "project-jupiter"
	})

	Describe("NewSubscriptionsVariables", func() {
		It("Returns a correctly populated instance of SubscriptionsVariables", func() {
			vars := NewSubscriptionsVariables(orgID)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QuerySubscriptions))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"orgId",
				"name",
				"uuid",
				"groups",
				"channelName",
				"channelUuid",
				"version",
			))
		})
	})

	Describe("Subcriptions", func() {

		var (
			token                 string
			subscriptionsResponse SubscriptionsResponse
			c                     SubscriptionService
			httpClient            *webfakes.FakeHTTPClient
			response              *http.Response
		)

		BeforeEach(func() {
			token = "monolith"

			subscriptionsResponse = SubscriptionsResponse{
				Data: &SubscriptionsResponseData{
					Subscriptions: types.SubscriptionList{
						{
							UUID:  "hal",
							OrgID: orgID,
							Name:  "subscription1",
						},
						{
							UUID:  "9000",
							OrgID: orgID,
							Name:  "subscription4",
						},
						{
							UUID:  "2001",
							OrgID: orgID,
							Name:  "subscription9",
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(subscriptionsResponse)
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
			_, err := c.Subscriptions(orgID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of subscriptions", func() {
			subscriptions, _ := c.Subscriptions(orgID, token)
			expected := subscriptionsResponse.Data.Subscriptions
			Expect(subscriptions).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("None whatsoever, Frank."))
			})

			It("Bubbles up the error", func() {
				_, err := c.Subscriptions(orgID, token)
				Expect(err).To(MatchError(MatchRegexp("None whatsoever, Frank.")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(SubscriptionsResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				subscriptions, err := c.Subscriptions(orgID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(subscriptions).To(BeNil())
			})
		})

	})

})
