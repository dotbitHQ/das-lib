package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/dotbitHQ/das-lib/witness"
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

func TestGenerateDigestListFromTx(t *testing.T) {
	cli, err := getClientTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	txJson := `{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x7b8ead25d97ab4ace27931a723f0c6ceb207eb9dad78976b26b29ffd9a64e2b7","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x538c1335196c88676b3287107c3d6cb53626f8381e9e6a5c4ac8f645ea69f54d","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x10b000eb473abf0847655a02ad8384ea808e8f9b88a3240f695475984f7e674d","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x9e2272d2ec864b2f22ff74dd19c12ecc82718ec0864c37324d31a5b09dd64a00","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x243d467d3a0c3355e64a03baab4f4850ebe2133b6ac34586e7f21b2248d997f0","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4c817bc81165aae004f0961d583492e95759212edaa210afc434766998ce2670","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0x5165a17d1a487ee19362c231fca1658e0cebcd4d3c3d57a79ca144ae45e9520c","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x35cca9f4c1bed9642a971777a656fe769abfa9e441d8b7f29b571e27b8239b63","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x8ffa409ba07d74f08f63c03f82b7428d36285fe75b2173fc2476c0f7b80c707a","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0xf6df36b34b535b3c1dc5fd93b9da072d8ae9958054d8761fec94ba6cfa46a7b9","index":"0x1"}}],"outputs":[{"capacity":"0x50775b0f0","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x05acd0d6ce063c8609414b511d05b682f8717f529a05acd0d6ce063c8609414b511d05b682f8717f529a"},"type":{"code_hash":"0x1106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b1","hash_type":"type","args":"0x"}}],"outputs_data":["0x4d18a43fa0e59a33cef987bb0a9bb183960dcbc21a6767be39e9e9a554419e84247a9b0c9d9dfd787861c9d7d928dc321124996f2484ac44dce340acb817023b056102d2bfc8cc14684625640000000031323333332e626974"],"witnesses":["0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","0x64617300000000210000000c0000001c0000000c000000656469745f7265636f7264730100000001","0x64617301000000cb02000010000000100000001b0100000b0100001000000014000000180000000000000003000000ef000000ef0000002c00000040000000c1000000c9000000d1000000d9000000e1000000e2000000e6000000e7000000247a9b0c9d9dfd787861c9d7d928dc321124996f81000000180000002d00000042000000570000006c000000150000000c00000010000000010000000100000031150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000033150000000c00000010000000010000000100000033e8124462000000000000000000000000000000000000000000000000000000000004000000000000000000000000b0010000100000001400000018000000000000000300000094010000940100002c00000040000000c1000000c9000000d1000000d9000000e1000000e20000008b0100008c010000247a9b0c9d9dfd787861c9d7d928dc321124996f81000000180000002d00000042000000570000006c000000150000000c00000010000000010000000100000031150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000033150000000c00000010000000010000000100000033e812446200000000000000000000000000000000000000006ee199630000000000a90000000c0000006b0000005f0000001800000023000000290000002d0000005b0000000700000061646472657373020000003630000000002a0000003078616364306436636530363363383630393431346235313164303562363832663837313766353239612c0100003e00000018000000260000002d000000320000003a0000000a000000637573746f6d5f6b657903000000616161010000006104000000616161762c010000000000000000000000","0x646173640000008c0000003c00000040000000480000005000000054000000580000005c000000640000006c000000740000007c0000008000000084000000880000002a000000000edbcb0400000000e1f50500000000008d27002c0100008813000010270000000000001027000000000000102700000000000010270000000000002c0100002c0100002c0100002c010000","0x6461736c0000007e030000616464726573732e61646100616464726573732e61746f6d00616464726573732e6176616c616e63686500616464726573732e62636800616464726573732e62736300616464726573732e62737600616464726573732e62746300616464726573732e63656c6f00616464726573732e636b6200616464726573732e6461736800616464726573732e6466696e69747900616464726573732e646f676500616464726573732e646f7400616464726573732e656f7300616464726573732e65746300616464726573732e65746800616464726573732e66696c00616464726573732e666c6f7700616464726573732e6865636f00616464726573732e696f737400616464726573732e696f746100616464726573732e6b736d00616464726573732e6c746300616464726573732e6e65617200616464726573732e706f6c79676f6e00616464726573732e736300616464726573732e736f6c00616464726573732e737461636b7300616464726573732e746572726100616464726573732e74727800616464726573732e76657400616464726573732e78656d00616464726573732e786c6d00616464726573732e786d7200616464726573732e78727000616464726573732e78747a00616464726573732e7a656300616464726573732e7a696c00647765622e6172776561766500647765622e6970667300647765622e69706e7300647765622e726573696c696f00647765622e736b796e65740070726f66696c652e6176617461720070726f66696c652e626568616e63650070726f66696c652e62696c6962696c690070726f66696c652e6465736372697074696f6e0070726f66696c652e646973636f72640070726f66696c652e6472696262626c650070726f66696c652e656d61696c0070726f66696c652e66616365626f6f6b0070726f66696c652e6769746875620070726f66696c652e696e7374616772616d0070726f66696c652e6a696b650070726f66696c652e6c696e6b6564696e0070726f66696c652e6d656469756d0070726f66696c652e6d6972726f720070726f66696c652e6e65787469640070726f66696c652e7265646469740070726f66696c652e74656c656772616d0070726f66696c652e74696b746f6b0070726f66696c652e747769747465720070726f66696c652e776562736974650070726f66696c652e776569626f0070726f66696c652e796f757475626500","0x64617368000000b90200001000000011000000c501000001b401000034000000540000007400000094000000b4000000d4000000f400000014010000340100005401000074010000940100001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4ef40000001c000000400000006400000088000000ac000000d0000000209b35208da7d20d882f0871f3979c68c53981bcc4caa71274c035449074d08200000000747411fb3914dd7ca5488a0762c6f4e76f56387e83bcbb24e3a01afef1d5a5b4000000000000000000000000000000000000000000000000000000000000000000000000000000008ffa409ba07d74f08f63c03f82b7428d36285fe75b2173fc2476c0f7b80c707a000000009e0823959e5b76bd010cc503964cced4f8ae84f3b03e94811b083f9765534ff100000000a706f46e58e355a6d29d7313f548add21b875639ea70605d18f682c1a08740d600000000"]}`
	var skipGroups []int
	digestList, err := txbuilder.GenerateDigestListFromTx(cli, txJson, skipGroups)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Digest:", digestList)
}
