package subscriptions_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/IBM/satcon-client-go/client/actions/subscriptions"
	"github.com/IBM/satcon-client-go/client/auth/iam"
)

var _ = Describe("ClusterClient", func() {
	Describe("NewClient", func() {
		var (
			iamClient *iam.IAMClient
			err       error
			h         *http.Client
			endpoint  string
		)

		BeforeEach(func() {
			endpoint = "https://satcon.foo"
			iamClient, err = iam.NewIAMClient("some_key")
			Expect(err).ToNot(HaveOccurred())
		})

		It("Creates a client using the default http client", func() {
			c, err := NewClient(endpoint, nil, iamClient.Client)
			Expect(c).NotTo(BeNil())
			Expect(c.(*Client).HTTPClient).To(Equal(http.DefaultClient))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("When a specific http client is supplied", func() {
			BeforeEach(func() {
				h = &http.Client{
					Timeout: time.Second * 3,
				}
			})

			It("Uses the supplied client", func() {
				c, err := NewClient(endpoint, h, iamClient.Client)
				Expect(c).NotTo(BeNil())
				Expect(c.(*Client).HTTPClient).To(Equal(h))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When the endpoint URL is empty", func() {
			BeforeEach(func() {
				endpoint = ""
			})

			It("Returns nil and an error", func() {
				c, err := NewClient(endpoint, nil, iamClient.Client)
				Expect(c).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
