package subscriptions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/subscriptions"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("SubscriptionIdsForCluster", func() {

	var (
		orgID, clusterID string
		fakeAuthClient   authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "project-jupiter"
		clusterID = "some-cluster"
	})

	Describe("NewSubscriptionIdsForClusterVariables", func() {
		It("Returns a correctly populated instance of SubscriptionIdsForClusterVariables", func() {
			vars := NewSubscriptionIdsForClusterVariables(orgID, clusterID)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QuerySubscriptionIdsForCluster))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"clusterId": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
			))
		})
	})

	Describe("SubscriptionIdsForCluster", func() {

		var (
			subscriptionsResponse SubscriptionIdsForClusterResponse
			subscriptionIds       []string
			c                     SubscriptionService
			httpClient            *webfakes.FakeHTTPClient
			response              *http.Response
		)

		BeforeEach(func() {

			subscriptionIds = []string{
				"hal",
				"9000",
				"2001",
			}
			responseDataField := make([]types.UuidOnly, len(subscriptionIds))
			for i := 0; i < len(subscriptionIds); i++ {
				responseDataField[i] = types.UuidOnly{UUID: subscriptionIds[i]}
			}
			subscriptionsResponse = SubscriptionIdsForClusterResponse{
				Data: &SubscriptionIdsForClusterResponseData{
					SubscriptionIdsForCluster: responseDataField,
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

			c, _ = NewClient("https://foo.bar", httpClient, &fakeAuthClient)

		})

		It("Makes a valid http request", func() {
			_, err := c.SubscriptionIdsForCluster(orgID, clusterID)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of subscriptions", func() {
			subscriptions, _ := c.SubscriptionIdsForCluster(orgID, clusterID)
			Expect(subscriptions).To(Equal(subscriptionIds))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("None whatsoever, Frank."))
			})

			It("Bubbles up the error", func() {
				_, err := c.SubscriptionIdsForCluster(orgID, clusterID)
				Expect(err).To(MatchError(MatchRegexp("None whatsoever, Frank.")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(SubscriptionIdsForClusterResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				subscriptions, err := c.SubscriptionIdsForCluster(orgID, clusterID)
				Expect(err).NotTo(HaveOccurred())
				Expect(subscriptions).To(BeNil())
			})
		})

	})

})
