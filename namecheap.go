// Package namecheap implements a client for the Namecheap API.
//
// In order to use this package you will need a Namecheap account and your API Token.
package namecheap

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.namecheap.com/xml.response"
)

type NamecheapClient struct {
	ApiUser    string
	ApiToken   string
	UserName   string
	HttpClient *http.Client

	// Base URL for API requests.
	// Defaults to the public Namecheap API, but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string
}

type ApiRequest struct {
	method  string
	command string
	params  url.Values
}

type ApiResponse struct {
	Status     string     `xml:"Status,attr"`
	Command    string     `xml:"RequestedCommand"'`
	Domains    []Domain   `xml:"CommandResponse>DomainGetListResult>Domain"`
	DomainInfo DomainInfo `xml:"CommandResponse>DomainGetInfoResult"`
	Errors     []ApiError `xml:"Errors>Error"`
}

type ApiError struct {
	Number int `xml:"Number,attr"`
	Message string `xml:",innerxml"`
}

func NewClient(apiUser, apiToken, userName string) *NamecheapClient {
	return &NamecheapClient{ApiUser: apiUser, ApiToken: apiToken, UserName: userName, HttpClient: &http.Client{}, BaseURL: defaultBaseURL}
}

func (client *NamecheapClient) get(request ApiRequest, resp interface{}) error {
	request.method = "GET"
	body, _, err := client.sendRequest(request, nil)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal([]byte(body), &resp); err != nil {
		return err
	}

	return nil
}

func (client *NamecheapClient) makeRequest(request ApiRequest, body io.Reader) (*http.Request, error) {
	url, err := url.Parse(client.BaseURL)
	if err != nil {
		return nil, err
	}
	p := request.params
	p.Set("ApiUser", client.ApiUser)
	p.Set("ApiKey", client.ApiToken)
	p.Set("UserName", client.UserName)
	p.Set("ClientIp", "127.0.0.1")
	p.Set("Command", request.command)
	url.RawQuery = p.Encode()

	urlString := fmt.Sprintf("%s?%s", client.BaseURL, url.RawQuery)
	req, err := http.NewRequest(request.method, urlString, body)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (client *NamecheapClient) sendRequest(request ApiRequest, body io.Reader) (string, int, error) {
	req, err := client.makeRequest(request, body)
	if err != nil {
		return "", 0, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(responseBytes), resp.StatusCode, nil
}
