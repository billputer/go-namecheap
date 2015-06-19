package namecheap

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	domainsGetList = "namecheap.domains.getList"
	domainsGetInfo = "namecheap.domains.getInfo"
	domainsCheck   = "namecheap.domains.check"
)

// DomainGetListResult represents the data returned by 'domains.getList'
type DomainGetListResult struct {
	ID         int    `xml:"ID,attr"`
	Name       string `xml:"Name,attr"`
	User       string `xml:"User,attr"`
	Created    string `xml:"Created,attr"`
	Expires    string `xml:"Expires,attr"`
	IsExpired  bool   `xml:"IsExpired,attr"`
	IsLocked   bool   `xml:"IsLocked,attr"`
	AutoRenew  bool   `xml:"AutoRenew,attr"`
	WhoisGuard string `xml:"WhoisGuard,attr"`
}

// DomainInfo represents the data returned by 'domains.getInfo'
type DomainInfo struct {
	ID         int        `xml:"ID,attr"`
	Name       string     `xml:"DomainName,attr"`
	Owner      string     `xml:"OwnerName,attr"`
	Created    string     `xml:"DomainDetails>CreatedDate"`
	Expires    string     `xml:"DomainDetails>ExpiredDate"`
	IsExpired  bool       `xml:"IsExpired,attr"`
	IsLocked   bool       `xml:"IsLocked,attr"`
	AutoRenew  bool       `xml:"AutoRenew,attr"`
	DNSDetails DNSDetails `xml:"DnsDetails"`
}

type DNSDetails struct {
	ProviderType  string   `xml:"ProviderType,attr"`
	IsUsingOurDNS bool     `xml:"IsUsingOurDNS"`
	Nameservers   []string `xml:"Nameserver"`
}

type DomainCheckStatus struct {
	Domain    string `xml:"Domain,attr"`
	Available bool   `xml:""`
}

func (client *Client) DomainsGetList() ([]DomainGetListResult, error) {
	resp := new(ApiResponse)
	requestInfo := &ApiRequest{
		command: domainsGetList,
		params:  url.Values{},
	}
	if err := client.get(requestInfo, resp); err != nil {
		return []DomainGetListResult{}, err
	}

	return resp.Domains, nil
}

func (client *Client) DomainGetInfo(domainName string) (*DomainInfo, error) {
	resp := new(ApiResponse)

	requestInfo := &ApiRequest{
		command: domainsGetInfo,
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	if err := client.get(requestInfo, resp); err != nil {
		return nil, err
	}

	if resp.Status == "ERROR" {
		errMsg := ""
		for _, apiError := range resp.Errors {
			errMsg += fmt.Sprintf("Error %d: %s\n", apiError.Number, apiError.Message)
		}
		return nil, errors.New(errMsg)
	}

	return resp.DomainInfo, nil
}
