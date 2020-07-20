package integration

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
)

type TestConfig struct {
	APIKey         string `json:"apiKey,omitempty"`
	Token          string `json:"token,omitempty"`
	IAMEndpoint    string `json:"iamEndpoint,omitempty"`
	SatConEndpoint string `json:"satconEndpoint,omitempty"`
	OrgID          string `json:"orgId,omitempty"`
}

func LoadConfig(configFile string) *TestConfig {
	var cfg TestConfig
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		return nil
	}

	return &cfg
}

type IAMTokenResponse struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func GetToken(apiKey, iamEndpoint string) (string, error) {
	var tokenResponse IAMTokenResponse

	payload := strings.NewReader("grant_type=urn%3Aibm%3Aparams%3Aoauth%3Agrant-type%3Aapikey&apikey=" + apiKey)

	req, _ := http.NewRequest("POST", iamEndpoint, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}
	return tokenResponse.Token, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
