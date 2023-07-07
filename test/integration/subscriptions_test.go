package integration_test

import (
	"encoding/base64"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	"github.com/IBM/satcon-client-go/client/actions/channels"
	"github.com/IBM/satcon-client-go/client/actions/versions"
	"github.com/IBM/satcon-client-go/client/auth/iam"
	"github.com/IBM/satcon-client-go/client/types"
	. "github.com/IBM/satcon-client-go/test/integration"
)

var _ = Describe("Subscriptions", func() {
	var (
		c         client.SatCon
		content   []byte
		iamClient *iam.Client
	)

	BeforeEach(func() {
		var err error
		iamClient, err = iam.NewIAMClient(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, iamClient.Client)
		Expect(c.Subscriptions).NotTo(BeNil())

		encodedContent := "YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IGludGVncmF0aW9uX3Rlc3QKc3BlYzoKICBjb250YWluZXJzOgogIC0gbmFtZTogaW50ZWdyYXRpb25fdGVzdAogICAgaW1hZ2U6IGh0dHBkOmFscGluZQo="
		content, err = base64.StdEncoding.DecodeString(encodedContent)
		Expect(err).NotTo(HaveOccurred())
		Expect(content).NotTo(BeNil())

	})

	Describe("Subscription Lifecycle", func() {
		var (
			channelName,
			version1Name,
			version2Name,
			subscriptionName,
			description string

			groupNames,
			groupIds,
			clusterNames,
			clusterIds []string

			groupCount,
			clusterCount int

			channelDetails                   *channels.AddChannelResponseDataDetails
			version1Details, version2Details *versions.AddChannelVersionResponseDataDetails
		)

		BeforeEach(func() {
			channelName = RandStringBytes(8)
			version1Name = RandStringBytes(8)
			version2Name = RandStringBytes(8)
			clusterCount = 3
			clusterNames = make([]string, clusterCount)
			clusterIds = make([]string, clusterCount)
			for i := 0; i < clusterCount; i++ {
				clusterNames[i] = RandStringBytes(8)
			}
			groupCount = 2
			groupNames = make([]string, groupCount)
			groupIds = make([]string, groupCount)
			for i := 0; i < groupCount; i++ {
				groupNames[i] = RandStringBytes(8)
			}
			subscriptionName = RandStringBytes(8)
			fmt.Println("Using channel name: ", channelName)
			fmt.Println("Using version1 name: ", version1Name)
			fmt.Println("Using version2 name: ", version2Name)
			fmt.Println("Using subscription name: ", subscriptionName)
			description = fmt.Sprintf("Integration test version: %s for channel: %s", version1Name, channelName)

			for i, name := range clusterNames {
				details, err := c.Clusters.RegisterCluster(testConfig.OrgID, types.Registration{Name: name})
				Expect(err).NotTo(HaveOccurred())
				Expect(details).NotTo(BeNil())
				for _, clusterId := range clusterIds {
					Expect(clusterId).NotTo(Equal(details.ClusterID))
				}
				clusterIds[i] = details.ClusterID
			}

			// create groups
			for i, groupName := range groupNames {
				newGroupDetails, err := c.Groups.AddGroup(testConfig.OrgID, groupName)
				Expect(newGroupDetails.UUID).NotTo(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
				groupIds[i] = newGroupDetails.UUID
			}

			// list groups again and find our new groups
			currentGroups, err := c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			currentGroupNames := make([]string, len(currentGroups))
			for i, group := range currentGroups {
				currentGroupNames[i] = group.Name
			}
			Expect(currentGroupNames).To(ContainElements(groupNames))

			// Group clusters: Clusters 0 and 1 in Group 0, Clusters 1 and 2 in Group 1
			for i, groupId := range groupIds {
				response, err := c.Groups.GroupClusters(testConfig.OrgID, groupId, clusterIds[i:i+2])
				Expect(response).NotTo(BeNil())
				Expect(response.Modified > 0).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
			}

			// Demonstrate channel version does not exist for the arguments of the current channelName and versionName
			version, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
			Expect(err).To(HaveOccurred())
			Expect(version).To(BeNil())

			version, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version2Name)
			Expect(err).To(HaveOccurred())
			Expect(version).To(BeNil())

			// Create a channel
			channelDetails, err = c.Channels.AddChannel(testConfig.OrgID, channelName)
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
			version1Details, err = c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, version1Name, content, description)
			Expect(err).NotTo(HaveOccurred())
			Expect(version1Details).NotTo(BeNil())

			// Verify that channel version exists
			getVersionDetails, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersionDetails).NotTo(BeNil())
			// Create a new channel version
			version2Details, err = c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, version2Name, content, description)
			Expect(err).NotTo(HaveOccurred())
			Expect(version2Details).NotTo(BeNil())

			// Verify that channel version exists
			getVersion2Details, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version2Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(getVersion2Details).NotTo(BeNil())

		})

		When("using both Subscriptions APIs", func() {
			It("Lists the subscriptions, creates our new subscription, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
				// Verify that our subscription does not already exist
				// Can't use SubscriptionIdsForCluster here, we don't have an id yet
				subscriptionList, err := c.Subscriptions.Subscriptions(testConfig.OrgID)
				Expect(err).NotTo(HaveOccurred())
				for _, e := range subscriptionList {
					Expect(e.Name).NotTo(Equal(subscriptionName))
				}

				// Create a new subscription using this channel and version
				subscriptionDetails, err := c.Subscriptions.AddSubscription(testConfig.OrgID, subscriptionName, channelDetails.UUID, version1Details.VersionUUID, groupNames)
				Expect(err).NotTo(HaveOccurred())
				Expect(subscriptionDetails).NotTo(BeNil())

				// Verify that our newly created subscription exists, first using Subscriptions API...
				subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID)
				Expect(err).NotTo(HaveOccurred())
				found := false
				for _, subscription := range subscriptionList {
					if subscription.UUID == subscriptionDetails.UUID {
						found = true
					}
				}
				Expect(found).To(BeTrue())
				// ...then with SubscriptionIdsForCluster API. Should be on all clusters
				for _, clusterId := range clusterIds {
					subscriptionIds, err := c.Subscriptions.SubscriptionIdsForCluster(testConfig.OrgID, clusterId)
					Expect(err).NotTo(HaveOccurred())
					Expect(subscriptionIds).To(ContainElement(subscriptionDetails.UUID))
					// fmt.Println(subscriptionDetails.UUID, " is present on cluster ", clusterId)
				}

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

				//Verify Subscription have been removed, first with Subscriptions API...
				subscriptionList, err = c.Subscriptions.Subscriptions(testConfig.OrgID)
				Expect(err).NotTo(HaveOccurred())
				found = false
				for _, subscription := range subscriptionList {
					if subscription.UUID == subscriptionDetails.UUID {
						found = true
					}
				}
				Expect(found).To(BeFalse())
				// ...then with SubscriptionIdsForCluster API. Should not be on any clusters
				for _, clusterId := range clusterIds {
					subscriptionIds, err := c.Subscriptions.SubscriptionIdsForCluster(testConfig.OrgID, clusterId)
					Expect(err).NotTo(HaveOccurred())
					Expect(subscriptionIds).NotTo(ContainElement(subscriptionDetails.UUID))
				}

			})
		})

		When("using SubscriptonIdsForCluster()", func() {})

		AfterEach(func() {
			// Remove channel version 1
			removeVersionDetails, err := c.Versions.RemoveChannelVersion(testConfig.OrgID, version1Details.VersionUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeVersionDetails).NotTo(BeNil())

			// Verify that channel version 1 has been removed
			getVersionDetails, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, version1Name)
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
			channelList, err := c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, channel := range channelList {
				if channel.UUID == channelDetails.UUID {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// Remove the groups
			for _, g := range groupNames {
				rg, err := c.Groups.RemoveGroupByName(testConfig.OrgID, g)
				Expect(rg.UUID).NotTo(BeEmpty())
				Expect(err).NotTo(HaveOccurred())
			}

			// Remove clusters
			for _, clusterId := range clusterIds {
				delDetails, err := c.Clusters.DeleteClusterByClusterID(testConfig.OrgID, clusterId)
				Expect(err).NotTo(HaveOccurred())
				Expect(delDetails.DeletedClusterCount).To(Equal(1))
			}
			clusterList, err := c.Clusters.ClustersByOrgID(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			for _, details := range clusterList {
				Expect(clusterIds).NotTo(ContainElement(details.ClusterID))
			}

		})
	})
})
