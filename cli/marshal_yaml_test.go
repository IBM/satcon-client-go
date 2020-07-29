package cli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.ibm.com/coligo/satcon-client/cli"
)

var _ = Describe("MarshalYaml", func() {

	var (
		fileName      string
		expectedBytes string
	)

	BeforeEach(func() {
		fileName = "clifakes/fake_pod.yml"
		expectedBytes = "YXBpVmVyc2lvbjogdjEKa2luZDogUG9kCm1ldGFkYXRhOgogIG5hbWU6IGNsaV90ZXN0CnNwZWM6CiAgY29udGFpbmVyczoKICAtIG5hbWU6IGNsaV90ZXN0CiAgICBpbWFnZTogaHR0cGQ6YWxwaW5lCg=="
	})

	Describe("MarshalYAMLFromFile", func() {

		Context("When given a properly formatted .yaml", func() {

			It("Successfully returns a []byte", func() {
				bytes, err := MarshalYAMLFromFile(fileName)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(bytes)).To(MatchRegexp(string(expectedBytes)))
			})

		})

		Context("When file does not exist", func() {

			BeforeEach(func() {
				fileName = "nonexistent.yml"
			})

			It("Returns an error", func() {
				bytes, err := MarshalYAMLFromFile(fileName)
				Expect(err).To(HaveOccurred())
				Expect(bytes).To(BeNil())
			})

		})

	})

})
