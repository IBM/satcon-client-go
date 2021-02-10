package apikey_test

import (
	"github.com/IBM/satcon-client-go/client/auth/apikey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("Client", func() {

	var apiKey string
	BeforeEach(func() {
		apiKey = "some_key"
	})

	It("returns a NewRazeeApiKeyClient", func() {
		local, err := apikey.NewClient(apiKey)
		Expect(err).NotTo(HaveOccurred())
		Expect(local).NotTo(BeNil())
	})

	Describe("ApiKeyClient testing", func() {
		It("executes token retrieval", func() {
			client, err := apikey.NewClient(apiKey)
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(BeNil())
			request := http.Request{
				Header: http.Header{},
			}
			request.Header.Add("content-type", "application/json")
			err = client.Authenticate(&request)
			Expect(err).NotTo(HaveOccurred())
			Expect(request.Header.Get(apikey.APIKeyHeader)).To(Equal(apiKey))
		})
	})
})
