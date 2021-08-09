package integration_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	"github.com/IBM/satcon-client-go/client/auth/iam"
	. "github.com/IBM/satcon-client-go/test/integration"
)

var _ = Describe("Channels", func() {
	var (
		c client.SatCon
	)

	BeforeEach(func() {
		var err error
		var iamClient *iam.Client
		iamClient, err = iam.NewIAMClient(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, iamClient.Client)
		Expect(c.Channels).NotTo(BeNil())
	})

	Describe("Channel Lifecycle", func() {
		var (
			channelName string
		)

		BeforeEach(func() {
			channelName = RandStringBytes(8)
			fmt.Println("Using channel name: ", channelName)
		})

		It("Lists the channels, creates our new channel, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			channelList, err := c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, channel := range channelList {
				if strings.Compare(channel.Name, channelName) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			details, err := c.Channels.AddChannel(testConfig.OrgID, channelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(details).NotTo(BeNil())
			for _, channel := range channelList {
				Expect(channel.UUID).NotTo(Equal(details.UUID))
			}

			channelList, err = c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, channel := range channelList {
				if channel.UUID == details.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			channel, err := c.Channels.Channel(testConfig.OrgID, details.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(channel).NotTo(BeNil())

			channelByName, err := c.Channels.ChannelByName(testConfig.OrgID, channelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(channelByName).NotTo(BeNil())

			rmDetails, err := c.Channels.RemoveChannel(testConfig.OrgID, details.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(rmDetails.Success).To(BeTrue())

			// Confirm channel has been removed using the channelName
			channelByName, err = c.Channels.ChannelByName(testConfig.OrgID, channelName)
			Expect(err).To(HaveOccurred())
			Expect(channelByName).To(BeNil())

			// Confirm channel has been removed using the channelUuid
			channel, err = c.Channels.Channel(testConfig.OrgID, details.UUID)
			Expect(err).To(HaveOccurred())
			Expect(channel).To(BeNil())

			channelList, err = c.Channels.Channels(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			for _, channel := range channelList {
				Expect(channel.UUID).NotTo(Equal(details.UUID))
			}
		})
	})
})
