package integration_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	. "github.com/IBM/satcon-client-go/test/integration"
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
			groupName1 string
			groupName2 string
		)

		BeforeEach(func() {
			groupName1 = RandStringBytes(8)
			groupName2 = RandStringBytes(8)
			fmt.Printf("groupName1 = %s\ngroupName2 = %s\n", groupName1, groupName2)
		})

		It("Lists the groups, creates our new group, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			// List the groups
			groups, err := c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// create a new group
			newGroupDetails, err := c.Groups.AddGroup(testConfig.OrgID, groupName1, token)
			Expect(newGroupDetails.UUID).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())

			// list groups again and find our new group
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// delete the group using RemoveGroupByName
			removeGroup1, err := c.Groups.RemoveGroupByName(testConfig.OrgID, groupName1, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeGroup1.UUID).To(MatchRegexp(newGroupDetails.UUID))

			// list groups again and prove group is gone
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// List groups again to show group does not exist
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName2) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// create a new group
			newGroupDetails, err = c.Groups.AddGroup(testConfig.OrgID, groupName2, token)
			Expect(newGroupDetails.UUID).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())

			// list groups again and find our new group
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.UUID, newGroupDetails.UUID) == 0 {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// delete the group using RemoveGroup
			removeGroup2, err := c.Groups.RemoveGroup(testConfig.OrgID, newGroupDetails.UUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeGroup2.UUID).To(MatchRegexp(newGroupDetails.UUID))

			// list groups again and prove group is gone
			groups, err = c.Groups.Groups(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.UUID, newGroupDetails.UUID) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())
		})

	})

})
