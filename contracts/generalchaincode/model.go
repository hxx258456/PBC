package main

import "time"

type PageResult struct {
	Count    int32       `json:"count"`    // 总数量
	Nextmark string      `json:"nextmark"` // 指向下一个的标记
	Data     interface{} `json:"data"`
}

// key log
type LogResult struct {
	Record    interface{} `json:"record"`
	TxId      string      `json:"txId"`
	Timestamp time.Time   `json:"timestamp"`
	IsDelete  bool        `json:"isDelete"`
}
