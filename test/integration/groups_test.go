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

var _ = Describe("Groups", func() {

	var (
		c         client.SatCon
		iamClient *iam.Client
	)

	BeforeEach(func() {
		var err error
		iamClient, err = iam.NewIAMClient(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, nil, iamClient.Client)
		Expect(c.Groups).NotTo(BeNil())
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
			groups, err := c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// create a new group
			newGroupDetails, err := c.Groups.AddGroup(testConfig.OrgID, groupName1)
			Expect(newGroupDetails.UUID).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())

			// list groups again and find our new group
			groups, err = c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// get group by name
			group1, err := c.Groups.GroupByName(testConfig.OrgID, groupName1)
			Expect(err).NotTo(HaveOccurred())
			Expect(group1).NotTo(BeNil())
			Expect(group1.Name).To(Equal(groupName1))

			// delete the group using RemoveGroupByName
			removeGroup1, err := c.Groups.RemoveGroupByName(testConfig.OrgID, groupName1)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeGroup1.UUID).To(MatchRegexp(newGroupDetails.UUID))

			// list groups again and prove group is gone
			groups, err = c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName1) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// List groups again to show group does not exist
			groups, err = c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.Name, groupName2) == 0 {
					found = true
				}
			}
			Expect(found).To(BeFalse())

			// create a new group
			newGroupDetails, err = c.Groups.AddGroup(testConfig.OrgID, groupName2)
			Expect(newGroupDetails.UUID).NotTo(BeEmpty())
			Expect(err).NotTo(HaveOccurred())

			// list groups again and find our new group
			groups, err = c.Groups.Groups(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, group := range groups {
				if strings.Compare(group.UUID, newGroupDetails.UUID) == 0 {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			// delete the group using RemoveGroup
			removeGroup2, err := c.Groups.RemoveGroup(testConfig.OrgID, newGroupDetails.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(removeGroup2.UUID).To(MatchRegexp(newGroupDetails.UUID))

			// list groups again and prove group is gone
			groups, err = c.Groups.Groups(testConfig.OrgID)
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
