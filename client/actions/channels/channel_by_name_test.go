package channels_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/channels"
	"github.ibm.com/coligo/satcon-client/client/types"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("ChannelByName", func() {

	var (
		orgID       string
		channelName string
	)

	BeforeEach(func() {
		orgID = "lemur"
		channelName = "gemini"
	})

	Describe("NewChannelVersionByNameVariables", func() {
		It("Returns a correctly populated instance of ChannelVersionByNameVariables", func() {
			vars := NewChannelByNameVariables(orgID, channelName)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryChannelByName))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.Name).To(Equal(channelName))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"orgId",
				"name",
				"created",
				"versions{uuid, name, description, location, created}",
				"subscriptions{uuid, orgId, name, groups, channelUuid, channelName, version, versionUuid, created, updated}",
			))
		})
	})

	Describe("ChannelByName", func() {
		var (
			token           string
			c               ChannelService
			h               *webfakes.FakeHTTPClient
			response        *http.Response
			channelResponse ChannelByNameResponse
		)

		BeforeEach(func() {
			token = "notreallyatoken"
			channelResponse = ChannelByNameResponse{
				Data: &ChannelByNameResponseData{
					Details: &types.Channel{
						UUID:    "asdf",
						OrgID:   orgID,
						Name:    "channel1",
						Created: "whenever",
						Versions: types.ChannelVersionList{
							{
								Name:        "version1",
								UUID:        "vesion-uuid-1",
								Location:    "location1",
								Description: "desc1",
								Created:     "then1",
							},
							{
								Name:        "version2",
								UUID:        "vesion-uuid-2",
								Location:    "location2",
								Description: "desc2",
								Created:     "then2",
							},
						},
						Subscriptions: []types.BasicChannelSubscription{
							{
								Groups:      []string{"group-1", "group-2"},
								Name:        "this-subscription",
								UUID:        "subscription-uuid-1",
								OrgID:       "userOrg",
								ChannelUUID: "subscription-channel-uuid-1",
								ChannelName: "subscription-channel-name",
								Version:     "someversion",
								Created:     "then",
								Updated:     "now",
							},
						},
					},
				},
			}

			respBodyBytes, err := json.Marshal(channelResponse)
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
			_, err := c.ChannelByName(orgID, channelName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the specified channel", func() {
			channel, _ := c.ChannelByName(orgID, channelName, token)
			expected := channelResponse.Data.Details
			Expect(channel).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Kablooie!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.ChannelByName(orgID, channelName, token)
				Expect(err).To(MatchError("Kablooie!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ChannelByNameResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				channel, err := c.ChannelByName(orgID, channelName, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(channel).To(BeNil())
			})
		})
	})

})
