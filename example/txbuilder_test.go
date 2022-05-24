package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/txbuilder"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestNewDasTxBuilderFormSystem(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	base := txbuilder.NewDasTxBuilderBase(context.Background(), dc, nil, "")
	builder := txbuilder.NewDasTxBuilderFromBase(base, nil)
	err = builder.BuildTransaction(&txbuilder.BuildTransactionParams{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.TxString())
}

func TestBuildMMJsonObj(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	base := txbuilder.NewDasTxBuilderBase(context.Background(), dc, nil, "")
	builder := txbuilder.NewDasTxBuilderFromBase(base, nil)

	tx, err := getEditAccountSaleTx(dc, "0x4da6fdb1295af7dc54c5374c463f134b5f91340110ece319acd09af45a200633")
	if err != nil {
		t.Fatal(err)
	}
	if err := builder.BuildTransaction(tx); err != nil {
		t.Fatal(err)
	}

	obj, err := builder.BuildMMJsonObj(0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("obj:", obj.String())

	signList, err := builder.GenerateDigestListFromTx([]int{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signList)
	//fmt.Println(builder.GetDasTxBuilderTransactionString())

	hash, err := builder.SendTransaction()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("hahs:", hash)

}

func getEditAccountSaleTx(dc *core.DasCore, hash string) (*txbuilder.BuildTransactionParams, error) {
	// sale cell
	res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash))
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}

	builder, err := witness.AccountSaleCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountSaleCellDataBuilderFromTx err: %s", err.Error())
	}

	// inputs
	var inputs []*types.CellInput
	inputs = append(inputs, &types.CellInput{PreviousOutput: &types.OutPoint{
		TxHash: types.HexToHash(hash),
		Index:  uint(builder.Index),
	}})

	// action witness
	var witnesses [][]byte
	actionWitness, err := witness.GenActionDataWitness(common.DasActionEditAccountSale, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	witnesses = append(witnesses, actionWitness)

	// sale cell witness
	dataWitness, accountSaleOutputData, _ := builder.GenWitness(&witness.AccountSaleCellParam{
		Price:       654 * 1e8,
		Description: "sa1",
		Action:      common.DasActionEditAccountSale,
	})
	witnesses = append(witnesses, dataWitness)

	// outputs
	var outputsData [][]byte
	outputsData = append(outputsData, accountSaleOutputData)

	fee := uint64(1e4)
	var outputs []*types.CellOutput
	outputs = append(outputs, &types.CellOutput{
		Capacity: res.Transaction.Outputs[builder.Index].Capacity - fee,
		Lock:     res.Transaction.Outputs[builder.Index].Lock,
		Type:     res.Transaction.Outputs[builder.Index].Type,
	})

	// height,time cell
	heightCell, err := dc.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}
	timeCell, err := dc.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	configCellMarket, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsSecondaryMarket)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	alwaysSuccessContract, err := core.GetDasContractInfo(common.DasContractNameAlwaysSuccess)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	fmt.Println("alwaysSuccessContract:", alwaysSuccessContract.OutPoint.TxHash.Hex())
	incomeCellType, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellTypeArgsAccount err: %s", err.Error())
	}

	configCellIncome, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsIncome)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellTypeArgsAccount err: %s", err.Error())
	}

	configCellProfit, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsProfitRate)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellTypeArgsAccount err: %s", err.Error())
	}

	so, _ := core.GetDasSoScript(common.DasAlgorithmIdCkb.ToSoScriptType())
	fmt.Println("so:", so.OutPoint.TxHash.Hex())
	// cell deps
	cellDeps := []*types.CellDep{
		heightCell.ToCellDep(),
		timeCell.ToCellDep(),
		configCellMarket.ToCellDep(),
		alwaysSuccessContract.ToCellDep(),
		incomeCellType.ToCellDep(),
		configCellIncome.ToCellDep(),
		configCellProfit.ToCellDep(),
		so.ToCellDep(),
	}

	tx := txbuilder.BuildTransactionParams{
		CellDeps:    cellDeps,
		Inputs:      inputs,
		Outputs:     outputs,
		OutputsData: outputsData,
		Witnesses:   witnesses,
	}

	return &tx, nil
}

func TestParam(t *testing.T) {
	data := common.Hex2Bytes("0x3c00000000000000050000000000000000")
	fmt.Println(len(data), data)
}

func TestGenerateMultiSignWitnessArgs(t *testing.T) {
	emptySignatures := make([][]byte, 3)
	addrList := []string{
		"ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqfdqkruhv2ac0z43yavczye39v457nq8vclg7xgl",
		"ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsq0nzujqmmmarw0azts6869ucjkn0xlt5esjs0cn0",
		"ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqdzc26ytd5dgz2f5uyc67v89yw50szgkwcp9sl0f",
		"ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqdvjum0ha4zr9k59w0k693gsvw563cgzjglua447",
		"ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqdugnv77uz0zlwfme8ah640qsw0slmdjegusx37n",
	}
	var sortArgsList [][]byte
	for _, v := range addrList {
		addrP, _ := address.Parse(v)
		sortArgsList = append(sortArgsList, addrP.Script.Args)
		fmt.Println(common.Bytes2Hex(addrP.Script.Args))
	}

	wa := txbuilder.GenerateMultiSignWitnessArgs(0, emptySignatures, sortArgsList)
	fmt.Println(common.Bytes2Hex(wa.Lock))
}
