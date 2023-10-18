package main

import (
	json "github.com/json-iterator/go"
)

const (
	paramCreateOrganization = ""
	paramCreateUser         = ""
	paramCreateNode         = ""
)

type Organization struct {
	OrgId       string   `json:"org_id"`       // 组织id,使用usci编码
	OrgName     string   `json:"org_name"`     // 组织名称
	Algorithm   int      `json:"algorithm"`    // 算法 0 sm2 HASH_TYPE_SHA256 1 ecdsa HASH_TYPE_SM3
	CaType      int      `json:"ca_type"`      // 证书模式 0 single 1 Double
	LeaderOrges []string `json:"leader_orges"` //

	Cert Cert `json:"cert"`
}

type User struct {
	UserID string `json:"user_id"` //用户标识id
	Name   string `json:"name"`    //用户名
	NameEN string `json:"name_en"` //用户名拼音
	OrgID  string `json:"org_id"`  //所属组织id
	Class  string `json:"class"`   //用户类型

	Cert Cert `json:"cert"`
}

type Node struct {
	NodeID string `json:"node_id"`
	OrgID  string `json:"org_id"`

	Cert Cert `json:"cert"`
}

type Cert struct {
	CertType int `json:"cert_type"` // 证书类型 6：用户light证书，
	// 5：普通节点证书，4：共识节点证书，3：用户client证书，2：用户admin证书, 1：ca证书， 0：根证书
	CertUse      int    `json:"cert_use"`       // 证书用途 1：tls，0：签名 2: pem
	Cert         string `json:"cert"`           //证书值
	PrivateKey   string `json:"private_key"`    //私钥值
	PublicKey    string `json:"public_key"`     //公钥值
	OrgId        string `json:"org_id"`         // 组织id
	OrgName      string `json:"org_name"`       // 组织名称
	CertUserName string `json:"cert_user_name"` // 证书用户名
	NodeName     string `json:"node_name"`      // 节点名
	Algorithm    int    `json:"algorithm"`      // 0:国密 1:非国密
	Addr         string `json:"addr"`           // 地址
	RemarkName   string `json:"remark_name"`    // 账户备注名
	ChainMode    string `json:"chain_mode"`     // 链账户类型
}

type ChainList struct {
	ChainId        string `json:"chain_id"`        //区块链id
	Name           string `json:"name"`            //区块链名
	ChainMode      string `json:"chain_mode"`      //链账户模式
	Algorithm      int    `json:"algorithm"`       // 0:国密 1:非国密
	ConsensusType  int    `json:"consensus_type"`  //共识类型 0-SOLO,1-TBFT,2-MBFT,3-MAXBFT,4-RAFT,10-POW
	ConsensusNodes []Node `json:"consensus_nodes"` //共识节点
	Remark         string `json:"remark"`
	Administrator  string `json:"administrator"`
}

type Param interface {
	paramInit(params map[string][]byte) error
}

func (org *Organization) paramInit(params []byte) error {
	err := json.Unmarshal(params, org)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) paramInit(params []byte) error {
	err := json.Unmarshal(params, user)
	if err != nil {
		return err
	}
	return nil
}

func (node *Node) paramInit(params []byte) error {
	err := json.Unmarshal(params, node)
	if err != nil {
		return err
	}
	return nil
}

func (chain *ChainList) paramInit(param []byte) error {
	return json.Unmarshal(param, chain)
}
