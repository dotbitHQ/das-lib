package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/DeAccountSystems/das-lib/witness"
	"testing"
)

func TestConfigCellDataBuilderByTypeArgs(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	configCell, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsMain)
	if err != nil {
		t.Fatal(err)
	}
	if res, err := dc.Client().GetTransaction(context.Background(), configCell.OutPoint.TxHash); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.ConfigCellDataBuilderByTypeArgs(res.Transaction, common.ConfigCellTypeArgsMain)
		if err != nil {
			t.Fatal(err)
		}
		status, err := molecule.Bytes2GoU8(builder.ConfigCellMain.Status().RawData())
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("status:", status)
	}
}

func TestGetDasConfigCellByBlockNumber(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsProfitRate)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(builder.ProfitRateInviter())
	fmt.Println(builder.ProfitRateSaleBuyerInviter())
	fmt.Println(builder.ProfitRateSaleBuyerChannel())
	fmt.Println(builder.ProfitRateSaleDas())
}

func TestGetOfferConfig(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsSecondaryMarket)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.OfferMessageBytesLimit())
}

func TestGetKeyNameConfig(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsRecordNamespace)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range builder.ConfigCellRecordKeys {
		fmt.Println(item)
	}

	//apply, err := core.GetDasContractInfo(common.DasContractNameApplyRegisterCellType)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//var o = types.CellOutput{
	//	Lock: common.GetNormalLockScript("0xc866479211cadf63ad115b9da50a6c16bd3d226d"),
	//	Type: apply.ToScript(nil),
	//}
	//data := common.Hex2Bytes("0x1b839f0eb8a356fdeb5b66d8e39779b0b31cfddeef592482bfd7b22e7b26140b82603a00000000007896c26100000000")
	//ca := o.OccupiedCapacity(data)
	//fmt.Println(ca)
}

func TestGetPrice(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPrice)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.AccountPrice(5))
}

func TestGetEmoji(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsUnavailable)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(builder.ConfigCellEmojis))
}
