package client

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestJSONClient_GetPeerList(t *testing.T) {
	jsonclient := NewJSONClient("", "http://123.60.25.80:8801")
	peerList, err := jsonclient.GetPeerList()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(peerList)
}

func TestJSONClient_SendTx(t *testing.T) {
	//tx := "0x0a0a757365722e77726974651a6d0801122103068eb379d486ad923f57c7983a1ed32227fc9fa09355412ed9709d7e9739bb6d1a46304402207f379436e2ee009effd43e67256102fc5fa5c97750aec989236824ec25d7ee9302205b949c77fe4011f2f3364011ac575b74d579a0df7b477cc0f6b8bd32808ecbb6208cb4ff9df6cea6f36330b88588cedd829ecd163a2231444e615344524739524431397335396d65416f654e34613246365248393766536f"
	jsonclient := NewJSONClient("", "http://123.60.25.80:8801")
	rep,err:=jsonclient.GetTotalTxCount(10)
	if err !=nil {
		t.Error(err)
	}
	t.Log(rep.TxCount)
	//txhash, err := jsonclient.SendTransaction(tx)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(txhash)
	//start := time.Now().UnixNano()
	//data, _ := common.FromHex(tx)
	//var tr types.Transaction
	//types.Decode(data, &tr)
	//end := time.Now().UnixNano()
	//fmt.Printf("执行消耗的时间为 %d纳秒", end-start)
	//fmt.Println(tr)
}

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(rand.Intn(1))
	}
	t.Log(strconv.FormatInt(10*1e9, 10))
}


func BenchmarkRandStringBytes(b *testing.B) {
	ch :=make(chan int,1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
        ch<-i
	}
	for i:=0;i<10;i++{
		go func() {
			for{
				select {
				case v:=<-ch:
					b.Log(v)

				}
			}
		}()
	}

}