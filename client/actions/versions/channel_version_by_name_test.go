package versions_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.ibm.com/coligo/satcon-client/client/actions"
	. "github.ibm.com/coligo/satcon-client/client/actions/versions"
	"github.ibm.com/coligo/satcon-client/client/types"
	"github.ibm.com/coligo/satcon-client/client/web/webfakes"
)

var _ = Describe("ChannelVersion", func() {

	var (
		orgID       string
		channelName string
		versionName string
	)

	BeforeEach(func() {
		orgID = "lemur"
		channelName = "wormhole"
		versionName = "foxtrot"
	})

	Describe("NewChannelVersionByNameVariables", func() {
		It("Returns a correctly populated instance of ChannelVersionByNameVariables", func() {
			vars := NewChannelVersionByNameVariables(orgID, channelName, versionName)
			Expect(vars.Type).To(Equal(actions.QueryTypeQuery))
			Expect(vars.QueryName).To(Equal(QueryChannelVersionByName))
			Expect(vars.OrgID).To(Equal(orgID))
			Expect(vars.ChannelName).To(Equal(channelName))
			Expect(vars.VersionName).To(Equal(versionName))
			Expect(vars.Returns).To(ConsistOf(
				"orgId",
				"uuid",
				"channelId",
				"channelName",
				"name",
				"type",
				"description",
				"content",
				"created",
			))

		})
	})

	Describe("ChannelVersionByName", func() {

		var (
			token                        string
			channelVersionByNameResponse ChannelVersionByNameResponse
			c                            VersionService
			httpClient                   *webfakes.FakeHTTPClient
			response                     *http.Response
		)

		BeforeEach(func() {
			token = "somefaketoken"
			channelVersionByNameResponse = ChannelVersionByNameResponse{
				&ChannelVersionByNameResponseData{
					Details: &types.DeployableVersion{
						OrgID:       "cyborgID",
						UUID:        "youyouID",
						ChannelID:   "chanID",
						ChannelName: "chanName",
						Name:        "somename",
						Type:        "sometype",
						Description: "somedescription",
						Content:     "somecontent",
						Created:     "createdToday",
					},
				},
			}

			respBodyBytes, err := json.Marshal(channelVersionByNameResponse)
			Expect(err).NotTo(HaveOccurred())
			response = &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader(respBodyBytes)),
			}

			httpClient = &webfakes.FakeHTTPClient{}
			Expect(httpClient.DoCallCount()).To(Equal(0))
			httpClient.DoReturns(response, nil)

			c, _ = NewClient("https://foo.bar", httpClient)
		})

		It("Makes a valid http request", func() {
			_, err := c.ChannelVersionByName(orgID, channelName, versionName, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})

		It("Returns the specified channel version", func() {
			channelVersion, _ := c.ChannelVersionByName(orgID, channelName, versionName, token)
			expected := channelVersionByNameResponse.Data.Details
			Expect(channelVersion).To(Equal(expected))
		})

		Context("When query execution errors", func() {
			BeforeEach(func() {
				httpClient.DoReturns(response, errors.New("None whatsoever, Frank."))
			})

			It("Bubbles up the error", func() {
				_, err := c.ChannelVersionByName(orgID, channelName, versionName, token)
				Expect(err).To(MatchError(MatchRegexp("None whatsoever, Frank.")))
			})
		})

		Context("When the response is empty for some reason", func() {
			BeforeEach(func() {
				respBodyBytes, _ := json.Marshal(ChannelVersionByNameResponse{})
				response.Body = ioutil.NopCloser(bytes.NewReader(respBodyBytes))
			})

			It("Returns nil", func() {
				channelVersion, err := c.ChannelVersionByName(orgID, channelName, versionName, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(channelVersion).To(BeNil())
			})
		})

	})

})
