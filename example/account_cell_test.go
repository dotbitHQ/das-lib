package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestAccountCellDataBuilderFromTx(t *testing.T) {
	bys := common.Blake2b(common.Hex2Bytes("0x110100001c00000030000000fc000000040100000c0100000d010000d1530ccd93bb35916be48e33eeb65eea6ef75c42cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000039150000000c00000010000000020000000100000062150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c000000100000000200000001000000610452d4600000000000000000000000000004000000"))
	fmt.Println(common.Bytes2Hex(bys), len(bys))
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x63254ca08bbbd3304809e062bb8843a172c5dca5ba50d342ac9911b2c0238c17"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMap, err := witness.AccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range builderMap {
			if v.Index == 0 {
				fmt.Println(v.Account, v.Version, v.AccountId)
				fmt.Println(common.Bytes2Hex(v.AccountCellDataV1.AsSlice()))
				//0x290100002400000038000000040100000c010000140100001c0100002401000025010000d1530ccd93bb35916be48e33eeb65eea6ef75c42cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000039150000000c00000010000000020000000100000062150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c000000100000000200000001000000610452d460000000000000000000000000000000000000000000000000000000000004000000
				//0x290100002400000038000000040100000c010000140100001c0100002401000025010000d1530ccd93bb35916be48e33eeb65eea6ef75c42cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000039150000000c00000010000000020000000100000062150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c000000100000000200000001000000610452d460000000000000000000000000000000000000000000000000000000000004000000
			}
			//fmt.Println(v.Account, v.Version,v.Index,common.Bytes2Hex(res.Transaction.OutputsData[v.Index][:32])
		}
		//fmt.Println("==")
		//builderMap, err = witness.AccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeOld)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//for _, v := range builderMap{
		//	fmt.Println(v.Account, v.Version,v.Index,common.Bytes2Hex(res.Transaction.OutputsData[v.Index][:32]))
		//}
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
	_, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	accountChars, _ := common.AccountToAccountChars("00😊001😊.0001.bit")
	//moleculeAccountChars := witness.ConvertToAccountChars(accountChars)
	//account := common.AccountCharsToAccount(moleculeAccountChars)
	fmt.Println(accountChars)

	//accountChars, _ = common.AccountToAccountChars("metavs.bit")
	//moleculeAccountChars = witness.ConvertToAccountChars(accountChars)
	//account = common.AccountCharsToAccount(moleculeAccountChars)
	//fmt.Println(account, accountChars)

}

func TestAccountCellGenWitness(t *testing.T) {
	accountId, _ := molecule.AccountIdFromSlice(common.Hex2Bytes("0xc475fcded6955abc8bf6e2f23e68c6912159505d"), true)
	accountCharSet, _ := common.AccountToAccountChars("7aaaaaaa.bit")
	records := witness.ConvertToCellRecords([]witness.Record{{
		Key:   "eth",
		Type:  "address",
		Label: "label",
		Value: "0xc9f53b1d85356B60453F867610888D89a0B667Ad",
		TTL:   300,
	}})

	accountCellData := molecule.NewAccountCellDataBuilder().
		Id(*accountId).
		Account(*common.ConvertToAccountChars(accountCharSet)).
		RegisteredAt(molecule.GoU64ToMoleculeU64(1624345781)).
		LastTransferAccountAt(molecule.Uint64Default()).
		LastEditManagerAt(molecule.Uint64Default()).
		LastEditRecordsAt(molecule.Uint64Default()).
		Status(molecule.GoU8ToMoleculeU8(0)).
		Records(*records).
		EnableSubAccount(molecule.GoU8ToMoleculeU8(1)).
		RenewSubAccountPrice(molecule.GoU64ToMoleculeU64(100000000)).
		//Dev1(molecule.GoString2MoleculeBytes("dev1")).
		//Dev2(molecule.GoString2MoleculeBytes("dev2")).
		Build()
	builder := witness.AccountCellDataBuilder{
		Version:         5,
		AccountCellData: &accountCellData,
	}
	wit, witHash, err := builder.GenWitness(&witness.AccountCellParam{
		OldIndex: 0,
		NewIndex: 0,
		Status:   0,
		Action:   common.DasActionRenewAccount,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(wit))
	fmt.Println(common.Bytes2Hex(witHash))
}

func TestAccountCellWitnessParser(t *testing.T) {
	witnessByte := common.Hex2Bytes("0x64617301000000cf010000100000001000000010000000bf0100001000000014000000180000000000000003000000a3010000a30100002c000000400000000c010000140100001c010000240100002c0100002d0100009a0100009b010000c475fcded6955abc8bf6e2f23e68c6912159505dcc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000037150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061b58cd16000000000000000000000000000000000000000000000000000000000006d000000080000006500000018000000230000002a0000003300000061000000070000006164647265737303000000657468050000006c6162656c2a0000003078633966353362316438353335364236303435334638363736313038383844383961304236363741642c0100000100e1f50500000000")
	b, err := json.Marshal(witness.ParserWitnessData(witnessByte))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestGetSatisfiedCapacityLiveCellWithOrder(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	dasLockScript := types.Script{
		CodeHash: types.HexToHash("0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd"),
		HashType: "type",
		Args:     common.Hex2Bytes("0x04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d"),
	}
	//dasTypeScript := types.Script{
	//	CodeHash: types.HexToHash("0x4ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c"),
	//	HashType: "type",
	//	Args:     nil,
	//}
	res, _, err := core.GetSatisfiedCapacityLiveCellWithOrder(dc.Client(), nil, &dasLockScript, nil, 1000, 100, indexer.SearchOrderDesc)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res {
		fmt.Println(v.BlockNumber, v.OutPoint.TxHash.String(), v.OutPoint.Index)
	}
}

func TestGetBalanceCells(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	dasLockScript := types.Script{
		CodeHash: types.HexToHash("0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd"),
		HashType: "type",
		Args:     common.Hex2Bytes("0x04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d"),
	}
	res, total, err := dc.GetBalanceCells(&core.ParamGetBalanceCells{
		DasCache:          nil,
		LockScript:        &dasLockScript,
		CapacityNeed:      0,
		CapacityForChange: 100,
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total:", total)
	for _, v := range res {
		fmt.Println(v.BlockNumber, v.OutPoint.TxHash.String(), v.OutPoint.Index)
	}
}
