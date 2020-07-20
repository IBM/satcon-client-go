package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client"
	"github.ibm.com/coligo/satcon-client/client/types"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

var _ = Describe("Clusters", func() {
	var (
		token string
		c     client.SatCon
	)

	BeforeEach(func() {
		var err error
		c, _ = client.New(testConfig.SatConEndpoint, nil)
		Expect(c.Clusters).NotTo(BeNil())

		token, err = GetToken(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).NotTo(BeZero())
	})

	Describe("Cluster Lifecycle", func() {
		var (
			clusterName string
		)

		BeforeEach(func() {
			clusterName = RandStringBytes(8)
			fmt.Println("Using cluster name: ", clusterName)
		})

		It("Lists the clusters, creates our new cluster, lists again and finds it, deletes it, and finally lists to see that it's gone", func() {
			clusterList, err := c.Clusters.ClustersByOrgID(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			// TODO: Iterate through the list and check that none of them is our cluster
			// for _, cluster := range clusterList {
			// 	Expect(cluster.)
			// }

			details, err := c.Clusters.RegisterCluster(testConfig.OrgID, types.Registration{Name: clusterName}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(details).NotTo(BeNil())
			for _, cluster := range clusterList {
				Expect(cluster.ClusterID).NotTo(Equal(details.ClusterID))
			}

			clusterList, err = c.Clusters.ClustersByOrgID(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, cluster := range clusterList {
				if cluster.ClusterID == details.ClusterID {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			delDetails, err := c.Clusters.DeleteClusterByClusterID(testConfig.OrgID, details.ClusterID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(delDetails.DeletedClusterCount).To(Equal(1))

			clusterList, err = c.Clusters.ClustersByOrgID(testConfig.OrgID, token)
			Expect(err).NotTo(HaveOccurred())
			for _, cluster := range clusterList {
				Expect(cluster.ClusterID).NotTo(Equal(details.ClusterID))
			}
		})
	})
})
