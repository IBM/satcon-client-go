package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

var _ = Describe("Versions", func() {
	var (
		token string
		c     client.SatCon
	)

	BeforeEach(func() {
		var err error
		c, _ = client.New(testConfig.SatConEndpoint, nil)
		Expect(c.Versions).NotTo(BeNil())

		token, err = GetToken(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).NotTo(BeZero())
	})

	Describe("Version Lifecycle", func() {
		var (
			versionName string
		)

		BeforeEach(func() {
			versionName = RandStringBytes(8)
			fmt.Println("Using version name: ", versionName)
		})

		//TODO WRITE THESE TESTS
		It("Lists the versions, creates our new version, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			//TODO need API to return a list of all channel versions
			/*
				versionList, err := c.Versions.Versions(testConfig.OrgID, token)
				Expect(err).NotTo(HaveOccurred())
				for _, e := range versionList {
					Expect(e.Name).NotTo(Equal(versionName))
				}
			*/

			// TODO before we can add a subscription, we need to be able to add a new channel and version so that we can pass channelUuid and versionUuid as arguments
		})
	})
})
