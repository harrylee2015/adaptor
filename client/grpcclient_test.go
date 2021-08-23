package client

import (
	"context"
	"github.com/33cn/chain33/types"
	"github.com/shimingyah/pool"
	"math/rand"
	"testing"
	"time"
)

func TestGrpcClient(t *testing.T) {
	client := types.NewChain33Client(newGrpcConn("123.60.25.80:8802"))
	header, err := client.GetLastHeader(context.Background(), &types.ReqNil{})
	if err != nil {
		t.Error(err)
	}
	t.Log(header.Height)
}

func TestSendTx(t *testing.T) {
	_, priv := genaddress()
	tx := writeTxPool.Get().(*types.Transaction)
	tx.To = execAddr
	tx.Fee = rand.Int63()
	tx.Nonce = time.Now().UnixNano()
	//tx.Expire = height + types.TxHeightFlag + types.LowAllowPackHeight
	tx.Payload = RandStringBytes(100)
	//交易签名
	tx.Sign(types.SECP256K1, priv)
	client := types.NewChain33Client(newGrpcConn("123.60.25.80:8802"))
	//reply, err := client.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(common.ToHex(reply.Msg))
	//time.Sleep(2*time.Second)
	//detail,err:=client.QueryTransaction(context.Background(),&types.ReqHash{Hash:reply.Msg})
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(common.ToHex(detail.Tx.Hash()))

	account,err:=client.GetBalance(context.Background(),&types.ReqBalance{Execer:"coins",Addresses:[]string{"13gEXsuQhzwfGado1dLmwmysPyFnat6Z11"}})
	if err != nil {
		t.Error(err)
	}
	t.Log(account.Acc[0].Balance)
}


func BenchmarkPoolConn(b *testing.B) {
	p, err := pool.New("123.60.25.80:8802", pool.DefaultOptions)
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := p.Get()
		if err != nil {
			b.Error(err)
		}
		client := types.NewChain33Client(conn.Value())
		client.GetLastHeader(context.Background(), &types.ReqNil{})

	}
}

func BenchmarkSingleConn(b *testing.B) {
	client := types.NewChain33Client(newGrpcConn("123.60.25.80:8802"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetLastHeader(context.Background(), &types.ReqNil{})
		if err != nil {
			b.Error(err)
		}
	}

}
