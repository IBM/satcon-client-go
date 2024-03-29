package clusters_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// "github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/clusters"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("DeleteClusterByClusterId", func() {
	var (
		// vars                              DeleteClusterByClusterIDVariables
		endpoint, orgID, clusterID string
		c                          ClusterService
		h                          *webfakes.FakeHTTPClient
		response                   *http.Response
		fakeAuthClient             authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		response = &http.Response{}
		h = &webfakes.FakeHTTPClient{}
		Expect(h.DoCallCount()).To(Equal(0))
		h.DoReturns(response, nil)
		endpoint = "http://foo.bar"
		orgID = "someorg"
		clusterID = "somecluster"

	})

	JustBeforeEach(func() {
		c, _ = NewClient(endpoint, h, &fakeAuthClient)
		Expect(c).NotTo(BeNil())

		// vars = NewDeleteClusterByClusterIDVariables(orgID, clusterID)
	})

	// XDescribe("Variable Template", func() {
	// 	payload, err := actions.BuildRequestBody(DeleteClusterByClusterIDVarTemplate, vars, nil)
	// 	Expect(err).NotTo(HaveOccurred())
	//
	// 	pbytes, _ := ioutil.ReadAll(payload)
	// 	Expect(len(pbytes)).To(BeNumerically(">", 0))
	// })

	Describe("DeleteClusterByClusterID", func() {
		var (
			delResponse DeleteClustersResponse
		)

		BeforeEach(func() {
			delResponse = DeleteClustersResponse{
				Data: &DeleteClustersResponseData{
					Details: &DeleteClustersResponseDataDetails{
						DeletedClusterCount:  1,
						DeletedResourceCount: 0,
					},
				},
			}

			respBodyBytes, err := json.Marshal(delResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends a correct request", func() {
			_, err := c.DeleteClusterByClusterID(orgID, clusterID)
			Expect(h.DoCallCount()).To(Equal(1))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Returns the response details", func() {
			details, _ := c.DeleteClusterByClusterID(orgID, clusterID)

			expected := delResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(nil, errors.New("Holy smokes!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.DeleteClusterByClusterID(orgID, clusterID)
				Expect(err).To(MatchError("Holy smokes!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(DeleteClustersResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.DeleteClusterByClusterID(orgID, clusterID)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
