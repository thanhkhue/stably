package fee

import (
	"errors"
	"fmt"
	// "log"
	"testing"
)

func main() {
	// defer runtime.Goexit()
	t := &testing.T{}
	// TestCalculateFees(t)
}

var NetworkFeeCharge = map[string]float64{
	"Bitcoin":  10,
	"Ethereum": 18,
	"Solana":   2,
}

type LiquidityProvider struct {
	Name string
	Fee  float64
	Type string // "ADD" , "MULTI"
}

var liquidityProviders = []LiquidityProvider{
	{Name: "Duck Provider", Fee: 4, Type: "ADD"},
	{Name: "Goose Provider", Fee: 2, Type: "ADD"},
	{Name: "Fox Provider", Fee: 0.1, Type: "MULTI"},
}

/*
//		FromAmount: "100",
//	  FromNetwork: "ACH",
//	  FromAsset: "USD",
//	  ToNetwork: "ethereum",
//	  ToAsset: "ETH",
*/
type CalcFeeRequest struct {
	FromAmount     float64 `json:"FromAmount"`
	FromNetwork    string  `json:"FromNetwork"`
	FromAsset      string  `json:"FromAsset"`
	ToNetwork      string  `json:"ToNetwork"`
	ToAsset        string  `json:"ToAsset"`
	FromCustomerID string
	IsTesting      bool
}

type CalcFeeInternalRequest struct {
	Req              *CalcFeeRequest
	NetworkChargeFee float64
	Tier             string
}

type FeeResonse struct {
	FeeUSD   float64 `json:"FeeUSD"`
	Provider string  `json:"Provider"`
	Error    error
}

// type service struct {
// 	httpClient *http.Client
// }

func CalculateFeesInternal(
	req *CalcFeeRequest, networkChargeFee float64, tier string,
) (*FeeResonse, error) {
	fiatPaymentNetworkFee, err := GetFiatPaymentNetwork(req.FromNetwork, req.FromAmount)
	if err != nil {
		// log.Printf("Get Fiat Payment Network Fee Error")
		return nil, err
		// return nil, err
	}

	CalcFeePerCustomerTier(
		float64(networkChargeFee), fiatPaymentNetworkFee, tier,
		req.FromNetwork, req.ToNetwork,
	)
	return &FeeResonse{
		Provider: "Duck",
		FeeUSD:   12,
	}, nil
}

func CalculateFees(req *CalcFeeRequest) (*FeeResonse, error) {
	networkChargeFee := GetNetworkChargedFees(req.ToNetwork)
	customerTier := GetCustomerTier(req.FromCustomerID)
	return CalculateFeesInternal(req, networkChargeFee, customerTier)
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func GetFiatPaymentNetwork(fiatPaymentNetwork string, amount float64) (float64, error) {
	switch fiatPaymentNetwork {
	case "ACH":
		return max(2.0, amount*0.015), nil
	case "Instant ACH":
		return max(3.0, amount*0.02), nil
	case "USD Balance":
		return 0.5, nil
	}
	return 0, errors.New("Invalid Fiat Payment Network")
}

func CalcFeePerCustomerTier(
	cryptorFee, fiatPaymentFee float64, tier, network, fiatPaymentNetwork string,
) (float64, float64) {
	switch tier {
	case "tier 1":
		return cryptorFee, fiatPaymentFee
	case "tier 2":
		return 0.25 * cryptorFee, 0.25 * fiatPaymentFee
	case "tier 3":
		return 0.5 * cryptorFee, 0.25 * fiatPaymentFee
	case "tier 4":
		if network == "Ethereum" {
			cryptorFee *= 0.5
		} else {
			cryptorFee = 0
		}
		if fiatPaymentNetwork == "ACH" {
			fiatPaymentFee *= 0.25
		} else {
			fiatPaymentFee *= 0.5
		}
		return cryptorFee, fiatPaymentFee
	}
	return cryptorFee, fiatPaymentFee
}

func GetCustomerTier(customerID string) string {
	return "tier 2"
}

func GetFiatNetworkChargedFees(currentFee float64) float64 {
	return 0.0
}

func GetNetworkChargedFees(network string) float64 {
	chargeFee, found := NetworkFeeCharge[network]
	if found {
		return chargeFee
	}

	// flag-service

	// call request 3rd party
	return GetNetworkChargedFeesFrom3rdParty()
}

// call request 3rd party
func GetNetworkChargedFeesFrom3rdParty() float64 {
	return 7
}

/*
// Transaction
{
  FromAmount: "100",
  FromNetwork: "ACH",
  FromAsset: "USD",
  ToNetwork: "ethereum",
  ToAsset: "ETH",
}


// Customer
{
  Tier: 3,
}


// AvailableProviders
[
  "Duck",
  "Goose",
]

*/

// Output
/*
{
  FeeUSD: 12.5,
  Provider: "Goose",
}

*/

type TestCase struct {
	Name      string
	Input     CalcFeeInternalRequest
	ExpectRes *FeeResonse
}
