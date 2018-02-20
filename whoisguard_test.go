package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestWhoisguardGetList(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <Warnings />
  <RequestedCommand>namecheap.whoisguard.getList</RequestedCommand>
  <CommandResponse Type="namecheap.whoisguard.getList">
    <WhoisguardGetListResult>
      <Whoisguard ID="34401" DomainName="" Created="12/18/2013" Expires="12/18/2014" Status="unused" />
      <Whoisguard ID="34400" DomainName="test.com" Created="12/26/2013" Expires="12/26/2014" Status="enabled" />
    </WhoisguardGetListResult>
    <Paging>
      <TotalItems>642</TotalItems>
      <CurrentPage>1</CurrentPage>
      <PageSize>20</PageSize>
    </Paging>
  </CommandResponse>
  <Server>API01</Server>
  <GMTTimeDifference>--5:00</GMTTimeDifference>
  <ExecutionTime>0.029</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.whoisguard.getList")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	list, err := client.WhoisguardGetList()
	if err != nil {
		t.Errorf("WhoisguardGetList returned error: %v", err)
	}

	// WhoisguardGetListResult we expect, given the respXML above
	want := []WhoisguardGetListResult{
		WhoisguardGetListResult{
			ID:      34401,
			Created: "12/18/2013",
			Expires: "12/18/2014",
			Status:  "unused",
		},
		WhoisguardGetListResult{
			ID:         34400,
			DomainName: "test.com",
			Created:    "12/26/2013",
			Expires:    "12/26/2014",
			Status:     "enabled",
		},
	}

	if !reflect.DeepEqual(list, want) {
		t.Errorf("WhoisguardGetList returned %+v, want %+v", list, want)
	}
}

func TestWhoisguardEnable(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <Warnings />
  <RequestedCommand>namecheap.whoisguard.enable</RequestedCommand>
  <CommandResponse Type="namecheap.whoisguard.enable">
    <WhoisguardEnableResult DomainName="domain1.com" IsSuccess="true" />
  </CommandResponse>
  <Server>API02</Server>
  <GMTTimeDifference>--5:00</GMTTimeDifference>
  <ExecutionTime>0.92</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.whoisguard.enable")
		correctParams.Set("WhoisguardID", "34401")
		correctParams.Set("ForwardedToEmail", "john@test.com")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	err := client.WhoisguardEnable(34401, "john@test.com")
	if err != nil {
		t.Errorf("WhoisguardEnable returned error: %v", err)
	}
}

func TestWhoisguardDisable(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <Warnings />
  <RequestedCommand>namecheap.whoisguard.disable</RequestedCommand>
  <CommandResponse Type="namecheap.whoisguard.disable">
    <WhoisguardDisableResult DomainName="domain1.com" IsSuccess="true" />
  </CommandResponse>
  <Server>API02</Server>
  <GMTTimeDifference>--5:00</GMTTimeDifference>
  <ExecutionTime>0.92</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.whoisguard.disable")
		correctParams.Set("WhoisguardID", "34401")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	err := client.WhoisguardDisable(34401)
	if err != nil {
		t.Errorf("WhoisguardDisable returned error: %v", err)
	}
}

func TestWhoisguardRenew(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse xmlns="http://api.namecheap.com/xml.response" Status="OK">
  <Errors />
  <Warnings />
  <RequestedCommand>namecheap.whoisguard.renew</RequestedCommand>
   <CommandResponse Type="namecheap.whoisguard.renew">
      <WhoisguardRenewResult WhoisguardId="38495" Years="1" Renew="true" OrderId="580938" TransactionId="884255" ChargedAmount="6.8000"/>
   </CommandResponse>
  <Server>API01</Server>
  <GMTTimeDifference>--5:00</GMTTimeDifference>
  <ExecutionTime>0.029</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.whoisguard.renew")
		correctParams.Set("WhoisguardID", "38495")
		correctParams.Set("Years", "1")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	result, err := client.WhoisguardRenew(38495, 1)
	if err != nil {
		t.Errorf("WhoisguardRenew returned error: %v", err)
	}

	// DomainCheckResult we expect, given the respXML above
	want := &WhoisguardRenewResult{
		WhoisguardID:  38495,
		Renewed:       true,
		ChargedAmount: 6.8,
		TransactionID: 884255,
		OrderID:       580938,
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("WhoisguardRenew returned %+v, want %+v", result, want)
	}
}
