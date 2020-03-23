package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestDomainsDNSGetHosts(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.dns.getHosts</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.getHosts">
    <DomainDNSGetHostsResult Domain="domain.com" IsUsingOurDNS="true">
      <host HostId="12" Name="@" Type="A" Address="1.2.3.4" MXPref="10" TTL="1800" />
      <host HostId="14" Name="www" Type="A" Address="122.23.3.7" MXPref="10" TTL="1800" />
    </DomainDNSGetHostsResult>
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.dns.getHosts")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	hosts, err := client.DomainsDNSGetHosts("domain", "com")
	if err != nil {
		t.Errorf("DomainsDNSGetHosts returned error: %v", err)
	}

	want := &DomainDNSGetHostsResult{
		Domain:        "domain.com",
		IsUsingOurDNS: true,
		Hosts: []DomainDNSHost{
			{
				ID:      12,
				Name:    "@",
				Type:    "A",
				Address: "1.2.3.4",
				MXPref:  10,
				TTL:     1800,
			},
			{
				ID:      14,
				Name:    "www",
				Type:    "A",
				Address: "122.23.3.7",
				MXPref:  10,
				TTL:     1800,
			},
		},
	}

	if !reflect.DeepEqual(hosts, want) {
		t.Errorf("DomainsDNSGetHosts returned %+v, want %+v", hosts, want)
	}
}

func TestDomainsDNSSetHosts(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.setHosts">
    <DomainDNSSetHostsResult Domain="domain51.com" IsSuccess="true" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.dns.setHosts")
		correctParams.Set("Address1", "http://www.namecheap.com")
		correctParams.Set("HostName1", "@")
		correctParams.Set("RecordType1", "URL")
		correctParams.Set("TTL1", "100")
		correctParams.Set("SLD", "domain51")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	hosts := []DomainDNSHost{
		{
			Name:    "@",
			Type:    "URL",
			Address: "http://www.namecheap.com",
			TTL:     100,
		},
	}

	result, err := client.DomainDNSSetHosts("domain51", "com", hosts)

	if err != nil {
		t.Errorf("DomainsDNSGetHosts returned error: %v", err)
	}

	want := &DomainDNSSetHostsResult{
		Domain:    "domain51.com",
		IsSuccess: true,
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("DomainsDNSSetHosts returned %+v, want %+v", hosts, want)
	}
}

func TestDomainsDNSSetCustom(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.dns.setCustom</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.setCustom">
    <DomainDNSSetCustomResult Domain="domain.com" Update="true" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.dns.setCustom")
		correctParams.Set("Nameservers", "dns1.name-servers.com,dns2.name-servers.com")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	result, err := client.DomainDNSSetCustom("domain", "com", "dns1.name-servers.com,dns2.name-servers.com")
	if err != nil {
		t.Errorf("DomainDNSSetCustom returned error: %v", err)
	}

	want := &DomainDNSSetCustomResult{
		Domain: "domain.com",
		Update: true,
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("DomainsDNSSetCustom returned %+v, want %+v", result, want)
	}
}

func TestDomainsDNSSetDefault(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.dns.setDefault</RequestedCommand>
  <CommandResponse Type="namecheap.domains.dns.setDefault">
    <DomainDNSSetDefaultResult Domain="domain.com" Updated="true" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.dns.setDefault")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	result, err := client.DomainDNSSetDefault("domain", "com")
	if err != nil {
		t.Errorf("DomainDNSSetDefault returned error: %v", err)
	}

	want := &DomainDNSSetDefaultResult{
		Domain: "domain.com",
		Update: true,
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("DomainsDNSSetDefault returned %+v, want %+v", result, want)
	}
}
