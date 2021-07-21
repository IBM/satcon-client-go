package resources_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResourcesByCluster", func() {

	var (
		orgID, clusterID, filter string
		limit                    int
	)
	BeforeEach(func() {
		orgID = "mercury-org"
		clusterID = "gemini-cluster"
		filter = "Pod"
		limit = 5
	})

	Describe("NewResourcesByClusterVariables", func() {
		It("Returns a correctly populated instance of ResourcesByClusterVariables", func() {
			vars := NewResourcesByClusterVariables(orgID, clusterID, filter, limit)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryResourcesByCluster))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.ClusterID).To(Equal(clusterID))
			Expect(vars.Filter).To(Equal(filter))
			Expect(vars.Limit).To(Equal(limit))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":     "String!",
				"clusterId": "String!",
				"filter":    "String",
				"limit":     "Int",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"count",
				"resources{id, orgId, clusterId, selfLink, searchableData, subscription{uuid, orgId, name, groups, channel{uuid, orgId, name, created}, version}}",
			))
		})
	})

	Describe("ResourcesByCluster", func() {

		var (
			r                 ResourceService
			h                 *webfakes.FakeHTTPClient
			response          *http.Response
			resourcesResponse ResourcesByClusterResponse
			fakeAuthClient    authfakes.FakeAuthClient
		)

		BeforeEach(func() {
			resourcesResponse = ResourcesByClusterResponse{
				Data: &ResourcesByClusterResponseData{
					ResourceList: &types.ResourceList{
						Count: limit,
						Resources: []types.Resource{
							{
								ID:        "indentify-yourself",
								OrgID:     "what-is-your-organization",
								ClusterID: "c7bc66fe-82e0-4d24-ad61-ac7773830ebc",
								Cluster: types.ClusterInfo{
									ClusterID: "cluster-ID",
									Name:      "cluster-name",
								},
								SelfLink:     "/api/v1/namespaces/razeedeploy/pods/watch-keeper-5dd8f8b5b8-k5t5h",
								Hash:         "bb5d00c8173bbb63704342f885385cfb1f5c3c25",
								Data:         "{\"kind\":\"Pod\",\"apiVersion\":\"v1\",\"metadata\":{\"name\":\"watch-keeper-abcdefg-xxxx\",\"generateName\":\"watch-keeper-abcdefg-\",\"namespace\":\"razeedeploy\",\"selfLink\":\"/api/v1/namespaces/razeedeploy/pods/watch-keeper-abcdefg-xxxx\",\"uid\":\"whatever-uid\",\"resourceVersion\":\"0000001\",\"creationTimestamp\":\"a-few-days-ago\",\"labels\":{\"app\":\"watch-keeper\",\"pod-template-hash\":\"1234hash\",\"razee/watch-resource\":\"lite\"},\"annotations\":{\"kubernetes.io/psp\":\"ibm-privileged-psp\"},\"ownerReferences\":[{\"apiVersion\":\"apps/v1\",\"kind\":\"ReplicaSet\",\"name\":\"watch-keeper-abcdefg\",\"uid\":\"some-other-uid\",\"controller\":true,\"blockOwnerDeletion\":true}]},\"status\":{\"phase\":\"Running\",\"conditions\":[{\"type\":\"Initialized\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"seconds-ago\"},{\"type\":\"Ready\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"not-long-ago\"},{\"type\":\"ContainersReady\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"yesterday\"},{\"type\":\"PodScheduled\",\"status\":\"True\",\"lastProbeTime\":null,\"lastTransitionTime\":\"who-know?\"}],\"hostIP\":\"some-host\",\"podIP\":\"some-pod-IP\",\"podIPs\":[{\"ip\":\"some-pod-IP\"}],\"startTime\":\"2020-06-30T17:52:02Z\",\"containerStatuses\":[{\"name\":\"watch-keeper\",\"state\":{\"running\":{\"startedAt\":\"beginning-of-time\"}},\"lastState\":{},\"ready\":true,\"restartCount\":0,\"image\":\"quay.io/razee/watch-keeper:tag\",\"imageID\":\"quay.io/razee/watch-keeper\",\"containerID\":\"containerd://1234567890\",\"started\":true}],\"qosClass\":\"Burstable\"}}",
								Deleted:      false,
								Created:      "some string",
								Updated:      "a little while ago",
								LastModified: "now",
								SearchableData: types.SearchableData{
									Kind:                 "Pod",
									Name:                 "watch-keeper-abcdefg-xxxx",
									Namespace:            "razeedeploy",
									APIVersion:           "v1",
									SearchableExpression: "Pod:watch-keeper-abcdefg-xxxx:razeedeploy:v1:ibm-privileged-psp:quay.io/razee/watch-keeper",
								},
								SearchableDataHash: "s34rchableH@$H",
								Subscription: types.ChannelSubscription{
									UUID:        "subscription-uuid",
									OrgID:       "subscription-org-id",
									Name:        "subscription-name",
									Groups:      []string{"subscription-group-1", "subscription-group-2"},
									ChannelUUID: "channel-uuid",
									ChannelName: "channel-name",
									Channel: &types.Channel{
										UUID:    "another-uuid",
										OrgID:   "the-orgID-again",
										Name:    "channel-name-once-again",
										Created: "some-time-today",
										Versions: []types.ChannelVersion{
											{
												UUID:        "version1uuid",
												Name:        "version1",
												Description: "version1 description",
												Location:    "kitchen",
												Created:     "yesterday",
											},
											{
												UUID:        "version2uuid",
												Name:        "version2",
												Description: "version2 description",
												Location:    "library",
												Created:     "day-before-yesterday",
											},
										},
									},
								},
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
			_, err := r.ResourcesByCluster(orgID, clusterID, filter, limit, nil)
			Expect(err).To(Not(HaveOccurred()))
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns resources for the specified cluster", func() {
			resources, _ := r.ResourcesByCluster(orgID, clusterID, filter, limit, nil)
			expected := resourcesResponse.Data.ResourceList
			Expect(resources).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Oh no, Something went wrong!"))
			})

			It("Bubbles up the error", func() {
				_, err := r.ResourcesByCluster(orgID, clusterID, filter, limit, nil)
				Expect(err).To(MatchError("Oh no, Something went wrong!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ResourcesByClusterResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				groups, err := r.ResourcesByCluster(orgID, clusterID, filter, limit, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(groups).To(BeNil())
			})
		})
	})

})
