package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDomainsGetList(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
    <?xml version="1.0" encoding="utf-8"?>
    <ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
      <Errors />
      <Warnings />
      <RequestedCommand>namecheap.domains.getList</RequestedCommand>
      <CommandResponse Type="namecheap.domains.getList">
        <DomainGetListResult>
          <Domain ID="57579" Name="example.com" User="anUser" Created="11/04/2014" Expires="11/04/2015" IsExpired="false" IsLocked="false" AutoRenew="false" WhoisGuard="ENABLED" />
        </DomainGetListResult>
        <Paging>
          <TotalItems>12</TotalItems>
          <CurrentPage>1</CurrentPage>
          <PageSize>20</PageSize>
        </Paging>
      </CommandResponse>
      <Server>WEB1-SANDBOX1</Server>
      <GMTTimeDifference>--5:00</GMTTimeDifference>
      <ExecutionTime>0.009</ExecutionTime>
    </ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.getList")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	domains, err := client.DomainsGetList()

	if err != nil {
		t.Errorf("DomainsGetList returned error: %v", err)
	}

	// DomainGetListResult we expect, given the respXML above
	want := []DomainGetListResult{{
		ID:         57579,
		Name:       "example.com",
		User:       "anUser",
		Created:    "11/04/2014",
		Expires:    "11/04/2015",
		IsExpired:  false,
		IsLocked:   false,
		AutoRenew:  false,
		WhoisGuard: "ENABLED",
	}}

	if !reflect.DeepEqual(domains, want) {
		t.Errorf("DomainsGetList returned %+v, want %+v", domains, want)
	}
}

func TestDomainGetInfo(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
  <Errors />
  <Warnings />
  <RequestedCommand>namecheap.domains.getInfo</RequestedCommand>
  <CommandResponse Type="namecheap.domains.getInfo">
    <DomainGetInfoResult Status="Ok" ID="57582" DomainName="example.com" OwnerName="anUser" IsOwner="true">
      <DomainDetails>
        <CreatedDate>11/04/2014</CreatedDate>
        <ExpiredDate>11/04/2015</ExpiredDate>
        <NumYears>0</NumYears>
      </DomainDetails>
      <LockDetails />
      <Whoisguard Enabled="True">
        <ID>53536</ID>
        <ExpiredDate>11/04/2015</ExpiredDate>
        <EmailDetails WhoisGuardEmail="08040e11d32d48ebb4346b02b98dda17.protect@whoisguard.com" ForwardedTo="billwiens@gmail.com" LastAutoEmailChangeDate="" AutoEmailChangeFrequencyDays="0" />
      </Whoisguard>
      <DnsDetails ProviderType="FREE" IsUsingOurDNS="true">
        <Nameserver>dns1.registrar-servers.com</Nameserver>
        <Nameserver>dns2.registrar-servers.com</Nameserver>
        <Nameserver>dns3.registrar-servers.com</Nameserver>
        <Nameserver>dns4.registrar-servers.com</Nameserver>
        <Nameserver>dns5.registrar-servers.com</Nameserver>
      </DnsDetails>
      <Modificationrights All="true" />
    </DomainGetInfoResult>
  </CommandResponse>
  <Server>WEB1-SANDBOX1</Server>
  <GMTTimeDifference>--5:00</GMTTimeDifference>
  <ExecutionTime>0.008</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.getInfo")
		correctParams.Set("DomainName", "example.com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	domain, err := client.DomainGetInfo("example.com")

	if err != nil {
		t.Errorf("DomainGetInfo returned error: %v", err)
	}

	// DomainInfo we expect, given the respXML above
	want := &DomainInfo{
		ID:        57582,
		Name:      "example.com",
		Owner:     "anUser",
		Created:   "11/04/2014",
		Expires:   "11/04/2015",
		IsExpired: false,
		IsLocked:  false,
		AutoRenew: false,
		DNSDetails: DNSDetails{
			ProviderType:  "FREE",
			IsUsingOurDNS: true,
			Nameservers: []string{
				"dns1.registrar-servers.com",
				"dns2.registrar-servers.com",
				"dns3.registrar-servers.com",
				"dns4.registrar-servers.com",
				"dns5.registrar-servers.com",
			},
		},
		Whoisguard: Whoisguard{
			Enabled:     true,
			ID:          53536,
			ExpiredDate: "11/04/2015",
		},
	}

	if !reflect.DeepEqual(domain, want) {
		t.Errorf("DomainGetInfo returned %+v, want %+v", domain, want)
	}
}

func TestDomainsCheck(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.check</RequestedCommand>
  <CommandResponse Type="namecheap.domains.check">
    <DomainCheckResult Domain="domain1.com" Available="true" />
    <DomainCheckResult Domain="availabledomain.com" Available="false" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.check")
		correctParams.Set("DomainList", "domain1.com,availabledomain.com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	domains, err := client.DomainsCheck("domain1.com", "availabledomain.com")
	if err != nil {
		t.Errorf("DomainsCheck returned error: %v", err)
	}

	// DomainCheckResult we expect, given the respXML above
	want := []DomainCheckResult{
		DomainCheckResult{
			Domain:    "domain1.com",
			Available: true,
		},
		DomainCheckResult{
			Domain:    "availabledomain.com",
			Available: false,
		},
	}

	if !reflect.DeepEqual(domains, want) {
		t.Errorf("DomainsCheck returned %+v, want %+v", domains, want)
	}
}

func TestDomainCreate(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
	<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
	  <Errors />
	  <RequestedCommand>namecheap.domains.create</RequestedCommand>
	  <CommandResponse Type="namecheap.domains.create">
	    <DomainCreateResult Domain="domain1.com" Registered="true" ChargedAmount="20.3600" DomainID="9007" OrderID="196074" TransactionID="380716" WhoisguardEnable="true" NonRealTimeDomain="false" />
	  </CommandResponse>
	  <Server>SERVER-NAME</Server>
	  <GMTTimeDifference>+5</GMTTimeDifference>
	  <ExecutionTime>0.078</ExecutionTime>
	</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		fillInfo := func(prefix string) {
			correctParams.Set(prefix+"FirstName", "John")
			correctParams.Set(prefix+"LastName", "Smith")
			correctParams.Set(prefix+"Address1", "8939 S.cross Blvd")
			correctParams.Set(prefix+"StateProvince", "CA")
			correctParams.Set(prefix+"PostalCode", "90045")
			correctParams.Set(prefix+"Country", "US")
			correctParams.Set(prefix+"Phone", "+1.6613102107")
			correctParams.Set(prefix+"EmailAddress", "john@gmail.com")
			correctParams.Set(prefix+"City", "CA")
		}
		correctParams.Set("Command", "namecheap.domains.create")
		correctParams.Set("DomainName", "domain1.com")
		correctParams.Set("Years", "2")
		correctParams.Set("AddFreeWhoisguard", "yes")
		correctParams.Set("WGEnabled", "yes")
		correctParams.Set("Nameservers", "ns1.test.com,ns2.test.com")
		fillInfo("AuxBilling")
		fillInfo("Tech")
		fillInfo("Admin")
		fillInfo("Registrant")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	client.NewRegistrant(
		"John", "Smith",
		"8939 S.cross Blvd", "",
		"CA", "CA", "90045", "US",
		"+1.6613102107", "john@gmail.com",
	)

	result, err := client.DomainCreate("domain1.com", 2, DomainCreateOption{
		AddFreeWhoisguard: true,
		WGEnabled:         true,
		Nameservers: []string{
			"ns1.test.com",
			"ns2.test.com",
		},
	})
	if err != nil {
		t.Fatalf("DomainCreate returned error: %v", nil)
	}

	// DomainCreateResult we expect, given the respXML above
	want := &DomainCreateResult{
		"domain1.com", true, 20.36, 9007, 196074, 380716, true, false,
	}

	if !reflect.DeepEqual(result, want) {
		t.Fatalf("DomainCreate returned\n%+v,\nwant\n%+v", result, want)
	}
}

func TestDomainsRenew(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.renew</RequestedCommand>
  <CommandResponse Type="namecheap.domains.renew">
		<DomainRenewResult DomainName="domain1.com" DomainID="151378" Renew="true" OrderID="109116" TransactionID="119569" ChargedAmount="650.0000">
			<DomainDetails>
        <ExpiredDate>4/30/2021 11:31:13 AM</ExpiredDate>
        <NumYears>0</NumYears>
			</DomainDetails>
		</DomainRenewResult>
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.renew")
		correctParams.Set("DomainName", "domain1.com")
		correctParams.Set("Years", "1")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	result, err := client.DomainRenew("domain1.com", 1)
	if err != nil {
		t.Errorf("DomainRenew returned error: %v", err)
	}

	// DomainCheckResult we expect, given the respXML above
	want := &DomainRenewResult{
		DomainID:      151378,
		Name:          "domain1.com",
		Renewed:       true,
		ChargedAmount: 650,
		TransactionID: 119569,
		OrderID:       109116,
		ExpireDate:    "4/30/2021 11:31:13 AM",
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("DomainRenew returned %+v, want %+v", result, want)
	}
}
