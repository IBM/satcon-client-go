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

var _ = Describe("RemoveGroupByName", func() {

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

	Describe("NewRemoveGroupByNameVariables", func() {
		It("Returns a correctly configured set of variables", func() {
			vars := NewRemoveGroupByNameVariables(orgID, name)
			Expect(vars.Type).To(Equal(actions.QueryTypeMutation))
			Expect(vars.QueryName).To(Equal(QueryRemoveGroupByName))
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

	Describe("RemoveGroupByName", func() {
		var (
			rgResponse RemoveGroupByNameResponse
		)

		BeforeEach(func() {
			rgResponse = RemoveGroupByNameResponse{
				Data: &RemoveGroupByNameResponseData{
					Details: &RemoveGroupByNameResponseDataDetails{
						UUID: "abcdabcd-abcd-abcd-abcd-abcdabcdabcd",
					},
				},
			}

			respBodyBytes, err := json.Marshal(rgResponse)
			Expect(err).NotTo(HaveOccurred())

			response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
		})

		It("Sends the http request", func() {
			_, err := c.RemoveGroupByName(orgID, name, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the removeGroupByName details", func() {
			details, _ := c.RemoveGroupByName(orgID, name, token)
			Expect(details).NotTo(BeNil())

			expected := rgResponse.Data.Details
			Expect(*details).To(Equal(*expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Excuse ME!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.RemoveGroupByName(orgID, name, token)
				Expect(err).To(MatchError("Excuse ME!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(RemoveGroupByNameResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				details, err := c.RemoveGroupByName(orgID, name, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
			})
		})
	})
})
