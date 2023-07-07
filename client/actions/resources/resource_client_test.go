package resources_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/IBM/satcon-client-go/client/actions/resources"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
)

var _ = Describe("ResourceClient", func() {
	Describe("NewClient", func() {
		var (
			h              *http.Client
			endpoint       string
			fakeAuthClient authfakes.FakeAuthClient
		)

		BeforeEach(func() {
			endpoint = "https://satcon.foo"
		})

		It("Creates a client using the default http client", func() {
			c, err := NewClient(endpoint, nil, &fakeAuthClient)
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
				c, err := NewClient(endpoint, h, &fakeAuthClient)
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
				c, err := NewClient(endpoint, nil, &fakeAuthClient)
				Expect(c).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
