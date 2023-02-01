package smt

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/parnurzeal/gorequest"
	"time"
)

const (
	GetSmtRoot        = "get_smt_root"
	DeleteSmt         = "delete_smt"
	UpdateMemorySmt   = "update_memory_smt"
	UpdateDbSmt       = "update_db_smt"
	UpdateDbSmtMiddle = "update_db_smt_middle"

	RetryNumber = 3
	RetryTime   = time.Second * 3
	TimeOut     = time.Second * 20
)

type SmtKvHex struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SmtKv struct {
	Key   H256
	Value H256
}

type SmtOpt struct {
	GetProof bool `json:"get_proof"`
	GetRoot  bool `json:"get_root"`
}

type UpdateSmtParam struct {
	Opt     SmtOpt     `json:"opt"`
	Data    []SmtKvHex `json:"data"`
	SmtName string     `json:"smt_name"`
}

type UpdateSmtOut struct {
	Root   H256              `json:"root"`
	Proofs map[string]string `json:"proofs"`
}

type UpdateMiddleSmtOut struct {
	Roots  map[string]H256   `json:"root"`
	Proofs map[string]string `json:"proofs"`
}

type smtServerReq struct {
	Id      int         `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UpdateResult struct {
	Root   string            `json:"root"`
	Proofs map[string]string `json:"proofs"`
}

type UpdateMiddleResult struct {
	Roots  map[string]string `json:"roots"`
	Proofs map[string]string `json:"proofs"`
}

type SmtServer struct {
	url     string
	smtName string
}

func NewSmtSrv(url, smtName string) *SmtServer {
	return &SmtServer{
		url:     url,
		smtName: smtName,
	}
}
func (s *SmtServer) GetSmtUrl() string {
	return s.url
}

func (s *SmtServer) GetSmtRoot() (H256, error) {
	var rpcReq smtServerReq

	rpcReq = newBasicReq(GetSmtRoot)
	params := make(map[string]string)
	params["smt_name"] = s.smtName
	rpcReq.Params = params
	body, err := sendAndCheck(s.url, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("UpdateSmt %s", err.Error())
	}

	rpcRep := struct {
		JsonRpc string       `json:"jsonrpc"`
		Result  string       `json:"result"`
		Error   JsonRpcError `json:"error"`
	}{}
	if err := json.Unmarshal([]byte(*body), &rpcRep); err != nil {
		reqByte, _ := json.Marshal(rpcReq)
		return nil, fmt.Errorf("GetSmtRoot Json Unmarshal err: %s body: %s, request: %s", err.Error(), *body, string(reqByte))
	}

	return common.Hex2Bytes(rpcRep.Result), nil
}

func (s *SmtServer) DeleteSmt() (bool, error) {
	return s.DeleteSmtWithTimeOut(TimeOut)
}

func (s *SmtServer) DeleteSmtWithTimeOut(timeout time.Duration) (bool, error) {
	var rpcReq smtServerReq

	rpcReq = newBasicReq(DeleteSmt)
	params := make(map[string]string)
	params["smt_name"] = s.smtName
	rpcReq.Params = params

	body, err := sendAndCheckWithTimeout(s.url, rpcReq, timeout)
	if err != nil {
		return false, fmt.Errorf("DeleteSmt %s", err.Error())
	}

	rpcRep := struct {
		JsonRpc string       `json:"jsonrpc"`
		Result  bool         `json:"result"`
		Error   JsonRpcError `json:"error"`
	}{}

	if err := json.Unmarshal([]byte(*body), &rpcRep); err != nil {
		reqByte, _ := json.Marshal(rpcReq)
		return false, fmt.Errorf("DeleteSmt Json Unmarshal err: %s body: %s, request:%s", err.Error(), *body, string(reqByte))
	}

	return rpcRep.Result, nil
}

func (s *SmtServer) UpdateSmt(kv []SmtKv, opt SmtOpt) (*UpdateSmtOut, error) {
	var (
		rpcReq smtServerReq
		out    UpdateSmtOut
		param  UpdateSmtParam
		kvHex  []SmtKvHex
	)

	for i, _ := range kv {
		kvHex = append(kvHex, SmtKvHex{
			Key:   common.Bytes2Hex(kv[i].Key)[2:],
			Value: common.Bytes2Hex(kv[i].Value)[2:],
		})
	}
	param.Opt = opt
	param.Data = kvHex

	if s.smtName == "" {
		rpcReq = newBasicReq(UpdateMemorySmt)
	} else {
		rpcReq = newBasicReq(UpdateDbSmt)
		param.SmtName = s.smtName
	}
	rpcReq.Params = param

	body, err := sendAndCheck(s.url, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("UpdateSmt %s", err.Error())
	}
	rpcRep := struct {
		JsonRpc string       `json:"jsonrpc"`
		Result  UpdateResult `json:"result"`
		Error   JsonRpcError `json:"error"`
	}{}
	if err := json.Unmarshal([]byte(*body), &rpcRep); err != nil {
		reqByte, _ := json.Marshal(rpcReq)
		return nil, fmt.Errorf("UpdateSmt Json Unmarshal err: %s body: %s, request: %s", err.Error(), *body, reqByte)
	}

	out.Root = common.Hex2Bytes(rpcRep.Result.Root)
	out.Proofs = make(map[string]string)

	for i, _ := range rpcRep.Result.Proofs {
		key := fmt.Sprintf("%s%s", common.HexPreFix, i)
		out.Proofs[key] = fmt.Sprintf("%s%s", common.HexPreFix, rpcRep.Result.Proofs[i])
	}

	return &out, nil
}

func (s *SmtServer) UpdateMiddleSmt(kv []SmtKv, opt SmtOpt) (*UpdateMiddleSmtOut, error) {
	var (
		rpcReq smtServerReq
		out    UpdateMiddleSmtOut
		param  UpdateSmtParam
		kvHex  []SmtKvHex
	)

	for i, _ := range kv {
		kvHex = append(kvHex, SmtKvHex{
			Key:   common.Bytes2Hex(kv[i].Key)[2:],
			Value: common.Bytes2Hex(kv[i].Value)[2:],
		})
	}
	param.Opt = opt
	param.Data = kvHex
	rpcReq = newBasicReq(UpdateDbSmtMiddle)
	param.SmtName = s.smtName
	rpcReq.Params = param

	body, err := sendAndCheck(s.url, rpcReq)
	if err != nil {
		return nil, fmt.Errorf("UpdateSmt %s", err.Error())
	}

	rpcRep := struct {
		JsonRpc string             `json:"jsonrpc"`
		Result  UpdateMiddleResult `json:"result"`
		Error   JsonRpcError       `json:"error"`
	}{}
	if err := json.Unmarshal([]byte(*body), &rpcRep); err != nil {
		reqByte, _ := json.Marshal(rpcReq)
		return nil, fmt.Errorf("UpdateMiddleSmt Json Unmarshal err: %s body: %s, request: %s", err.Error(), *body, string(reqByte))
	}

	out.Roots = make(map[string]H256)
	out.Proofs = make(map[string]string)
	for i, _ := range rpcRep.Result.Roots {
		key := fmt.Sprintf("%s%s", common.HexPreFix, i)
		out.Roots[key] = common.Hex2Bytes(rpcRep.Result.Roots[i])
	}
	for i, _ := range rpcRep.Result.Proofs {
		key := fmt.Sprintf("%s%s", common.HexPreFix, i)
		out.Proofs[key] = fmt.Sprintf("%s%s", common.HexPreFix, rpcRep.Result.Proofs[i])
	}

	return &out, nil
}

func sendAndCheck(url string, req smtServerReq) (*string, error) {
	return sendAndCheckWithTimeout(url, req, TimeOut)
}

func sendAndCheckWithTimeout(url string, req smtServerReq, timeout time.Duration) (*string, error) {
	rpcReq := req
	reqByte, _ := json.Marshal(rpcReq)
	_, body, err := gorequest.New().Post(url).Retry(RetryNumber, RetryTime).Timeout(timeout).SendStruct(&rpcReq).End()
	if err != nil {
		return nil, fmt.Errorf("Smt server request error: %v, %s, request:%s", err, body, string(reqByte))
	}

	repTemp := struct {
		JsonRpc string
		Error   JsonRpcError
	}{}
	if err := json.Unmarshal([]byte(body), &repTemp); err != nil {
		return nil, fmt.Errorf("json Unmarshal err: %s body: %s, request: %s", err.Error(), body, string(reqByte))

	}
	if repTemp.Error.Code != 0 {
		return nil, fmt.Errorf("rpc error: %s, request: %s", repTemp.Error.Message, string(reqByte))
	}

	return &body, nil
}

func newBasicReq(method string) smtServerReq {
	return smtServerReq{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  method,
	}
}
