package namecheap

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	domainsGetList          = "namecheap.domains.getList"
	domainsGetInfo          = "namecheap.domains.getInfo"
	domainsCheck            = "namecheap.domains.check"
	domainsCreate           = "namecheap.domains.create"
	domainsTLDList          = "namecheap.domains.getTldList"
	domainsRenew            = "namecheap.domains.renew"
	domainsSetContacts      = "namecheap.domains.setContacts"
	domainsGetRegistrarLock = "namecheap.domains.getRegistrarLock"
	domainsSetRegistrarLock = "namecheap.domains.setRegistrarLock"
)

// DomainGetListResult represents the data returned by 'domains.getList'
type DomainGetListResult struct {
	ID            int       `xml:"ID,attr"`
	Name          string    `xml:"Name,attr"`
	User          string    `xml:"User,attr"`
	Created       time.Time `xml:"-"`
	CreatedString string    `xml:"Created,attr"`
	Expires       time.Time `xml:"-"`
	ExpiresString string    `xml:"Expires,attr"`
	IsExpired     bool      `xml:"IsExpired,attr"`
	IsLocked      bool      `xml:"IsLocked,attr"`
	AutoRenew     bool      `xml:"AutoRenew,attr"`
	WhoisGuard    string    `xml:"WhoisGuard,attr"`
}

// DomainInfo represents the data returned by 'domains.getInfo'
type DomainInfo struct {
	ID            int        `xml:"ID,attr"`
	Name          string     `xml:"DomainName,attr"`
	Owner         string     `xml:"OwnerName,attr"`
	Created       time.Time  `xml:"-"`
	CreatedString string     `xml:"DomainDetails>CreatedDate"`
	Expires       time.Time  `xml:"-"`
	ExpiresString string     `xml:"DomainDetails>ExpiredDate"`
	IsExpired     bool       `xml:"IsExpired,attr"`
	IsLocked      bool       `xml:"IsLocked,attr"`
	AutoRenew     bool       `xml:"AutoRenew,attr"`
	DNSDetails    DNSDetails `xml:"DnsDetails"`
	Whoisguard    Whoisguard `xml:"Whoisguard"`
}

type DNSDetails struct {
	ProviderType  string   `xml:"ProviderType,attr"`
	IsUsingOurDNS bool     `xml:"IsUsingOurDNS,attr"`
	Nameservers   []string `xml:"Nameserver"`
}

type WhoisguardEmailDetails struct {
	Email                        string `xml:"WhoisGuardEmail,attr"`
	ForwardedTo                  string `xml:"ForwardedTo,attr"`
	LastAutoEmailChangeDate      string `xml:"LastAutoEmailChangeDate,attr"`
	AutoEmailChangeFrequencyDays int    `xml:"AutoEmailChangeFrequencyDays,attr"`
}
type Whoisguard struct {
	RawEnabled        string                 `xml:"Enabled,attr"`
	Enabled           bool                   `xml:"-"`
	ID                int64                  `xml:"ID"`
	ExpiredDate       time.Time              `xml:"-"`
	ExpiredDateString string                 `xml:"ExpiredDate"`
	EmailDetails      WhoisguardEmailDetails `xml:"EmailDetails"`
}

type DomainCheckResult struct {
	Domain                   string  `xml:"Domain,attr"`
	Available                bool    `xml:"Available,attr"`
	IsPremiumName            bool    `xml:"IsPremiumName,attr"`
	PremiumRegistrationPrice float64 `xml:"PremiumRegistrationPrice,attr"`
	PremiumRenewalPrice      float64 `xml:"PremiumRenewalPrice,attr"`
	PremiumRestorePrice      float64 `xml:"PremiumRestorePrice,attr"`
	PremiumTransferPrice     float64 `xml:"PremiumTransferPrice,attr"`
	IcannFee                 float64 `xml:"IcannFee,attr"`
}

type TLDListResult struct {
	Name string `xml:"Name,attr"`
}

type DomainCreateResult struct {
	Domain            string  `xml:"Domain,attr"`
	Registered        bool    `xml:"Registered,attr"`
	ChargedAmount     float64 `xml:"ChargedAmount,attr"`
	DomainID          int     `xml:"DomainID,attr"`
	OrderID           int     `xml:"OrderID,attr"`
	TransactionID     int     `xml:"TransactionID,attr"`
	WhoisguardEnable  bool    `xml:"WhoisguardEnable,attr"`
	NonRealTimeDomain bool    `xml:"NonRealTimeDomain,attr"`
}

type DomainRenewResult struct {
	DomainID         int     `xml:"DomainID,attr"`
	Name             string  `xml:"DomainName,attr"`
	Renewed          bool    `xml:"Renew,attr"`
	ChargedAmount    float64 `xml:"ChargedAmount,attr"`
	OrderID          int     `xml:"OrderID,attr"`
	TransactionID    int     `xml:"TransactionID,attr"`
	ExpireDate       time.Time
	ExpireDateString string `xml:"DomainDetails>ExpiredDate"`
}

type DomainSetContactsResult struct {
	Name      string `xml:"Domain,attr"`
	IsSuccess bool   `xml:"IsSuccess,attr"`
}

type DomainRegistrarLockStatusResult struct {
	Name     string `xml:"Domain,attr"`
	IsLocked bool   `xml:"RegistrarLockStatus,attr"`
}

type DomainSetRegistrarLockResult struct {
	Name      string `xml:"Domain,attr"`
	IsSuccess bool   `xml:"IsSuccess,attr"`
}
type DomainCreateOption struct {
	AddFreeWhoisguard bool
	WGEnabled         bool
	Nameservers       []string
}

func (client *Client) DomainsGetList() ([]DomainGetListResult, error) {
	requestInfo := &ApiRequest{
		command: domainsGetList,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("PageSize", "100")

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	for _, domain := range resp.Domains {
		if domain.CreatedString != "" {
			created, createdErr := resp.ParseDate(domain.CreatedString)
			if createdErr != nil {
				return nil, fmt.Errorf("error when parsing Created date to time.Time: %w", createdErr)
			}

			domain.Created = created
		}

		if domain.ExpiresString != "" {
			expires, expiresErr := resp.ParseDate(domain.ExpiresString)
			if expiresErr != nil {
				return nil, fmt.Errorf("error when parsing Expires date to time.Time: %w", expiresErr)
			}

			domain.Expires = expires
		}
	}

	return resp.Domains, nil
}

func (client *Client) DomainGetInfo(domainName string) (*DomainInfo, error) {
	requestInfo := &ApiRequest{
		command: domainsGetInfo,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainName", domainName)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	if resp.DomainInfo != nil {
		if strings.EqualFold(resp.DomainInfo.Whoisguard.RawEnabled, "true") {
			resp.DomainInfo.Whoisguard.Enabled = true
		}

		if resp.DomainInfo.CreatedString != "" {
			created, createdErr := resp.ParseDate(resp.DomainInfo.CreatedString)
			if createdErr != nil {
				return nil, fmt.Errorf("error when parsing domain Created date to time.Time: %w", createdErr)
			}

			resp.DomainInfo.Created = created
		}

		if resp.DomainInfo.ExpiresString != "" {
			expires, expiresErr := resp.ParseDate(resp.DomainInfo.ExpiresString)
			if expiresErr != nil {
				return nil, fmt.Errorf("error when parsing domain Expires date to time.Time: %w", expiresErr)
			}

			resp.DomainInfo.Expires = expires
		}

		if resp.DomainInfo.Whoisguard.ExpiredDateString != "" {
			expired, expiredErr := resp.ParseDate(resp.DomainInfo.Whoisguard.ExpiredDateString)
			if expiredErr != nil {
				return nil, fmt.Errorf("error when parsing WhoisGuard Expired date to time.Time: %w", expiredErr)
			}

			resp.DomainInfo.Whoisguard.ExpiredDate = expired
		}

	}

	return resp.DomainInfo, nil
}

func (client *Client) DomainsCheck(domainNames ...string) ([]DomainCheckResult, error) {
	requestInfo := &ApiRequest{
		command: domainsCheck,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainList", strings.Join(domainNames, ","))
	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainsCheck, nil
}

func (client *Client) DomainsTLDList() ([]TLDListResult, error) {
	requestInfo := &ApiRequest{
		command: domainsTLDList,
		method:  "POST",
		params:  url.Values{},
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.TLDList, nil
}

func (client *Client) DomainCreate(domainName string, years int, options ...DomainCreateOption) (*DomainCreateResult, error) {
	if client.Registrant == nil {
		return nil, errors.New("Registrant information on client cannot be empty")
	}

	requestInfo := &ApiRequest{
		command: domainsCreate,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("DomainName", domainName)
	requestInfo.params.Set("Years", strconv.Itoa(years))
	for _, opt := range options {
		if opt.AddFreeWhoisguard {
			requestInfo.params.Set("AddFreeWhoisguard", "yes")
		}
		if opt.WGEnabled {
			requestInfo.params.Set("WGEnabled", "yes")
		}
		if len(opt.Nameservers) > 0 {
			requestInfo.params.Set("Nameservers", strings.Join(opt.Nameservers, ","))
		}
	}
	if err := client.Registrant.addValues(requestInfo.params); err != nil {
		return nil, err
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainCreate, nil
}

func (client *Client) DomainRenew(domainName string, years int) (*DomainRenewResult, error) {
	requestInfo := &ApiRequest{
		command: domainsRenew,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	requestInfo.params.Set("Years", strconv.Itoa(years))

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	if resp.DomainRenew.ExpireDateString != "" {
		expired, expiredErr := resp.ParseDate(resp.DomainRenew.ExpireDateString)
		if expiredErr != nil {
			return nil, fmt.Errorf("error when parsing Expired date to time.Time: %w", expiredErr)
		}

		resp.DomainRenew.ExpireDate = expired
	}

	return resp.DomainRenew, nil
}

func (client *Client) DomainSetContacts(domainName string) (*DomainSetContactsResult, error) {
	requestInfo := &ApiRequest{
		command: domainsSetContacts,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	if err := client.Registrant.addValues(requestInfo.params); err != nil {
		return nil, err
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainSetContacts, nil
}

func (client *Client) DomainGetRegistrarLock(domainName string) (*DomainRegistrarLockStatusResult, error) {
	requestInfo := &ApiRequest{
		command: domainsGetRegistrarLock,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainRegistrarLockStatus, nil
}

func (client *Client) DomainSetRegistrarLock(domainName string, lock bool) (*DomainSetRegistrarLockResult, error) {
	requestInfo := &ApiRequest{
		command: domainsSetRegistrarLock,
		method:  "POST",
		params:  url.Values{},
	}
	requestInfo.params.Set("DomainName", domainName)
	if !lock {
		requestInfo.params.Set("LockAction", "UNLOCK")
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.DomainSetRegistrarLock, nil
}
