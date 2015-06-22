// Package namecheap implements a client for the Namecheap API.
//
// In order to use this package you will need a Namecheap account and your API Token.
package namecheap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://api.namecheap.com/xml.response"

// Client represents a client used to make calls to the Namecheap API.
type Client struct {
	ApiUser    string
	ApiToken   string
	UserName   string
	HttpClient *http.Client

	// Base URL for API requests.
	// Defaults to the public Namecheap API,
	// but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string
}

type ApiRequest struct {
	method  string
	command string
	params  url.Values
}

type ApiResponse struct {
	Status            string                   `xml:"Status,attr"`
	Command           string                   `xml:"RequestedCommand"`
	Domains           []DomainGetListResult    `xml:"CommandResponse>DomainGetListResult>Domain"`
	DomainInfo        *DomainInfo              `xml:"CommandResponse>DomainGetInfoResult"`
	DomainDNSHosts    *DomainDNSGetHostsResult `xml:"CommandResponse>DomainDNSGetHostsResult"`
	DomainDNSSetHosts *DomainDNSSetHostsResult `xml:"CommandResponse>DomainDNSSetHostsResult"`
	DomainsCheck      []DomainCheckResult      `xml:"CommandResponse>DomainCheckResult"`
	Errors            []ApiError               `xml:"Errors>Error"`
}

type ApiError struct {
	Number  int    `xml:"Number,attr"`
	Message string `xml:",innerxml"`
}

func (err *ApiError) Error() string {
	return err.Message
}

func NewClient(apiUser, apiToken, userName string) *Client {
	return &Client{
		ApiUser:    apiUser,
		ApiToken:   apiToken,
		UserName:   userName,
		HttpClient: http.DefaultClient,
		BaseURL:    defaultBaseURL,
	}
}

func (client *Client) get(request *ApiRequest) (*ApiResponse, error) {
	request.method = "GET"
	body, _, err := client.sendRequest(request)
	if err != nil {
		return nil, err
	}

	resp := new(ApiResponse)
	if err = xml.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	if resp.Status == "ERROR" {
		errMsg := ""
		for _, apiError := range resp.Errors {
			errMsg += fmt.Sprintf("Error %d: %s\n", apiError.Number, apiError.Message)
		}
		return nil, errors.New(errMsg)
	}

	return resp, nil
}

func (client *Client) makeRequest(request *ApiRequest) (*http.Request, error) {
	url, err := url.Parse(client.BaseURL)
	if err != nil {
		return nil, err
	}
	p := request.params
	p.Set("ApiUser", client.ApiUser)
	p.Set("ApiKey", client.ApiToken)
	p.Set("UserName", client.UserName)
	// This param is required by the API, but not actually used.
	p.Set("ClientIp", "127.0.0.1")
	p.Set("Command", request.command)
	url.RawQuery = p.Encode()

	// UGH
	//
	// Need this for the domain name part of the domains.check endpoint
	url.RawQuery = strings.Replace(url.RawQuery, "%2C", ",", -1)

	urlString := fmt.Sprintf("%s?%s", client.BaseURL, url.RawQuery)
	req, err := http.NewRequest(request.method, urlString, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (client *Client) sendRequest(request *ApiRequest) ([]byte, int, error) {
	req, err := client.makeRequest(request)
	if err != nil {
		return nil, 0, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return buf, resp.StatusCode, nil
}
