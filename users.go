package namecheap

import (
	"net/url"
)

const (
	usersGetPricing  = "namecheap.users.getPricing"
	usersGetBalances = "namecheap.users.getBalances"
)

type UsersGetPricingResult struct {
	ProductType     string `xml:"Name,attr"`
	ProductCategory []struct {
		Name    string `xml:"Name,attr"`
		Product []struct {
			Name  string `xml:"Name,attr"`
			Price []struct {
				Duration     int     `xml:"Duration,attr"`
				DurationType string  `xml:"DurationType,attr"`
				Price        float64 `xml:"Price,attr"`
				RegularPrice float64 `xml:"RegularPrice,attr"`
				YourPrice    float64 `xml:"YourPrice,attr"`
				CouponPrice  float64 `xml:"CouponPrice,attr"`
				Currency     string  `xml:"Currency,attr"`
			} `xml:"Price"`
		} `xml:"Product"`
	} `xml:"ProductCategory"`
}

type UsersGetBalancesResult struct {
	Currency                  string  `xml:"Currency,attr"`
	AvailableBalance          float64 `xml:"AvailableBalance,attr"`
	AccountBalance            float64 `xml:"AccountBalance,attr"`
	EarnedAmount              float64 `xml:"EarnedAmount,attr"`
	WithdrawableAmount        float64 `xml:"WithdrawableAmount,attr"`
	FundsRequiredForAutoRenew float64 `xml:"FundsRequiredForAutoRenew,attr"`
}

func (client *Client) UsersGetPricing(productType string) ([]UsersGetPricingResult, error) {
	requestInfo := &ApiRequest{
		command: usersGetPricing,
		method:  "POST",
		params:  url.Values{},
	}

	requestInfo.params.Set("ProductType", productType)
	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.UsersGetPricing, nil
}

func (client *Client) UsersGetBalances() ([]UsersGetBalancesResult, error) {
	requestInfo := &ApiRequest{
		command: usersGetBalances,
		method:  "POST",
		params:  url.Values{},
	}

	resp, err := client.do(requestInfo)
	if err != nil {
		return nil, err
	}

	return resp.UsersGetBalances, nil
}
