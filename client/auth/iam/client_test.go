package iam_test

import (
	"github.com/IBM/satcon-client-go/client/auth/iam"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	var apiKey string
	BeforeEach(func() {
		apiKey = "some_key"
	})

	It("returns a new IAMClient", func() {

		iamClient, err := iam.NewIAMClient(apiKey)
		Expect(iamClient.Client).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

	})

	Describe("errors", func() {

		BeforeEach(func() {

			apiKey = ""
		})

		It("returns an error", func() {

			iamClient, err := iam.NewIAMClient(apiKey)
			Expect(err).To(HaveOccurred())
			Expect(iamClient).To(BeNil())

		})

	})
})
