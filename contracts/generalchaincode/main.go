package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

func main() {
	err := sandbox.Start(new(pyramidContract))
	if err != nil {
		sdk.Instance.Errorf(err.Error())
	}
}
