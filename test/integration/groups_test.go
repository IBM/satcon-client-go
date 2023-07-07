package integration_test

import (
	"fmt"
	"strings"

	"github.com/IBM/satcon-client-go/client/types"

	. "github.com/onsi/ginkgo/v2"
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
		c, _ = client.New(testConfig.SatConEndpoint, iamClient.Client)
		Expect(c.Groups).NotTo(BeNil())
	})

	Describe("Group Lifecycle", func() {

		var (
			groupName1  string
			groupName2  string
			clusterName string
		)

		BeforeEach(func() {
			groupName1 = RandStringBytes(8)
			groupName2 = RandStringBytes(8)
			clusterName = RandStringBytes(8)
			fmt.Printf("groupName1 = %s\ngroupName2 = %s\nclusterName = %s\n", groupName1, groupName2, clusterName)
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

		It("Assign group and check if servers are listed", func() {
			// List the groups to check that group does not exist
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

			// create new cluster
			newClusterDetails, err := c.Clusters.RegisterCluster(testConfig.OrgID, types.Registration{Name: clusterName})
			Expect(err).NotTo(HaveOccurred())
			Expect(newClusterDetails).NotTo(BeNil())

			// add new cluster to the group
			_, err = c.Groups.GroupClusters(testConfig.OrgID, newGroupDetails.UUID, []string{newClusterDetails.ClusterID})
			Expect(err).NotTo(HaveOccurred())

			// get group by name and verify that the new cluster is part of it
			group1, err := c.Groups.GroupByName(testConfig.OrgID, groupName1)
			Expect(err).NotTo(HaveOccurred())
			Expect(group1).NotTo(BeNil())
			Expect(group1.Name).To(Equal(groupName1))
			Expect(group1.Clusters).To(HaveLen(1))
			Expect(group1.Clusters[0].ClusterID).To(Equal(newClusterDetails.ClusterID))
			Expect(group1.Clusters[0].Name).To(Equal(clusterName))

			//Remove cluster from group
			_, err = c.Groups.UnGroupClusters(testConfig.OrgID, newGroupDetails.UUID, []string{newClusterDetails.ClusterID})
			Expect(err).NotTo(HaveOccurred())

			group1, err = c.Groups.GroupByName(testConfig.OrgID, groupName1)
			Expect(err).NotTo(HaveOccurred())
			Expect(group1).ToNot(BeNil())
			Expect(group1.Clusters).To(HaveLen(0))

			// delete cluster
			delClusterDetails, err := c.Clusters.DeleteClusterByClusterID(testConfig.OrgID, newClusterDetails.ClusterID)
			Expect(err).NotTo(HaveOccurred())
			Expect(delClusterDetails.DeletedClusterCount).To(Equal(1))

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

			// get group by name and verify that it does not exist
			group1, err = c.Groups.GroupByName(testConfig.OrgID, groupName1)
			Expect(err).To(HaveOccurred())
			Expect(group1).To(BeNil())
		})

	})

})
