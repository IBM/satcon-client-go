package auth_test

import (
	"github.com/IBM/satcon-client-go/client/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	var apiKey string
	BeforeEach(func() {

		apiKey = "some_key"
	})

	It("returns a new IAMClient", func() {

		iamClient, err := auth.NewIAMClient(apiKey)
		Expect(iamClient.Client).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

	})

	Describe("errors", func() {

		BeforeEach(func() {

			apiKey = ""
		})

		It("returns an error", func() {

			iamClient, err := auth.NewIAMClient(apiKey)
			Expect(err).To(HaveOccurred())
			Expect(iamClient).To(BeNil())

		})

	})

})