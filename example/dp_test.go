package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
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
		DPSplitCount:       2,
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

func TestDPOrderInfo(t *testing.T) {
	info := witness.DPOrderInfo{
		OrderId: "aaa",
		Action:  witness.DPActionTransferDeposit,
	}
	wit, data, err := witness.GenDPOrderInfoWitness(info)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(wit)
	fmt.Println(data)
	orderInfo, err := witness.ConvertDPOrderInfoWitness(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(orderInfo)
}

func TestContract(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	dpContract, err := core.GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dpContract.ContractName, dpContract.ContractTypeId, dpContract.OutPoint.TxHash)
	dpConfigCell, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dpConfigCell.Name, dpConfigCell.OutPoint.TxHash)
	mapT, err := dc.GetDPointTransferWhitelist()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("mapT:")
	for k, v := range mapT {
		fmt.Println(k, common.Bytes2Hex(v.Args))
	}
	mapC, err := dc.GetDPointCapacityRecycleWhitelist()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("mapC:")
	for k, v := range mapC {
		fmt.Println(k, common.Bytes2Hex(v.Args))
	}
}
