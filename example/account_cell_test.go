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
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x71141fd011235ef06b8bb6640ac14c23afe7d0ed657b2771fb828d320a21fc80"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.AccountCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Version, builder.Account)
		fmt.Println(builder.RecordList())
		fmt.Println(builder.NextAccountId)
	}
}

func TestAccountCellDataBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x3cefd87b4c0102e3679ea456ac3766df6028296ba7e2d51185ccc5a29399ec49"
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
