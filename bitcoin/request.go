package bitcoin

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

// BaseRequest
type BaseRequest struct {
	RpcUrl   string
	User     string
	Password string
	Proxy    string
}

// Error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BaseResponse
type BaseResponse struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   Error       `json:"error"`
}

//
type ReqJsonRpc struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  RpcMethod   `json:"method"`
	Params  interface{} `json:"params"`
	Id      string      `json:"id"`
}

type RpcMethod string

const (
	RpcMethodGetBlockChainInfo    RpcMethod = "getblockchaininfo"
	RpcMethodGetBlockHash         RpcMethod = "getblockhash"
	RpcMethodGetBlock             RpcMethod = "getblock"
	RpcMethodGetRawTransaction    RpcMethod = "getrawtransaction"
	RpcMethodListUnspent          RpcMethod = "listunspent"
	RpcMethodSendRawTransaction   RpcMethod = "sendrawtransaction"
	RpcMethodEstimateFee          RpcMethod = "estimatefee"
	RpcMethodEstimateSmartFee     RpcMethod = "estimatesmartfee"
	RpcMethodDecodeRawTransaction RpcMethod = "decoderawtransaction"
)

func (b *BaseRequest) Request(method RpcMethod, params []interface{}, result interface{}) error {
	var req ReqJsonRpc
	req.Jsonrpc = "2.0"
	req.Method = method
	req.Id = "1"
	req.Params = params

	engine := gorequest.New().Timeout(time.Second * 30)
	if b.User != "" && b.Password != "" {
		engine = engine.SetBasicAuth(b.User, b.Password)
	}
	if b.Proxy != "" {
		engine = engine.Proxy(b.Proxy)
	}

	var resp BaseResponse
	resp.Result = result

	res, body, errs := engine.Post(b.RpcUrl).
		Set("Content-Type", "application/json").
		Send(&req).EndStruct(&resp)

	if len(errs) > 0 {
		return fmt.Errorf("req errs: %v", errs)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http code: %d, [%s]", res.StatusCode, body)
	}
	if result == nil {
		fmt.Println("body:", string(body))
	}
	return nil
}

//  '{"jsonrpc":"2.0","id":"0","method":"getblockchaininfo"}' -H 'Content-Type: application/json'
type BlockChainInfo struct {
	BestBlockHash string `json:"bestblockhash"`
	Blocks        uint64 `json:"blocks"`
	Chain         string `json:"chain"`
	ChainWork     string `json:"chainwork"`
	Headers       uint64 `json:"headers"`
	MedianTime    uint64 `json:"mediantime"`
}

//  '{"jsonrpc":"2.0","id":"0","method":"getblockhash","params":[4600472]}' -H 'Content-Type: application/json'

//  '{"jsonrpc":"2.0","id":"0","method":"getblock","params":["5d0954672b3d7bc9becbfa017f7cb47714c39ef74ab99c969217ee2af0d40a82"]}' -H 'Content-Type: application/json'
type BlockInfo struct {
	ChainWork         string   `json:"chainwork"`
	Confirmations     uint64   `json:"confirmations"`
	Hash              string   `json:"hash"`
	Height            uint64   `json:"height"`
	MedianTime        uint64   `json:"mediantime"`
	NextBlockHash     string   `json:"nextblockhash"`
	PreviousBlockHash string   `json:"previousblockhash"`
	Time              uint64   `json:"time"`
	Tx                []string `json:"tx"`
}

//  '{"jsonrpc":"2.0","id":"0","method":"getrawtransaction","params":["c9b477a5afabbd6ff7afea9a2b0dce9687e1dc56a452b72e336b2961126fe411",true]}' -H 'Content-Type: application/json'
