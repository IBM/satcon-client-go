package groups_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGroups(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Groups Suite")
}
