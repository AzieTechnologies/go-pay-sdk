package main

import (
	"fmt"

	"github.com/dev-pay/go-pay-sdk/payclient"
)

func main() {
	client := payclient.New(&payclient.Config{
		ShareableKey: "Shareable Key",
		Secret:       "Secret Key",
		AccountId:    "Accoutn ID",
		Sandbox:      true})

	var expMap = map[string]string{
		"month": "10",
		"year":  "2024",
	}
	var metaData = map[string]string{
		"client": "golang-sdk",
	}

	intent, err := client.ConfirmPayment(payclient.PaymentDetail{
		Amount:   103,
		Currency: payclient.USD,
		Card: payclient.Card{CardNum: "4037111111000000",
			CardExpiry: expMap, Cvv: "102"},
		BillingAddress: payclient.BillingAddress{Country: "US",
			Zip:    "38138",
			State:  "TN",
			City:   "Memphis",
			Street: "123 ABC Lane"},
		Name:     "Jnix - Golang",
		MetaData: metaData,
	})

	if err != nil {
		fmt.Println("Error - " + err.Error())
	}

	if intent != nil {
		fmt.Printf("\nStatus - %s", intent.Status)
		fmt.Printf("\nAmount - %d", intent.Amount)
	}

}
