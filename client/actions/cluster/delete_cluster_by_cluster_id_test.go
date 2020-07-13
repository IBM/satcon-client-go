package cluster_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	)

	BeforeEach(func() {
		HTTPClient = &clientfakes.FakeHTTPClient{}
		Expect(HTTPClient.DoCallCount()).To(Equal(0))
		HTTPClient.DoReturns(&http.Response{}, nil)
		endpoint = "http://foo.bar"
		orgID = "someorg"
		clusterID = "somecluster"
		token = "sometoken"
	})

	JustBeforeEach(func() {
		client, _ = NewClient(endpoint, HTTPClient)
		Expect(client).NotTo(BeNil())

		vars = NewDeleteClusterByClusterIDVariables(orgID, clusterID)
		Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
	})

	Describe("Variable Template", func() {

	})

	FIt("Sends a correct request", func() {
		client.DeleteClusterByClusterID(orgID, clusterID, token)
		Expect(HTTPClient.DoCallCount()).To(Equal(1))
		req := HTTPClient.DoArgsForCall(0)

		bodyBytes, err := ioutil.ReadAll(req.Body)
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintln(os.Stderr, string(bodyBytes))
	})
})
