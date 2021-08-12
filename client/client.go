package client

import (
	"fmt"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const fee = 1e6
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-=_+=/<>!@#$%^&"

var r *rand.Rand
var count uint64
var execAddr = address.ExecAddress("user.write")
var txPool = sync.Pool{
	New: func() interface{} {
		tx := &types.Transaction{Execer: []byte("user.write")}
		return tx
	},
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func genaddress() (string, crypto.PrivKey) {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	privto, err := cr.GenKey()
	if err != nil {
		panic(err)
	}
	addrto := address.PubKeyToAddress(privto.PubKey().Bytes())
	fmt.Println("addr:", addrto.String())
	return addrto.String(), privto
}

// RandStringBytes ...
func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

type Client struct {
	LoaderBalancer map[int]*JSONClient
	sync.RWMutex
	len int
}

func NewClient(urls []string) *Client {
	balancer := make(map[int]*JSONClient)
	for i, url := range urls {
		balancer[i] = NewJSONClient("", url)
	}
	return &Client{LoaderBalancer: balancer, len: len(urls)}
}
func (c *Client) GetLBClient() *JSONClient {
	c.RLock()
	jsonClient := c.LoaderBalancer[rand.Intn(c.len)]
	c.RUnlock()
	return jsonClient
}

// 获取最新区块高度
func (c *Client) GetBlockHeight(w rest.ResponseWriter, r *rest.Request) {
	header, err := c.GetLBClient().GetLastHeader()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteJson(&ReplyHeight{Result: strconv.FormatInt(header.Height, 10)})
}

// 节点总数
func (c *Client) GetNodeCount(w rest.ResponseWriter, r *rest.Request) {
	peerList, err := c.GetLBClient().GetPeerList()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteJson(&ReplyNodeCount{Result: strconv.FormatInt(int64(len(peerList.Peers)), 10)})
}

// 往链交易池推送交易总数
func (c *Client) GetTxAccepted(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(&ReplyAcceptedTxCount{Result: strconv.FormatUint(atomic.LoadUint64(&count),10)})
}

// 已经打包确认交易总数
func (c *Client) GetTxConfirmed(w rest.ResponseWriter, r *rest.Request) {
	header, err := c.GetLBClient().GetLastHeader()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	txfee, err := c.GetLBClient().GetTotalTxCount(header.Height)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&ReplyConfirmedTxCount{Result: strconv.FormatInt(txfee.TxCount,10)})
}

// 获取交易信息
func (c *Client) GetTxInfo(w rest.ResponseWriter, r *rest.Request) {
	//TODO   测试工具接口暂时不可用
}

// 获取区块信息
func (c *Client) GetBlockInfo(w rest.ResponseWriter, r *rest.Request) {

	height, err := strconv.ParseInt(r.Request.FormValue("height"), 10, 64)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blockInfo, err := c.GetLBClient().GetBlockByHeight(height)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&BlockInfo{
		Height:  int(blockInfo.Head.Height),
		TxCount: int(blockInfo.Head.TxCount),
		Hash:    blockInfo.Head.Hash,
		PreHash: blockInfo.Head.ParentHash,
		//时间直接字符串化处理
		CreateTime: strconv.FormatInt(blockInfo.Head.BlockTime, 10),
		TxHashList: blockInfo.TxHashes,
	})
}

// 构建交易，本地构建
func (c *Client) CreateTx(w rest.ResponseWriter, r *rest.Request) {
	var size TxSize
	if err := r.DecodeJsonPayload(&size); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s,_:=strconv.ParseInt(size.Size,64,10)
	_, priv := genaddress()
	//构造存证交易
	tx := txPool.Get().(*types.Transaction)
	fmt.Println("execAddr:",execAddr)
	tx.To = execAddr
	tx.Fee = rand.Int63()
	tx.Nonce = time.Now().UnixNano()
	//tx.Expire = height + types.TxHeightFlag + types.LowAllowPackHeight
	tx.Payload = RandStringBytes(int(s))
	//交易签名
	tx.Sign(types.SECP256K1, priv)
	fmt.Println(common.ToHex(types.Encode(tx)))
	//w.WriteJson(common.ToHex(types.Encode(tx)))
	w.WriteJson(&ReplyTx{
		TxContent: common.ToHex(types.Encode(tx)),
	})

}

// 发送交易
func (c *Client) SendTx(w rest.ResponseWriter, r *rest.Request) {
	//var requestTx RequestTx
	//if err := r.DecodeJsonPayload(&requestTx); err != nil {
	//	rest.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	tx,err :=ioutil.ReadAll(r.Body)
	if err !=nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	fmt.Println("tx:",string(tx))
	hash, err := c.GetLBClient().SendTransaction(string(tx))
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//计数器，区块链系统接收了多少条数据
	atomic.AddUint64(&count, 1)
	w.WriteJson(&ReplyTxHash{ID: hash})

}
