package integration_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

var _ = Describe("Groups", func() {

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

	Describe("Group Lifecycle", func() {

		var (
			groupName string
		)

		BeforeEach(func() {
			groupName = RandStringBytes(8)
			fmt.Println("Using group name: ", groupName)
		})

		It("Lists the groups, creates our new group, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			// List the groups
			groups, err := c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// create a new group
			newGroupDetails, err := c.Groups.AddGroup(testConfig.OrgID, groupName, token)
			Expect(newGroupDetails.UUID).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())

			// list groups again and find our new group
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName) == 0 {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// delete the group
			removeGroup, err := c.Groups.RemoveGroupByName(testConfig.OrgID, groupName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeGroup.UUID).To(MatchRegexp(newGroupDetails.UUID))

			// list groups again and prove group is gone
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

		})

	})

})
