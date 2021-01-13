package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	"github.com/IBM/satcon-client-go/client/auth"
)

var _ = Describe("Users", func() {
	var (
		c         client.SatCon
		iamClient *auth.IAMClient
	)

	BeforeEach(func() {
		var err error
		iamClient, err = auth.NewIAMClient(testConfig.APIKey)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, nil, iamClient.Client)
		Expect(c.Clusters).NotTo(BeNil())
	})

	Describe("User Details", func() {
		BeforeEach(func() {
			fmt.Println("Using IAM user: ")
		})

		It("List user details", func() {
			me, err := c.Users.Me()
			Expect(err).NotTo(HaveOccurred())
			Expect(me).NotTo(BeNil())
			Expect(me.OrgId).To(Equal(testConfig.OrgID))
			Expect(me.Type).To(Equal("iam"))
		})
	})
})
