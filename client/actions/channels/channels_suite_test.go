package channels_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChannels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Channels Suite")
}
