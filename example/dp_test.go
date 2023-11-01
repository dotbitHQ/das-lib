package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestSplitDPCell(t *testing.T) {

	outputs, outputsData, normalCell, err := core.SplitDPCell(&core.ParamSplitDPCell{
		FromLock: &types.Script{
			CodeHash: types.Hash{},
			HashType: "1",
			Args:     nil,
		},
		ToLock: &types.Script{
			CodeHash: types.Hash{},
			HashType: "2",
			Args:     nil,
		},
		DPLiveCell: []*indexer.LiveCell{
			&indexer.LiveCell{
				BlockNumber: 0,
				OutPoint:    nil,
				Output:      nil,
				OutputData:  nil,
				TxIndex:     0,
			},
		},
		DPLiveCellCapacity: 1000,
		DPTotalAmount:      1000,
		DPTransferAmount:   100,
		DPBaseCapacity:     200,
		DPContract:         &core.DasContractInfo{},
		SplitCount:         2,
		DPSplitAmount:      100,
		NormalCellLock: &types.Script{
			CodeHash: types.Hash{},
			HashType: "3",
			Args:     nil,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range outputs {
		fmt.Println(v.Lock.HashType)
		fmt.Println(outputsData[i])
	}
	fmt.Println(normalCell)
}
