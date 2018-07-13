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
	want := &SslCreateResult{
		IsSuccess:     true,
		OrderId:       1234567,
		TransactionId: 1234567,
		ChargedAmount: 908.1600,
		SSLCertificate: []SSLCertificate{{
			CertificateID: 123456,
			SSLType:       "PositiveSSL",
			Created:       "02/20/2018",
			Years:         2,
			Status:        "NewPurchase",
		}},
	}

	if !reflect.DeepEqual(certificates, want) {
		t.Errorf("SslCreate returned %+v, want %+v", certificates, want)
	}
}

func TestSslActivateHttpDCValidation(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
	<?xml version="1.0" encoding="UTF-8"?>
	<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
	<Errors/>
	<Warnings/>
	<RequestedCommand>namecheap.ssl.activate</RequestedCommand>
	<CommandResponse Type="namecheap.ssl.activate">
		<SSLActivateResult ID="953413" IsSuccess="true">
		<HttpDCValidation ValueAvailable="true">
			<DNS domain="test.example.org">
				<FileName><![CDATA[4E3324A380B58813D5A2F32AA13A96F0.txt]]></FileName>
				<FileContent><![CDATA[6694010FAC8ED8F806F1EAD56A1A0478DE6620A256BB8C356A8DD2146B00E884 comodoca.com 5a955211b1f8c]]></FileContent>
			</DNS>
		</HttpDCValidation>
		</SSLActivateResult>
	</CommandResponse>
	<Server>5eda89c931f6</Server>
	<GMTTimeDifference>--5:00</GMTTimeDifference>
	<ExecutionTime>2.227</ExecutionTime>
	</ApiResponse>`

	csrContent := `-----BEGIN CERTIFICATE REQUEST-----
	MIICyjCCAbICAQAwgYQxCzAJBgNVBAYTAkRFMQkwBwYDVQQIEwAxETAPBgNVBAcT
	CGlyZ2VuZHdvMQ8wDQYDVQQKEwZNQVggQUcxCTAHBgNVBAsTADEZMBcGA1UEAxMQ
	dGVzdC5leGFtcGxlLm9yZzEgMB4GCSqGSIb3DQEJARYRYWRtaW5AZXhhbXBsZS5v
	cmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDsYAf6QorCNP4+bbyX
	RoVHcx5zq37Qc7SzRH3Jus9i/zjINT+2Yq0rAKgyiJ2Z1duBl3fNoDS64KRNB15a
	v/d1aH5XBk3motdVxuPcX3/3a6yEepfew6eb2gWI/1J0v9OC3bPzNQB+EEXs0P4E
	wKhdG3+Qxp2XV8EHvdoh0da+kE9mvxlTyqSnkI/03Awu/iHJq7UChNgG3ElmM3qV
	ybqItYnzvi1iZ/gU0l5RrCkj3/uCc8ODnrMM6QeTM3FbVKtEF3b6O+iTRn4uz0LJ
	dKODzxSok9fUD8/FKSzHKwAxo4gmYpR1yIvbuHRPhekoP+bdelhySn5JeZnR1iEb
	dfNBAgMBAAGgADANBgkqhkiG9w0BAQsFAAOCAQEA6Xei1GBkTxqBqzu6QDft9d48
	J5ID4TU3U2piLJVkbjUDBPpkk5TRZWkUG/0PKZopd0c5ujzBJCx37ipsyU+T9g5i
	BEcoEzCPE+zlg9nTsMpNZVR17sBoM2xNkyHdytormrCYrAtu/E43Fymg8Fp8ygqQ
	/UvEww4vnadnLxNYitb7HeaG0QN+XlP3vt3uXW2HxZL9fpsQV93TQXZ5w5+B3mg4
	nnS+Y+N/O3nd4fcsQlIt7//mb5Ikd+txuAUYJRdm7bQMn1MN/Jef4slw4tP0KZA1
	v5DDv8p49Ae+08d0TTFRViMBI6sTHJ+AqF5vep0R4GWOsbdUjG/wiJhpyMLOGQ==
	-----END CERTIFICATE REQUEST-----`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.ssl.activate")
		correctParams.Set("CertificateID", "953413")
		correctParams.Set("CSR", csrContent)
		correctParams.Set("AdminEmailAddress", "admin@example.org")
		correctParams.Set("WebServerType", "nginx")
		correctParams.Set("HTTPDCValidation", "true")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	sslActivateparams := SslActivateParams{
		CertificateId:      953413,
		Csr:                csrContent,
		AdminEmailAddress:  "admin@example.org",
		WebServerType:      "nginx",
		IsHTTPDCValidation: true,
	}

	certificates, err := client.SslActivate(sslActivateparams)

	if err != nil {
		t.Errorf("SslActivate returned error: %v", err)
	}

	// SslActivateResult we expect, given the respXML above
	want := &SslActivateResult{
		ID:        953413,
		IsSuccess: true,
		HttpDCValidation: SslDcValidation{
			ValueAvailable: true,
			Dns: SslDns{
				Domain:      "test.example.org",
				FileName:    "4E3324A380B58813D5A2F32AA13A96F0.txt",
				FileContent: "6694010FAC8ED8F806F1EAD56A1A0478DE6620A256BB8C356A8DD2146B00E884 comodoca.com 5a955211b1f8c",
			},
		},
	}

	if !reflect.DeepEqual(certificates, want) {
		t.Errorf("SslActivate returned %+v, want %+v", certificates, want)
	}
}

func TestSslActivateDNSDCValidation(t *testing.T) {
	setup()
	defer teardown()

	respXML := `
	<?xml version="1.0" encoding="UTF-8"?>
	<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
	<Errors/>
	<Warnings/>
	<RequestedCommand>namecheap.ssl.activate</RequestedCommand>
	<CommandResponse Type="namecheap.ssl.activate">
		<SSLActivateResult ID="953413" IsSuccess="true">
		<DNSDCValidation ValueAvailable="true">
			<DNS domain="test.example.org">
				<HostName><![CDATA[_4E3324A380B58813D5A2F32AA13A96F0.test.example.org]]></HostName>
				<Target><![CDATA[6694010FAC8ED8F806F1EAD56A1A0478.DE6620A256BB8C356A8DD2146B00E884.5a955211b1f8c.comodoca.com]]></Target>
			</DNS>
		</DNSDCValidation>
		</SSLActivateResult>
	</CommandResponse>
	<Server>5eda89c931f6</Server>
	<GMTTimeDifference>--5:00</GMTTimeDifference>
	<ExecutionTime>2.227</ExecutionTime>
	</ApiResponse>`

	csrContent := `-----BEGIN CERTIFICATE REQUEST-----
	MIICyjCCAbICAQAwgYQxCzAJBgNVBAYTAkRFMQkwBwYDVQQIEwAxETAPBgNVBAcT
	CGlyZ2VuZHdvMQ8wDQYDVQQKEwZNQVggQUcxCTAHBgNVBAsTADEZMBcGA1UEAxMQ
	dGVzdC5leGFtcGxlLm9yZzEgMB4GCSqGSIb3DQEJARYRYWRtaW5AZXhhbXBsZS5v
	cmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDsYAf6QorCNP4+bbyX
	RoVHcx5zq37Qc7SzRH3Jus9i/zjINT+2Yq0rAKgyiJ2Z1duBl3fNoDS64KRNB15a
	v/d1aH5XBk3motdVxuPcX3/3a6yEepfew6eb2gWI/1J0v9OC3bPzNQB+EEXs0P4E
	wKhdG3+Qxp2XV8EHvdoh0da+kE9mvxlTyqSnkI/03Awu/iHJq7UChNgG3ElmM3qV
	ybqItYnzvi1iZ/gU0l5RrCkj3/uCc8ODnrMM6QeTM3FbVKtEF3b6O+iTRn4uz0LJ
	dKODzxSok9fUD8/FKSzHKwAxo4gmYpR1yIvbuHRPhekoP+bdelhySn5JeZnR1iEb
	dfNBAgMBAAGgADANBgkqhkiG9w0BAQsFAAOCAQEA6Xei1GBkTxqBqzu6QDft9d48
	J5ID4TU3U2piLJVkbjUDBPpkk5TRZWkUG/0PKZopd0c5ujzBJCx37ipsyU+T9g5i
	BEcoEzCPE+zlg9nTsMpNZVR17sBoM2xNkyHdytormrCYrAtu/E43Fymg8Fp8ygqQ
	/UvEww4vnadnLxNYitb7HeaG0QN+XlP3vt3uXW2HxZL9fpsQV93TQXZ5w5+B3mg4
	nnS+Y+N/O3nd4fcsQlIt7//mb5Ikd+txuAUYJRdm7bQMn1MN/Jef4slw4tP0KZA1
	v5DDv8p49Ae+08d0TTFRViMBI6sTHJ+AqF5vep0R4GWOsbdUjG/wiJhpyMLOGQ==
	-----END CERTIFICATE REQUEST-----`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.ssl.activate")
		correctParams.Set("CertificateID", "953413")
		correctParams.Set("CSR", csrContent)
		correctParams.Set("AdminEmailAddress", "admin@example.org")
		correctParams.Set("WebServerType", "nginx")
		correctParams.Set("DNSDCValidation", "true")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	sslActivateparams := SslActivateParams{
		CertificateId:     953413,
		Csr:               csrContent,
		AdminEmailAddress: "admin@example.org",
		WebServerType:     "nginx",
		IsDNSDCValidation: true,
	}

	certificates, err := client.SslActivate(sslActivateparams)

	if err != nil {
		t.Errorf("SslActivate returned error: %v", err)
	}

	// SslActivateResult we expect, given the respXML above
	want := &SslActivateResult{
		ID:        953413,
		IsSuccess: true,
		DNSDCValidation: SslDcValidation{
			ValueAvailable: true,
			Dns: SslDns{
				Domain:   "test.example.org",
				HostName: "_4E3324A380B58813D5A2F32AA13A96F0.test.example.org",
				Target:   "6694010FAC8ED8F806F1EAD56A1A0478.DE6620A256BB8C356A8DD2146B00E884.5a955211b1f8c.comodoca.com",
			},
		},
	}

	if !reflect.DeepEqual(certificates, want) {
		t.Errorf("SslActivate returned %+v, want %+v", certificates, want)
	}
}
