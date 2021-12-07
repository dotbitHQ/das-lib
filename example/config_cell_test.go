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
	builder, err := dc.GetDasConfigCellByBlockNumber(1948397, common.ConfigCellTypeArgsProfitRate)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.ProfitRateInviter())
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
	for i, item := range builder.ConfigCellRecordKeys {
		fmt.Println("i: ", i)
		fmt.Println("key: ", item)
	}
}