package namecheap

import "net/url"

const (
	sslGetList = "namecheap.ssl.getList"
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
