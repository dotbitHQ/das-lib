package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

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

func TestTxDPInfo(t *testing.T) {
	//dc, err := getNewDasCoreTestnet2()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//hash := "0xa7e780250f5db774f2fcae028e7b6f44bb6e7b04ed8b2cb1beb6f4f7e969295c"
	//tx, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//molecule.Bytes2GoU64()
	data := common.Hex2Bytes("0x080000006400000000000000")
	res, err := witness.ConvertBysToDPData(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Value)
	bys, err := witness.ConvertDPDataToBys(res)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
	//res, err := dc.GetOutputsDPInfo(tx.Transaction)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for k, v := range res {
	//	fmt.Println(k, toolib.JsonString(v))
	//}
}

func TestGetDpCells(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	addr := "ckt1qsexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6pgy55pufwt8rrg5v6vg08d2dm2wekv7djc9qjjs839evuvdz3nf3pua4fhdfmxenektlczhev"
	pAddr, err := address.Parse(addr)
	if err != nil {
		t.Fatal(err)
	}
	AmountNeed := uint64(10) * common.UsdRateBase
	list, dpAmount, capacityAmount, err := dc.GetDpCells(&core.ParamGetDpCells{
		DasCache:           nil,
		LockScript:         pAddr.Script,
		AmountNeed:         AmountNeed,
		CurrentBlockNumber: 0,
		SearchOrder:        indexer.SearchOrderDesc,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(list), dpAmount, capacityAmount)

	//

	normal := "ckt1qyqre7f5hpeujdlq5q9xvj59f6qq5nemar8qv73xan"
	nAddr, _ := address.Parse(normal)
	outputs, outputsData, normalCell, err := dc.SplitDPCell(&core.ParamSplitDPCell{
		FromLock:           pAddr.Script,
		ToLock:             pAddr.Script,
		DPLiveCell:         list,
		DPLiveCellCapacity: capacityAmount,
		DPTotalAmount:      dpAmount,
		DPTransferAmount:   15 * common.UsdRateBase,
		DPSplitCount:       2,
		DPSplitAmount:      25 * common.UsdRateBase,
		NormalCellLock:     nAddr.Script,
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range outputs {
		fmt.Println(common.Bytes2Hex(v.Lock.Args), v.Lock.OccupiedCapacity())
		fmt.Println(common.Bytes2Hex(outputsData[i]))
	}
	fmt.Println(len(outputs), normalCell)

}

func TestGetOutputsDPInfo(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xafc34f958c620d5d8f5f2ef1df95eb655a831bc97012ac9418e9c3f820202cc8"
	tx, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash))
	if err != nil {
		t.Fatal(err)
	}
	res, err := dc.GetOutputsDPInfo(tx.Transaction)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range res {
		fmt.Println("res:", k, v.AlgId, v.SubAlgId, v.Payload, v.AmountDP)
	}
}
