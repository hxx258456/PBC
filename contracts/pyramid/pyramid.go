package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
)

const (
	invokeMethodQuery  = "query"
	invokeMethodInsert = "insert"

	queryParamKeyOrganization = "organization"
	queryParamKeyNode         = "node"
	queryParamKeyUser         = "user"
	queryParamKeyChain        = "chain"

	keyOrganizationID = "OID"
	keyUserID         = "UID"
	keyUserOrgID      = "UOID"
	keyNodeID         = "NID"
	keyNodeOrgID      = "NOID"

	keyChainID = "CID"
)

const (
	methodTypeOrganization = "0"
	methodTypeNode         = "1"
	methodTypeUser         = "2"
	methodTypeChain        = "3"
)

type Log interface {
	Info(format string, val ...string)
	Error(format string, val ...string)
}

type pyramidContract struct {
	//log Log
}

func (p *pyramidContract) InitContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (p *pyramidContract) UpgradeContract() protogo.Response {
	return sdk.Success([]byte("Init contract success"))
}

func (p *pyramidContract) InvokeContract(method string) protogo.Response {
	switch method {
	case "save":
		return p.save()
	case "query":
		return p.query()
	default:
		return sdk.Error("invalid method")
	}
}

func (p *pyramidContract) save() protogo.Response {
	params := sdk.Instance.GetArgs()
	t := string(params["type"])
	switch t {
	case methodTypeOrganization:
		paramOrg, ok := params[queryParamKeyOrganization]
		if ok && len(paramOrg) > 0 {
			return p.saveOrganization(paramOrg)
		}
		return sdk.Error(fmt.Sprintf("invalid param(organization:%v)", paramOrg))
	case methodTypeNode:
		paramNode, ok := params[queryParamKeyNode]
		if ok && len(paramNode) > 0 {
			return p.saveNode(paramNode)
		}
		return sdk.Error(fmt.Sprintf("invalid param(node:%v)", paramNode))
	case methodTypeUser:
		paramUser, ok := params[queryParamKeyUser]
		if ok && len(paramUser) > 0 {
			return p.saveUser(paramUser)
		}
		return sdk.Error(fmt.Sprintf("invalid param(user:%v)", paramUser))
	case methodTypeChain:
		paramChain, ok := params[queryParamKeyChain]
		if ok && len(paramChain) > 0 {
			return p.saveChain(paramChain)
		}
		return sdk.Error(fmt.Sprintf("invalid param(chain:%v)", paramChain))
	default:
		return sdk.Error(fmt.Sprintf("invalid param(type:%v)", t))
	}
}

func (p *pyramidContract) query() protogo.Response {
	params := sdk.Instance.GetArgs()
	t := string(params["type"])
	switch t {
	case methodTypeOrganization:
		orgId := string(params["org_id"])
		return p.findOrganizationById(orgId)
	case methodTypeNode:
		oid := string(params["org_id"])
		nid := string(params["node_id"])
		if oid != "" {
			return p.queryNodesByOrgID(oid)
		}
		if nid != "" {
			return p.queryNodesByID(nid)
		}
		return sdk.Error(fmt.Sprintf("invalid param"))
	case methodTypeUser:
		oid := string(params["org_id"])
		uid := string(params["user_id"])
		if oid != "" {
			return p.queryUsersByOrgID(oid)
		}
		if uid != "" {
			return p.queryUsersByID(uid)
		}
		return sdk.Error(fmt.Sprintf("invalid param"))
	case methodTypeChain:
		paramChain, ok := params[queryParamKeyChain]
		if ok && len(paramChain) > 0 {
			return p.saveChain(paramChain)
		}
		return sdk.Error(fmt.Sprintf("invalid param(chain:%v)", paramChain))
	default:
		return sdk.Error(fmt.Sprintf("invalid param(type:%v)", t))
	}
}

func (p *pyramidContract) delete() protogo.Response {
	params := sdk.Instance.GetArgs()
	orgId := string(params["org_id"])
	return p.deleteOrganizationById(orgId)
}

func (p *pyramidContract) getState(key, field string) ([]byte, error) {
	result, err := sdk.Instance.GetStateByte(key, field)
	return result, err
}

func (p *pyramidContract) getStateExist(key, field string) (bool, error) {
	_, exist, err := sdk.Instance.GetStateWithExists(key, field)
	return exist, err
}

func (p *pyramidContract) insertState(key, field string, metadata interface{}) error {
	b, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errors.New("metadata is empty")
	}
	sdk.Instance.Infof("[save] value=" + string(b))
	return sdk.Instance.PutStateByte(key, field, b)
}

func (p *pyramidContract) delState(key, field string) error {
	return sdk.Instance.DelState(key, field)
}
