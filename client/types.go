package client

import "encoding/json"

// RawParm defines raw parameter command
type RawParm struct {
	Token string `json:"token"`
	Data  string `json:"data"`
}

// CreateTxIn create tx input
type CreateTxIn struct {
	Execer     string          `json:"execer"`
	ActionName string          `json:"actionName"`
	Payload    json.RawMessage `json:"payload"`
}

// QueryParm Query parameter
type QueryParm struct {
	Hash string `json:"hash"`
}

// 无参
type ParmNil struct {
}

// Signature parameter
type Signature struct {
	Ty        int32  `json:"ty"`
	Pubkey    string `json:"pubkey"`
	Signature string `json:"signature"`
}

// Transaction parameter
type Transaction struct {
	Execer     string          `json:"execer"`
	Payload    json.RawMessage `json:"payload"`
	RawPayload string          `json:"rawPayload"`
	Signature  *Signature      `json:"signature"`
	Fee        int64           `json:"fee"`
	FeeFmt     string          `json:"feefmt"`
	Expire     int64           `json:"expire"`
	Nonce      int64           `json:"nonce"`
	From       string          `json:"from,omitempty"`
	To         string          `json:"to"`
	Amount     int64           `json:"amount,omitempty"`
	AmountFmt  string          `json:"amountfmt,omitempty"`
	GroupCount int32           `json:"groupCount,omitempty"`
	Header     string          `json:"header,omitempty"`
	Next       string          `json:"next,omitempty"`
	Hash       string          `json:"hash,omitempty"`
}

// ReceiptDataResult receipt data result
type ReceiptDataResult struct {
	Ty     int32               `json:"ty"`
	TyName string              `json:"tyName"`
	Logs   []*ReceiptLogResult `json:"logs"`
}

// ReceiptLogResult receipt log result
type ReceiptLogResult struct {
	Ty     int32           `json:"ty"`
	TyName string          `json:"tyName"`
	Log    json.RawMessage `json:"log"`
	RawLog string          `json:"rawLog"`
}

// Asset asset
type Asset struct {
	Exec   string `json:"exec"`
	Symbol string `json:"symbol"`
	Amount int64  `json:"amount"`
}

// TxProof :
type TxProof struct {
	Proofs   []string `json:"proofs"`
	Index    uint32   `json:"index"`
	RootHash string   `json:"rootHash"`
}

// TransactionDetail transaction detail
type TransactionDetail struct {
	Tx         *Transaction       `json:"tx"`
	Receipt    *ReceiptDataResult `json:"receipt"`
	Proofs     []string           `json:"proofs"`
	Height     int64              `json:"height"`
	Index      int64              `json:"index"`
	Blocktime  int64              `json:"blockTime"`
	Amount     int64              `json:"amount"`
	Fromaddr   string             `json:"fromAddr"`
	ActionName string             `json:"actionName"`
	Assets     []*Asset           `json:"assets"`
	TxProofs   []*TxProof         `json:"txProofs"`
	FullHash   string             `json:"fullHash"`
}

// Query4Jrpc query jrpc
type Query4Jrpc struct {
	Execer   string          `json:"execer"`
	FuncName string          `json:"funcName"`
	Payload  json.RawMessage `json:"payload"`
}

// Header header parameter
type Header struct {
	Version    int64      `json:"version"`
	ParentHash string     `json:"parentHash"`
	TxHash     string     `json:"txHash"`
	StateHash  string     `json:"stateHash"`
	Height     int64      `json:"height"`
	BlockTime  int64      `json:"blockTime"`
	TxCount    int64      `json:"txCount"`
	Hash       string     `json:"hash"`
	Difficulty uint32     `json:"difficulty"`
	Signature  *Signature `json:"signature,omitempty"`
}

// Peer  information
type Peer struct {
	Addr           string  `json:"addr"`
	Port           int32   `json:"port"`
	Name           string  `json:"name"`
	MempoolSize    int32   `json:"mempoolSize"`
	Self           bool    `json:"self"`
	Header         *Header `json:"header"`
	Version        string  `json:"version,omitempty"`
	LocalDBVersion string  `json:"localDBVersion,omitempty"`
	StoreDBVersion string  `json:"storeDBVersion,omitempty"`
}

// PeerList peer list
type PeerList struct {
	Peers []*Peer `json:"peers"`
}

// WalletAccounts Wallet Module
type WalletAccounts struct {
	Wallets []*WalletAccount `json:"wallets"`
}

// WalletAccount  wallet account
type WalletAccount struct {
	Acc   *Account `json:"acc"`
	Label string   `json:"label"`
}

// Account account information
type Account struct {
	Currency int32  `json:"currency"`
	Balance  int64  `json:"balance"`
	Frozen   int64  `json:"frozen"`
	Addr     string `json:"addr"`
}

// WalletStatus wallet status
type WalletStatus struct {
	IsWalletLock bool `json:"isWalletLock"`
	IsAutoMining bool `json:"isAutoMining"`
	IsHasSeed    bool `json:"isHasSeed"`
	IsTicketLock bool `json:"isTicketLock"`
}

// p2p get peer req
type P2PGetPeerReq struct {
	P2PType              string   `protobuf:"bytes,1,opt,name=p2pType,proto3" json:"p2pType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type ReqWalletSendToAddress struct {
	From                 string   `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	To                   string   `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Amount               int64    `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Note                 string   `protobuf:"bytes,4,opt,name=note,proto3" json:"note,omitempty"`
	IsToken              bool     `protobuf:"varint,5,opt,name=isToken,proto3" json:"isToken,omitempty"`
	TokenSymbol          string   `protobuf:"bytes,6,opt,name=tokenSymbol,proto3" json:"tokenSymbol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type RequestTx struct {
	Tx string `json:"tx"`
}

type RequestH struct {
	Height string `json:"height"`
}
type TxType struct {
	Size       string `json:"size"`
	IsTransfer bool   `json:"istransfer"`
}

// ReplyHash reply hash string json
type ReplyHash struct {
	Hash string `json:"hash"`
}

// block height
type ReplyHeight struct {
	Result string `json:"height"`
}

// Peer  information
type ReplyNodeCount struct {
	Result string `json:"result"`
}

//接收交易总数
type ReplyAcceptedTxCount struct {
	Result string `json:"result"`
}

//接收交易总数
type ReplyConfirmedTxCount struct {
	Result string `json:"result"`
}

//发送交易还回hash
type ReplyTxHash struct {
	ID string `json:"ID"`
}

//return sign tx
type ReplyTx struct {
	TxContent string `json:"txcontent"`
}

//type BlockInfo struct {
//	Height     int      `json:"height"`
//	TxCount    int      `json:"txcount"`
//	Hash       string   `json:"hash"`
//	PreHash    string   `json:"prehash"`
//	CreateTime string   `json:"createtime"`
//	TxHashList []string `json:"txhashlist"` // 交易hash列表
//}
