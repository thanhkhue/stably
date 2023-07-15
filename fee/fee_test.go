package fee

import (
	"errors"
	"fmt"
	// "log"
	"testing"
)

func TestCalculateFees(t *testing.T) {
	testcases := []TestCase{
		{
			Name: "test with invalid fiat payment network",
			Input: CalcFeeInternalRequest{
				Req: &CalcFeeRequest{
					FromAmount:  100,
					FromNetwork: "Stripe",
					FromAsset:   "USD",
					ToNetwork:   "ethereum",
					ToAsset:     "ETH",
				},
				NetworkChargeFee: 1,
				Tier:             "Tier 1",
			},
			ExpectRes: &FeeResonse{
				Error: errors.New("Invalid Fiat Payment Network"),
			},
		},
		{
			Name: "Case buy Ethereum with ACH",
			Input: CalcFeeInternalRequest{
				Req: &CalcFeeRequest{
					FromAmount:  100,
					FromNetwork: "ACH",
					FromAsset:   "USD",
					ToNetwork:   "ethereum",
					ToAsset:     "ETH",
				},
				NetworkChargeFee: 1,
				Tier:             "Tier 1",
			},
			ExpectRes: &FeeResonse{
				FeeUSD:   12.5,
				Provider: "Goose",
			},
		},
	}
	for _, tc := range testcases {
		res, err := CalculateFeesInternal(tc.Input.Req, tc.Input.NetworkChargeFee, tc.Input.Tier)
		if tc.ExpectRes.Error != nil && tc.ExpectRes.Error.Error() != err.Error() {
			// t.Error("Expect %+v, but got %+v\n", tc.ExpectRes.Error, err)
			fmt.Printf("Expect %+v, but got %+v\n", tc.ExpectRes.Error, err)
			t.Fatalf("Expect %+v, but got %+v\n", tc.ExpectRes.Error, err)
			continue
		}
		if res != nil && res.Provider != tc.ExpectRes.Provider {
			// t.Fatalf()
			fmt.Printf("Expect %s, but got %s\n", tc.ExpectRes.Provider, res.Provider)
			t.Fatalf("Expect %s, but got %s", tc.ExpectRes.Provider, res.Provider)
		}
		if res != nil && res.FeeUSD != tc.ExpectRes.FeeUSD {
			fmt.Printf("Expect %f, but got %f\n", tc.ExpectRes.FeeUSD, res.FeeUSD)
			t.Fatalf("Expect %f, but got %f", tc.ExpectRes.FeeUSD, res.FeeUSD)
		}
	}
}
