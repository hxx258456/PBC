package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"fmt"
	json "github.com/json-iterator/go"
	"time"
)

type pyramidContract struct {
	//log Log
}

var (
	docType = "notary"
)

const (
	ArgNumErrorf   = "Incorrect number of arguments. Expecting "
	ArgInvalied    = "Invalid parameter"
	DATA_NOT_FOUND = "Data does not exist"
	CCNOTINIT      = "chaincode not init"
	MAXCOUNT       = 1000
)

func (p *pyramidContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (p *pyramidContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (p *pyramidContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "initChaincode":
		return p.initChaincode()
	case "add":
		return p.add()
	case "query":
		return p.query()
	case "queryAll":
		return p.queryAll()
	case "update":
		return p.update()
	case "delete":
		return p.delete()
	case "querysByPagination":
		return p.querysByPagination()
	case "queryLog":
		return p.queryLog()
	case "check":
		return p.check()
	default:
		return sdk.Error("invalid method")
	}
}

func (p *pyramidContract) initChaincode() protogo.Response {
	params := sdk.Instance.GetArgs()
	t := string(params["doc_type"])
	if t == "" {
		return sdk.Error("Invalid parameter")
	}
	docType = t
	return sdk.SuccessResponse
}

func (p *pyramidContract) add() protogo.Response {
	params := sdk.Instance.GetArgs()
	key := string(params["key"])
	value := string(params["value"])
	b, _ := json.Marshal(value)
	err := sdk.Instance.PutStateByte(docType, key, b)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}
func (p *pyramidContract) query() protogo.Response {
	params := sdk.Instance.GetArgs()
	k := string(params["key"])
	b, err := sdk.Instance.GetStateByte(docType, k)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success(b)
}

func (p *pyramidContract) queryAll() protogo.Response {
	params := sdk.Instance.GetArgs()
	key := string(params["key"])
	iter, err := sdk.Instance.NewIteratorPrefixWithKeyField(docType, key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	result := make([][]byte, 1024)
	for iter.HasNext() {
		_, field, value, err := iter.Next()
		sdk.Instance.Infof("[allRes] iter get next done,err=%v,field=%v,value=%v", err, field, value)
		if err != nil || len(value) == 0 {
			break
		}
		result = append(result, value)
	}
	res, _ := json.Marshal(result)
	return sdk.Success(res)
}

func (p *pyramidContract) update() protogo.Response {
	params := sdk.Instance.GetArgs()
	key := string(params["key"])
	value := string(params["value"])
	b, _ := json.Marshal(value)
	_, ok, err := sdk.Instance.GetStateWithExists(docType, key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	if !ok {
		return sdk.Error(DATA_NOT_FOUND)
	}
	err = sdk.Instance.PutStateByte(docType, key, b)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}

func (p *pyramidContract) delete() protogo.Response {
	params := sdk.Instance.GetArgs()
	key := string(params["key"])
	err := sdk.Instance.DelState(docType, key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}

func (p *pyramidContract) querysByPagination() protogo.Response {
	//params := sdk.Instance.GetArgs()
	//t := string(params["type"])
	return sdk.SuccessResponse
}

func (p *pyramidContract) queryLog() protogo.Response {
	params := sdk.Instance.GetArgs()
	key := string(params["key"])
	iter, err := sdk.Instance.NewHistoryKvIterForKey(docType, key)
	if err != nil {
		return sdk.Error(err.Error())
	}
	result := make([]LogResult, 1024)
	for iter.HasNext() {
		km, err := iter.Next()
		if err != nil {

		}
		loc, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation("2006-01-02 15:04:05", km.Timestamp, loc)
		if err != nil {
			fmt.Println("km.timestamp parse to time failed:", km.Timestamp)
		}
		result = append(result, LogResult{
			Record:    string(km.Value),
			TxId:      km.TxId,
			Timestamp: t,
			IsDelete:  km.IsDelete,
		})
	}
	res, _ := json.Marshal(result)
	return sdk.Success(res)
}

func (p *pyramidContract) check() protogo.Response {
	params := sdk.Instance.GetArgs()
	paramKey := string(params["key"])
	paramCompares := params["compares"]
	paramContent := params["content"]
	paramCheckKey := string(params["checkKey"])

	txId, _ := sdk.Instance.GetTxId()

	compares := map[string]string{}
	if err := json.Unmarshal(paramCompares, &compares); err != nil {
		return sdk.Error(err.Error())
	}
	content := map[string]interface{}{}
	if err := json.Unmarshal(paramContent, &content); err != nil {
		return sdk.Error(err.Error())
	}
	contentByte, err := sdk.Instance.GetStateByte(docType, paramKey)
	if err != nil {
		return sdk.Error(err.Error())
	}
	ccContent := map[string]interface{}{} //链上记录
	if err := json.Unmarshal(contentByte, &ccContent); err != nil {
		return sdk.Error(err.Error())
	}
	content["txId"] = txId
	content["checkResult"] = "1"
	for k, v := range compares {
		if content[k] != ccContent[v] {
			content["checkResult"] = "0"
		}
	}
	dataByte, _ := json.Marshal(&content)

	err = sdk.Instance.PutStateByte(docType, paramCheckKey, dataByte)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}
