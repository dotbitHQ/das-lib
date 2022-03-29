package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"strings"
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
		script, err := core.GetDasSoScript(common.SoScriptTypeTron)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
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
	fmt.Println(builder.PriceConfigMap)
	for k, v := range builder.PriceConfigMap {
		fmt.Println(k)
		fmt.Println(molecule.Bytes2GoU64(v.New().RawData()))
	}
	fmt.Println(builder.AccountPrice(5))
}

func TestGetEmoji(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(
		common.ConfigCellTypeArgsPreservedAccount00,
		common.ConfigCellTypeArgsPreservedAccount01,
		common.ConfigCellTypeArgsPreservedAccount02,
		common.ConfigCellTypeArgsPreservedAccount03,
		common.ConfigCellTypeArgsPreservedAccount04,
		common.ConfigCellTypeArgsPreservedAccount05,
		common.ConfigCellTypeArgsPreservedAccount06,
		common.ConfigCellTypeArgsPreservedAccount07,
		common.ConfigCellTypeArgsPreservedAccount08,
		common.ConfigCellTypeArgsPreservedAccount09,
		common.ConfigCellTypeArgsPreservedAccount10,
		common.ConfigCellTypeArgsPreservedAccount11,
		common.ConfigCellTypeArgsPreservedAccount12,
		common.ConfigCellTypeArgsPreservedAccount13,
		common.ConfigCellTypeArgsPreservedAccount14,
		common.ConfigCellTypeArgsPreservedAccount15,
		common.ConfigCellTypeArgsPreservedAccount16,
		common.ConfigCellTypeArgsPreservedAccount17,
		common.ConfigCellTypeArgsPreservedAccount18,
		common.ConfigCellTypeArgsPreservedAccount19,
		common.ConfigCellTypeArgsUnavailable,
	)
	if err != nil {
		t.Fatal(err)
	}
	str := ``
	//fmt.Println(len(builder.ConfigCellEmojis))
	list := strings.Split(str, "\n")
	for _, v := range list {
		bys, _ := blake2b.Blake160([]byte(v))
		accountHashIndex := uint32(bys[0] % 20)
		tmp := common.Bytes2Hex(common.Blake2b([]byte(v))[:20])

		fmt.Println(accountHashIndex, v, tmp)
		if _, ok := builder.ConfigCellPreservedAccountMap[tmp]; ok {
			fmt.Println(v, true)
		}
		if _, ok := builder.ConfigCellUnavailableAccountMap[tmp]; ok {
			fmt.Println(v, true)
		}
	}
	fmt.Println(len(builder.ConfigCellUnavailableAccountMap), len(builder.ConfigCellPreservedAccountMap))
	//byStr:=common.Bytes2Hex(builder.ConfigCellUnavailableAccount)
	//fmt.Println(byStr,len(builder.ConfigCellUnavailableAccount))
}

func TestIncomeCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	configCell, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsPreservedAccount03)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(configCell.OutPoint.TxHash, configCell.OutPoint.Index)
	res, err := dc.Client().GetTransaction(context.Background(), configCell.OutPoint.TxHash)
	for _, v := range res.Transaction.Witnesses {
		fmt.Println(len(v))
	}
	//builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPreservedAccount03)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(builder.IncomeBasicCapacity())
}

func TestBasicCapacityFromOwnerDasAlgorithmId(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsAccount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.BasicCapacityFromOwnerDasAlgorithmId("0x04"))
}

func TestConfigCellReleaseLuckyNumber(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsRelease)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.LuckyNumber())
}
