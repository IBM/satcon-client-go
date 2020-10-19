package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/groups"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Adding a Group", func() {
	var (
		orgID, name, token string
		c                  GroupService
		h                  *webfakes.FakeHTTPClient
		response           *http.Response
	)

	BeforeEach(func() {
		orgID = "someorg"
		name = "somegroup"
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

	Describe("NewAddGroupVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewAddGroupVariables(orgID, name)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryAddGroup))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Name).To(Equal(name))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"name":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
			))
		})
	})

	Describe("AddGroup", func() {
		var (
			agResponse AddGroupResponse
		)

		BeforeEach(func() {
			agResponse = AddGroupResponse{
				Data: &AddGroupResponseData{
					Details: &AddGroupResponseDataDetails{
						UUID: "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
					},
				},
			}

			respBodyBytes, err := json.Marshal(agResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.AddGroup(orgID, name, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the add group details", func() {
			details, _ := c.AddGroup(orgID, name, token)
			Expect(details).NotTo(BeNil())

			expected := agResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Fart Monkeys!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.AddGroup(orgID, name, token)
				Expect(err).To(MatchError("Fart Monkeys!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(AddGroupResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.AddGroup(orgID, name, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
