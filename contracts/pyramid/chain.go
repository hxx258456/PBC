package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

func (p *pyramidContract) saveChain(paramVal []byte) protogo.Response {
	var chain = ChainList{}
	err := chain.paramInit(paramVal)
	if err != nil {
		return sdk.Error("params unmarshal failed")
	}
	err = p.insertState(keyChainID, chain.ChainId, chain)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(chain.ChainId))
}
