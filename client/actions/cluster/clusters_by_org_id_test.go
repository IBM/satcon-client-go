package cluster_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/cluster"
	"github.ibm.com/coligo/satcon-client/client/clientfakes"
)

var _ = Describe("ClustersByOrgId", func() {
	var (
		orgID string
	)

	BeforeEach(func() {
		orgID = "someorg"
	})

	Describe("NewClustersByOrgIDVariables", func() {
		It("Returns a correctly populated instance of ClustersByOrgIDVariables", func() {
			vars := NewClustersByOrgIDVariables(orgID)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryClustersByOrgID))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Args).To(Equal(map[string]string{
				"org_id": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"_id",
				"org_id",
				"cluster_id",
				"metadata",
			))
		})
	})

	Describe("ClustersByOrgID", func() {
		var (
			token           string
			client          ClusterService
			httpClient      *clientfakes.FakeHTTPClient
			response        *http.Response
			clusterResponse ClustersByOrgIDResponse
		)

		BeforeEach(func() {
			token = "notreallyatoken"
			clusterResponse = ClustersByOrgIDResponse{
				Data: &ClustersByOrgIDResponseData{
					Clusters: ClusterList{
						{
							ID:        "asdf",
							OrgID:     orgID,
							ClusterID: "cluster1",
						},
						{
							ID:        "qwer",
							OrgID:     orgID,
							ClusterID: "cluster2",
						},
						{
							ID:        "xzcv",
							OrgID:     orgID,
							ClusterID: "cluster3",
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(clusterResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			httpClient = &clientfakes.FakeHTTPClient{}
			Expect(httpClient.DoCallCount()).To(Equal(0))
			httpClient.DoReturns(response, nil)

			client, _ = NewClient("https://foo.bar", httpClient)
		})

		It("Makes a valid http request", func() {
			_, err := client.ClustersByOrgID(orgID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of clusters", func() {
			clusters, _ := client.ClustersByOrgID(orgID, token)
			expected := clusterResponse.Data.Clusters
			Expect(clusters).To(Equal(expected))
		})
	})
})
