package namecheap

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	whoisguardGetList = "namecheap.whoisguard.getList"
	whoisguardEnable  = "namecheap.whoisguard.enable"
	whoisguardDisable = "namecheap.whoisguard.disable"
	whoisguardRenew   = "namecheap.whoisguard.renew"
)

type WhoisguardGetListResult struct {
	ID         int64  `xml:"ID,attr"`
	DomainName string `xml:"DomainName,attr"`
	Created    string `xml:"Created,attr"`
	Expires    string `xml:"Expires,attr"`
	Status     string `xml:"Status,attr"`
}

type whoisguardEnableResult struct {
	Domain    string `xml:"Domain,attr"`
	IsSuccess bool   `xml:"IsSuccess,attr"`
}

type whoisguardDisableResult struct {
	Domain    string `xml:"Domain,attr"`
	IsSuccess bool   `xml:"IsSuccess,attr"`
}

type WhoisguardRenewResult struct {
	WhoisguardID  int64   `xml:"WhoisguardId,attr"`
	Renewed       bool    `xml:"Renew,attr"`
	ChargedAmount float64 `xml:"ChargedAmount,attr"`
	OrderID       int     `xml:"OrderId,attr"`
	TransactionID int     `xml:"TransactionId,attr"`
}

func (client *Client) WhoisguardGetList() ([]WhoisguardGetListResult, error) {
	requestInfo := &ApiRequest{
		command: whoisguardGetList,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("PageSize", "100")

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.WhoisguardList, nil
}

func (client *Client) WhoisguardEnable(id int64, email string) error {
	requestInfo := &ApiRequest{
		command: whoisguardEnable,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("WhoisguardID", strconv.FormatInt(id, 10))
	requestInfo.params.Set("ForwardedToEmail", email)
	resp, err := client.do(requestInfo)
	if err == nil && !resp.WhoisguardEnable.IsSuccess {
		err = errors.New("IsSuccess was false")
	}

	return err
}

func (client *Client) WhoisguardDisable(id int64) error {
	requestInfo := &ApiRequest{
		command: whoisguardDisable,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("WhoisguardID", strconv.FormatInt(id, 10))
	resp, err := client.do(requestInfo)
	if err == nil && !resp.WhoisguardDisable.IsSuccess {
		err = errors.New("IsSuccess was false")
	}

	return err
}

func (client *Client) WhoisguardRenew(id int64, years int) (*WhoisguardRenewResult, error) {
	requestInfo := &ApiRequest{
		command: whoisguardRenew,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("WhoisguardID", strconv.FormatInt(id, 10))
	requestInfo.params.Set("Years", strconv.Itoa(years))
	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.WhoisguardRenew, nil
}
