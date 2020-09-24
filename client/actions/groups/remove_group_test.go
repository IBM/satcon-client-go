package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/groups"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("RemoveGroup", func() {

	var (
		orgID, uuid, token string
		c                  GroupService
		h                  *webfakes.FakeHTTPClient
		response           *http.Response
	)

	BeforeEach(func() {
		orgID = "someorg"
		uuid = "somelongstringofcharactersgoeshere"
		token = "thisissupposedtobeatoken"

		h = &webfakes.FakeHTTPClient{}
		response = &http.Response{}
		h.DoReturns(response, nil)
	})

	JustBeforeEach(func() {
		c, _ = NewClient("https://foo.bar", h)
		Expect(c).NotTo(BeNil())

		Expect(h.DoCallCount()).To(Equal(0))
	})

	Describe("NewRemoveGroupVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewRemoveGroupVariables(orgID, uuid)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryRemoveGroup))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"uuid":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
			))
		})
	})

	Describe("RemoveGroup", func() {
		var (
			rgResponse RemoveGroupResponse
		)

		BeforeEach(func() {
			rgResponse = RemoveGroupResponse{
				Data: &RemoveGroupResponseData{
					Details: &RemoveGroupResponseDataDetails{
						UUID: "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
					},
				},
			}

			respBodyBytes, err := json.Marshal(rgResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.RemoveGroup(orgID, uuid, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the removeGroup details", func() {
			details, _ := c.RemoveGroup(orgID, uuid, token)
			Expect(details).NotTo(BeNil())

			expected := rgResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoStub = func(r *http.Request) (*http.Response, error) {
					err := errors.New("KABOOM!")
					return response, err
				}
			})

			It("Bubbles up the error", func() {
				_, err := c.RemoveGroup(orgID, uuid, token)
				Expect(err).To(MatchError("KABOOM!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(RemoveGroupResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.RemoveGroup(orgID, uuid, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
