package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
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
	fmt.Println(configCell.OutPoint.TxHash.Hex())
	if res, err := dc.Client().GetTransaction(context.Background(), configCell.OutPoint.TxHash); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(common.Bytes2Hex(res.Transaction.Witnesses[len(res.Transaction.Witnesses)-1]))
		//builder, err := witness.ConfigCellDataBuilderByTypeArgs(res.Transaction, common.ConfigCellTypeArgsMain)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//status, err := molecule.Bytes2GoU8(builder.ConfigCellMain.Status().RawData())
		//if err != nil {
		//	t.Fatal(err)
		//}
		//fmt.Println("status:", status)
		//script, err := core.GetDasSoScript(common.SoScriptTypeTron)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
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
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsRelease)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.LuckyNumber())
}

func TestConfigCellDataBuilderRefByTypeArgs(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSubAccountWhiteList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(builder.ConfigCellSubAccountWhiteListMap))
}

func TestConfigCellTypeArgsCharSetEn(t *testing.T) {
	_, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(common.CharSetTypeEmojiMap,len(common.CharSetTypeEmojiMap))
	fmt.Println(common.CharSetTypeDigitMap)
	fmt.Println(common.CharSetTypeEnMap)
	//fmt.Println(common.CharSetTypeHanSMap)
	//fmt.Println(common.CharSetTypeHanTMap)
	//fmt.Println(common.CharSetTypeJaMap)
	//fmt.Println(common.CharSetTypeKoMap)
	//fmt.Println(common.CharSetTypeViMap)
	//fmt.Println(common.CharSetTypeRuMap)
	//fmt.Println(common.CharSetTypeThMap)
	//fmt.Println(common.CharSetTypeTrMap)
	//builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsCharSetEmoji)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(builder.ConfigCellEmojis)
}

func TestSubAccountConfigCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSubAccount)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(molecule.Bytes2GoU32(builder.ConfigCellSubAccount.NewSubAccountCustomPriceDasProfitRate().RawData()))
	fmt.Println(molecule.Bytes2GoU32(builder.ConfigCellSubAccount.RenewSubAccountCustomPriceDasProfitRate().RawData()))
}

func TestAccountConfigCell(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.ExpirationGracePeriod())

}

func TestConfigCellProp(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	res, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsProposal)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(molecule.Bytes2GoU32(res.ConfigCellProposal.ProposalMaxPreAccountContain().RawData()))

}

func TestConfigCellTypeArgsSystemStatus(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	res, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSystemStatus)
	if err != nil {
		t.Fatal(err)
	}
	contractStatus, err := res.GetContractStatus(common.DasContractNameAccountCellType)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("status:", contractStatus.Status, "version:", contractStatus.Version)
	fmt.Println(dc.CheckContractVersion(common.DasContractNameAccountCellType))
	fmt.Println(dc.CheckContractStatusOK(common.DasContractNameAccountCellType))

	list := []common.DasContractName{
		common.DasContractNameApplyRegisterCellType,
		common.DasContractNamePreAccountCellType,
		common.DasContractNameProposalCellType,
		common.DasContractNameConfigCellType,
		common.DasContractNameAccountCellType,
		common.DasContractNameAccountSaleCellType,
		common.DASContractNameSubAccountCellType,
		common.DASContractNameOfferCellType,
		common.DasContractNameBalanceCellType,
		common.DasContractNameIncomeCellType,
		common.DasContractNameReverseRecordCellType,
		common.DASContractNameEip712LibCellType,
		common.DasContractNameReverseRecordRootCellType,
	}
	sysStatus, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsSystemStatus)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range list {
		fmt.Println(v)
		fmt.Println(dc.CheckContractVersionV2(sysStatus, v))
	}
}
