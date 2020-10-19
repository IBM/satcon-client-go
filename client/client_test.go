package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/IBM/satcon-client-go/client"
)

var _ = Describe("Client", func() {
	Describe("New", func() {
		var (
			endpointURL string
		)

		BeforeEach(func() {
			endpointURL = "https://foo.bar"
		})

		It("Creates a new SatCon client", func() {
			s, err := New(endpointURL, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.Channels).NotTo(BeNil())
			Expect(s.Clusters).NotTo(BeNil())
			Expect(s.Groups).NotTo(BeNil())
			Expect(s.Resources).NotTo(BeNil())
			Expect(s.Subscriptions).NotTo(BeNil())
			Expect(s.Versions).NotTo(BeNil())
		})

		It("Errors when endpointURL is empty", func() {
			s, err := New("", nil)
			Expect(err).To(HaveOccurred())
			Expect(s.Channels).To(BeNil())
			Expect(s.Clusters).To(BeNil())
			Expect(s.Groups).To(BeNil())
			Expect(s.Resources).To(BeNil())
			Expect(s.Subscriptions).To(BeNil())
			Expect(s.Versions).To(BeNil())

		})
	})
})
