package actions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Actions Suite")
}
