package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNSCreate(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.ns.create</RequestedCommand>
  <CommandResponse Type="namecheap.domains.ns.create">
    <DomainNSCreateResult Domain="domain.com" Nameserver="ns1.domain.com" IP="12.23.23.23" IsSuccess="true" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.ns.create")
		correctParams.Set("Nameserver", "ns1.domain.com")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		correctParams.Set("IP", "12.23.23.23")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	ns, err := client.NSCreate("domain", "com", "ns1.domain.com", "12.23.23.23")
	if err != nil {
		t.Errorf("NSCreate returned error: %v", err)
	}
	want := &DomainNSCreateResult{
		Domain:     "domain.com",
		Nameserver: "ns1.domain.com",
		IP:         "12.23.23.23",
		IsSuccess:  true,
	}
	if !reflect.DeepEqual(ns, want) {
		t.Errorf("NSCreate returned %+v, want %+v", ns, want)
	}
}

func TestNSDelete(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.ns.delete</RequestedCommand>
  <CommandResponse Type="namecheap.domains.ns.delete">
    <DomainNSDeleteResult Domain="domain.com" Nameserver="ns1.domain.com" IsSuccess="true" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.ns.delete")
		correctParams.Set("Nameserver", "ns1.domain.com")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	ns, err := client.NSDelete("domain", "com", "ns1.domain.com")
	if err != nil {
		t.Errorf("NSDelete returned error: %v", err)
	}
	want := &DomainNSDeleteResult{
		Domain:     "domain.com",
		Nameserver: "ns1.domain.com",
		IsSuccess:  true,
	}
	if !reflect.DeepEqual(ns, want) {
		t.Errorf("NSDelete returned %+v, want %+v", ns, want)
	}
}

func TestNSGetInfo(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.ns.getInfo</RequestedCommand>
  <CommandResponse Type="namecheap.domains.ns.getInfo">
    <DomainNSInfoResult Domain="domain.com" Nameserver="ns1.domain.com" IP="12.23.23.23">
      <NameserverStatuses>
        <Status>OK</Status>
        <Status>Linked</Status>
      </NameserverStatuses>
    </DomainNSInfoResult>
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.ns.getInfo")
		correctParams.Set("Nameserver", "ns1.domain.com")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	ns, err := client.NSGetInfo("domain", "com", "ns1.domain.com")
	if err != nil {
		t.Errorf("NSGetInfo returned error: %v", err)
	}
	want := &DomainNSInfoResult{
		Domain:     "domain.com",
		Nameserver: "ns1.domain.com",
		IP:         "12.23.23.23",
		Statuses:   []string{"OK", "Linked"},
	}
	if !reflect.DeepEqual(ns, want) {
		t.Errorf("NSGetInfo returned %+v, want %+v", ns, want)
	}
}

func TestNSUpdate(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <RequestedCommand>namecheap.domains.ns.update</RequestedCommand>
  <CommandResponse Type="namecheap.domains.ns.update">
    <DomainNSUpdateResult Domain="domain.com" Nameserver="ns1.domain.com" IsSuccess="true" />
  </CommandResponse>
  <Server>SEVER-ONE</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>32.76</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.domains.ns.update")
		correctParams.Set("Nameserver", "ns1.domain.com")
		correctParams.Set("SLD", "domain")
		correctParams.Set("TLD", "com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	ns, err := client.NSUpdate("domain", "com", "ns1.domain.com")
	if err != nil {
		t.Errorf("NSUpdate returned error: %v", err)
	}
	want := &DomainNSUpdateResult{
		Domain:     "domain.com",
		Nameserver: "ns1.domain.com",
		IsSuccess:  true,
	}
	if !reflect.DeepEqual(ns, want) {
		t.Errorf("NSUpdate returned %+v, want %+v", ns, want)
	}
}
