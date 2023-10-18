package main

import (
	"chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
	"pyramid/common/code"

	"container/list"
	json "github.com/json-iterator/go"
)

func (p *pyramidContract) saveOrganization(paramVal []byte) protogo.Response {
	var org = Organization{}
	err := org.paramInit(paramVal)
	if err != nil {
		return sdk.Error("params unmarshal failed")
	}
	//check org
	if !p.dagCheck(org) {
		return sdk.Error("dag check failed")
	}
	err = p.insertState(keyOrganizationID, org.OrgId, org)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.Success([]byte(org.OrgId))
}

func (p *pyramidContract) queryOrganizationByOID(oid string) protogo.Response {
	return p.findOrganizationById(oid)
}

func (p *pyramidContract) findOrganizationById(id string) protogo.Response {
	organization := Organization{}
	result, err := sdk.Instance.GetStateByte(keyOrganizationID, id)
	if err != nil {
		return sdk.Error(err.Error())
	}
	sdk.Instance.Infof("result:" + string(result))
	err = json.Unmarshal(result, &organization)
	if err != nil {
		return sdk.Error("json unmarshal failed,result:" + string(result))
	}
	return sdk.Success(result)
}

func (p *pyramidContract) deleteOrganizationById(id string) protogo.Response {
	err := p.delState(keyOrganizationID, id)
	if err != nil {
		return sdk.Error(err.Error())
	}
	return sdk.SuccessResponse
}

func (p *pyramidContract) allRes() map[string][]byte {
	result := make(map[string][]byte)
	iter, err := sdk.Instance.NewIteratorPrefixWithKey(keyOrganizationID)
	if err != nil {
		return nil
	}
	for iter.HasNext() {
		_, field, value, err := iter.Next()
		sdk.Instance.Infof("[allRes] iter get next done,err=%v,field=%v,value=%v", err, field, value)
		if err != nil || len(value) == 0 {
			break
		}
		result[field] = value
	}
	return result
}

func (p *pyramidContract) getAllOrganization() map[string]Organization {
	result := p.allRes()
	m := make(map[string]Organization)
	for oid, val := range result {
		var o Organization
		err := json.Unmarshal(val, &o)
		if err != nil {
			sdk.Instance.Infof("[getAllOrganization] data unmarshal failed,key=%v ;val=%v", oid, string(val))
			continue
		}
		m[oid] = o
	}
	return m
}

func (p *pyramidContract) dagCheck(newOrg Organization) bool {
	//1.获取全量节点
	orgMap := p.getAllOrganization()

	//2.循环,用map存储,key为组织id,val为该组织的下级列表
	inDegreeMap := make(map[string]int)
	leaderMap := make(map[string][]string)
	for orgID, org := range orgMap {
		leaderMap[orgID] = org.LeaderOrges
		if _, ok := inDegreeMap[orgID]; !ok {
			inDegreeMap[orgID] = 0
		}
		for _, leaderOID := range org.LeaderOrges {
			inDegreeMap[leaderOID]++
		}
	}
	//3.map中add新增的组织
	leaderMap[newOrg.OrgId] = newOrg.LeaderOrges
	for _, leaderOID := range newOrg.LeaderOrges {
		if _, ok := orgMap[leaderOID]; !ok {
			continue
		}
		inDegreeMap[leaderOID]++
	}
	//4.循环map,使用拓扑排序观察是否能排序完成,能完成既合规
	queue := list.New()
	for orgID, num := range inDegreeMap {
		if num == 0 {
			queue.PushBack(orgID)
		}
	}
	total := len(inDegreeMap)
	for {
		if queue.Len() == 0 {
			break
		}
		cur := queue.Remove(queue.Front())
		total--
		orgID := cur.(string)

		for _, leaderOID := range leaderMap[orgID] {
			_, ok := inDegreeMap[leaderOID]
			if !ok {
				continue
			}
			inDegreeMap[leaderOID]--
			if inDegreeMap[leaderOID] == 0 {
				queue.PushBack(leaderOID)
			}
		}
	}
	return total == 0
}

func (p *pyramidContract) codeCheck(newOrg Organization) bool {
	c := code.NewCode(newOrg.OrgId)
	return c.IsValid()
}
