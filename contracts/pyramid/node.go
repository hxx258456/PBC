package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"fmt"
	json "github.com/json-iterator/go"
)

func (p *pyramidContract) saveNode(paramVal []byte) protogo.Response {
	var node = Node{}
	err := node.paramInit(paramVal)
	if err != nil {
		return sdk.Error("params unmarshal failed")
	}
	err = p.insertState(keyNodeID, node.NodeID, node)
	if err != nil {
		return sdk.Error(err.Error())
	}
	err = p.insertState(keyNodeOrgID, fmt.Sprintf("%v_%v", node.OrgID, node.NodeID), node)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(node.NodeID))
}

func (p *pyramidContract) queryNodesByOrgID(oid string) protogo.Response {
	nodes, err := p.getNodes(oid)
	if err != nil {
		return sdk.Error(err.Error())
	}
	result, _ := json.Marshal(nodes)
	return sdk.Success(result)
}

func (p *pyramidContract) getNodes(oid string) ([]Node, error) {
	iter, err := sdk.Instance.NewIteratorPrefixWithKeyField(keyNodeOrgID, oid)
	if err != nil {
		return nil, err
	}
	nodes := make([]Node, 0, 20)
	for iter.HasNext() {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var node Node
		json.Unmarshal(value, &node)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (p *pyramidContract) queryNodesByID(id string) protogo.Response {
	node := Node{}
	result, err := sdk.Instance.GetStateByte(keyNodeID, id)
	if err != nil {
		return sdk.Error(err.Error())
	}
	sdk.Instance.Infof("result:" + string(result))
	err = json.Unmarshal(result, &node)
	if err != nil {
		return sdk.Error("json unmarshal failed,result:" + string(result))
	}
	return sdk.Success(result)
}
