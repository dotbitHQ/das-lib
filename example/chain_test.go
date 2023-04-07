package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/chain/chain_evm"
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
