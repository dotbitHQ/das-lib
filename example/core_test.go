package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"testing"
)

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
	script, err := core.GetDasSoScript(common.SoScriptTypeEth)
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
	tc, err := dc.GetTimeCell()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tc.Timestamp(), tc.LiveCell.OutPoint.TxHash.Hex(), tc.LiveCell.OutPoint.Index)

	hc, err := dc.GetHeightCell()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hc.BlockNumber(), hc.LiveCell.OutPoint.TxHash.Hex(), hc.LiveCell.OutPoint.Index)

	qc, err := dc.GetQuoteCell()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(qc.Quote(), qc.LiveCell.OutPoint.TxHash.Hex(), qc.LiveCell.OutPoint.Index)
}

func TestParseDasLockArgs(t *testing.T) {
	args := "0x053919a8eb619ccae32fba88d333829929db2f432405c9f53b1d85356b60453f867610888d89a0b667ad"
	fmt.Println(core.FormatDasLockToHexAddress(common.Hex2Bytes(args)))
	fmt.Println(core.FormatNormalCkbLockToAddress(common.DasNetTypeMainNet, common.Hex2Bytes(args)))
}

func TestGetAccountCellOnChain(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	liveCell, err := dc.GetAccountCellOnChain(common.ChainTypeEth, common.ChainTypeEth, "0xc82ee26529193afd4252592c585178d8baf07545", "0xc82ee26529193afd4252592c585178d8baf07545", "asdsadsada.bit")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(liveCell.BlockNumber, liveCell.OutPoint.TxHash.Hex(), liveCell.OutPoint.Index)
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

func TestGetSatisfiedCapacityLiveCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	dasLock, dasType, err := dc.FormatAddressToDasLockScript(common.ChainTypeTron, "TQoLh9evwUmZKxpD1uhFttsZk3EBs8BksV", true)
	if err != nil {
		t.Fatal(err)
	}
	cells, total, err := core.GetSatisfiedCapacityLiveCell(dc.Client(), nil, dasLock, dasType, 0, 116*common.OneCkb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(total)
	for _, v := range cells {
		fmt.Println(len(v.OutputData))
		fmt.Println(v.BlockNumber, v.OutPoint.TxHash, v.OutPoint.Index)
	}
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

func TestFormatOwnerManagerAddressToArgs(t *testing.T) {
	oCT, mCT, oA, mA := common.ChainTypeEth, common.ChainTypeTron, "0xc9f53b1d85356B60453F867610888D89a0B667Ad", "TEooRfPxhqJ7AJfmsRg5hZWEX95VeNxvtX"
	args := core.FormatOwnerManagerAddressToArgs(oCT, mCT, oA, mA)
	fmt.Println(common.Bytes2Hex(args))
}
