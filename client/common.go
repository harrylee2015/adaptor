package client

import (
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	cty "github.com/33cn/chain33/system/dapp/coins/types"
	"github.com/33cn/chain33/types"
	"math/rand"
	"time"
)
//随机生成一对密钥
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
	return addrto.String(), privto
}
//根据四要字符串获取私钥
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
//构造存证交易
func createWriteTx(size int) *types.Transaction {
	_, priv := genaddress()
	tx := writeTxPool.Get().(*types.Transaction)
	tx.To = execAddr
	tx.Fee = rand.Int63()
	tx.Nonce = time.Now().UnixNano()
	//tx.Expire = height + types.TxHeightFlag + types.LowAllowPackHeight
	tx.Payload = RandStringBytes(size)
	//交易签名
	tx.Sign(types.SECP256K1, priv)
	return tx
}

//构造转账交易
func createTransferTx(priv crypto.PrivKey) *types.Transaction {
	addr, _ := genaddress()
	tx := transferTxPool.Get().(*types.Transaction)
	tx.Fee = fee
	tx.Nonce = time.Now().UnixNano()
	action := cty.CoinsAction{Value: &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{
		Cointoken: "",
		Amount:    1e6,
		To:        addr,
	}},
		Ty: cty.CoinsActionTransfer}
	tx.Payload = types.Encode(&action)
	tx.To = coinsExecAddr
	//交易签名
	tx.Sign(types.SECP256K1, priv)
	return tx
}

// 随机生成指定字节大小的字符串
func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}