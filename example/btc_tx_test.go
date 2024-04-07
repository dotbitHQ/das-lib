package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/shopspring/decimal"
	"testing"
)

var (
	node = ""
)

func getBtcClient(node string) (*rpcclient.Client, error) {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         node,
		User:         "root",
		Pass:         "root",
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   false, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	return client, err
}

func TestBtcRpc(t *testing.T) {
	client, err := getBtcClient(node)
	if err != nil {
		t.Fatal(err)
	}
	blockCount, err := client.GetBlockCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("blockCount:", blockCount)

	fee, err := client.EstimateSmartFee(10, &btcjson.EstimateModeConservative)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("fee:", *fee.FeeRate)

	h, _ := chainhash.NewHashFromStr("811b0bdc3047cc2c4af1f11015889cad5c1787e27ecd06a22b5a7df726cc0e41")
	tx, err := client.GetRawTransaction(h)
	if err != nil {
		t.Fatal(err)
	}
	netParams := bitcoin.GetBTCTestNetParams()
	_, addrList, _, err := txscript.ExtractPkScriptAddrs(tx.MsgTx().TxOut[0].PkScript, &netParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(addrList), addrList[0])
	decValue := decimal.NewFromInt(tx.MsgTx().TxOut[0].Value)
	fmt.Println("decValue:", decValue)

	pkScript, err := txscript.ComputePkScript(tx.MsgTx().TxIn[0].SignatureScript, tx.MsgTx().TxIn[0].Witness)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := pkScript.Address(&netParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("addr:", addr)

	var orderId string
	txOut := tx.MsgTx().TxOut[len(tx.MsgTx().TxOut)-1]
	if txscript.IsNullData(txOut.PkScript) && len(txOut.PkScript) > 32 {
		orderId = hex.EncodeToString(txOut.PkScript[len(txOut.PkScript)-32:])
	} else {
		orderId = string(txOut.PkScript[2:])
	}
	data := []byte("test")
	fmt.Println("txOut.PkScript:", len(txOut.PkScript), len(data), string(txOut.PkScript))
	fmt.Println("orderId:", orderId)

	//for _, v := range tx.MsgTx().TxOut {
	//	fmt.Println("PkScript:", hex.EncodeToString(v.PkScript))
	//}
	//for _, v := range tx.MsgTx().TxIn {
	//	fmt.Println(len(v.Witness))
	//	for _, j := range v.Witness {
	//		fmt.Println("witness:", hex.EncodeToString(j))
	//	}
	//}
	//witness: 304402206e0bce0676946ced5e02afbea79a5c40f8d7a14d6bae5fada627c6b6d0024c90022013ddb28b0b6ddc5021a2802135c3b7f5015539b41f8bcd95c9a785733616b19b01
	//witness: 0262c6eb28bc42cc168f61319dfa54fa64267bc3626ab05094cd1195fdf49a3009

	//// Get the current block count.
	//blockCount, err := client.GetBlockCount()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("blockCount:", blockCount)
	//
	////
	//blockHash, err := client.GetBlockHash(836783) //836782)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("blockHash:", blockHash)
	//
	////
	//blockHash, _ := chainhash.NewHashFromStr("000000000000000000012bd821d9d4baa773aa456b9f33571f7abe2d5d1c26ba")
	//block, err := client.GetBlock(blockHash)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("PrevBlock:", block.Header.PrevBlock.String())
	//fmt.Println("BlockHash:", block.Header.BlockHash().String())
	//fmt.Println("BlockHash:", block.Header.Timestamp)
	//fmt.Println("Transactions:", len(block.Transactions))
	//netParams := bitcoin.GetBTCMainNetParams()
	//for _, tx := range block.Transactions {
	//	//fmt.Println(tx.TxHash())
	//	if tx.TxHash().String() == "36e1d97cc0fe7b439aaaa3b9a0237d9c9d29a7846e30148ab888caa3384a79a3" {
	//		fmt.Println("=========")
	//		for _, txOut := range tx.TxOut {
	//			_, addrList, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, &netParams)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			fmt.Println(addrList[0], txOut.Value)
	//		}
	//		fmt.Println("=========")
	//		for _, txIn := range tx.TxIn {
	//			fmt.Println("==")
	//			pkScript, err := txscript.ComputePkScript(txIn.SignatureScript, txIn.Witness)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			addr, err := pkScript.Address(&netParams)
	//			if err != nil {
	//				t.Fatal(err)
	//			}
	//			fmt.Println(addr)
	//		}
	//
	//		break
	//	}
	//}
}

func TestNewBTCTx(t *testing.T) {
	client, err := getBtcClient(node)
	if err != nil {
		t.Fatal(err)
	}

	txTool := bitcoin.TxTool{
		RpcClient:    nil,
		Ctx:          context.Background(),
		DustLimit:    bitcoin.DustLimitBtc,
		Params:       bitcoin.GetBTCTestNetParams(),
		RpcClientBTC: client,
	}

	//var uos []bitcoin.UnspentOutputs
	// get uos
	addr := "tb1qg56gqsanept494plpa3hdhhmtl2ejyrq5sfw42"
	privateKey := ""
	url := "https://btcbook-testnet.nownodes.io/api/v2/utxo"
	apiKey := ""
	_, uos, err := bitcoin.GetUnspentOutputsBtc(addr, privateKey, url, apiKey, 1000)
	if err != nil {
		t.Fatal(err)
	}
	// 001445348043b3c85752d43f0f6376defb5fd5991060
	//
	//PkScript: 001445348043b3c85752d43f0f6376defb5fd5991060
	//PkScript: 0014e6c61a595983dafb9282d3eb510d517d024160be

	toAddr := "tb1qumrp5k2es0d0hy5z6044zr2305pyzc978qz0ju"
	tx, err := txTool.NewBTCTx(uos, []string{toAddr}, []int64{1000}, "test")
	if err != nil {
		t.Fatal(err)
	}

	res, err := txTool.LocalSignTxWithWitness(tx, uos)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)
	for _, v := range tx.TxOut {
		fmt.Println("txout:", v.Value, hex.EncodeToString(v.PkScript))
	}
	hash, err := txTool.SendBTCTx(tx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("hash:", hash)
}

func TestBTCUTXO(t *testing.T) {
	client, err := getBtcClient(node)
	if err != nil {
		t.Fatal(err)
	}

	min := 1
	max := 100
	addr := []string{"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM"}
	cmd := btcjson.NewListUnspentCmd(&min, &max, &addr)

	var res rpcclient.FutureListUnspentResult
	res = client.SendCmd(cmd)
	list, err := res.Receive()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range list {
		fmt.Println(i, v.Account, v.TxID)
	}
}

func TestGetUnspentOutputsBtc(t *testing.T) {
	addr := "tb1qg56gqsanept494plpa3hdhhmtl2ejyrq5sfw42"
	apiKey := ""
	url := "https://btcbook-testnet.nownodes.io/api/v2/utxo"
	value := int64(1000)
	total, utxoList, err := bitcoin.GetUnspentOutputsBtc(addr, "", url, apiKey, value)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(total)
	for _, v := range utxoList {
		fmt.Println(v.Hash, v.Index, v.Value)
	}
}

func TestExtractPkScriptAddrs(t *testing.T) {
	netParams := bitcoin.GetBTCMainNetParams()
	addr := "bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf"
	decodeAddress, err := btcutil.DecodeAddress(addr, &netParams)
	if err != nil {
		t.Fatal(err)
	}
	script, err := txscript.PayToAddrScript(decodeAddress)
	if err != nil {
		t.Fatal(err)
	}
	c, addrList, n, err := txscript.ExtractPkScriptAddrs(script, &netParams)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c, addrList, n)
}

func Test(t *testing.T) {
	w1, _ := hex.DecodeString("3045022100e1b1243a0960dc06c8c72355626e33d93c400d0aa7725fad2ad5eadc517b986c022036f427360d6c66b94daea673661d5c9baf4b5879618a7df6a5d0c613b194611201")
	w2, _ := hex.DecodeString("03a957d526a85ebe4cd4e0411119dbaeeb9d244831a801462c27876a3ae6e31463")
	w := wire.TxWitness{w1, w2}
	pkScript, err := txscript.ComputePkScript(nil, w)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pkScript.String(), hex.EncodeToString(pkScript.Script()))
	net := bitcoin.GetBTCTestNetParams()
	addr, _ := pkScript.Address(&net)
	fmt.Println(addr.EncodeAddress())
}
