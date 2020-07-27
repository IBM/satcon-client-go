package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.ibm.com/coligo/satcon-client/cli"
	"github.ibm.com/coligo/satcon-client/client/actions/clusters"
	"github.ibm.com/coligo/satcon-client/client/types"
)

var _ = Describe("Print", func() {

	var (
		testInterface interface{}
	)

	BeforeEach(func() {
		testInterface = clusters.RegisterClusterResponse{
			Data: &clusters.RegisterClusterResponseData{
				Details: &clusters.RegisterClusterResponseDataDetails{
					URL:       "https://over.there",
					OrgID:     "fake-orgID",
					OrgKey:    "whatshouldakeylooklike",
					ClusterID: "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
					RegState:  "Faaaabulous!",
					Registration: types.Registration{
						Name: "my_cluster",
					},
				},
			},
		}
	})

	Describe("Print", func() {

		It("Can print a valid interface", func() {
			err := Print(testInterface)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("When printing returns errors", func() {

			BeforeEach(func() {
				// Provide something that cannot be marshalled / unmarshalled
				testInterface = make(chan interface{}, 1)
			})
			It("Returns an error", func() {
				err := Print(testInterface)
				Expect(err).To(HaveOccurred())

			})
		})
	})

})
