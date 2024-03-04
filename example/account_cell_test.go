package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestAccountCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x8821fdd6f30f7018c3c9dd9627bf9f321e82830e831bacfa5f3679de08608b5b"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMap, err := witness.AccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range builderMap {
			//list := common.ConvertToAccountCharSets(v.AccountChars)
			fmt.Println(v.Account, v.ExpiredAt)
			//var resMap = make(map[common.AccountCharType]struct{})
			//common.GetAccountCharType(resMap, list)
			//var num uint64
			//for k, _ := range resMap {
			//	numTmp := common.AccountCharTypeToUint64(k)
			//	num += numTmp
			//}
			//fmt.Println(num)
			//fmt.Println(common.Uint64ToAccountCharType(num))
		}
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
	acc := "ぁ0001.ぁぁ123ぁぁ.bit"
	accountChars, _ := common.AccountToAccountChars(acc[:strings.Index(acc, ".")])
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
	addr, _ := address.Parse("ckt1qyqre7f5hpeujdlq5q9xvj59f6qq5nemar8qv73xan")
	res, total, err := dc.GetBalanceCellsFilter(&core.ParamGetBalanceCells{
		DasCache:           nil,
		LockScript:         addr.Script,
		CapacityNeed:       0,
		CapacityForChange:  0,
		SearchOrder:        indexer.SearchOrderDesc,
		OutputDataLenRange: nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("total:", total)
	for _, v := range res {
		fmt.Println(v.BlockNumber, v.OutPoint.TxHash.String(), v.OutPoint.Index, v.Output.Capacity, len(v.OutputData))
	}
}

func TestAccount(t *testing.T) {
	fmt.Println(common.GetAccountLength("ให้บริการ.bit"))
	fmt.Println(len([]rune("ให้บริการ")), utf8.RuneCountInString("ให้บริการ"))
}

func TestLockAccount(t *testing.T) {
	str := ``
	list := strings.Split(str, "\n")
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range list {
		ou := common.String2OutPointStruct(v)
		if res, err := dc.Client().GetTransaction(context.Background(), ou.TxHash); err != nil {
			t.Fatal(err)
		} else {
			acc, err := witness.AccountCellDataBuilderFromTx(res.Transaction, common.DataTypeOld)
			if err != nil {
				t.Fatal(err)
			}
			accTx, err := dc.Client().GetTransaction(context.Background(), res.Transaction.Inputs[acc.Index].PreviousOutput.TxHash)
			if err != nil {
				t.Fatal(err)
			}
			ownerScript := accTx.Transaction.Outputs[res.Transaction.Inputs[acc.Index].PreviousOutput.Index].Lock
			owner, _, err := dc.Daf().ScriptToHex(ownerScript)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(fmt.Sprintf("%s,%s", acc.Account, owner.AddressHex))
		}
	}
}

func TestBalance(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	asContract, err := core.GetDasContractInfo(common.DasContractNameAlwaysSuccess)
	if err != nil {
		t.Fatal(err)
	}
	preContract, err := core.GetDasContractInfo(common.DasContractNamePreAccountCellType)
	if err != nil {
		t.Fatal(err)
	}

	searchKey := indexer.SearchKey{
		Script:     asContract.ToScript(nil),
		ScriptType: indexer.ScriptTypeLock,
		ArgsLen:    0,
		Filter: &indexer.CellsFilter{
			Script:              preContract.ToScript(nil),
			OutputDataLenRange:  nil,
			OutputCapacityRange: nil,
			BlockRange:          nil,
		},
	}
	searchKey.Filter.BlockRange = &[2]uint64{8303666, 8303668}
	res, err := dc.Client().GetCells(context.Background(), &searchKey, indexer.SearchOrderAsc, 200, "")

	//res, err := dc.Client().GetCells(context.Background(), &indexer.SearchKey{
	//	Script: &types.Script{
	//		CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
	//		HashType: types.HashTypeType,
	//		Args:     common.Hex2Bytes("0x75bebe0707641658cb8020b9233de32c20c3e172"),
	//	},
	//	ScriptType: indexer.ScriptTypeLock,
	//	ArgsLen:    0,
	//	Filter:     &indexer.CellsFilter{
	//		Script:              nil,
	//		OutputDataLenRange:  nil,
	//		OutputCapacityRange: nil,
	//		BlockRange:          nil,
	//	},
	//}, indexer.SearchOrderDesc, indexer.SearchLimit, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res.Objects {
		fmt.Println(v.OutPoint.TxHash.String(), v.OutPoint.Index)
	}
}

func TestRenewIncome(t *testing.T) {
	str := ``
	list := strings.Split(str, "\n")
	fmt.Println(len(list))

	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}

	var mapIncome = make(map[int]uint64)
	var mapCount = make(map[int]int)
	for i, v := range list {
		tx, err := dc.Client().GetTransaction(context.Background(), common.String2OutPointStruct(v).TxHash)
		if err != nil {
			t.Log(err, v)
			continue
		}
		accNew, err := witness.AccountCellDataBuilderFromTx(tx.Transaction, common.DataTypeNew)
		if err != nil {
			t.Log(err, v)
			continue
		}
		accO, err := witness.AccountCellDataBuilderFromTx(tx.Transaction, common.DataTypeOld)
		if err != nil {
			t.Log(err, v)
			continue
		}
		txOld, err := dc.Client().GetTransaction(context.Background(), tx.Transaction.Inputs[accO.Index].PreviousOutput.TxHash)
		if err != nil {
			t.Log(err, v)
			continue
		}
		accOldMap, err := witness.AccountIdCellDataBuilderFromTx(txOld.Transaction, common.DataTypeNew)
		if err != nil {
			t.Log(err, v)
			continue
		}
		accOld := accOldMap[accNew.AccountId]
		renewYears := (accNew.ExpiredAt - accOld.ExpiredAt) / uint64(common.OneYearSec)
		fmt.Println(i, renewYears, accOld.AccountChars.Len(), accOld.Account, accNew.ExpiredAt, accOld.ExpiredAt)
		if accOld.AccountChars.Len() == 4 {
			mapCount[4]++
			mapIncome[4] += renewYears
		} else {
			mapCount[5]++
			mapIncome[5] += renewYears
		}
	}
	for k, v := range mapCount {
		fmt.Println(k, v)
	}
	for k, v := range mapIncome {
		fmt.Println(k, v)
	}
}

func TestAccountApprovalFromSlice(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	lock, _, err := daf.HexToScript(core.DasAddressHex{
		DasAlgorithmId: common.ChainTypeEth.ToDasAlgorithmId(true),
		AddressHex:     "0xe58673b9bF0a57398e0C8A1BDAe01EEB730177C8",
		IsMulti:        false,
		ChainType:      common.ChainTypeEth,
	})
	if err != nil {
		t.Fatal(err)
	}
	accountApproval := &witness.AccountApproval{
		Action: witness.AccountApprovalActionTransfer,
		Params: witness.AccountApprovalParams{
			Transfer: witness.AccountApprovalParamsTransfer{
				PlatformLock:     lock,
				ProtectedUntil:   1691644224,
				SealedUntil:      1691722652,
				DelayCountRemain: 1,
				ToLock:           lock,
			},
		},
	}
	accountApprovalMol, err := accountApproval.GenToMolecule()
	if err != nil {
		t.Fatal(err)
	}

	accApproval, err := witness.AccountApprovalFromSlice(accountApprovalMol.AsSlice())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accApproval)
}

func TestParseAccountApproval(t *testing.T) {
	accApproval, err := witness.AccountApprovalFromSlice(common.Hex2Bytes("0x030100000c00000018000000080000007472616e73666572e7000000e700000018000000770000007f00000087000000880000005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a00000003deefc10a42cd84c072f2b0e2fa99061a74a0698c03deefc10a42cd84c072f2b0e2fa99061a74a0698c4b93d464000000009ca3d56400000000005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a0000000352045950a5b582e9b426ad89296c8970c96d09d90352045950a5b582e9b426ad89296c8970c96d09d9"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(accApproval)
}

func TestRecycle(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	tx, err := dc.Client().GetTransaction(context.Background(), types.HexToHash("0x90d30bee4fb34dbab3be1cda84c38d4f44fc2d5658e38e733b7fb2ef9e798efc"))
	if err != nil {
		t.Fatal(err)
	}
	accMap, err := witness.AccountIdCellDataBuilderFromTx(tx.Transaction, common.DataTypeNew)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range accMap {
		fmt.Println(k, v.Account)
	}
	accbuilder := accMap["0xfffe71bd9c8a662eac263ba9846b3ed55a4479ec"]

	preParam := &witness.AccountCellParam{
		OldIndex:  0,
		NewIndex:  0,
		Action:    common.DasActionRecycleExpiredAccount,
		SubAction: "previous",
	}
	_, newDataHash, err := accbuilder.GenWitness(preParam)
	newData := append([]byte{}, newDataHash...)
	fmt.Println("newDataHash:", common.Bytes2Hex(newDataHash), len(newDataHash))
	//// change the next account id
	originData := tx.Transaction.OutputsData[0]
	newData = append(newData, originData[32:52]...)
	fmt.Println("originData:", common.Bytes2Hex(originData[32:52]), len(originData[32:52]))
	fmt.Println("newData:", common.Bytes2Hex(newData))
}
