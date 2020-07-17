package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.ibm.com/coligo/satcon-client/client"
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
		})
	})
})
