package example

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/rpcclient"
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
	//baseRep := bitcoin.BaseRequest{
	//	RpcUrl:   node,
	//	User:     "root",
	//	Password: "root",
	//	Proxy:    "socks5://127.0.0.1:8838",
	//}
	//
	//var data bitcoin.BlockChainInfo
	//err := baseRep.Request(bitcoin.RpcMethodGetBlockChainInfo, nil, &data)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(data)

	client, err := getBtcClient(node)
	if err != nil {
		t.Fatal(err)
	}

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("blockCount:", blockCount)

	//
	blockHash, err := client.GetBlockHash(836782)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("blockHash:", blockHash)

	//
	block, err := client.GetBlock(blockHash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PrevBlock:", block.Header.PrevBlock.String())
	fmt.Println("BlockHash:", block.Header.BlockHash().String())
	fmt.Println("BlockHash:", block.Header.Timestamp)
	fmt.Println("Transactions:", len(block.Transactions))

	client.ListUnspent()
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
