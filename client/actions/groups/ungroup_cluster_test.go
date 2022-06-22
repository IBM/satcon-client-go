package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/satcon-client-go/client/actions"
	"github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
)

var _ = Describe("UngroupCluster", func() {
	var (
		orgID, uuid    string
		clusters       []string
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "fakeOrgId"
		uuid = "fakeuuid"

		clusters = []string{
			"cluster1",
			"cluster2",
			"cluster3",
		}
	})

	Describe("NewUnGroupClusterVariables", func() {
		It("Returns a correctly configured instance of GroupClusterVariables", func() {
			vars := groups.NewUnGroupClustersVariables(orgID, uuid, clusters)
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Clusters).To(Equal(clusters))
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(groups.QueryUnGroupClusters))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":    "String!",
				"uuid":     "String!",
				"clusters": "[String!]!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"modified",
			))
		})
	})
	Describe("GroupClustersVarTemplate", func() {
		var (
			vars groups.UnGroupClustersVariables
		)

		BeforeEach(func() {
			vars = groups.NewUnGroupClustersVariables(orgID, uuid, clusters)
		})

		It("Processes the variables", func() {
			payload, err := actions.BuildRequestBody(groups.UnGroupClustersVarTemplate, vars, nil)
			Expect(err).NotTo(HaveOccurred())

			b, _ := ioutil.ReadAll(payload)
			Expect(b).To(MatchRegexp(fmt.Sprintf("\"orgId\":\"%s\"", vars.OrgID)))
			Expect(b).To(MatchRegexp(fmt.Sprintf("\"uuid\":\"%s\"", vars.UUID)))
			Expect(b).To(MatchRegexp(`"clusters":`))
		})
	})

	Describe("GroupClusters", func() {
		var (
			c                       groups.GroupService
			h                       *webfakes.FakeHTTPClient
			response                *http.Response
			unGroupClustersResponse groups.UnGroupClustersResponse
		)
		BeforeEach(func() {
			unGroupClustersResponse := groups.UnGroupClustersResponse{
				Data: &groups.UnGroupClustersResponseData{
					Details: &groups.UnGroupClustersResponseDataDetails{
						Modified: 5,
					},
				},
			}
			respBodyBytes, err := json.Marshal(unGroupClustersResponse)
			Expect(err).NotTo(HaveOccurred())

			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			h.DoReturns(response, nil)

			c, _ = groups.NewClient("https://foo.bar", h, &fakeAuthClient)
		})
		It("Does not error", func() {
			_, err := c.UnGroupClusters(orgID, uuid, clusters)
			Expect(err).NotTo(HaveOccurred())
		})
		It("Returns the response details", func() {
			details, _ := c.UnGroupClusters(orgID, uuid, clusters)
			Expect(details).NotTo(BeNil())
			expected := unGroupClustersResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})
		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("some error!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.GroupClusters(orgID, uuid, clusters)
				Expect(err).To(MatchError("some error"))
			})
		})
		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(groups.UnGroupClustersResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.UnGroupClusters(orgID, uuid, clusters)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})

	})
})
