package clusters_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/clusters"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
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
				"orgId": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"id",
				"orgId",
				"clusterId",
				"metadata",
			))
		})
	})

	Describe("ClustersByOrgID", func() {
		var (
			token           string
			c               ClusterService
			httpClient      *webfakes.FakeHTTPClient
			response        *http.Response
			clusterResponse ClustersByOrgIDResponse
		)

		BeforeEach(func() {
			token = "notreallyatoken"
			clusterResponse = ClustersByOrgIDResponse{
				Data: &ClustersByOrgIDResponseData{
					Clusters: types.ClusterList{
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

			httpClient = &webfakes.FakeHTTPClient{}
			Expect(httpClient.DoCallCount()).To(Equal(0))
			httpClient.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", httpClient)
		})

		It("Makes a valid http request", func() {
			_, err := c.ClustersByOrgID(orgID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of clusters", func() {
			clusters, _ := c.ClustersByOrgID(orgID, token)
			expected := clusterResponse.Data.Clusters
			Expect(clusters).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.ClustersByOrgID(orgID, token)
				Expect(err).To(MatchError(MatchRegexp("Fart Monkeys!")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ClustersByOrgIDResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				clusters, err := c.ClustersByOrgID(orgID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(clusters).To(BeNil())
			})
		})
	})
})
