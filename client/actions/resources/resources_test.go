package resources_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Resources", func() {

	var (
		params         types.ResourcesParams
		fakeAuthClient authfakes.FakeAuthClient
		orgID          string
	)

	BeforeEach(func() {
		orgID = "some-cybORG"
		params = types.ResourcesParams{OrgID: orgID}
	})

	Describe("NewResourcesVariables", func() {
		vars := NewResourcesVariables(params)
		Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
		Expect(vars.QueryName).To(Equal(QueryResources))
		Expect(vars.OrgID).To(Equal(orgID))
		Expect(vars.Args).To(Equal(map[string]string{
			"orgId":              "String!",
			"filter":             "String",
			"mongoQuery":         "JSON",
			"fromDate":           "Date",
			"toDate":             "Date",
			"limit":              "Int",
			"kinds":              "[String!]",
			"sort":               "[SortObj!]",
			"subscriptionsLimit": "Int",
		}))
		Expect(vars.Returns).To(ConsistOf(
			"count",
			"resources{id, orgId, clusterId, selfLink, searchableData, created, deleted, subscription{uuid, orgId, name, groups, channel{uuid, orgId, name, created}, version}}",
		))
	})

	Describe("Resources", func() {

		var (
			r                 ResourceService
			h                 *webfakes.FakeHTTPClient
			response          *http.Response
			resourcesResponse ResourcesResponse
		)

		BeforeEach(func() {
			resourcesResponse = ResourcesResponse{
				Data: &ResourcesResponseData{
					ResourceList: &types.ResourceList{
						Count: 1,
						Resources: []types.Resource{
							{
								ID:        "indentify-yourself",
								OrgID:     "what-is-your-organization",
								ClusterID: "c7bc66fe-82e0-4d24-ad61-ac7773830ebc",
								Cluster: types.ClusterInfo{
									ClusterID: "cluster-ID",
									Name:      "cluster-name",
								},
								SelfLink: "/api/v1/namespaces/razeedeploy/pods/watch-keeper-5dd8f8b5b8-k5t5h",
							},
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(resourcesResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			Expect(h.DoCallCount()).To(Equal(0))
			h.DoReturns(response, nil)

			r, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		})

		It("Makes a valid http request", func() {
			_, err := r.Resources(params)
			Expect(err).To(Not(HaveOccurred()))
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns resources for the specified orgID", func() {
			resources, _ := r.Resources(params)
			expected := resourcesResponse.Data.ResourceList
			Expect(resources).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Oh no, Something went wrong!"))
			})

			It("Bubbles up the error", func() {
				_, err := r.Resources(params)
				Expect(err).To(MatchError("Oh no, Something went wrong!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ResourcesResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				groups, err := r.Resources(params)
				Expect(err).NotTo(HaveOccurred())
				Expect(groups).To(BeNil())
			})
		})

	})

})
