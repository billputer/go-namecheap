package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestSslGetList(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
	<?xml version="1.0" encoding="UTF-8"?>
	<ApiResponse Status="OK">
	  <Errors />
	  <RequestedCommand>namecheap.ssl.getList</RequestedCommand>
	  <CommandResponse Type="namecheap.ssl.getList">
		<SSLListResult>
		  <SSL CertificateID="52556" HostName="domainxy.com" SSLType="SSLCertificate3" PurchaseDate="10/17/2006" ExpireDate="10/17/2008" ActivationExpireDate="12/31/2009" IsExpiredYN="false" Status="new" />
		</SSLListResult>
		<Paging>
		  <TotalItems>3</TotalItems>
		  <CurrentPage>1</CurrentPage>
		  <PageSize>20</PageSize>
		</Paging>
	  </CommandResponse>
	  <Server>SERVER</Server>
	  <GMTTimeDifference>+5:30</GMTTimeDifference>
	  <ExecutionTime>1.094</ExecutionTime>
	</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.ssl.getList")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	certificates, err := client.SslGetList()

	if err != nil {
		t.Errorf("SslGetList returned error: %v", err)
	}

	// DomainGetListResult we expect, given the respXML above
	want := []SslGetListResult{{
		CertificateID:        52556,
		HostName:             "domainxy.com",
		SSLType:              "SSLCertificate3",
		PurchaseDate:         "10/17/2006",
		ExpireDate:           "10/17/2008",
		ActivationExpireDate: "12/31/2009",
		IsExpired:            false,
		Status:               "new",
	}}

	if !reflect.DeepEqual(certificates, want) {
		t.Errorf("SslGetList returned %+v, want %+v", certificates, want)
	}
}

func TestSslCreate(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
	<?xml version="1.0" encoding="UTF-8"?>
	<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
	<Errors/>
	<Warnings/>
	<RequestedCommand>namecheap.ssl.create</RequestedCommand>
	<CommandResponse Type="namecheap.ssl.create">
		<SSLCreateResult IsSuccess="true" OrderId="1234567" TransactionId="1234567" ChargedAmount="908.1600">
			<SSLCertificate CertificateID="123456" Created="02/20/2018" SSLType="PositiveSSL" Years="2" Status="NewPurchase"/>
		</SSLCreateResult>
	</CommandResponse>
	<Server>202005e9484c</Server>
	<GMTTimeDifference>--5:00</GMTTimeDifference>
	<ExecutionTime>2.608</ExecutionTime>
	</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.ssl.create")
		correctParams.Set("Type", "PositiveSSL")
		correctParams.Set("Years", "2")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	certificates, err := client.SslCreate("PositiveSSL", 2)

	if err != nil {
		t.Errorf("SslCreate returned error: %v", err)
	}

	// SslCreateResult we expect, given the respXML above
	wantCerts := []SSLCertificate{{
		CertificateID: 123456,
		SSLType:       "PositiveSSL",
		Created:       "02/20/2018",
		Years:         2,
		Status:        "NewPurchase",
	}}

	want := &SslCreateResult{
		IsSuccess:      true,
		OrderId:        1234567,
		TransactionId:  1234567,
		ChargedAmount:  908.1600,
		SSLCertificate: wantCerts,
	}

	if !reflect.DeepEqual(certificates, want) {
		t.Errorf("SslCreate returned %+v, want %+v", certificates, want)
	}
}
