package subscriptions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSubscriptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Subscriptions Suite")
}
