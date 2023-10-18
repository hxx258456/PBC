package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	json "github.com/json-iterator/go"
)

func (p *pyramidContract) saveUser(paramVal []byte) protogo.Response {
	var user = User{}
	err := user.paramInit(paramVal)
	if err != nil {
		return sdk.Error("params unmarshal failed")
	}
	err = p.insertState(keyUserID, user.UserID, user)
	if err != nil {
		return sdk.Error(err.Error())
	}
	err = p.insertState(keyUserOrgID, user.OrgID, user)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(user.UserID))
}

func (p *pyramidContract) queryUsersByOrgID(uid string) protogo.Response {
	users, err := p.getUsers(uid)
	if err != nil {
		return sdk.Error(err.Error())
	}
	result, _ := json.Marshal(users)
	return sdk.Success(result)
}

func (p *pyramidContract) getUsers(uid string) ([]User, error) {
	iter, err := sdk.Instance.NewIteratorPrefixWithKeyField(keyUserOrgID, uid)
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, 20)
	for iter.HasNext() {
		_, _, value, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var node User
		json.Unmarshal(value, &node)
		users = append(users, node)
	}
	return users, nil
}

func (p *pyramidContract) queryUsersByID(id string) protogo.Response {
	node := User{}
	result, err := sdk.Instance.GetStateByte(keyUserID, id)
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
