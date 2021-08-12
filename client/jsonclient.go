// Package jsonclient 实现JSON rpc客户端请求功能
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/33cn/chain33/common"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"

	"io/ioutil"
	"net/http"
	"strings"
)

// JSONClient a object of jsonclient
type JSONClient struct {
	url       string
	prefix    string
	tlsVerify bool
	client    *http.Client
}

func addPrefix(prefix, name string) string {
	if strings.Contains(name, ".") {
		return name
	}
	return prefix + "." + name
}

// NewJSONClient produce a json object
func NewJSONClient(prefix, url string) *JSONClient {
	return new(prefix, url, false)
}

// New produce a jsonclient by perfix and url
func new(prefix, url string, tlsVerify bool) *JSONClient {
	httpcli := http.DefaultClient
	if strings.Contains(url, "https") { //暂不校验tls证书
		httpcli = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerify}}}
	}
	return &JSONClient{
		url:       url,
		prefix:    prefix,
		tlsVerify: tlsVerify,
		client:    httpcli,
	}
}

type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	ID     uint64         `json:"id"`
}

type clientResponse struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

func (client *JSONClient) Call(method string, params, resp interface{}) error {
	method = addPrefix(client.prefix, method)
	req := &clientRequest{}
	req.Method = method
	req.Params[0] = params
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	postresp, err := client.client.Post(client.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer postresp.Body.Close()
	b, err := ioutil.ReadAll(postresp.Body)
	if err != nil {
		return err
	}
	cresp := &clientResponse{}
	err = json.Unmarshal(b, &cresp)
	if err != nil {
		return err
	}
	if cresp.Error != nil {
		x, ok := cresp.Error.(string)
		if !ok {
			return fmt.Errorf("invalid error %v", cresp.Error)
		}
		if x == "" {
			x = "unspecified error"
		}
		return fmt.Errorf(x)
	}
	if cresp.Result == nil {
		return errors.New("Empty result")
	}
	if msg, ok := resp.(proto.Message); ok {
		var str json.RawMessage
		err = json.Unmarshal(*cresp.Result, &str)
		if err != nil {
			return err
		}
		b, err := str.MarshalJSON()
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, msg)
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		return nil
	}
	return json.Unmarshal(*cresp.Result, resp)
}

type ParseFunc func(result json.RawMessage) (interface{}, error)

//回调函数，用于自定义解析返回得result数据
func (client *JSONClient) CallBack(method string, params interface{}, parseFunc ParseFunc) (interface{}, error) {
	method = addPrefix(client.prefix, method)
	req := &clientRequest{}
	req.Method = method
	req.Params[0] = params
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	postresp, err := client.client.Post(client.url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer postresp.Body.Close()
	b, err := ioutil.ReadAll(postresp.Body)
	if err != nil {
		return nil, err
	}
	cresp := &clientResponse{}
	err = json.Unmarshal(b, &cresp)
	if err != nil {
		return nil, err
	}
	if cresp.Error != nil {
		x, ok := cresp.Error.(string)
		if !ok {
			return nil, fmt.Errorf("invalid error %v", cresp.Error)
		}
		if x == "" {
			x = "unspecified error"
		}
		return nil, fmt.Errorf(x)
	}
	if cresp.Result == nil {
		return nil, errors.New("Empty result")
	}
	return parseFunc(*cresp.Result)
}

// 发送交易
func (client *JSONClient) SendTransaction(signedTx string) (string, error) {
	var res string
	send := &RawParm{
		Token: "BTY",
		Data:  signedTx,
	}
	err := client.Call("Chain33.SendTransaction", send, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// 查询交易
func (client *JSONClient) QueryTransaction(hash string) (*TransactionDetail, error) {
	query := QueryParm{
		Hash: hash,
	}
	var detail TransactionDetail
	err := client.Call("Chain33.QueryTransaction", query, &detail)
	if err != nil {
		return nil, err
	}

	return &detail, nil
}

//获取节点列表
func (client *JSONClient) GetPeerList() (*PeerList, error) {
	var peerList PeerList
	err := client.Call("Chain33.GetPeerInfo", &P2PGetPeerReq{}, &peerList)
	if err != nil {
		return nil, err
	}
	return &peerList, nil
}

//获取最新区块
func (client *JSONClient) GetLastHeader() (*Header, error) {
	var header Header
	err := client.Call("Chain33.GetLastHeader", &ParmNil{}, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
func (client *JSONClient) QueryBlockHash(height int64) (*ReplyHash, error) {
	var res ReplyHash
	params := types.ReqInt{
		Height: height,
	}
	err := client.Call("Chain33.GetBlockHash", &params, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (client *JSONClient) GetBlockByHeight(height int64) (*rpctypes.BlockOverview, error) {
	reply, err := client.QueryBlockHash(height)
	if err != nil {
		return nil, err
	}
	var blockOverview rpctypes.BlockOverview
	err = client.Call("Chain33.GetBlockOverview", &rpctypes.QueryParm{
		Hash: reply.Hash,
	}, &blockOverview)
	if err != nil {
		return nil, err
	}
	return &blockOverview, nil
}

// 获取上链交易总数
func (client *JSONClient) GetTotalTxCount(height int64) (*types.TotalFee, error) {
	reply, err := client.QueryBlockHash(height)
	if err != nil {
		return nil, err
	}
	hash, err := common.FromHex(reply.Hash)
	if err != nil {
		return nil, err
	}
	hash = append([]byte("TotalFeeKey:"), hash...)
	params := types.LocalDBGet{Keys: [][]byte{hash[:]}}
	res := &types.TotalFee{}

	err = client.Call("Chain33.QueryTotalFee", &params, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
