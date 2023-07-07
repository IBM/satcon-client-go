package integration_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client"
	"github.com/IBM/satcon-client-go/client/auth/iam"
	"github.com/IBM/satcon-client-go/client/types"
	. "github.com/IBM/satcon-client-go/test/integration"
)

var _ = Describe("Clusters", func() {
	var (
		c         client.SatCon
		iamClient *iam.Client
	)

	BeforeEach(func() {
		var err error
		iamClient, err = iam.NewIAMClient(testConfig.APIKey, testConfig.IAMEndpoint)
		Expect(err).ToNot(HaveOccurred())
		c, _ = client.New(testConfig.SatConEndpoint, iamClient.Client)
		Expect(c.Clusters).NotTo(BeNil())
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
			clusterList, err := c.Clusters.ClustersByOrgID(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())

			details, err := c.Clusters.RegisterCluster(testConfig.OrgID, types.Registration{Name: clusterName})
			Expect(err).NotTo(HaveOccurred())
			Expect(details).NotTo(BeNil())
			for _, cluster := range clusterList {
				Expect(cluster.ClusterID).NotTo(Equal(details.ClusterID))
			}

			clusterList, err = c.Clusters.ClustersByOrgID(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, cluster := range clusterList {
				if cluster.ClusterID == details.ClusterID && cluster.Name == clusterName {
					found = true
				}
			}
			Expect(found).To(BeTrue())

			cluster, err := c.Clusters.ClusterByName(testConfig.OrgID, clusterName)
			Expect(err).NotTo(HaveOccurred())
			Expect(cluster).NotTo(BeNil())
			Expect(cluster.ClusterID).To(Equal(details.ClusterID))
			Expect(cluster.Name).To(Equal(clusterName))

			delDetails, err := c.Clusters.DeleteClusterByClusterID(testConfig.OrgID, details.ClusterID)
			Expect(err).NotTo(HaveOccurred())
			Expect(delDetails.DeletedClusterCount).To(Equal(1))

			clusterList, err = c.Clusters.ClustersByOrgID(testConfig.OrgID)
			Expect(err).NotTo(HaveOccurred())
			for _, cluster := range clusterList {
				Expect(cluster.ClusterID).NotTo(Equal(details.ClusterID))
			}
		})
	})
})
