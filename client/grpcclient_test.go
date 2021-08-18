package client

import (
	"context"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/types"
	"math/rand"
	"testing"
	"time"
)

func TestGrpcClient(t  *testing.T){
	client:=types.NewChain33Client(newGrpcConn("123.60.25.80:8802"))
	header,err:=client.GetLastHeader(context.Background(),&types.ReqNil{})
	if err !=nil {
		t.Error(err)
	}
	t.Log(header.Height)
}

func TestSendTx(t  *testing.T){
	_, priv := genaddress()
	tx := txPool.Get().(*types.Transaction)
	tx.To = execAddr
	tx.Fee = rand.Int63()
	tx.Nonce = time.Now().UnixNano()
	//tx.Expire = height + types.TxHeightFlag + types.LowAllowPackHeight
	tx.Payload = RandStringBytes(100)
	//交易签名
	tx.Sign(types.SECP256K1, priv)
	client:=types.NewChain33Client(newGrpcConn("123.60.25.80:8802"))
	reply,err:=client.SendTransaction(context.Background(),tx)
	if err !=nil {
		t.Error(err)
	}
	t.Log(common.ToHex(reply.Msg))
}
