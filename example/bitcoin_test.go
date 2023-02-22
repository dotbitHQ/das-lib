package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
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

func TestDecodeWIF(t *testing.T) {
	wif, err := btcutil.DecodeWIF("")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(wif.PrivKey.Serialize()))
}

func TestHexToPrivateKey(t *testing.T) {
	bys, _, err := bitcoin.HexToPrivateKey(bitcoin.GetDogeMainNetParams(), "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bys))
}

func TestCreateDogeWallet(t *testing.T) {
	if err := bitcoin.CreateDogeWallet(); err != nil {
		t.Fatal(err)
	}
}

func TestFormatDogeAddress(t *testing.T) {
	res, v, err := base58.CheckDecode("D8tA4yZjXexxXTDLDPkUUe2fwd4a2FU77T")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(res), v)

	fmt.Println(base58.CheckEncode(res, common.DogeCoinBase58Version))
}

func TestRpcMethodEstimateFee(t *testing.T) {
	baseRep := getRpcClient()
	var fee float64
	err := baseRep.Request(bitcoin.RpcMethodEstimateFee, []interface{}{10}, &fee)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fee) //0.01003342
}

func TestRpcMethodEstimateSmartFee(t *testing.T) {
	baseRep := getRpcClient()
	//var fee float64
	err := baseRep.Request(bitcoin.RpcMethodEstimateSmartFee, []interface{}{10}, nil)
	if err != nil {
		t.Fatal(err)
	}
	//{"feerate":0.01003339,"blocks":10}
	//fmt.Println(fee) //0.01003342
}

func TestGetUnspentOutputsDoge(t *testing.T) {
	var txTool bitcoin.TxTool

	uos, err := txTool.GetUnspentOutputsDoge("DMjVFBqbqZGAyTXgkt7fTuqihhCCVuLwZ6", 7700000000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(uos)
}

func TestNewTx(t *testing.T) {
	baseRep := getRpcClient()

	//client, err := rpc.Dial("")
	//if err != nil {
	//	t.Fatal(err)
	//}

	txTool := bitcoin.TxTool{
		RpcClient: baseRep,
		Ctx:       context.Background(),
		//RemoteSignClient: client,
		DustLimit:  bitcoin.DustLimitDoge,
		Params:     bitcoin.GetDogeMainNetParams(),
		PrivateKey: "", // note
	}

	//var uos []bitcoin.UnspentOutputs
	// get uos
	uos, err := txTool.GetUnspentOutputsDoge("", 3400000000)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txTool.NewTx(uos, []string{""}, []int64{1000000000})
	if err != nil {
		t.Fatal(err)
	}

	res, err := txTool.LocalSignTx(tx, uos)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)

	hash, err := txTool.SendTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("hash:", hash)
}

func TestRpcMethodDecodeRawTransaction(t *testing.T) {
	baseRep := getRpcClient()
	raw := "01000000017fc48d5025a50bb55359a0eca2cedf83cf2af44e71ff49aae0183f26c1de0e23020000006b483045022100d50313af8dff46f014d58e60c8c37c0f2bcb961f18692a707746af72aa01249b0220250bcb5719597032afd69015860a88966b2e9435072aad9f79bc0570f83ff2ab012102a18b81e15f6d7739683c1e39628419ec04ae32d221d6f8bd6fcdad0a9ff07340ffffffff0200ca9a3b000000001976a914b6031be679d6bfa9ce6db1e3bf61b6e6552423be88ac14886089010000001976a914b6031be679d6bfa9ce6db1e3bf61b6e6552423be88ac00000000"

	err := baseRep.Request(bitcoin.RpcMethodDecodeRawTransaction, []interface{}{raw}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBalanceDoge(t *testing.T) {
	res, err := bitcoin.GetBalanceDoge("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Balance, res.Confirmed, res.Unconfirmed)
}
