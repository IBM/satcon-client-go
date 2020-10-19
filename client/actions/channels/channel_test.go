package channels_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/IBM/satcon-client-go/client/actions"
	. "github.com/IBM/satcon-client-go/client/actions/channels"
	"github.com/IBM/satcon-client-go/client/types"
	"github.com/IBM/satcon-client-go/client/web/webfakes"
)

var _ = Describe("Channels", func() {
	var (
		orgID string
		uuid  string
	)

	BeforeEach(func() {
		orgID = "somecybORG"
		uuid = "atlas-v"
	})

	Describe("NewChannelVariables", func() {
		It("Returns a correctly populated instance of ChannelVariables", func() {
			vars := NewChannelVariables(orgID, uuid)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryChannel))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.UUID).To(Equal(uuid))
			Expect(vars.Args).To(Equal(map[string]string{
				"orgId": "String!",
				"uuid":  "String!",
			}))
			Expect(vars.Returns).To(ConsistOf(
				"uuid",
				"orgId",
				"name",
				"created",
				"versions{uuid, name, location}",
				"subscriptions{uuid, orgId, name, groups}",
			))
		})
	})

	Describe("Channel", func() {
		var (
			token           string
			c               ChannelService
			h               *webfakes.FakeHTTPClient
			response        *http.Response
			channelResponse ChannelResponse
		)

		BeforeEach(func() {
			token = "notreallyatoken"
			channelResponse = ChannelResponse{
				Data: &ChannelResponseData{
					Channel: &types.Channel{
						UUID:    "asdf",
						OrgID:   orgID,
						Name:    "channel1",
						Created: "Now",
						Versions: types.ChannelVersionList{
							{
								UUID:        "version1uuid",
								Name:        "version1",
								Description: "version1 description",
								Location:    "kitchen",
								Created:     "yesterday",
							},
							{
								UUID:        "version2uuid",
								Name:        "version2",
								Description: "version2 description",
								Location:    "library",
								Created:     "day-before-yesterday",
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
			_, err := c.Channel(orgID, uuid, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(h.DoCallCount()).To(Equal(1))
		})

		It("Returns the specified channel", func() {
			groups, _ := c.Channel(orgID, uuid, token)
			expected := channelResponse.Data.Channel
			Expect(groups).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				h.DoReturns(response, errors.New("Kablooie!"))
			})

			It("Bubbles up the error", func() {
				_, err := c.Channel(orgID, uuid, token)
				Expect(err).To(MatchError("Kablooie!"))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ChannelsResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				groups, err := c.Channel(orgID, uuid, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(groups).To(BeNil())
			})
		})
	})
})
