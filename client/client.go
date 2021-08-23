package client

import (
	. "adaptor/common"
	"context"
	"fmt"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/shimingyah/pool"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

const fee = 1e6
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-=_+=/<>!@#$%^&"

var (
	r     *rand.Rand
	count uint64

	execAddr    = address.ExecAddress("user.write")
	writeTxPool = sync.Pool{
		New: func() interface{} {
			tx := &types.Transaction{Execer: []byte("user.write")}
			return tx
		},
	}
	coinsExecAddr  = address.ExecAddress("coins")
	transferTxPool = sync.Pool{
		New: func() interface{} {
			tx := &types.Transaction{Execer: []byte("coins")}
			return tx
		},
	}

	//未发送成功的交易
	UnSendedTxMap sync.Map
	//限流器
	Limiter *ChannelLimiter

	//转账签名账户
	privkey crypto.PrivKey
)

func InitLimiter(limit int) *ChannelLimiter {
	Limiter = NewChannelLimiter(limit)
	return Limiter
}

func InitPrivKey(key string) {
	privkey = getprivkey(key)
}

func LoadPrivKey(str string) (crypto.PrivKey, error) {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		return nil, err
	}
	data, err := common.FromHex(str)
	if err != nil {
		return nil, err
	}
	return cr.PrivKeyFromBytes(data)
}

type Client struct {
	JrpcBalancer    []*JSONClient
	GrpcBalancer    []types.Chain33Client
	GrpcConnectPool []pool.Pool
	Async           bool
	sync.RWMutex
}

func NewClient(jsonurls, grpcurls []string,async bool) *Client {
	jrpcBalancer := make([]*JSONClient, len(jsonurls))
	grpcBalancer := make([]types.Chain33Client, len(grpcurls))
	grpcConnectPool := make([]pool.Pool, len(grpcurls))
	for i, url := range jsonurls {
		jrpcBalancer[i] = NewJSONClient("", url)
	}
	for i, url := range grpcurls {
		gcli := types.NewChain33Client(newGrpcConn(url))
		grpcBalancer[i] = gcli
		p, err := pool.New(url, pool.DefaultOptions)
		if err != nil {
			panic(err)
		}
		grpcConnectPool[i] = p
	}
	return &Client{JrpcBalancer: jrpcBalancer, GrpcBalancer: grpcBalancer, GrpcConnectPool: grpcConnectPool,Async:async}
}
func (c *Client) GetGrpcClient() types.Chain33Client {
	c.RLock()
	grpcClient := c.GrpcBalancer[rand.Intn(len(c.GrpcBalancer))]
	c.RUnlock()
	return grpcClient
}

func (c *Client) GetJrpcClient() *JSONClient {
	c.RLock()
	jrcpClient := c.JrpcBalancer[rand.Intn(len(c.JrpcBalancer))]
	c.RUnlock()
	return jrcpClient
}

//func (c *Client) GetGrpcConnectPool() pool.Pool {
//	c.RLock()
//	p := c.GrpcConnectPool[rand.Intn(c.len)]
//	c.RUnlock()
//	return p
//}
func (c *Client) GetClient() (types.Chain33Client, error) {
	c.RLock()
	p := c.GrpcConnectPool[rand.Intn(len(c.GrpcConnectPool))]
	c.RUnlock()
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}
	return types.NewChain33Client(conn.Value()), nil
}

// 获取最新区块高度
func (c *Client) GetBlockHeight(w rest.ResponseWriter, r *rest.Request) {
	header, err := c.GetGrpcClient().GetLastHeader(context.Background(), &types.ReqNil{})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteJson(&ReplyHeight{Result: strconv.FormatInt(header.Height, 10)})
}

// 节点总数
func (c *Client) GetNodeCount(w rest.ResponseWriter, r *rest.Request) {
	peerList, err := c.GetGrpcClient().GetPeerInfo(context.Background(), &types.P2PGetPeerReq{})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteJson(&ReplyNodeCount{Result: strconv.FormatInt(int64(len(peerList.Peers)), 10)})
}

// 往链交易池推送交易总数
func (c *Client) GetTxAccepted(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(&ReplyAcceptedTxCount{Result: strconv.FormatUint(atomic.LoadUint64(&count), 10)})
}

// 已经打包确认交易总数
func (c *Client) GetTxConfirmed(w rest.ResponseWriter, r *rest.Request) {
	header, err := c.GetGrpcClient().GetLastHeader(context.Background(), &types.ReqNil{})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	txfee, err := c.GetJrpcClient().GetTotalTxCount(header.Height)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&ReplyConfirmedTxCount{Result: strconv.FormatInt(txfee.TxCount, 10)})
}
// 获取交易信息
func (c *Client) GetTxInfo(w rest.ResponseWriter, r *rest.Request) {
	//TODO   测试工具接口暂时不可用
	tx, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := common.FromHex(string(tx))
	var tr types.Transaction
	types.Decode(data, &tr)
	client, err := c.GetClient()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	detail,err:=client.QueryTransaction(context.Background(),&types.ReqHash{Hash:tr.Hash()})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&ReplyTxHash{ID: common.ToHex(detail.Tx.Hash())})
}
// 获取存证信息
func (c *Client) GetWriteInfo(w rest.ResponseWriter, r *rest.Request) {
	//TODO   测试工具接口暂时不可用
	tx, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := common.FromHex(string(tx))
	var tr types.Transaction
	types.Decode(data, &tr)
	client, err := c.GetClient()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	detail,err:=client.QueryTransaction(context.Background(),&types.ReqHash{Hash:tr.Hash()})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(common.ToHex(detail.Tx.Payload))
}

// 获取账户余额
func (c *Client) GetBalance(w rest.ResponseWriter, r *rest.Request) {
	//TODO 这个接口不可用
	tx, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := common.FromHex(string(tx))
	var tr types.Transaction
	types.Decode(data, &tr)
	client, err := c.GetClient()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addr:=tr.From()
	detail,err:=client.GetBalance(context.Background(),&types.ReqBalance{Execer:"coins",Addresses:[]string{addr}})
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(detail)

}

// 获取区块信息
func (c *Client) GetBlockInfo(w rest.ResponseWriter, r *rest.Request) {
	var req RequestH
	if err := r.DecodeJsonPayload(&req); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	height, err := strconv.ParseInt(req.Height, 10, 64)
	//height, err := strconv.ParseInt(r.Request.FormValue("height"), 10, 64)
	fmt.Println("hegiht:", height)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blockInfo, err := c.GetJrpcClient().GetBlockByHeight(height)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&BlockInfo{
		Height:  int(blockInfo.Head.Height),
		TxCount: int(blockInfo.Head.TxCount),
		Hash:    blockInfo.Head.Hash,
		PreHash: blockInfo.Head.ParentHash,
		//时间直接字符串化处理,时间戳19位
		CreateTime: strconv.FormatInt(blockInfo.Head.BlockTime*1e9, 10),
		TxHashList: blockInfo.TxHashes,
	})
}

// 构建交易，本地构建
func (c *Client) CreateTx(w rest.ResponseWriter, r *rest.Request) {
	var txType TxType
	if err := r.DecodeJsonPayload(&txType); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tx *types.Transaction
	if txType.IsTransfer {
		tx = createTransferTx()
	} else {
		s, _ := strconv.ParseInt(txType.Size, 10, 64)
		tx = createWriteTx(int(s))
	}
	//w.WriteJson(common.ToHex(types.Encode(tx)))
	w.WriteJson(&ReplyTx{
		TxContent: common.ToHex(types.Encode(tx)),
	})
}

// 发送交易
func (c *Client) SendTx(w rest.ResponseWriter, r *rest.Request) {
	if c.Async{
		c.asyncSendTx(w,r)
		return
	}
	c.syncSendTx(w,r)
}
// 同步发送
func (c *Client) syncSendTx(w rest.ResponseWriter, r *rest.Request) {
	tx, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := common.FromHex(string(tx))
	var tr types.Transaction
	types.Decode(data, &tr)
	grpcClient, err := c.GetClient()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reply, err := grpcClient.SendTransaction(context.Background(), &tr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if reply.IsOk {
		//计数器，区块链系统接收了多少条数据
		atomic.AddUint64(&count, 1)
		w.WriteJson(&ReplyTxHash{ID: common.ToHex(reply.Msg)})
	} else {
		rest.Error(w, fmt.Errorf("The service did not handle the request properly!").Error(), http.StatusInternalServerError)
		return
	}
}
// 异步发送
func (c *Client) asyncSendTx(w rest.ResponseWriter, r *rest.Request) {
	tx, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, _ := common.FromHex(string(tx))
	var tr types.Transaction
	types.Decode(data, &tr)
	//计数器，区块链系统接收了多少条数据
	atomic.AddUint64(&count, 1)
	w.WriteJson(&ReplyTxHash{ID: common.ToHex(tr.Hash())})
	//异步转发处理
	cli, err := c.GetClient()
	if err != nil {
		UnSendedTxMap.Store(common.ToHex(tr.Hash()), tr)
		return
	}
	go func(tx *types.Transaction, grpcClient types.Chain33Client) {
		reply, err := cli.SendTransaction(context.Background(), tx)
		if err != nil || !reply.IsOk {
			UnSendedTxMap.Store(common.ToHex(tx.Hash()), *tx)
		}
	}(&tr, cli)
}
