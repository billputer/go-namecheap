package namecheap

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestUsersGetBalances(t *testing.T) {
	setup()
	defer teardown()
	respXML := `<?xml version="1.0" encoding="UTF-8"?>
<ApiResponse Status="OK">
  <Errors />
  <RequestedCommand>namecheap.users.getBalances</RequestedCommand>
  <CommandResponse Type="namecheap.users.getBalances">
    <UserGetBalancesResult Currency="USD" AvailableBalance="4932.96" AccountBalance="4932.96" EarnedAmount="381.70" WithdrawableAmount="1243.36" FundsRequiredForAutoRenew="0.00" />
  </CommandResponse>
  <Server>SERVER-NAME</Server>
  <GMTTimeDifference>+5</GMTTimeDifference>
  <ExecutionTime>0.024</ExecutionTime>
</ApiResponse>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		correctParams := fillDefaultParams(url.Values{})
		correctParams.Set("Command", "namecheap.users.getBalances")
		testBody(t, r, correctParams)
		testMethod(t, r, "POST")
		fmt.Fprint(w, respXML)
	})

	balances, err := client.UsersGetBalances()

	if err != nil {
		t.Errorf("UsersGetBalances returned error: %v", err)
	}

	// Balances response we want, given the above respXML
	want := &UsersGetBalancesResult{
		Currency:                  "USD",
		AvailableBalance:          4932.96,
		AccountBalance:            4932.96,
		EarnedAmount:              381.70,
		WithdrawableAmount:        1243.36,
		FundsRequiredForAutoRenew: 0.00,
	}

	if !reflect.DeepEqual(balances, want) {
		t.Errorf("UsersGetBalances returned %+v, want %+v", balances, want)
	}
}
