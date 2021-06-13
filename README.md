# go-pay-sdk
A Golang SDK for Devpay Payment Gateway  Get your API Keys at https://devpay.io

# Install
```go
go mod download github.com/dev-pay/go-pay-sdk/payclient
```

# Usage 

```golang

import "github.com/dev-pay/go-pay-sdk/payclient"

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
	"property1": "value1",
}

intent, err := client.ConfirmPayment(payclient.PaymentDetail{
	Amount:   Amount,
	Currency: payclient.USD,
	Card: payclient.Card{CardNum: "Card_Number",
		CardExpiry: expMap, Cvv: "CVV"},
	BillingAddress: payclient.BillingAddress{Country: "US",
		Zip:    "38138",
		State:  "TN",
		City:   "Memphis",
		Street: "123 ABC Lane"},
	Name:     "John",
	MetaData: metaData,
})


```

# Demo
Please refer example code [here](https://github.com/dev-pay/go-pay-sdk/tree/main/example), follow below steps to run the example code
1. Download the code
2. cd to `go-pay-sdk/example`
3. Run go mod download
4. Update inputs in maing.go 
5. Run the file `gp run main.go`

# License
Refer [LICENSE](https://github.com/dev-pay/go-pay-sdk/blob/main/LICENSE) file
