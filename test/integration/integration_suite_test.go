package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.ibm.com/coligo/satcon-client/test/integration"
)

const (
	integrationConfigFile = "integration.json"
)

var testConfig *TestConfig

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	testConfig = LoadConfig(integrationConfigFile)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	Expect(testConfig).NotTo(BeNil())
})
