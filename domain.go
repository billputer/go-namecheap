package namecheap

import (
	"net/url"
)

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
