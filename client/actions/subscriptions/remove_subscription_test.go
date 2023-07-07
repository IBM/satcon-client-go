package subscriptions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/subscriptions"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Removing a Subscription", func() {
	var (
		orgID, uuid    string
		c              SubscriptionService
		h              *webfakes.FakeHTTPClient
		response       *http.Response
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "someorg"
		uuid = "somesubscription"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {

		c, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewRemoveSubscriptionVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewRemoveSubscriptionVariables(orgID, uuid)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryRemoveSubscription))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"uuid":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"success",
			))
		})
	})

	Describe("RemoveSubscription", func() {
		var (
			rcResponse RemoveSubscriptionResponse
		)

		BeforeEach(func() {
			rcResponse = RemoveSubscriptionResponse{
				Data: &RemoveSubscriptionResponseData{
					Details: &RemoveSubscriptionResponseDataDetails{
						UUID:    "abacab-is-a-genesis-album",
						Success: true,
					},
				},
			}

			respBodyBytes, err := json.Marshal(rcResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.RemoveSubscription(orgID, uuid)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the details of subscription removal", func() {
			details, _ := c.RemoveSubscription(orgID, uuid)
			Expect(details).NotTo(BeNil())

			expected := rcResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("robots will win!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.RemoveSubscription(orgID, uuid)
				Expect(err).To(MatchError("robots will win!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(RemoveSubscriptionResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.RemoveSubscription(orgID, uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
