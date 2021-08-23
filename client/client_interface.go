package client

import "github.com/ant0ine/go-json-rest/rest"

type TxInfo struct {
	BlockHeight int    `json:"blockHeight"`
	TxHash      string `json:"txhash"`
	ConfirmTime string `json:"confirmtime"`
}

type BlockInfo struct {
	Height     int      `json:"height"`
	TxCount    int      `json:"txcount"`
	Hash       string   `json:"hash"`
	PreHash    string   `json:"prehash"`
	CreateTime string   `json:"createtime"`
	TxHashList []string `json:"txhashlist"` // 交易hash列表
}

type SendTxInfo struct {
	TxHash       string `json:"txhash"`
	RequestTime  string `json:"requesttime"`
	ResponseTime string `json:"responsetime"`
	Invalid      string `json:"invalid"`
}

type TxReportInfo struct {
	TxHash       string `json:"txhash"`
	BlockHeight  int    `json:"blockheight"`
	CreateTime   string `json:"createtime"`
	ResponseTime string `json:"responsetime"`
	ConfirmTime  string `json:"confirmtime"`
}

type Tasks interface {
	// 获取区块高度
	GetBlockHeight(w rest.ResponseWriter, r *rest.Request)

	// 获取区块链网络的节点数量
	GetNodeCount(w rest.ResponseWriter, r *rest.Request)

	// 获取已经接收到的交易数
	GetTxAccepted(w rest.ResponseWriter, r *rest.Request)

	// 获取已经确认（落块）的交易数
	GetTxConfirmed(w rest.ResponseWriter, r *rest.Request)

	// 获取交易信息
	// 一般为 Json 字符串，然后根据配置文件中的键值映射获取，从返回结果中获取需要的值
	GetTxInfo(w rest.ResponseWriter, r *rest.Request)

	// 获取区块信息
	// 一般为 Json 字符串，然后根据配置文件中的键值映射获取，从返回结果中获取需要的值
	GetBlockInfo(w rest.ResponseWriter, r *rest.Request)

	// 创建交易
	// 用于先在本地生成一批准备好的交易
	CreateTx(w rest.ResponseWriter, r *rest.Request)

	// 发送交易
	SendTx(w rest.ResponseWriter, r *rest.Request)

	// 查询余额
	GetBalance(w rest.ResponseWriter, r *rest.Request)
}
