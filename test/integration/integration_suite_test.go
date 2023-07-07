package integration_test

import (
	"testing"

	. "github.com/IBM/satcon-client-go/test/integration"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
