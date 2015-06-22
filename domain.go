package namecheap

import (
	"net/url"
	"strings"
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

type DomainCheckResult struct {
	Domain    string `xml:"Domain,attr"`
	Available bool   `xml:"Available,attr"`
}

func (client *Client) DomainsGetList() ([]DomainGetListResult, error) {
	resp := new(ApiResponse)
	requestInfo := &ApiRequest{
		command: domainsGetList,
		params:  url.Values{},
	}
	if err := client.get(requestInfo, resp); err != nil {
		return nil, err
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

	return resp.DomainInfo, nil
}

func (client *Client) DomainsCheck(domainNames ...string) ([]DomainCheckResult, error) {
	resp := new(ApiResponse)
	requestInfo := &ApiRequest{
		command: domainsCheck,
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainList", strings.Join(domainNames, ","))
	if err := client.get(requestInfo, resp); err != nil {
		return nil, err
	}

	return resp.DomainsCheck, nil
}
