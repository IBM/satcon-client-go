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

var _ = Describe("Groups", func() {
	var (
		orgID string
	)

	BeforeEach(func() {
		orgID = "someorg"
	})

	Describe("NewGroupsVariables", func() {
		It("Returns a correctly populated instance of GroupsVariables", func() {
			vars := NewGroupsVariables(orgID)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryGroups))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"orgId",
				"name",
				"created",
			))
		})
	})

	Describe("Groups", func() {
		var (
			token          string
			c              GroupService
			h              *webfakes.FakeHTTPClient
			response       *http.Response
			groupsResponse GroupsResponse
		)

		BeforeEach(func() {
			token = "notreallyatoken"
			groupsResponse = GroupsResponse{
				Data: &GroupsResponseData{
					Groups: GroupList{
						{
							UUID:  "asdf",
							OrgID: orgID,
							Name:  "cluster1",
						},
						{
							UUID:  "qwer",
							OrgID: orgID,
							Name:  "cluster2",
						},
						{
							UUID:  "xzcv",
							OrgID: orgID,
							Name:  "cluster3",
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(groupsResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			h = &webfakes.FakeHTTPClient{}
			Expect(h.DoCallCount()).To(Equal(0))
			h.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", h)
		})

		It("Makes a valid http request", func() {
			_, err := c.Groups(orgID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the list of clusters", func() {
			groups, _ := c.Groups(orgID, token)
			expected := groupsResponse.Data.Groups
			Expect(groups).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Kablooie!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.Groups(orgID, token)
				Expect(err).To(MatchError("Kablooie!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(GroupsResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				groups, err := c.Groups(orgID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(groups).To(BeNil())
			})
		})
	})
})
