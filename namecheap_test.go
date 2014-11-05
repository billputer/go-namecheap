package namecheap

import (
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("anApiUser", "anToken", "anUser")

	if c.BaseURL != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL, defaultBaseURL)
	}
}

// Verify that the MakeRequest function assembles the correct API URL
func TestMakeRequest(t *testing.T) {
	c := NewClient("anApiUser", "anToken", "anUser")
	c.BaseURL = "https://fake-api-server/"
	requestInfo := ApiRequest{
		method:  "GET",
		command: "namecheap.domains.getList",
		params:  url.Values{},
	}
	req, _ := c.makeRequest(requestInfo, nil)

	// correctly assembled URL
	outURL := "https://fake-api-server/?ApiKey=anToken&ApiUser=anApiUser&ClientIp=127.0.0.1&Command=namecheap.domains.getList&UserName=anUser"

	// test that URL was correctly assembled
	if req.URL.String() != outURL {
		t.Errorf("NewRequest() URL = %v, want %v", req.URL, outURL)
	}
}
