package example

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
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
	hash := "670a62465d46d3088832a009dbcbe4c1b584a68b958eaec664954fc23c7080ae"
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
	for _, v := range data.Vout {
		//fmt.Println(v.ScriptPubKey.Addresses,v.ScriptPubKey)
		fmt.Println("hex:", common.Hex2Bytes("0x6a"), common.Hex2Bytes(v.ScriptPubKey.Hex)[2:], []byte("test"))
	}

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
	//res, v, err := base58.CheckDecode("")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(hex.EncodeToString(res), v)
	//
	//fmt.Println(base58.CheckEncode(res, v))
	addr, err := common.Base58CheckDecode("", common.DogeCoinBase58Version)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr)
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
	_, uos, err := txTool.GetUnspentOutputsDoge(addr, privateKey, 1000000000)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txTool.NewTx(uos, []string{addr}, []int64{1000000000}, "test")
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
	msg := "0x4ee3835331f9bb84db2dd2cf674db221ad42484e2c6aa96609911e1f810d38a5"
	privateKey := "0000000000000000000000000000000000000000000000000000000000000001"
	bys, err := sign.DogeSignature([]byte(msg), privateKey, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(bys), bys)
	fmt.Println(hex.EncodeToString(bys))

	payload := "751e76e8199196d454941c45d1b3a323f1433bd6" // com
	////payload = "91b24bf9f5288532960ac687abb035127b1d28a5"
	//payload = "500de0c9a7c7777e02ab8e0e86c9f55bda5df756"
	////payload = "03da55778c9d441ea212cdbfb4f8f1c1cebd6f9e"
	//
	//bys = common.Hex2Bytes("0x3fe95c1ed28fa2e80728450024943fe59be28c995370511878d1c08ff376543352e70c8570b84a0c68da92b5fcf42abeb21550c21727850c51352a489d5a71530000")
	//data := common.Hex2Bytes("0xb47e70087d6fde994af0d1852e4a0fa558139a6b562a03c2eb00696d4060644a")
	//payload = "b6031be679d6bfa9ce6db1e3bf61b6e6552423be"
	fmt.Println(sign.VerifyDogeSignature(bys, []byte(msg), payload))
	// 6eb037f8db51a81be05d6c1813181384f286b0d52febbb0b28e0a90e6f0c118c08a248925efd5d6bb51678246f46a25fbcbd5f64791e34e5ca2bbe62c9d0a1370101
}

func TestDogeSig(t *testing.T) {
	str := "H83e/zo4/m1MtX55jc//gp0yyMGUDgK0bmkpylRPbCNyF53kLwmGQhyowkTz9JhpDUO+xyH0R3xRPx/HWxz7hKM="
	str = "G6k+dZwJ8oOei3PCSpdj603fDvhlhQ+sqaFNIDvo/bI+Xh6zyIKGzZpyud6YhZ1a5mcrwMVtTWL+VXq/hC5Zj7s="
	str = "IG6wN/jbUagb4F1sGBMYE4TyhrDVL+u7CyjgqQ5vDBGMCKJIkl79XWu1Fngkb0aiX7y9X2R5HjTlyiu+YsnQoTc="
	str = "H0xeoog3VBjVl98NCBxC+Y/zApnjtTyeQNf92BD6IIBSYUFAlnXXfAk/u6ygQZPbWi4fP9+MFZRAp/jJHU01s0U="
	//str = "IOiJaq+w2tcbkKRO2XiCG6p6UkHNyvaGP21V1aVKelmYQNfoyRO/kFB54gH+mmaPtt0YpuC9ca4ZI4IecUf2aL4="
	res, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println([]byte("abc123"))
	//0xe8896aafb0dad71b90a44ed978821baa7a5241cdcaf6863f6d55d5a54a7a599840d7e8c913bf905079e201fe9a668fb6dd18a6e0bd71ae1923821e7147f668be0101
	//0x4c3265a41e9ca4b178a6ebe617c27f47f4ddb4983267d351136229e6f155f6f53d8f556f8f1eb82f2554c3cd5bd7aa5b413485d360830066b27b341f92ebf09f0101

	//0x4c3265a41e9ca4b178a6ebe617c27f47f4ddb4983267d351136229e6f155f6f53d8f556f8f1eb82f2554c3cd5bd7aa5b413485d360830066b27b341f92ebf09f0100

	//0x4c5ea288375418d597df0d081c42f98ff30299e3b53c9e40d7fdd810fa2080526141409675d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b3450001
	//0x4c5ea288375418d597df0d081c42f98ff30299e3b53c9e40d7fdd810fa2080526141409675d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b3450001
	//fmt.Println(len(res))
	si, err := sign.DecodeSignature(res)
	if err != nil {
		t.Fatal(err)
	}
	//0x4c5ea288375418d597df0d081c42f98ff30299e3b53c9e40d7fdd810fa2080526141409675d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b3450001
	//0x4c5ea288375418d523df0d081c42f98ff30299e3b53c9e40d7fdd810fa2080526141402275d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b3450001
	//4c5ea288375418d597df0d081c4298ff30299e3b53c9e40d7fdd810fa2080526141409675d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b345
	fmt.Println(si.Compressed, si.Recovery, si.SegwitType, hex.EncodeToString(si.Signature))
	fmt.Println(hex.EncodeToString(si.ToSig()))
	//fmt.Println(hex.EncodeToString(res))
	// false 0 <nil> a93e759c09f2839e8b73c24a9763eb4ddf0ef865850faca9a14d203be8fdb23e5e1eb3c88286cd9a72b9de98859d5ae6672bc0c56d4d62fe557abf842e598fbb
	// true 0 <nil> a93e759c09f2839e8b73c24a9763eb4ddf0ef865850faca9a14d203be8fdb23e5e1eb3c88286cd9a72b9de98859d5ae6672bc0c56d4d62fe557abf842e598fbb
	// a93e759c09f2839e8b73c24a9763eb4ddf0ef865850faca9a14d203be8fdb23e5e1eb3c88286cd9a72b9de98859d5ae6672bc0c56d4d62fe557abf842e598fbb000100
}

func TestVerifyDogeSignature(t *testing.T) {
	payload, err := common.Base58CheckDecode("D9YnEkJGK5HTmRAtf61uyXTYeXNPkhceCg", common.DogeCoinBase58Version)
	if err != nil {
		t.Fatal(err)
	}

	sig := common.Hex2Bytes("4c5ea288375418d597df0d081c42f98ff30299e3b53c9e40d7fdd810fa2080526141409675d77c093fbbaca04193db5a2e1f3fdf8c159440a7f8c91d4d35b3450001")
	data := []byte("0x4ee3835331f9bb84db2dd2cf674db221ad42484e2c6aa96609911e1f810d38a5")
	fmt.Println(payload, len(sig), sig[65:66])
	fmt.Println(sign.VerifyDogeSignature(sig, data, payload))
}
