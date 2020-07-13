package cluster_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/cluster"
	"github.ibm.com/coligo/satcon-client/client/clientfakes"
)

var _ = Describe("Registering a Cluster", func() {
	var (
		orgID, token string
		reg          Registration
		c            ClusterService
		HTTPClient   *clientfakes.FakeHTTPClient
		response     *http.Response
	)

	BeforeEach(func() {
		orgID = "someorg"
		reg = Registration{
			Name: "my_cluster",
		}
		token = "thisissupposedtobeatoken"

		HTTPClient = &clientfakes.FakeHTTPClient{}
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

	It("Sends a correct request", func() {
		_, err := c.RegisterCluster(orgID, reg, token)
		Expect(err).NotTo(HaveOccurred())
		Expect(HTTPClient.DoCallCount()).To(Equal(1))

		req := HTTPClient.DoArgsForCall(0)

		Expect(req.Header).To(HaveKeyWithValue(
			MatchRegexp(`[cC]ontent-[tT]ype`),
			ContainElement("application/json"),
		))
		Expect(req.Header).To(HaveKeyWithValue(
			MatchRegexp("[aA]uthorization"),
			ContainElement(fmt.Sprintf("Bearer %s", token)),
		))
	})
})
