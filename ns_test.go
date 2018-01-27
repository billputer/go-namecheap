package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

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
