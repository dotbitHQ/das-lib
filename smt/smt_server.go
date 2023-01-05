package smt

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/parnurzeal/gorequest"
	"time"
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

type GetSmtRootRep struct {
	JsonRpc string       `json:"jsonrpc"`
	Result  string       `json:"result"`
	Error   JsonRpcError `json:"error"`
}

type DeleteSmtRootRep struct {
	JsonRpc string       `json:"jsonrpc"`
	Result  bool         `json:"result"`
	Error   JsonRpcError `json:"error"`
}

type UpdateResult struct {
	Root   string            `json:"root"`
	Proofs map[string]string `json:"proofs"`
}

type UpdateSmtRep struct {
	JsonRpc string
	Result  UpdateResult `json:"result"`
	Error   JsonRpcError
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
	var (
		rpcReq smtServerReq
		rpcRep GetSmtRootRep
	)

	rpcReq = newBasicReq("get_smt_root")
	params := make(map[string]string)
	params["smt_name"] = s.smtName
	rpcReq.Params = params
	reqByte, _ := json.Marshal(rpcReq)
	_, body, errs := gorequest.New().Post(s.url).Timeout(time.Minute * 20).SendStruct(&rpcReq).End()
	if errs != nil {
		return nil, fmt.Errorf("GetSmtRoot Smt server request error: %v, %s, reqest: %s", errs, body, string(reqByte))
	}
	err := json.Unmarshal([]byte(body), &rpcRep)
	if err != nil {
		return nil, fmt.Errorf("GetSmtRoot Json Unmarshal err: %s body: %s, request: %s", err.Error(), body, string(reqByte))
	}

	if rpcRep.Error.Code != 0 && rpcRep.Error.Message != "" {
		return nil, fmt.Errorf("GetSmtRoot Rpc error: %s, request : %s", rpcRep.Error.Message, string(reqByte))
	}
	return common.Hex2Bytes(rpcRep.Result), nil

}

func (s *SmtServer) DeleteSmt() (bool, error) {
	var (
		rpcReq smtServerReq
		rpcRep DeleteSmtRootRep
	)

	rpcReq = newBasicReq("delete_smt")
	params := make(map[string]string)
	params["smt_name"] = s.smtName
	rpcReq.Params = params
	reqByte, _ := json.Marshal(rpcReq)
	_, body, errs := gorequest.New().Post(s.url).Timeout(time.Minute * 20).SendStruct(&rpcReq).End()
	if errs != nil {
		return false, fmt.Errorf("DeleteSmt Smt server request error: %v, %s, request:%s", errs, body, string(reqByte))
	}
	err := json.Unmarshal([]byte(body), &rpcRep)
	if err != nil {
		return false, fmt.Errorf("DeleteSmt Json Unmarshal err: %s body: %s, request:%s", err.Error(), body, string(reqByte))
	}
	if rpcRep.Error.Code != 0 && rpcRep.Error.Message != "" {
		return false, fmt.Errorf("DeleteSmt Rpc error: %s, request:%s", rpcRep.Error.Message, string(reqByte))
	}
	return true, nil
}

func (s *SmtServer) UpdateSmt(kv []SmtKv, opt SmtOpt) (*UpdateSmtOut, error) {
	var (
		rpcReq smtServerReq
		rpcRep UpdateSmtRep
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
		rpcReq = newBasicReq("update_memory_smt")
	} else {
		rpcReq = newBasicReq("update_db_smt")
		param.SmtName = s.smtName
	}
	rpcReq.Params = param
	reqByte, _ := json.Marshal(rpcReq)

	_, body, errs := gorequest.New().Post(s.url).Timeout(time.Minute * 20).SendStruct(&rpcReq).End()
	if errs != nil {
		return nil, fmt.Errorf("UpdateSmt Smt server request error: %v, %s, request:%s", errs, body, string(reqByte))
	}

	err := json.Unmarshal([]byte(body), &rpcRep)
	if err != nil {
		return nil, fmt.Errorf("UpdateSmt Json Unmarshal err: %s body: %s, request: %s", err.Error(), body, string(reqByte))
	}
	if rpcRep.Error.Code != 0 && rpcRep.Error.Message != "" {
		return nil, fmt.Errorf("UpdateSmt Rpc error: %s, request: %s", rpcRep.Error.Message, string(reqByte))
	}
	out.Root = common.Hex2Bytes(rpcRep.Result.Root)
	out.Proofs = make(map[string]string)

	for i, _ := range rpcRep.Result.Proofs {
		key := fmt.Sprintf("%s%s", common.HexPreFix, i)
		out.Proofs[key] = fmt.Sprintf("%s%s", common.HexPreFix, rpcRep.Result.Proofs[i])
	}
	return &out, nil
}

func newBasicReq(method string) smtServerReq {
	return smtServerReq{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  method,
	}
}
