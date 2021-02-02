package integration_test

import (
	"encoding/base64"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	"github.com/IBM/satcon-client-go/client/auth/iam"
	. "github.com/IBM/satcon-client-go/test/integration"
)

var _ = Describe("Subscriptions", func() {
	var (
		c         client.SatCon
		content   []byte
		iamClient *iam.IAMClient
	)

	BeforeEach(func() {
		var err error
		iamClient, err = iam.NewIAMClient(testConfig.APIKey)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, nil, iamClient.Client)
		Expect(c.Subscriptions).NotTo(BeNil())

		encodedContent := "YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IGludGVncmF0aW9uX3Rlc3QKc3BlYzoKICBjb250YWluZXJzOgogIC0gbmFtZTogaW50ZWdyYXRpb25fdGVzdAogICAgaW1hZ2U6IGh0dHBkOmFscGluZQo="
		content, err = base64.StdEncoding.DecodeString(encodedContent)
		Expect(err).NotTo(HaveOccurred())
		Expect(content).NotTo(BeNil())

	})

	Describe("Subscription Lifecycle", func() {
		var (
			channelName      string
			version1Name     string
			version2Name     string
			subscriptionName string
			description      string
			groups           []string
		)

		BeforeEach(func() {
			channelName = RandStringBytes(8)
			version1Name = RandStringBytes(8)
			version2Name = RandStringBytes(8)
			group1 := RandStringBytes(8)
			group2 := RandStringBytes(8)
			subscriptionName = RandStringBytes(8)
			fmt.Println("Using channel name: ", channelName)
			fmt.Println("Using version1 name: ", version1Name)
			fmt.Println("Using version2 name: ", version2Name)
			fmt.Println("Using subscription name: ", subscriptionName)
			description = fmt.Sprintf("Integration test version: %s for channel: %s", version1Name, channelName)
			groups = []string{group1, group2}
		})

		It("Lists the subscriptions, creates our new subscription, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			// Verify that our subscription does not already exist
			subscriptionList, err := c.Subscriptions.Subscriptions(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			for _, e := range subscriptionList {
				Expect(e.Name).NotTo(Equal(subscriptionName))
			}

			// TODO before we can add a subscription, we need to be able to add a new channel and version so that we can pass channelUuid and versionUuid as arguments
			// Demonstrate channel version does not exist for the arguments of the current channelName and versionName
			version, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
			Expect(err).To(HaveOccurred())
			Expect(version).NotTo(Equal(version1Name))

			version, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version2Name)
			Expect(err).To(HaveOccurred())
			Expect(version).NotTo(Equal(version2Name))

			// Create a channel
			channelDetails, err := c.Channels.AddChannel(testConfig.OrgID, channelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(channelDetails).NotTo(BeNil())

			// Prove that channel exists
			channelList, err := c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, channel := range channelList {
				if channel.UUID == channelDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// Create a channel version using the previously created channel
			version1Details, err := c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, version1Name, content, description)
			Expect(err).NotTo(HaveOccurred())
			Expect(version1Details).NotTo(BeNil())

			// Verify that channel version exists
			getVersionDetails, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersionDetails).NotTo(BeNil())

			// Create a new subscription using this channel and version
			subscriptionDetails, err := c.Subscriptions.AddSubscription(testConfig.OrgID, subscriptionName, channelDetails.UUID, version1Details.VersionUUID, groups)
			Expect(err).NotTo(HaveOccurred())
			Expect(subscriptionDetails).NotTo(BeNil())

			// Verify that our newly created subscription exists
			subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, subscription := range subscriptionList {
				if subscription.UUID == subscriptionDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// Create a new channel version
			version2Details, err := c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, version2Name, content, description)
			Expect(err).NotTo(HaveOccurred())
			Expect(version2Details).NotTo(BeNil())

			// Verify that channel version exists
			getVersion2Details, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version2Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersion2Details).NotTo(BeNil())

			// Set existing subscription to new version
			setSubscriptionDetail, err := c.Subscriptions.SetSubscription(testConfig.OrgID, subscriptionDetails.UUID, version2Details.VersionUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(setSubscriptionDetail).NotTo(BeNil())

			// Verify that our newly updated subscription is up to date
			subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, subscription := range subscriptionList {
				if subscription.UUID == setSubscriptionDetail.UUID && subscription.Version == version2Name {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// Remove Subscription
			removeSubscriptionDetails, err := c.Subscriptions.RemoveSubscription(testConfig.OrgID, subscriptionDetails.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeSubscriptionDetails).NotTo(BeNil())

			//Verify Subscription have been removed
			subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, subscription := range subscriptionList {
				if subscription.UUID == subscriptionDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// Remove channel version 1
			removeVersionDetails, err := c.Versions.RemoveChannelVersion(testConfig.OrgID, version1Details.VersionUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeVersionDetails).NotTo(BeNil())

			// Verify that channel version 1 has been removed
			getVersionDetails, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
			Expect(err).To(HaveOccurred())
			Expect(getVersionDetails).To(BeNil())

			// Remove channel version 2
			removeVersionDetails, err = c.Versions.RemoveChannelVersion(testConfig.OrgID, version2Details.VersionUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeVersionDetails).NotTo(BeNil())

			// Verify that channel version 2 has been removed
			getVersionDetails, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version2Name)
			Expect(err).To(HaveOccurred())
			Expect(getVersionDetails).To(BeNil())

			// Delete channel
			removeChannelDetails, err := c.Channels.RemoveChannel(testConfig.OrgID, channelDetails.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeChannelDetails).NotTo(BeNil())

			// Verify that channel has been removed
			channelList, err = c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, channel := range channelList {
				if channel.UUID == channelDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// Remove the groups
			for _, g := range groups {
				rg, err := c.Groups.RemoveGroupByName(testConfig.OrgID, g)
				Expect(rg.UUID).NotTo(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
			}

		})
	})
})
