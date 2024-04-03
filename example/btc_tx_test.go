package example

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/dotbitHQ/das-lib/bitcoin"
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
	blockHash, _ := chainhash.NewHashFromStr("000000000000000000012bd821d9d4baa773aa456b9f33571f7abe2d5d1c26ba")
	block, err := client.GetBlock(blockHash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PrevBlock:", block.Header.PrevBlock.String())
	fmt.Println("BlockHash:", block.Header.BlockHash().String())
	fmt.Println("BlockHash:", block.Header.Timestamp)
	fmt.Println("Transactions:", len(block.Transactions))
	netParams := bitcoin.GetBTCMainNetParams()
	for _, tx := range block.Transactions {
		//fmt.Println(tx.TxHash())
		if tx.TxHash().String() == "36e1d97cc0fe7b439aaaa3b9a0237d9c9d29a7846e30148ab888caa3384a79a3" {
			fmt.Println("=========")
			for _, txOut := range tx.TxOut {
				_, addrList, _, err := txscript.ExtractPkScriptAddrs(txOut.PkScript, &netParams)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(addrList[0], txOut.Value)
			}
			fmt.Println("=========")
			for _, txIn := range tx.TxIn {
				fmt.Println("==")
				pkScript, err := txscript.ComputePkScript(txIn.SignatureScript, txIn.Witness)
				if err != nil {
					t.Fatal(err)
				}
				addr, err := pkScript.Address(&netParams)
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(addr)
			}

			break
		}
	}
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
		Params:       bitcoin.GetBTCMainNetParams(),
		RpcClientBTC: client,
	}

	//var uos []bitcoin.UnspentOutputs
	// get uos
	addr := ""
	privateKey := ""
	_, uos, err := bitcoin.GetUnspentOutputsBtc(addr, privateKey, "", "", 100000)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txTool.NewBTCTx(uos, []string{addr}, []int64{50000}, "test")
	if err != nil {
		t.Fatal(err)
	}

	res, err := txTool.LocalSignTx(tx, uos)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)

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
	addr := ""
	apiKey := ""
	url := "https://btc.nownodes.io/api/v2/utxo"
	value := int64(7323360)
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
