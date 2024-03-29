package example

import (
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
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
}
