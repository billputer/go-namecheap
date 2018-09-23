package namecheap

import (
	"net/url"
	"strconv"
)

const (
	sslActivate             = "namecheap.ssl.activate"
	sslCreate               = "namecheap.ssl.create"
	sslGetList              = "namecheap.ssl.getList"
	sslGetApproverEmailList = "namecheap.ssl.getApproverEmailList"
	sslGetInfo              = "namecheap.ssl.getInfo"
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

type SslGetApproverEmailListResult struct {
	DomainEmails  []string `xml:"Domainemails>email"`
	GenericEmails []string `xml:"Genericemails>email"`
	ManualEmails  []string `xml:"Manualemails>email"`
}

type SslGetInfoResult struct {
	Status               string                `xml:"Status,attr"`
	StatusDescription    string                `xml:"StatusDescription,attr"`
	Type                 string                `xml:"Type,attr"`
	IssuedOn             string                `xml:"IssuedOn,attr"`
	Expires              string                `xml:"Expires,attr"`
	ActivationExpireDate string                `xml:"ActivationExpireDate,attr"`
	OrderID              int                   `xml:"OrderId,attr"`
	ReplacedBy           int                   `xml:"ReplacedBy,attr"`
	SANSCount            int                   `xml:"SANSCount,attr"`
	CertificateDetails   SslCertificateDetails `xml:"CertificateDetails"`
	Provider             SslProvider           `xml:"Provider"`
}

type SslCertificateDetails struct {
	CSR                string                 `xml:"CSR"`
	ApproverEmail      string                 `xml:"ApproverEmail"`
	CommonName         string                 `xml:"CommonName"`
	AdministratorName  string                 `xml:"AdministratorName"`
	AdministratorEmail string                 `xml:"AdministratorEmail"`
	Certificates       SslGetInfoCertificates `xml:"Certificates"`
}

type SslGetInfoCertificates struct {
	CertificateReturned bool     `xml:"CertificateReturned,attr"`
	ReturnType          string   `xml:"ReturnType,attr"`
	Certificate         []string `xml:"Certificate"`

	CACertificates SslCACertificates `xml:"CaCertificates"`
}

type SslCACertificates struct {
	Certificate []SslCACertificate `xml:"Certificate"`
}

type SslCACertificate struct {
	Type        string `xml:"Type,attr"`
	Certificate string `xml:"Certificate"`
}

type SslProvider struct {
	OrderID int    `xml:"OrderID"`
	Name    string `xml:"Name"`
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

// SslGetApproverEmailList returns email addresses that can be used for domain
// control validation
func (client *Client) SslGetApproverEmailList(domainName, certificateType string) (*SslGetApproverEmailListResult, error) {
	requestInfo := &ApiRequest{
		command: sslGetApproverEmailList,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	requestInfo.params.Set("CertificateType", certificateType)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.SslGetApproverEmailList, nil
}

// SslGetInfo returns information about a purchased certificate
func (client *Client) SslGetInfo(certificateID int) (*SslGetInfoResult, error) {
	requestInfo := &ApiRequest{
		command: sslGetInfo,
		method:  "GET",
		params:  url.Values{},
	}

	requestInfo.params.Set("CertificateID", strconv.Itoa(certificateID))
	requestInfo.params.Set("Returncertificate", "true")
	requestInfo.params.Set("Returntype", "Individual")

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.SslGetInfo, nil
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
