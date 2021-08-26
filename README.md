# Fipper-go-sdk
A client library for Golang (SDK)

Fipper.io - a feature toggle (aka feature flags) software. More info https://fipper.io

## Install
> go get github.com/Fipper-io/fipper-go-sdk

## Example
#### main.go
```go
package main

import (
	"fipper_go_sdk"
	"fmt"
)

func main() {
	client := fipper_go_sdk.FipperClient{Rate: fipper_go_sdk.Rarely}
	config, err := client.GetConfig(
		"production",
		"* place your token here *",
		12345)

	type myJsonScheme struct {
		Test int
	}
	mySchema := myJsonScheme{}

	if err != nil {
		fmt.Printf("Errors: %v", err)
		return
	}

	if flag := config.Flags["bool_flag"]; flag.State {
		val, _ := flag.GetBool()
		fmt.Printf("bool_flag: %v\n", val)
	}

	if flag := config.Flags["int_flag"]; flag.State {
		val, _ := flag.GetInt()
		fmt.Printf("int_flag: %v\n", val)
	}

	if flag := config.Flags["string_flag"]; flag.State {
		val, _ := flag.GetString()
		fmt.Printf("string_flag: %v\n", val)
	}

	if flag := config.Flags["json_flag"]; flag.State {
		flag.GetJson(&mySchema)
		fmt.Printf("json_flag: %v\n", mySchema)
	}
}
```

More information and more client libraries: https://docs.fipper.io