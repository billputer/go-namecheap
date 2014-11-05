package namecheap

import (
	"errors"
	"fmt"
	"net/url"
)

// domain type returned by 'domains.getList'
type Domain struct {
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

// domain type returned by 'domains.getInfo'
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

func (client *NamecheapClient) Domains() ([]Domain, error) {
	resp := ApiResponse{}
	requestInfo := ApiRequest{
		command: "namecheap.domains.getList",
		params:  url.Values{},
	}
	if err := client.get(requestInfo, &resp); err != nil {
		return []Domain{}, err
	}
	return resp.Domains, nil
}

func (client *NamecheapClient) Domain(domainName string) (DomainInfo, error) {
	resp := ApiResponse{}

	requestInfo := ApiRequest{
		command: "namecheap.domains.getInfo",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	if err := client.get(requestInfo, &resp); err != nil {
		return DomainInfo{}, err
	}

	if resp.Status == "ERROR" {
		err_message := ""
		for _, apiError := range resp.Errors {
			err_message += fmt.Sprintf("Error %d: %s\n", apiError.Number, apiError.Message)
		}
		return DomainInfo{}, errors.New(err_message)
	}

	return resp.DomainInfo, nil
}
