package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestSoScript(t *testing.T) {
	_, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	script, err := core.GetDasSoScript(common.SoScriptTypeDogeCoin)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
	script, err = core.GetDasSoScript(common.SoScriptTypeEd25519)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
}

func TestNewDasCore(t *testing.T) {
	_, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	// contract
	cont, err := core.GetDasContractInfo(common.DasContractNameConfigCellType)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cont.ContractName, cont.ContractTypeId, cont.OutPoint.TxHash.Hex())
	// config cell
	cc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsMain)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cc.Name, cc.OutPoint.TxHash.Hex(), cc.OutPoint.Index)
	// so script
	script, err := core.GetDasSoScript(common.SoScriptTypeEd25519)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
}

func TestTHQCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	heightCell, err := dc.GetHeightCell()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(heightCell.LiveCell.OutPoint.TxHash.String(), heightCell.LiveCell.OutPoint.Index)
	//tc, err := dc.GetTimeCell()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(tc.Timestamp(), tc.LiveCell.OutPoint.TxHash.Hex(), tc.LiveCell.OutPoint.Index)
	//
	//hc, err := dc.GetHeightCell()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(hc.BlockNumber(), hc.LiveCell.OutPoint.TxHash.Hex(), hc.LiveCell.OutPoint.Index)
	//
	//qc, err := dc.GetQuoteCell()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(qc.Quote(), qc.LiveCell.OutPoint.TxHash.Hex(), qc.LiveCell.OutPoint.Index)
}

func TestGetAccountCellOnChainByAlgorithmId(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	liveCell, err := dc.GetAccountCellOnChainByAlgorithmId(common.DasAlgorithmIdEth, common.DasAlgorithmIdEth, "0xad63e52c73397ef5c0d38445e83dd6673cc60ebb", "0xad63e52c73397ef5c0d38445e83dd6673cc60ebb", "345435dsfsfg.bit")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(liveCell.BlockNumber, liveCell.OutPoint.TxHash.Hex(), liveCell.OutPoint.Index)
}

func TestGetCells(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	pre, _ := core.GetDasContractInfo(common.DasContractNamePreAccountCellType)
	searchKey := &indexer.SearchKey{
		Script:     pre.ToScript(nil),
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			OutputDataLenRange: &[2]uint64{52, 53},
		},
	}
	res, _ := dc.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderDesc, 100, "")

	for _, v := range res.Objects {
		fmt.Println(v.OutPoint.TxHash)
	}
}

func TestConfigCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	conf, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(conf.BasicCapacity())
}

func TestGetLiveCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	res, err := dc.Client().GetLiveCell(context.Background(), common.String2OutPointStruct("0x80ed13d2f0b1192e49f6130d5802044c96c2baff34496bc2d04a3e47572be015-1"), true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Status, res.Cell.Output)
}

func TestGoU64ToBytes(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	heightCell, _ := dc.GetHeightCell()
	fmt.Println(molecule.Go64ToBytes(heightCell.BlockNumber()))
}

func TestConvertScriptToAddress(t *testing.T) {
	// ckb1qj0n46hjl3pe2jwtepcvv5ehf9p6l94qvk9addgmarvfsvvrum6j7qwm3ct4htjlnv9fzz6js78jmkpk4veqdpcq0qxqzc
	// ckb1qn3yze8zyp8enzcg3ysyqh0vu0w06hqlhjer4m8uujea8m03fzyfwsxh8uxnc4sle2hrxr4tcqcd3kt2n590xmgv2y2gsdjc5dgvh83m5yqqqqqsqqqqqvqqqqqfjqqqqpvgk3ep8wj9g3yfdm5ff5c07c6uepr3sswuwg5x7d7q8z70fl22u6gqqqqpqqqqqqcqqqqqxyqqqqq4vvyq696mlrwagjjgapgvancvpdzhtq6h2m44ll2n4kpsjvdelyqngqqqqpqdw0cd83tplj4wxv82hspsmrvk48g27dksc5g53qm93g6sew0rhkuwzad6uhump2gsk5583ukasd4txgrgwq9rqgqqqqqqcq5qvmhg

	addr := "ckb1qyqyz7atfywzrldrllhqe4jswuxd7ge4a7mstd5ekq"
	parseAddress, err := address.Parse(addr)
	if err != nil {
		t.Fatal(err)
	}
	resAddr, err := common.ConvertScriptToAddress(address.Mainnet, parseAddress.Script)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resAddr)

	//addr1, err := address.ConvertToBech32mFullAddress(addr)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(addr1)

}

func TestInitEnv(t *testing.T) {
	fmt.Println(core.InitEnv(common.DasNetTypeMainNet))
	fmt.Println(core.InitEnv(common.DasNetTypeTestnet2))
	fmt.Println(core.InitEnv(common.DasNetTypeTestnet3))
	fmt.Println()

	fmt.Println(core.InitEnvOpt(common.DasNetTypeMainNet, common.DasContractNameAccountCellType, common.DasContractNameAccountSaleCellType))
	fmt.Println(core.InitEnvOpt(common.DasNetTypeTestnet2, common.DasContractNameAccountCellType, common.DasContractNameAccountSaleCellType))
	fmt.Println(core.InitEnvOpt(common.DasNetTypeTestnet3, common.DasContractNameAccountCellType, common.DasContractNameAccountSaleCellType))
}

func TestArgs(t *testing.T) {
	fmt.Println(common.ConvertScriptToAddress(address.Testnet, &types.Script{
		CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes("0xa897829e60ee4e3fb0e4abe65549ec4a5ddafad7"),
	}))
	fmt.Println(common.ConvertScriptToAddress(address.Testnet, &types.Script{
		CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes("0xa897829e60ee4e3fb0e4abe65549ec4a5ddafad7"),
	}))
}
