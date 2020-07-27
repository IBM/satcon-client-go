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
		Expect(c.Channels).NotTo(BeNil())

		token, err = GetToken(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).NotTo(BeZero())
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
			channelList, err := c.Channels.Channels(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			// TODO: Iterate through the list and check that none of them is our channel
			// for _, channel := range channelList {
			// 	Expect(channel.)
			// }

			details, err := c.Channels.AddChannel(testConfig.OrgID, channelName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(details).NotTo(BeNil())
			for _, channel := range channelList {
				Expect(channel.UUID).NotTo(Equal(details.UUID))
			}

			channelList, err = c.Channels.Channels(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, channel := range channelList {
				if channel.UUID == details.UUID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			rmDetails, err := c.Channels.RemoveChannel(testConfig.OrgID, details.UUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(rmDetails.Success).To(BeTrue())

			channelList, err = c.Channels.Channels(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			for _, channel := range channelList {
				Expect(channel.UUID).NotTo(Equal(details.UUID))
			}
		})
	})
})
