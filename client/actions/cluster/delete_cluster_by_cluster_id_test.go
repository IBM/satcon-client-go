package cluster_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/cluster"
	"github.ibm.com/coligo/satcon-client/client/clientfakes"
)

var _ = Describe("DeleteClusterByClusterId", func() {
	var (
		vars                              DeleteClusterByClusterIDVariables
		endpoint, orgID, clusterID, token string
		client                            ClusterService
		HTTPClient                        *clientfakes.FakeHTTPClient
		response                          *http.Response
	)

	BeforeEach(func() {
		response = &http.Response{}
		HTTPClient = &clientfakes.FakeHTTPClient{}
		Expect(HTTPClient.DoCallCount()).To(Equal(0))
		HTTPClient.DoReturns(response, nil)
		endpoint = "http://foo.bar"
		orgID = "someorg"
		clusterID = "somecluster"
		token = "sometoken"
	})

	JustBeforeEach(func() {
		client, _ = NewClient(endpoint, HTTPClient)
		Expect(client).NotTo(BeNil())

		vars = NewDeleteClusterByClusterIDVariables(orgID, clusterID)
	})

	Describe("Variable Template", func() {
		payload, err := actions.BuildRequestBody(DeleteClusterByClusterIDVarTemplate, vars, nil)
		Expect(err).NotTo(HaveOccurred())

		pbytes, _ := ioutil.ReadAll(payload)
		Expect(len(pbytes)).To(BeNumerically(">", 0))
	})

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
			_, err := client.DeleteClusterByClusterID(orgID, clusterID, token)
			Expect(HTTPClient.DoCallCount()).To(Equal(1))
			Expect(err).NotTo(HaveOccurred())
		})

		It("Returns the response details", func() {
			details, _ := client.DeleteClusterByClusterID(orgID, clusterID, token)

			expected := delResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When the http call errors", func() {
			BeforeEach(func() {
				HTTPClient.DoReturns(nil, errors.New("Holy smokes!"))
			})

			It("Bubbles up the error", func() {
				_, err := client.DeleteClusterByClusterID(orgID, clusterID, token)
				Expect(err).To(MatchError("Holy smokes!"))
			})
		})

		Context("When the response cannot be unmarshalled", func() {
			BeforeEach(func() {
				response.Body = ioutil.NopCloser(bytes.NewBufferString("thisis}not{validjson"))
			})

			It("Bubbles up the error", func() {
				_, err := client.DeleteClusterByClusterID(orgID, clusterID, token)
				Expect(err.Error()).To(HavePrefix("Unable to unmarshal response"))
			})
		})
	})
})
