package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

var _ = Describe("Channels", func() {
	var (
		token string
		c     client.SatCon
	)

	BeforeEach(func() {
		var err error
		c, _ = client.New(testConfig.SatConEndpoint, nil)
		Expect(c.Subscriptions).NotTo(BeNil())

		token, err = GetToken(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).NotTo(BeZero())
	})

	Describe("Subscription Lifecycle", func() {
		var (
			subscriptionName string
		)

		BeforeEach(func() {
			subscriptionName = RandStringBytes(8)
			fmt.Println("Using subscription name: ", subscriptionName)
		})

		It("Lists the subscriptions, creates our new subscription, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			subscriptionList, err := c.Subscriptions.Subscriptions(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			for _, e := range subscriptionList {
				Expect(e.Name).NotTo(Equal(subscriptionName))
			}

			// TODO before we can add a subscription, we need to be able to add a new channel and version so that we can pass channelUuid and versionUuid as arguments
		})
	})
})
