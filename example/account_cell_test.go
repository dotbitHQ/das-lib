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
	record := molecule.NewRecordBuilder().
		RecordKey(molecule.GoString2MoleculeBytes("eth")).
		RecordType(molecule.GoString2MoleculeBytes("address")).
		RecordLabel(molecule.GoString2MoleculeBytes("label")).
		RecordValue(molecule.GoString2MoleculeBytes("0xc9f53b1d85356B60453F867610888D89a0B667Ad")).
		RecordTtl(molecule.GoU32ToMoleculeU32(300)).Build()
	records := molecule.NewRecordsBuilder().Push(record).Build()

	accountCellData := molecule.NewAccountCellDataBuilder().
		Id(*accountId).
		Account(*common.ConvertToAccountChars(accountCharSet)).
		RegisteredAt(molecule.GoU64ToMoleculeU64(1624345781)).
		LastTransferAccountAt(molecule.Uint64Default()).
		LastEditManagerAt(molecule.Uint64Default()).
		LastEditRecordsAt(molecule.Uint64Default()).
		Status(molecule.GoU8ToMoleculeU8(0)).
		Records(records).
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
