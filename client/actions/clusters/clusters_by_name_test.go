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
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("ClusterByName", func() {
	var (
		orgID          string
		clusterName    string
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "someorg"
		clusterName = "somename"
	})

	Describe("NewClustersByNameVariables", func() {
		It("Returns a correctly populated instance of ClustersByNameVariables", func() {
			vars := NewClusterByNameVariables(orgID, clusterName)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryClusterByName))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.ClusterName).To(Equal(clusterName))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":       "String!",
				"clusterName": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"id",
				"orgId",
				"clusterId",
				"name",
				"metadata",
			))
		})
	})

	Describe("ClustersByName", func() {
		var (
			c               ClusterService
			httpClient      *webfakes.FakeHTTPClient
			response        *http.Response
			clusterResponse ClusterByNameResponse
		)

		BeforeEach(func() {
			clusterResponse = ClusterByNameResponse{
				Data: &ClusterByNameResponseData{
					&types.Cluster{
						ID:        "asdf",
						OrgID:     orgID,
						ClusterID: "cluster1",
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

			c, _ = NewClient("https://foo.bar", httpClient, &fakeAuthClient)
		})

		It("Makes a valid http request", func() {
			_, err := c.ClusterByName(orgID, clusterName)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of clusters", func() {
			clusters, _ := c.ClusterByName(orgID, clusterName)
			expected := clusterResponse.Data.Cluster
			Expect(clusters).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("None whatsoever, Frank."))
			})

			It("Bubbles up the error", func() {
				_, err := c.ClusterByName(orgID, clusterName)
				Expect(err).To(MatchError(MatchRegexp("None whatsoever, Frank.")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ClustersByOrgIDResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				clusters, err := c.ClusterByName(orgID, clusterName)
				Expect(err).NotTo(HaveOccurred())
				Expect(clusters).To(BeNil())
			})
		})
	})
})
