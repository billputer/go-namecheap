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
	"strconv"
	"strings"
)

const defaultBaseURL = "https://api.namecheap.com/xml.response"

// Client represents a client used to make calls to the Namecheap API.
type Client struct {
	ApiUser    string
	ApiToken   string
	UserName   string
	HttpClient *http.Client
	ClientIp   string

	// Base URL for API requests.
	// Defaults to the public Namecheap API,
	// but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	*Registrant
}

type ApiRequest struct {
	method  string
	command string
	params  url.Values
}

type ApiResponse struct {
	Status              string                     `xml:"Status,attr"`
	Command             string                     `xml:"RequestedCommand"`
	TLDList             []TLDListResult            `xml:"CommandResponse>Tlds>Tld"`
	Domains             []DomainGetListResult      `xml:"CommandResponse>DomainGetListResult>Domain"`
	DomainInfo          *DomainInfo                `xml:"CommandResponse>DomainGetInfoResult"`
	DomainDNSHosts      *DomainDNSGetHostsResult   `xml:"CommandResponse>DomainDNSGetHostsResult"`
	DomainDNSSetHosts   *DomainDNSSetHostsResult   `xml:"CommandResponse>DomainDNSSetHostsResult"`
	DomainCreate        *DomainCreateResult        `xml:"CommandResponse>DomainCreateResult"`
	DomainRenew         *DomainRenewResult         `xml:"CommandResponse>DomainRenewResult"`
	DomainsCheck        []DomainCheckResult        `xml:"CommandResponse>DomainCheckResult"`
	DomainNSInfo        *DomainNSInfoResult        `xml:"CommandResponse>DomainNSInfoResult"`
	DomainDNSSetCustom  *DomainDNSSetCustomResult  `xml:"CommandResponse>DomainDNSSetCustomResult"`
	DomainDNSSetDefault *DomainDNSSetDefaultResult `xml:"CommandResponse>DomainDNSSetDefaultResult"`
	DomainSetContacts   *DomainSetContactsResult   `xml:"CommandResponse>DomainSetContactResult"`
	SslActivate         *SslActivateResult         `xml:"CommandResponse>SSLActivateResult"`
	SslCreate           *SslCreateResult           `xml:"CommandResponse>SSLCreateResult"`
	SslCertificates     []SslGetListResult         `xml:"CommandResponse>SSLListResult>SSL"`
	UsersGetPricing     []UsersGetPricingResult    `xml:"CommandResponse>UserGetPricingResult>ProductType"`
	WhoisguardList      []WhoisguardGetListResult  `xml:"CommandResponse>WhoisguardGetListResult>Whoisguard"`
	WhoisguardEnable    whoisguardEnableResult     `xml:"CommandResponse>WhoisguardEnableResult"`
	WhoisguardDisable   whoisguardDisableResult    `xml:"CommandResponse>WhoisguardDisableResult"`
	WhoisguardRenew     *WhoisguardRenewResult     `xml:"CommandResponse>WhoisguardRenewResult"`
	Errors              ApiErrors                  `xml:"Errors>Error"`
}

// ApiError is the format of the error returned in the api responses.
type ApiError struct {
	Number  int    `xml:"Number,attr"`
	Message string `xml:",innerxml"`
}

func (err *ApiError) Error() string {
	return err.Message
}

// ApiErrors holds multiple ApiError's but implements the error interface
type ApiErrors []ApiError

func (errs ApiErrors) Error() string {
	errMsg := ""
	for _, apiError := range errs {
		errMsg += fmt.Sprintf("Error %d: %s\n", apiError.Number, apiError.Message)
	}
	return errMsg
}

func NewClient(apiUser, apiToken, userName string) *Client {
	return &Client{
		ApiUser:    apiUser,
		ApiToken:   apiToken,
		UserName:   userName,
		HttpClient: http.DefaultClient,
		BaseURL:    defaultBaseURL,
		ClientIp:   "127.0.0.1",
	}
}

// NewRegistrant associates a new registrant with the
func (client *Client) NewRegistrant(
	firstName, lastName,
	addr1, addr2,
	city, state, postalCode, country,
	phone, email string,
) {
	client.Registrant = newRegistrant(
		firstName, lastName,
		addr1, addr2,
		city, state, postalCode, country,
		phone, email,
	)
}

func (client *Client) do(request *ApiRequest) (*ApiResponse, error) {
	if request.method == "" {
		return nil, errors.New("request method cannot be blank")
	}

	body, status, err := client.sendRequest(request)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from api: %d", status)
	}

	resp := new(ApiResponse)
	if err = xml.Unmarshal(body, resp); err != nil {
		return nil, err
	}

	if resp.Status == "" {
		return nil, errors.New("failed to parse xml from api")
	}
	if resp.Status == "ERROR" {
		return nil, resp.Errors
	}

	return resp, nil
}

func (client *Client) makeRequest(request *ApiRequest) (*http.Request, error) {
	p := request.params
	p.Set("ApiUser", client.ApiUser)
	p.Set("ApiKey", client.ApiToken)
	p.Set("UserName", client.UserName)
	p.Set("ClientIp", client.ClientIp)
	p.Set("Command", request.command)

	b := p.Encode()
	req, err := http.NewRequest(request.method, client.BaseURL, strings.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(b)))
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
