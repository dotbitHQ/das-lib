package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/chain/chain_evm"
	"github.com/dotbitHQ/das-lib/chain/chain_tron"
	"testing"
)

func TestEVM(t *testing.T) {
	node := ""
	chainEVM, err := chain_evm.NewChainEvm(context.Background(), node, 0)
	if err != nil {
		t.Fatal(err)
	}

	block, err := chainEVM.GetBlockByNumber(16930167)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range block.Transactions {
		fmt.Println(v.Hash, v.Value, v.To, v.From)
	}
}

func TestTron(t *testing.T) {
	chainTron, err := chain_tron.NewChainTron(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	balance, err := chainTron.GetBalance("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance)
}
