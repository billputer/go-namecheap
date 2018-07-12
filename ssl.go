package namecheap

import (
	"net/url"
	"strconv"
)

const (
	sslGetList = "namecheap.ssl.getList"
	sslCreate  = "namecheap.ssl.create"
)

// SslGetListResult represents the data returned by 'domains.getList'
type SslGetListResult struct {
	CertificateID        int    `xml:"CertificateID,attr"`
	HostName             string `xml:"HostName,attr"`
	SSLType              string `xml:"SSLType,attr"`
	PurchaseDate         string `xml:"PurchaseDate,attr"`
	ExpireDate           string `xml:"ExpireDate,attr"`
	ActivationExpireDate string `xml:"ActivationExpireDate,attr"`
	IsExpired            bool   `xml:"IsExpiredYN,attr"`
	Status               string `xml:"Status,attr"`
}

type SslCreateResult struct {
	IsSuccess      bool             `xml:"IsSuccess,attr"`
	OrderId        int              `xml:"OrderId,attr"`
	TransactionId  int              `xml:"TransactionId,attr"`
	ChargedAmount  float64          `xml:"ChargedAmount,attr"`
	SSLCertificate []SSLCertificate `xml:"SSLCertificate"`
}

type SSLCertificate struct {
	CertificateID int    `xml:"CertificateID,attr"`
	SSLType       string `xml:"SSLType,attr"`
	Created       string `xml:"Created,attr"`
	Years         int    `xml:"Years,attr"`
	Status        string `xml:"Status,attr"`
}

// SslGetList gets a list of SSL certificates for a particular user
func (client *Client) SslGetList() ([]SslGetListResult, error) {
	requestInfo := &ApiRequest{
		command: sslGetList,
		method:  "POST",
		params:  url.Values{},
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.SslCertificates, nil
}

// SslCreate creates a new SSL certificate by purchasing it using the account funds
func (client *Client) SslCreate(productType string, years int) (*SslCreateResult, error) {
	requestInfo := &ApiRequest{
		command: sslCreate,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("Type", productType)
	requestInfo.params.Set("Years", strconv.Itoa(years))

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.SslCreate, nil
}
