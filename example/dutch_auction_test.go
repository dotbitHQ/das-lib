package example

import (
	"context"
	"das_register_server/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strconv"
	"testing"
	"time"
)

func TestPremium(t *testing.T) {
	//exp := 1690447973

	nowTime := time.Now().Unix()
	//s := common.Premium(int64(exp), nowTime)
	//fmt.Println("nowTime: ", nowTime, " premium: ", s)
	s1 := common.Premium(nowTime-90*24*3600-1, nowTime)
	num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", s1), 64)

	fmt.Println("nowTime: ", nowTime, " premium: ", num)

}

func TestAccId(t *testing.T) {
	account := "dutch-auction-test1"
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	fmt.Println(accountId)
}

func TestGetPeriod(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builderConfigCell, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		fmt.Printf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return
	}
	gracePeriodTime, err := builderConfigCell.ExpirationGracePeriod()
	if err != nil {
		fmt.Printf("ExpirationGracePeriod err: %s", err.Error())
		return
	}
	fmt.Println("gracePeriodTime: ", gracePeriodTime)
	auctionPeriodTime, err := builderConfigCell.ExpirationAuctionPeriod()
	if err != nil {
		fmt.Printf("ExpirationAuctionPeriod err: %s", err.Error())
	}
	fmt.Println("auctionPeriodTime: ", auctionPeriodTime)

	deliverPeriodTime, err := builderConfigCell.ExpirationDeliverPeriod()
	if err != nil {
		fmt.Printf("ExpirationDeliverPeriod err: %s", err.Error())
	}
	fmt.Println("deliverPeriodTime: ", deliverPeriodTime)

}
func TestEditExpiredAt(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	var txParams txbuilder.BuildTransactionParams
	account := "auctiontest10.bit"
	accountOutpoint := "0x7b28c68e50b254e5ab108a34e0b0064c7b74793b5753d5f5436b74048063f617-1"
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(account))
	// config cell
	quoteCell, err := dc.GetQuoteCell()
	if err != nil {
		fmt.Printf("GetQuoteCell err: %s", err.Error())
		return
	}
	quote := quoteCell.Quote()
	priceBuilder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPrice)
	if err != nil {
		fmt.Printf("ConfigCellDataBuilderByTypeArgsList err: %s", err.Error())
		return
	}

	// inputs
	accOutpoint := common.String2OutPointStruct(accountOutpoint)
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		PreviousOutput: accOutpoint,
	})

	// outputs
	accTx, err := dc.Client().GetTransaction(context.Background(), accOutpoint.TxHash)
	if err != nil {
		fmt.Printf("GetTransaction err: %s", err.Error())
		return
	}
	mapAcc, err := witness.AccountIdCellDataBuilderFromTx(accTx.Transaction, common.DataTypeNew)
	if err != nil {
		fmt.Printf("AccountCellDataBuilderMapFromTx err: %s", err.Error())
		return
	}
	accBuilder, ok := mapAcc[accountId]
	if !ok {
		fmt.Printf("account map builder is nil [%s]", accountOutpoint)
		return
	}

	accountLength := uint8(accBuilder.AccountChars.Len())

	_, renewPrice, _ := priceBuilder.AccountPrice(accountLength)
	priceCapacity := (renewPrice / quote) * common.OneCkb
	priceCapacity = priceCapacity * uint64(1)
	fmt.Println("buildOrderRenewTx:", priceCapacity, renewPrice, 1, quote)

	// renew years 90 27 3
	//newExpiredAt := int64(accBuilder.ExpiredAt) + int64(p.renewYears)*common.OneYearSec
	newExpiredAt := time.Now().Unix() - 109*24*3600
	byteExpiredAt := molecule.Go64ToBytes(newExpiredAt)

	accWitness, accData, err := accBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex: 0,
		NewIndex: 0,
		Action:   common.DasActionRenewAccount,
	})
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accTx.Transaction.Outputs[accBuilder.Index].Capacity,
		Lock:     accTx.Transaction.Outputs[accBuilder.Index].Lock,
		Type:     accTx.Transaction.Outputs[accBuilder.Index].Type,
	})

	accData = append(accData, accTx.Transaction.OutputsData[accBuilder.Index][32:]...)
	accData1 := accData[:common.ExpireTimeEndIndex-common.ExpireTimeLen]
	accData2 := accData[common.ExpireTimeEndIndex:]
	newAccData := append(accData1, byteExpiredAt...)
	newAccData = append(newAccData, accData2...)
	txParams.OutputsData = append(txParams.OutputsData, newAccData) // change expired_at
	// witness
	actionWitness, err := witness.GenActionDataWitness(common.DasActionRenewAccount, nil)
	if err != nil {
		fmt.Printf("GenActionDataWitness err: %s", err.Error())
		return
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// cell deps
	dasLockContract, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		fmt.Printf("GetDasContractInfo err: %s", err.Error())
		return
	}
	accContract, err := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		fmt.Printf("GetDasContractInfo err: %s", err.Error())
		return
	}
	incomeContract, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		fmt.Printf("GetDasContractInfo err: %s", err.Error())
		return
	}
	timeCell, err := dc.GetTimeCell()
	if err != nil {
		fmt.Printf("GetTimeCell err: %s", err.Error())
		return
	}
	heightCell, err := dc.GetHeightCell()
	if err != nil {
		fmt.Printf("GetHeightCell err: %s", err.Error())
		return
	}
	accountConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		fmt.Printf("GetDasConfigCellInfo err: %s", err.Error())
		return
	}
	priceConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsPrice)
	if err != nil {
		fmt.Printf("GetDasConfigCellInfo err: %s", err.Error())
		return
	}
	incomeConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsIncome)
	if err != nil {
		fmt.Printf("GetDasConfigCellInfo err: %s", err.Error())
		return
	}
	txParams.CellDeps = append(txParams.CellDeps,
		dasLockContract.ToCellDep(),
		accContract.ToCellDep(),
		incomeContract.ToCellDep(),
		timeCell.ToCellDep(),
		heightCell.ToCellDep(),
		quoteCell.ToCellDep(),
		accountConfig.ToCellDep(),
		priceConfig.ToCellDep(),
		incomeConfig.ToCellDep(),
	)
	base := txbuilder.NewDasTxBuilderBase(context.Background(), dc, nil, "")
	txBuilder := txbuilder.NewDasTxBuilderFromBase(base, nil)
	if err := txBuilder.BuildTransaction(&txParams); err != nil {
		fmt.Printf("BuildTransaction err: %s", err.Error())
		return
	}
	sizeInBlock, _ := txBuilder.Transaction.SizeInBlock()
	changeCapacity := txBuilder.Transaction.Outputs[0].Capacity
	//tx fee
	changeCapacity = changeCapacity - sizeInBlock
	txBuilder.Transaction.Outputs[0].Capacity = changeCapacity

	if hash, err := txBuilder.SendTransaction(); err != nil {
		fmt.Printf("SendTransaction err: %s", err.Error())

	} else {
		fmt.Println("SendTransaction ok:", tables.TxActionRenewAccount, hash)

	}
}
