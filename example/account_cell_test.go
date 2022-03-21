package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestAccountCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xb00f8f1e78723d6e0bdde33838c424fed04e11dc9a59789fcf5483d68e2a7c64"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.AccountCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Version, builder.Account)
		fmt.Println(builder.RecordList())
		fmt.Println(builder.NextAccountId, builder.ExpiredAt)
	}
}

func TestAccountCellDataBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xb1b7a83bc35bc2d3721e612f182ccec88aa8a6de3fd531cb9fa6adb7b01d8979"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMap, err := witness.AccountIdCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range builderMap {
			fmt.Println(k, v.Index)
			//bys := res.Transaction.OutputsData[v.Index]
			//fmt.Println(v.Index, common.Bytes2Hex(bys))
			//fmt.Println(k, v.Version, v.Status, v.AccountId, v.RegisteredAt, v.ExpiredAt)
		}
		tmp := builderMap["0x0000000000000000000000000000000000000000"]
		_, _, _ = tmp.GenWitness(&witness.AccountCellParam{})
	}
}

func TestAccountCellVersionV1(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	acc, _ := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	searchKey := &indexer.SearchKey{
		Script:     acc.ToScript(nil),
		ScriptType: indexer.ScriptTypeType,
	}
	liveCells, _ := dc.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderDesc, 10000, "")

	for k, v := range liveCells.Objects {
		res, _ := dc.Client().GetTransaction(context.Background(), v.OutPoint.TxHash)
		builders, _ := witness.AccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeNew)
		for _, builder := range builders {
			if builder.Index == uint32(v.TxIndex) {
				if builder.Version == 1 {
					fmt.Println("--------------------------------------------")
					fmt.Println(builder.Version, builder.Account)
					fmt.Println(k, v.OutPoint.TxHash)
					fmt.Println()
				}
			}
		}
	}
}

func TestTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	applyTx, err := dc.Client().GetTransaction(context.Background(), types.HexToHash("0x6cb507b9c9eb2a4b794dab9cbb42d5ab6eeefd820aa8d8fd4ed1a007abd00f30"))
	if err != nil {
		t.Fatal(err)
	}
	applyCapacity := applyTx.Transaction.Outputs[0].Capacity
	fmt.Println(applyCapacity)
}

func TestAccountToAccountChars(t *testing.T) {
	accountChars, _ := common.AccountToAccountChars("meta🆚.bit")
	moleculeAccountChars := witness.ConvertToAccountChars(accountChars)
	account := common.AccountCharsToAccount(moleculeAccountChars)
	fmt.Println(account, accountChars)

	accountChars, _ = common.AccountToAccountChars("metavs.bit")
	moleculeAccountChars = witness.ConvertToAccountChars(accountChars)
	account = common.AccountCharsToAccount(moleculeAccountChars)
	fmt.Println(account, accountChars)

}
