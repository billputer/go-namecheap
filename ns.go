package namecheap

import "net/url"

const (
	nsCreate  = "namecheap.domains.ns.create"
	nsDelete  = "namecheap.domains.ns.delete"
	nsGetInfo = "namecheap.domains.ns.getInfo"
	nsUpdate  = "namecheap.domains.ns.update"
)

type DomainNSCreateResult struct {
	Domain     string `xml:"Domain,attr"`
	Nameserver string `xml:"Nameserver,attr"`
	IP         string `xml:"IP,attr"`
	IsSuccess  bool   `xml:"IsSuccess,attr"`
}

type DomainNSDeleteResult struct {
	Domain     string `xml:"Domain,attr"`
	Nameserver string `xml:"Nameserver,attr"`
	IsSuccess  bool   `xml:"IsSuccess,attr"`
}

type DomainNSInfoResult struct {
	Domain     string   `xml:"Domain,attr"`
	Nameserver string   `xml:"Nameserver,attr"`
	IP         string   `xml:"IP,attr"`
	Statuses   []string `xml:"NameserverStatuses>Status"`
}

type DomainNSUpdateResult struct {
	Domain     string `xml:"Domain,attr"`
	Nameserver string `xml:"Nameserver,attr"`
	IsSuccess  bool   `xml:"IsSuccess,attr"`
}

func (client *Client) NSCreate(sld, tld, nameserver, ip string) (*DomainNSCreateResult, error) {
	requestInfo := &ApiRequest{
		command: nsCreate,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("SLD", sld)
	requestInfo.params.Set("TLD", tld)
	requestInfo.params.Set("Nameserver", nameserver)
	requestInfo.params.Set("IP", ip)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainNSCreate, nil
}

func (client *Client) NSDelete(sld, tld, nameserver string) (*DomainNSDeleteResult, error) {
	requestInfo := &ApiRequest{
		command: nsDelete,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("SLD", sld)
	requestInfo.params.Set("TLD", tld)
	requestInfo.params.Set("Nameserver", nameserver)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainNSDelete, nil
}

func (client *Client) NSGetInfo(sld, tld, nameserver string) (*DomainNSInfoResult, error) {
	requestInfo := &ApiRequest{
		command: nsGetInfo,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("SLD", sld)
	requestInfo.params.Set("TLD", tld)
	requestInfo.params.Set("Nameserver", nameserver)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainNSInfo, nil
}

func (client *Client) NSUpdate(sld, tld, nameserver string) (*DomainNSUpdateResult, error) {
	requestInfo := &ApiRequest{
		command: nsUpdate,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("SLD", sld)
	requestInfo.params.Set("TLD", tld)
	requestInfo.params.Set("Nameserver", nameserver)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainNSUpdate, nil
}
