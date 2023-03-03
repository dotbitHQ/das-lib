package example

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/sign"
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
	var data btcjson.TxRawResult
	hash := "c76c114ecedf7c006be2d93ab6477558973bcf917dc0fa4719625affeb6aca28"
	err := baseRep.Request(bitcoin.RpcMethodGetRawTransaction, []interface{}{hash, true}, &data)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%+v", data)
	bys, err := json.Marshal(&data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bys))
	fmt.Println(bitcoin.VinScriptSigToAddress(data.Vin[0].ScriptSig, bitcoin.GetDogeMainNetParams()))

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
	//WIF: QNcdLVw8fHkixm6NNyN6nVwxKek4u7qrioRbQmjxac5TVoTtZuot 0000000000000000000000000000000000000000000000000000000000000001
	//PubKey: 751e76e8199196d454941c45d1b3a323f1433bd6
	//PubKey: DFpN6QqFfUm3gKNaxN6tNcab1FArL9cZLE
	//	=======================
	//WIF: 6J8csdv3eDrnJcpSEb4shfjMh2JTiG9MKzC1Yfge4Y4GyUsjdM6 0000000000000000000000000000000000000000000000000000000000000001
	//PubKey: 91b24bf9f5288532960ac687abb035127b1d28a5
	//PubKey: DJRU7MLhcPwCTNRZ4e8gJzDebtG1H5M7pc
	//
	addr := "DFpN6QqFfUm3gKNaxN6tNcab1FArL9cZLE"
	privateKeyHex := "0000000000000000000000000000000000000000000000000000000000000001"
	bys, prvKey, compress, err := bitcoin.HexPrivateKeyToScript(addr, bitcoin.GetDogeMainNetParams(), privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bys), hex.EncodeToString(prvKey.Serialize()), compress)
	//
	addr = "DJRU7MLhcPwCTNRZ4e8gJzDebtG1H5M7pc"
	bys, prvKey, compress, err = bitcoin.HexPrivateKeyToScript(addr, bitcoin.GetDogeMainNetParams(), privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bys), hex.EncodeToString(prvKey.Serialize()), compress)
	//PubKey: c541b148bf600efe206e9b3116dcfbd7f8dc6d16
	//PubKey: DP86MSmWjEZw8GKotxcvAaW5D4e3qoEh6f
}

func TestCreateDogeWallet(t *testing.T) {
	if err := bitcoin.CreateDogeWallet(true); err != nil {
		t.Fatal(err)
	}
}

func TestFormatDogeAddress(t *testing.T) {
	res, v, err := base58.CheckDecode("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(res), v)

	fmt.Println(base58.CheckEncode(res, v))
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

func TestGetUnspentOutputsDoge(t *testing.T) {
	var txTool bitcoin.TxTool

	total, uos, err := txTool.GetUnspentOutputsDoge("", "", 7700000000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(total, uos)
}

func TestNewTx(t *testing.T) {
	baseRep := getRpcClient()

	txTool := bitcoin.TxTool{
		RpcClient: baseRep,
		Ctx:       context.Background(),
		DustLimit: bitcoin.DustLimitDoge,
		Params:    bitcoin.GetDogeMainNetParams(),
	}

	//var uos []bitcoin.UnspentOutputs
	// get uos
	addr := ""
	privateKey := ""
	_, uos, err := txTool.GetUnspentOutputsDoge(addr, privateKey, 3000000000)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txTool.NewTx(uos, []string{"D9YnEkJGK5HTmRAtf61uyXTYeXNPkhceCg"}, []int64{3000000000})
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

func TestDogeSignature(t *testing.T) {

	privateKey := "0000000000000000000000000000000000000000000000000000000000000001"
	decodePrvKey, err := hex.DecodeString(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	prvKey, _ := btcec.PrivKeyFromBytes(decodePrvKey)
	params := bitcoin.GetDogeMainNetParams()

	wif, err := btcutil.NewWIF(prvKey, &params, true)
	if err != nil {
		t.Fatal(err)
	}
	res, err := btcutil.DecodeWIF(wif.String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("WIF:", wif.String(), hex.EncodeToString(res.PrivKey.Serialize()))
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PubKey:", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
	fmt.Println("PubKey:", addressPubKey.EncodeAddress())

	fmt.Println("=======================")
	wif, err = btcutil.NewWIF(prvKey, &params, false)
	if err != nil {
		t.Fatal(err)
	}
	res, err = btcutil.DecodeWIF(wif.String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("WIF:", wif.String(), hex.EncodeToString(res.PrivKey.Serialize()))
	addressPubKey, err = btcutil.NewAddressPubKey(wif.SerializePubKey(), &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PubKey:", hex.EncodeToString(addressPubKey.AddressPubKeyHash().Hash160()[:]))
	fmt.Println("PubKey:", addressPubKey.EncodeAddress())
}

func TestDogeSignature2(t *testing.T) {
	res := common.Hex2Bytes("0xd97772883d1e0e65f2017cf9b5a05a01a3f4f0b93ba4aa970aa78fadaacaf291")
	fmt.Println(string(res), res)
	//base64.

	msg := "vires is numeris"
	privateKey := "0000000000000000000000000000000000000000000000000000000000000001"
	bys, err := sign.DogeSignature([]byte(msg), privateKey, common.Hex2Bytes("1b"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(bys))

	payload := "751e76e8199196d454941c45d1b3a323f1433bd6"
	// 91b24bf9f5288532960ac687abb035127b1d28a5
	fmt.Println(sign.VerifyDogeSignature(bys, []byte(msg), payload))
}
