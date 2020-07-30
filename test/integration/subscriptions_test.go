package integration_test

import (
	"encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

var _ = Describe("Subscriptions", func() {
	var (
		token   string
		c       client.SatCon
		content []byte
	)

	BeforeEach(func() {
		var err error
		c, _ = client.New(testConfig.SatConEndpoint, nil)
		Expect(c.Subscriptions).NotTo(BeNil())

		token, err = GetToken(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).NotTo(BeZero())

		encodedContent := "YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IGludGVncmF0aW9uX3Rlc3QKc3BlYzoKICBjb250YWluZXJzOgogIC0gbmFtZTogaW50ZWdyYXRpb25fdGVzdAogICAgaW1hZ2U6IGh0dHBkOmFscGluZQo="
		content, err = base64.StdEncoding.DecodeString(encodedContent)
		Expect(err).NotTo(HaveOccurred())
		Expect(content).NotTo(BeNil())

	})

	Describe("Subscription Lifecycle", func() {
		var (
			channelName      string
			versionName      string
			subscriptionName string
			description      string
			groups           []string
		)

		BeforeEach(func() {
			channelName = RandStringBytes(8)
			versionName = RandStringBytes(8)
			subscriptionName = RandStringBytes(8)
			fmt.Println("Using channel name: ", channelName)
			fmt.Println("Using version name: ", versionName)
			fmt.Println("Using subscription name: ", subscriptionName)
			description = fmt.Sprintf("Integration test version: %s for channel: %s", versionName, channelName)
			groups = []string{"integration-1", "integration-2"}
		})

		It("Lists the subscriptions, creates our new subscription, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			// Verify that our subscription does not already exist
			subscriptionList, err := c.Subscriptions.Subscriptions(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			for _, e := range subscriptionList {
				Expect(e.Name).NotTo(Equal(subscriptionName))
			}

			// TODO before we can add a subscription, we need to be able to add a new channel and version so that we can pass channelUuid and versionUuid as arguments
			// Demonstrate channel version does not exist for the arguments of the current channelName and versionName
			version, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(version).NotTo(Equal(versionName))

			// Create a channel
			channelDetails, err := c.Channels.AddChannel(testConfig.OrgID, channelName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(channelDetails).NotTo(BeNil())

			// Prove that channel exists
			channelList, err := c.Channels.Channels(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, channel := range channelList {
				if channel.UUID == channelDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// Create a channel version using the previously created channel
			versionDetails, err := c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, versionName, content, description, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(versionDetails).NotTo(BeNil())

			// Verify that channel version exists
			getVersionDetails, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersionDetails).NotTo(BeNil())

			// Create a new subscription using this channel and version
			subscriptionDetails, err := c.Subscriptions.AddSubscription(testConfig.OrgID, subscriptionName, channelDetails.UUID, versionDetails.VersionUUID, groups, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(subscriptionDetails).NotTo(BeNil())

			// Verify that our newly created subscription exists
			subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, subscription := range subscriptionList {
				if subscription.UUID == subscriptionDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// Remove Subscription
			removeSubscriptionDetails, err := c.Subscriptions.RemoveSubscription(testConfig.OrgID, subscriptionDetails.UUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeSubscriptionDetails).NotTo(BeNil())

			//Verify Subscription have been removed
			subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, subscription := range subscriptionList {
				if subscription.UUID == subscriptionDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// Remove channel version
			removeVersionDetails, err := c.Versions.RemoveChannelVersion(testConfig.OrgID, versionDetails.VersionUUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeVersionDetails).NotTo(BeNil())

			// Verify that channel version has been removed
			getVersionDetails, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersionDetails).To(BeNil())

			// Delete channel
			removeChannelDetails, err := c.Channels.RemoveChannel(testConfig.OrgID, channelDetails.UUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeChannelDetails).NotTo(BeNil())

			// Verify that channel has been removed
			channelList, err = c.Channels.Channels(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, channel := range channelList {
				if channel.UUID == channelDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeFalse())

		})
	})
})
