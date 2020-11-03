package versions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/versions"
	"github.com/IBM/satcon-client-go/client/auth/authfakes"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("RemoveChannelVersion", func() {

	var (
		orgID, uuid    string
		c              VersionService
		h              *webfakes.FakeHTTPClient
		response       *http.Response
		fakeAuthClient authfakes.FakeAuthClient
	)

	BeforeEach(func() {
		orgID = "someorg"
		uuid = "somechannel"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		c, _ = NewClient("https://foo.bar", h, &fakeAuthClient)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewRemoveChannelVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewRemoveChannelVersionVariables(orgID, uuid)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryRemoveChannelVersion))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"uuid":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"success",
			))
		})
	})

	Describe("RemoveChannelVersion", func() {
		var (
			rcResponse RemoveChannelVersionResponse
		)

		BeforeEach(func() {
			rcResponse = RemoveChannelVersionResponse{
				Data: &RemoveChannelVersionResponseData{
					Details: &RemoveChannelVersionResponseDataDetails{
						UUID:    "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
						Success: true,
					},
				},
			}

			respBodyBytes, err := json.Marshal(rcResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.RemoveChannelVersion(orgID, uuid)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the add group details", func() {
			details, _ := c.RemoveChannelVersion(orgID, uuid)
			Expect(details).NotTo(BeNil())

			expected := rcResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("YIKES!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.RemoveChannelVersion(orgID, uuid)
				Expect(err).To(MatchError("YIKES!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(RemoveChannelVersionResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.RemoveChannelVersion(orgID, uuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})

})
