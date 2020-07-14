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
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("Registering a Cluster", func() {
	var (
		orgID, token string
		reg          Registration
		c            ClusterService
		HTTPClient   *webfakes.FakeHTTPClient
		response     *http.Response
	)

	BeforeEach(func() {
		orgID = "someorg"
		reg = Registration{
			Name: "my_cluster",
		}
		token = "thisissupposedtobeatoken"

		HTTPClient = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		HTTPClient.DoReturns(response, nil)
	})

	JustBeforeEach(func() {

		c, _ = NewClient("https://foo.bar", HTTPClient)
		Expect(c).NotTo(BeNil())

		Expect(HTTPClient.DoCallCount()).To(Equal(0))
	})

	Describe("NewRegisterClusterVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewRegisterClusterVariables(orgID, reg)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryRegisterCluster))
			Expect(vars.OrgID).To(Equal(orgID))
			regBytes, _ := json.Marshal(reg)
			Expect(vars.Registration).To(Equal(regBytes))
			Expect(vars.Args).To(Equal(map[string]string{
				"org_id":       "String!",
				"registration": "JSON!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"url",
				"org_id",
				"orgKey",
				"clusterId",
				"regState",
				"registration",
			))
		})
	})

	Describe("RegisterCluster", func() {
		var (
			regResponse RegisterClusterResponse
		)

		BeforeEach(func() {
			regResponse = RegisterClusterResponse{
				Data: &RegisterClusterResponseData{
					Details: &RegisterClusterResponseDataDetails{
						URL:          "https://over.there",
						OrgID:        orgID,
						OrgKey:       "whatshouldakeylooklike",
						ClusterID:    "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
						RegState:     "Faaaabulous!",
						Registration: reg,
					},
				},
			}

			respBodyBytes, err := json.Marshal(regResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.RegisterCluster(orgID, reg, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(HTTPClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the cluster registration details", func() {
			details, _ := c.RegisterCluster(orgID, reg, token)
			Expect(details).NotTo(BeNil())

			expected := regResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When the http call errors", func() {
			BeforeEach(func() {
				HTTPClient.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.RegisterCluster(orgID, reg, token)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})
	})
})
