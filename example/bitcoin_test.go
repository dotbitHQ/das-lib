package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func getRpcClient() *bitcoin.BaseRequest {
	baseRep := bitcoin.BaseRequest{
		RpcUrl:   "",
		User:     "",
		Password: "",
		Proxy:    "socks5://127.0.0.1:8838",
	}
	return &baseRep
}

func TestRpcGetBlockChainInfo(t *testing.T) {
	baseRep := getRpcClient()
	var data bitcoin.BlockChainInfo
	err := baseRep.Request(bitcoin.RpcMethodGetBlockChainInfo, nil, &data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}

func TestRpcGetBlockHash(t *testing.T) {
	baseRep := getRpcClient()
	var blockHash string
	err := baseRep.Request(bitcoin.RpcMethodGetBlockHash, []interface{}{4600472}, &blockHash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(blockHash)
}

func TestRpcGetBlock(t *testing.T) {
	baseRep := getRpcClient()
	var data bitcoin.BlockInfo
	hash := "5d0954672b3d7bc9becbfa017f7cb47714c39ef74ab99c969217ee2af0d40a82"
	err := baseRep.Request(bitcoin.RpcMethodGetBlock, []interface{}{hash}, &data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}

func TestRpcGetRawTransaction(t *testing.T) {
	baseRep := getRpcClient()
	var data bitcoin.RawTransaction
	hash := "24eab97067999aab06b3f95854b8aef653db029ad1ef706e5b05bbf21d4b3f3c"
	err := baseRep.Request(bitcoin.RpcMethodGetRawTransaction, []interface{}{hash, true}, &data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", data)
}

func TestRpcListUnspent(t *testing.T) {
	baseRep := getRpcClient()
	//var data bitcoin.RawTransaction
	req := []interface{}{1, 9999999, []string{"AC8Q9Z4i4sXcbW7TV1jqrjG1JEWMdLyzcy"}, 0}
	err := baseRep.Request(bitcoin.RpcMethodListUnspent, req, nil)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%+v", data)
}

func TestRpcSendRawTransaction(t *testing.T) {
	baseRep := getRpcClient()
	req := []interface{}{"11", false}
	var data string
	err := baseRep.Request(bitcoin.RpcMethodSendRawTransaction, req, &data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(data)
}

func TestCreateDogecoinWallet(t *testing.T) {
	if err := bitcoin.CreateDogecoinWallet(); err != nil {
		t.Fatal(err)
	}
	//PubKey: 290d35c7ec8193604a44bc6d1b96cac0e1ce4dd3
	//PubKey: D8tA4yZjXexxXTDLDPkUUe2fwd4a2FU77T
	//WIF: QTLxZ1Td7U3i74yV21cpcQoFVABjsgVvMDswwAYvMKTYAQfNmQDt
	//PriKey: 8d205c955cabfe0f4c931f34cbc1ca13df515c5440f1b2af4af0318bf4c29396
}

func TestFormatDogecoinAddress(t *testing.T) {
	payload, err := bitcoin.FormatAddressToPayload("D8tA4yZjXexxXTDLDPkUUe2fwd4a2FU77T")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(payload)

	addr, err := bitcoin.FormatPayloadToAddress(common.DasAlgorithmIdDogecoin, payload)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr)
}
