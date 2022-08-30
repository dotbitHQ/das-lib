package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"testing"
)

func TestFixIncomeCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	//dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	asContract, err := core.GetDasContractInfo(common.DasContractNameAlwaysSuccess)
	if err != nil {
		t.Fatal(err)
	}
	inContract, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		t.Fatal(err)
	}

	searchKey := &indexer.SearchKey{
		Script:     inContract.ToScript(nil),
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script: asContract.ToScript(nil),
		},
	}
	cells, err := dc.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderAsc, indexer.SearchLimit, "")
	if err != nil {
		t.Fatal(err)
	}
	var list []string
	fmt.Println("cells:", len(cells.Objects))
	for i, v := range cells.Objects {
		fmt.Println("cell:", i)
		res, err := dc.Client().GetTransaction(context.Background(), v.OutPoint.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		_, err = witness.IncomeCellDataBuilderListFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			list = append(list, v.OutPoint.TxHash.String())
			//fmt.Println("IncomeCellDataBuilderListFromTx: ", err.Error(), v.OutPoint.TxHash.String(), v.OutPoint.Index)
		}
	}
	fmt.Println(len(list))
	fmt.Println(list)
}
