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

var _ = Describe("Versions", func() {
	var (
		description string
		c           client.SatCon
		content     []byte
		iamClient   *iam.Client
	)

	BeforeEach(func() {
		var err error
		iamClient, err = iam.NewIAMClient(testConfig.APIKey)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, nil, iamClient.Client)
		Expect(c.Versions).NotTo(BeNil())

		encodedContent := "YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IGludGVncmF0aW9uX3Rlc3QKc3BlYzoKICBjb250YWluZXJzOgogIC0gbmFtZTogaW50ZWdyYXRpb25fdGVzdAogICAgaW1hZ2U6IGh0dHBkOmFscGluZQo="
		content, err = base64.StdEncoding.DecodeString(encodedContent)
		Expect(err).NotTo(HaveOccurred())
		Expect(content).NotTo(BeNil())
	})

	Describe("Version Lifecycle", func() {
		var (
			channelName string
			versionName string
		)

		BeforeEach(func() {
			channelName = RandStringBytes(8)
			versionName = RandStringBytes(8)
			description = fmt.Sprintf("Integration test version: %s for channel: %s", versionName, channelName)
			fmt.Println("Using channel name: ", channelName)
			fmt.Println("Using version name: ", versionName)
		})

		It("Gets our channel version by name to show it does not exist, creates a channel, creates a "+
			"channel version, gets the versionby name, removes the version, removes the channel, then tries "+
			"to get the version by name again to show it no longer exists", func() {
			// Demonstrate channel version does not exist for the arguments of the current channelName and versionName
			version, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName)
			Expect(err).To(HaveOccurred())
			Expect(version).NotTo(Equal(versionName))

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
			versionDetails, err := c.Versions.AddChannelVersion(testConfig.OrgID, channelDetails.UUID, versionName, content, description)
			Expect(err).NotTo(HaveOccurred())
			Expect(versionDetails).NotTo(BeNil())

			// Verify that channel version exists (query by name)
			channelVersionByNameDetails, err := c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName)
			Expect(err).NotTo(HaveOccurred())
			Expect(channelVersionByNameDetails).NotTo(BeNil())

			// Verify the channel version exists (query by UUID)
			channelVersionDetails, err := c.Versions.ChannelVersion(testConfig.OrgID, channelVersionByNameDetails.ChannelID, channelVersionByNameDetails.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(channelVersionDetails).NotTo(BeNil())
			Expect(channelVersionDetails.Name).To(MatchRegexp(versionName))

			// Remove channel version
			removeVersionDetails, err := c.Versions.RemoveChannelVersion(testConfig.OrgID, versionDetails.VersionUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeVersionDetails).NotTo(BeNil())

			// Verify that channel version has been removed
			channelVersionByNameDetails, err = c.Versions.ChannelVersionByName(testConfig.OrgID, channelName, versionName)
			Expect(err).To(HaveOccurred())
			Expect(channelVersionByNameDetails).To(BeNil())

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

		})
	})
})
