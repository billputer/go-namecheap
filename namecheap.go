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
	ApiUser     string
	ApiToken    string
	UserName    string
	HttpClient  *http.Client

	// Base URL for API requests.
	// Defaults to the public Namecheap API, but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string
}

func NewClient(apiUser, apiToken, userName string) *NamecheapClient {
	return &NamecheapClient{ApiUser: apiUser, ApiToken: apiToken, UserName: userName, HttpClient: &http.Client{}, BaseURL: defaultBaseURL}
}

func (client *NamecheapClient) get(command string, val interface{}) error {
	body, _, err := client.sendRequest("GET", command, nil)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal([]byte(body), &val); err != nil {
		return err
	}

	return nil
}

func (client *NamecheapClient) makeRequest(method, command string, body io.Reader) (*http.Request, error) {
	url, err := url.Parse(client.BaseURL)
	if err != nil {
		return nil, err
	}
	q := url.Query()
	q.Set("ApiUser", client.ApiUser)
	q.Set("ApiKey", client.ApiToken)
	q.Set("UserName", client.UserName)
	q.Set("ClientIp", "127.0.0.1")
	q.Set("Command", command)
	url.RawQuery = q.Encode()

	urlString := fmt.Sprintf("%s?%s", client.BaseURL, q.Encode())
	req, err := http.NewRequest(method, urlString, body)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (client *NamecheapClient) sendRequest(method, command string, body io.Reader) (string, int, error) {
	req, err := client.makeRequest(method, command, body)
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
