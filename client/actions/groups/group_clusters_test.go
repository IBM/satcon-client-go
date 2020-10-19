package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("GroupClusters", func() {
	var (
		orgID, uuid, token string
		clusters           []string
	)

	BeforeEach(func() {
		orgID = "thisistheidofanorg"
		uuid = "thisist-heuu-idof-agro-upabcdabcdab"
		token = "reallylookslikeatoken"

		clusters = []string{
			"cluster1",
			"cluster2",
			"cluster3",
		}
	})

	Describe("NewGroupClusterVariables", func() {
		It("Returns a correctly configured instance of GroupClusterVariables", func() {
			vars := NewGroupClustersVariables(orgID, uuid, clusters)
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Clusters).To(Equal(clusters))
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryGroupClusters))
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
			vars GroupClustersVariables
		)

		BeforeEach(func() {
			vars = NewGroupClustersVariables(orgID, uuid, clusters)
		})

		It("Processes the variables", func() {
			payload, err := actions.BuildRequestBody(GroupClustersVarTemplate, vars, nil)
			Expect(err).NotTo(HaveOccurred())

			b, _ := ioutil.ReadAll(payload)
			Expect(b).To(MatchRegexp(fmt.Sprintf("\"orgId\":\"%s\"", vars.OrgID)))
			Expect(b).To(MatchRegexp(fmt.Sprintf("\"uuid\":\"%s\"", vars.UUID)))
			Expect(b).To(MatchRegexp(`"clusters":`))
		})
	})

	Describe("GroupClusters", func() {
		var (
			c                     GroupService
			h                     *webfakes.FakeHTTPClient
			response              *http.Response
			groupClustersResponse GroupClustersResponse
		)

		BeforeEach(func() {
			groupClustersResponse = GroupClustersResponse{
				Data: &GroupClustersResponseData{
					Details: &GroupClustersResponseDataDetails{
						Modified: 5,
					},
				},
			}

			respBodyBytes, err := json.Marshal(groupClustersResponse)
			Expect(err).NotTo(HaveOccurred())

			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			h.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", h)
		})

		It("Does not error", func() {
			_, err := c.GroupClusters(orgID, uuid, clusters, token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Returns the response details", func() {
			details, _ := c.GroupClusters(orgID, uuid, clusters, token)
			Expect(details).NotTo(BeNil())
			expected := groupClustersResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.GroupClusters(orgID, uuid, clusters, token)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(GroupClustersResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.GroupClusters(orgID, uuid, clusters, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
