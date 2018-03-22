package namecheap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

// This method of testing http client APIs is borrowed from
// Will Norris's work in go-github @ https://github.com/google/go-github
func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("anApiUser", "anToken", "anUser")
	client.BaseURL = server.URL + "/"
}

func fillDefaultParams(p url.Values) url.Values {
	p.Set("ApiKey", "anToken")
	p.Set("ApiUser", "anApiUser")
	p.Set("ClientIp", "127.0.0.1")
	p.Set("UserName", "anUser")
	return p
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func testBody(t *testing.T, r *http.Request, p url.Values) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}
	if p.Encode() != string(b) {
		t.Errorf("Body:\n %v\nwant:\n %v", string(b), p.Encode())
	}
}

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
	requestInfo := &ApiRequest{
		method:  "POST",
		command: "namecheap.domains.getList",
		params:  url.Values{},
	}
	req, _ := c.makeRequest(requestInfo)

	// correctly assembled URL
	outURL := "https://fake-api-server/"
	// test that URL was correctly assembled
	if req.URL.String() != outURL {
		t.Errorf("NewRequest() URL = %v, want %v", req.URL, outURL)
	}

	correctParams := fillDefaultParams(url.Values{})
	correctParams.Set("Command", "namecheap.domains.getList")
	testBody(t, req, correctParams)
}

// Verify that Do correctly handles errors and invalid returns
func TestDo(t *testing.T) {
	setup()
	defer teardown()

	state := "500"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch state {
		case "500":
			w.WriteHeader(http.StatusInternalServerError)
		case "invalid":
			fmt.Fprintln(w, "<invalid />")
		case "error":
			fmt.Fprintln(w, `<?xml version="1.0" encoding="utf-8"?><ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
				<Errors>
					<Error Number="1">Some Error</Error>
				</Errors>
			</ApiResponse>`)
		default:
			fmt.Fprintln(w, `<?xml version="1.0" encoding="utf-8"?><ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response"></ApiResponse>`)
		}
	})

	requestInfo := &ApiRequest{
		method:  "POST",
		command: "namecheap.domains.getList",
		params:  url.Values{},
	}
	_, err := client.do(requestInfo)
	if err == nil {
		t.Errorf("Expected error for non-200 response, got %v", err)
	}

	state = "invalid"
	_, err = client.do(requestInfo)
	if err == nil {
		t.Errorf("Expected error for invalid response, got %v", err)
	}

	state = "error"
	_, err = client.do(requestInfo)
	if err == nil || err.Error() != "Error 1: Some Error\n" {
		t.Errorf("Expected error for error response, got %v", err)
	}

	state = "ok"
	resp, err := client.do(requestInfo)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if resp == nil {
		t.Errorf("Expected non-nil response, got %v", resp)
	}
}
