package resources_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("ResourceContent", func() {

	var (
		orgID, clusterID, resourceSelfLink string
		resourceContentResponse            string
		r                                  ResourceService
		response                           *http.Response
		h                                  *webfakes.FakeHTTPClient
		fakeAuthClient                     authfakes.FakeAuthClient
		responseStruct                     ResourceContentResponse
	)

	Describe("NewResourceContentVariables", func() {

		BeforeEach(func() {
			orgID = "utah-monolith"
			clusterID = "romania-monolith"
			resourceSelfLink = "/apis/apps/v1/namespace/blahblah/california-monolith"
		})

		It("Returns new resourceContent vars", func() {
			vars := NewResourceContentVariables(orgID, clusterID, resourceSelfLink)
			Expect(vars).NotTo(BeNil())
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.ClusterID).To(Equal(clusterID))
			Expect(vars.ResourceSelfLink).To(Equal(resourceSelfLink))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId":            "String!",
				"clusterId":        "String!",
				"resourceSelfLink": "String!",
			}))
			Expect(vars.Returns).To(Equal([]string{
				"id",
				"histId",
				"content",
				"updated",
			}))
		})
	})

	Describe("ResourceContent", func() {
		BeforeEach(func() {
			resourceContentResponse = `{
  "data": {
    "resourceContent": {
      "id": "somefakeid1234567",
      "content": "{\"kind\":\"Deployment\",\"apiVersion\":\"apps/v1\",\"metadata\":{\"message\":\"Deployment has minimum availability.\"},{\"type\":\"Progressing\",\"status\":\"True\",\"lastUpdateTime\":\"2020-12-15T21:08:07Z\",\"lastTransitionTime\":\"2020-12-11T17:05:17Z\",\"reason\":\"NewReplicaSetAvailable\",\"message\":\"ReplicaSet \\\"some-pod-123456asdfg\\\" has successfully progressed.\"}]}}"
    }
  }
}`

			err := json.Unmarshal([]byte(resourceContentResponse), &responseStruct)
			Expect(err).NotTo(HaveOccurred())

			respBodyBytes, err := json.Marshal(responseStruct)
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
			_, err := r.ResourceContent(orgID, clusterID, resourceSelfLink)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns resource content for the specified cluster", func() {
			content, err := r.ResourceContent(orgID, clusterID, resourceSelfLink)
			Expect(err).NotTo(HaveOccurred())
			expected := responseStruct.Data.ResourceContent
			Expect(content).To(Equal(expected))
		})

		Context("When query execution fails", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Oh no, Something went wrong!"))
			})

			It("Bubbles up the error", func() {
				_, err := r.ResourceContent(orgID, clusterID, resourceSelfLink)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("When response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ResourceContentResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				content, err := r.ResourceContent(orgID, clusterID, resourceSelfLink)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(BeNil())
			})
		})
	})
})
