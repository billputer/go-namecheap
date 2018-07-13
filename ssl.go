package namecheap

import (
	"net/url"
	"strconv"
)

const (
	sslActivate = "namecheap.ssl.activate"
	sslCreate   = "namecheap.ssl.create"
	sslGetList  = "namecheap.ssl.getList"
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

type SslActivateParams struct {
	CertificateId      int
	Csr                string
	AdminEmailAddress  string
	WebServerType      string
	ApproverEmail      string
	IsHTTPDCValidation bool
	IsDNSDCValidation  bool
}

type SslActivateResult struct {
	ID               int             `xml:"ID,attr"`
	IsSuccess        bool            `xml:"IsSuccess,attr"`
	HttpDCValidation SslDcValidation `xml:"HttpDCValidation"`
	DNSDCValidation  SslDcValidation `xml:"DNSDCValidation"`
}

type SslDcValidation struct {
	ValueAvailable bool   `xml:"ValueAvailable,attr"`
	Dns            SslDns `xml:"DNS"`
}

type SslDns struct {
	Domain      string `xml:"domain,attr"`
	FileName    string `xml:"FileName,omitempty"`
	FileContent string `xml:"FileContent,omitempty"`
	HostName    string `xml:"HostName,omitempty"`
	Target      string `xml:"Target,omitempty"`
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

// SslActivate activates a purchased and non-activated SSL certificate
func (client *Client) SslActivate(params SslActivateParams) (*SslActivateResult, error) {
	requestInfo := &ApiRequest{
		command: sslActivate,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("CertificateID", strconv.Itoa(params.CertificateId))
	requestInfo.params.Set("CSR", params.Csr)
	requestInfo.params.Set("AdminEmailAddress", params.AdminEmailAddress)
	requestInfo.params.Set("WebServerType", params.WebServerType)

	if params.IsHTTPDCValidation {
		requestInfo.params.Set("HTTPDCValidation", "true")
	}

	if params.IsDNSDCValidation {
		requestInfo.params.Set("DNSDCValidation", "true")
	}

	if len(params.ApproverEmail) > 0 {
		requestInfo.params.Set("ApproverEmail", params.ApproverEmail)
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.SslActivate, nil
}
