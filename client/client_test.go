package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/IBM/satcon-client-go/client"

	"github.com/IBM/satcon-client-go/client/actions/channels/channelsfakes"
	"github.com/IBM/satcon-client-go/client/actions/clusters/clustersfakes"
	"github.com/IBM/satcon-client-go/client/actions/groups/groupsfakes"
	"github.com/IBM/satcon-client-go/client/actions/resources/resourcesfakes"
	"github.com/IBM/satcon-client-go/client/actions/subscriptions/subscriptionsfakes"
	"github.com/IBM/satcon-client-go/client/actions/versions/versionsfakes"
	"github.com/IBM/satcon-client-go/client/auth"
)

var _ = Describe("Client", func() {
	Describe("New", func() {
		var (
			endpointURL string
			iamClient   *auth.IAMClient
			err         error
		)

		BeforeEach(func() {
			endpointURL = "https://foo.bar"
			iamClient, err = auth.NewIAMClient("some_key")
			Expect(err).NotTo(HaveOccurred())
		})

		It("Creates a new SatCon client", func() {
			s, err := New(endpointURL, nil, iamClient.Client)
			Expect(err).NotTo(HaveOccurred())
			Expect(s.Channels).NotTo(BeNil())
			Expect(s.Clusters).NotTo(BeNil())
			Expect(s.Groups).NotTo(BeNil())
			Expect(s.Resources).NotTo(BeNil())
			Expect(s.Subscriptions).NotTo(BeNil())
			Expect(s.Versions).NotTo(BeNil())
		})

		It("Errors when endpointURL is empty", func() {
			s, err := New("", nil, iamClient.Client)
			Expect(err).To(HaveOccurred())
			Expect(s.Channels).To(BeNil())
			Expect(s.Clusters).To(BeNil())
			Expect(s.Groups).To(BeNil())
			Expect(s.Resources).To(BeNil())
			Expect(s.Subscriptions).To(BeNil())
			Expect(s.Versions).To(BeNil())

		})
	})

	Describe("NewTesting", func() {
		var (
			ch *channelsfakes.FakeChannelService
			cl *clustersfakes.FakeClusterService
			gr *groupsfakes.FakeGroupService
			re *resourcesfakes.FakeResourceService
			su *subscriptionsfakes.FakeSubscriptionService
			ve *versionsfakes.FakeVersionService
		)

		BeforeEach(func() {
			ch = &channelsfakes.FakeChannelService{}
			cl = &clustersfakes.FakeClusterService{}
			gr = &groupsfakes.FakeGroupService{}
			re = &resourcesfakes.FakeResourceService{}
			su = &subscriptionsfakes.FakeSubscriptionService{}
			ve = &versionsfakes.FakeVersionService{}
		})

		It("Creates a client containing only fakes", func() {
			sc := NewTesting("foo", nil)
			Expect(sc.Channels).To(BeAssignableToTypeOf(ch))
			Expect(sc.Clusters).To(BeAssignableToTypeOf(cl))
			Expect(sc.Groups).To(BeAssignableToTypeOf(gr))
			Expect(sc.Resources).To(BeAssignableToTypeOf(re))
			Expect(sc.Subscriptions).To(BeAssignableToTypeOf(su))
			Expect(sc.Versions).To(BeAssignableToTypeOf(ve))
		})
	})
})
