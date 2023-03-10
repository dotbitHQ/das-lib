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
	txJson := `{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x6e2d0ea0543984d3baf13bd25038dd6d3222baa20b4f1872f4864e4d7eb6c827","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x277d21646ca48621008cad78c5c8e089ecbabfb3150402093f80470f5f830456","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf56b13151f259dcb9074a019276cd1b223b96a79e16323e07347aa137eb01304","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xd4d5023b35db6e1f1212e658b114778c49ccdbade6aeb966e6ed7b0b89b6eb8f","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x25e88427f146a2e6419db3e64777549cef653932d4d56d11ca6c71b43e348b74","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xd8994603c0b7b5d114f9b6f5b3f2107105b5628b3775c00c93c3f4d62c865477","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x19901f755e4ac85bb505dc760b2fe63d8fae4551c55ad9df64dd2d0ded91bae7","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4a1ad60d8cdb1d41d5579e9ec8e4ef6fd09fa09083651375adca99d3fd0654b1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x29a278a2498ca40d386c1432843e32c04ee996c99e7e16db7e09bfb60f83f47e","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x5eeacd8e0f093e7366cf7c8b2a4cda1b01eb15b11e9a0c535dd2bd63c481e56f","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x448880888c693ba05c3835d28069089aedfd432d65e4629fb507c31688319aa9","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0xf10d683189f9fec5406e632b8d21271ea58e20728b9f3f3bbb4e88d7d1a73c16","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xb5da406fba558d217d50b8aa30fbdd269c94d8483518105e70dc2e1d0aca7f6d","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xe7b9978b125e35d8cf4358bba5b08b655d0f9abf577bc0acfcbd386b0622740a","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0x07dfcb9511725027968c3aa0c8c80300be907b7c5c1ac41c4c261d943778b0bb","index":"0x0"}}],"outputs":[{"capacity":"0x5195767ad","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x07b6031be679d6bfa9ce6db1e3bf61b6e6552423be07b6031be679d6bfa9ce6db1e3bf61b6e6552423be"},"type":{"code_hash":"0x1106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b1","hash_type":"type","args":"0x"}}],"outputs_data":["0xff661c5b25525805fd11e11021106dbbce0b70f12aa77eb89e1777bf904be9eedc081ac3ed5f4c8c44282550ad62823ad43d50f5dc30c14f2b9738ecefc1638df02e47af1322456abb71e0650000000032303233303330312e626974"],"witnesses":["0x56000000100000005600000056000000420000002e40b5ec1ef32e6f830cadecf4a39c157cf9d6c228323206539836c3ad4762ce79daa77f550891574fc33f59ee7a04eb5073a6e6645c268456b6662c666a34600100","0x64617300000000210000000c0000001c0000000c000000656469745f7265636f7264730100000001","0x64617301000000160300001000000010000000660100005601000010000000140000001800000000000000030000003a0100003a0100002c000000400000000c010000140100001c010000240100002c0100002d0100003101000032010000dc081ac3ed5f4c8c44282550ad62823ad43d50f5cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000032150000000c00000010000000010000000100000030150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000030150000000c00000010000000010000000100000033150000000c00000010000000010000000100000030150000000c000000100000000100000001000000313b3eff63000000000000000000000000000000000000000000000000000000000004000000000000000000000000b0010000100000001400000018000000000000000300000094010000940100002c000000400000000c010000140100001c010000240100002c0100002d0100008b0100008c010000dc081ac3ed5f4c8c44282550ad62823ad43d50f5cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000032150000000c00000010000000010000000100000030150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000030150000000c00000010000000010000000100000033150000000c00000010000000010000000100000030150000000c000000100000000100000001000000313b3eff6300000000000000000000000000000000000000005ff5096400000000005e00000008000000560000001800000023000000280000002c00000052000000070000006164647265737301000000330000000022000000444d6a5646427162715a4741795458676b743766547571696868434356754c775a362c010000000000000000000000","0x646173640000008c0000003c00000040000000480000005000000054000000580000005c000000640000006c000000740000007c0000008000000084000000880000002a000000000edbcb0400000000e1f50500000000010000002c0100008813000010270000000000001027000000000000102700000000000010270000000000002c0100002c0100002c0100002c010000","0x6461736c0000008c030000616464726573732e61646100616464726573732e61746f6d00616464726573732e6176616c616e63686500616464726573732e62636800616464726573732e62736300616464726573732e62737600616464726573732e62746300616464726573732e63656c6f00616464726573732e636b6200616464726573732e6461736800616464726573732e6466696e69747900616464726573732e646f676500616464726573732e646f7400616464726573732e656f7300616464726573732e65746300616464726573732e65746800616464726573732e66696c00616464726573732e666c6f7700616464726573732e6865636f00616464726573732e696f737400616464726573732e696f746100616464726573732e6b736d00616464726573732e6c746300616464726573732e6e65617200616464726573732e706f6c79676f6e00616464726573732e736300616464726573732e736f6c00616464726573732e737461636b7300616464726573732e746572726100616464726573732e74727800616464726573732e76657400616464726573732e78656d00616464726573732e786c6d00616464726573732e786d7200616464726573732e78727000616464726573732e78747a00616464726573732e7a656300616464726573732e7a696c00647765622e6172776561766500647765622e6970667300647765622e69706e7300647765622e726573696c696f00647765622e736b796e65740070726f66696c652e6176617461720070726f66696c652e626568616e63650070726f66696c652e62696c6962696c690070726f66696c652e6465736372697074696f6e0070726f66696c652e646973636f72640070726f66696c652e6472696262626c650070726f66696c652e656d61696c0070726f66696c652e66616365626f6f6b0070726f66696c652e6769746875620070726f66696c652e696e7374616772616d0070726f66696c652e6a696b650070726f66696c652e6c696e6b6564696e0070726f66696c652e6d656469756d0070726f66696c652e6d6972726f720070726f66696c652e6e65787469640070726f66696c652e6e6f7374720070726f66696c652e7265646469740070726f66696c652e74656c656772616d0070726f66696c652e74696b746f6b0070726f66696c652e747769747465720070726f66696c652e776562736974650070726f66696c652e776569626f0070726f66696c652e796f757475626500","0x64617368000000bd0300001400000015000000ed010000e102000001d801000038000000580000007800000098000000b8000000d8000000f80000001801000038010000580100007801000098010000b80100001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4e8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5ff40000001c000000400000006400000088000000ac000000d0000000448880888c693ba05c3835d28069089aedfd432d65e4629fb507c31688319aa900000000ccf7216a4af8aad8d8872e9660f56edc21a9712016551fcb679051ec6b5e6ee600000000000000000000000000000000000000000000000000000000000000000000000000000000322d876ac47f8d901d39d6ece8d64cf7d91adbeef76e72532581207f415bb266000000004a1ad60d8cdb1d41d5579e9ec8e4ef6fd09fa09083651375adca99d3fd0654b10000000029a278a2498ca40d386c1432843e32c04ee996c99e7e16db7e09bfb60f83f47e00000000dc0000001c0000003c0000005c0000007c0000009c000000bc000000c9fc9f3dc050f8bf11019842a2426f48420f79da511dd169ee243f455e9f84ed991bcf61b6d7a26e6c27bda87d5468313d99ef0cd37113eee9e16c2680fa4532ebb79383a2947f36a095b434dd4f7c670dec6c2a53d925fb5c5f949104e59a6fe3ce977f83fd46cdacc6bccef4d8b045ecf245b39c29558cad2c1405220b6914f8f6b58d548231bc6fe19c1a1ceafa3a429f54c21a458b211097ebe564b146157ab1b06d51c579d528395d7f472582bf1d3dce45ba96c2bff2c19e30f0d90281"]}`
	//txJson=`{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x25e88427f146a2e6419db3e64777549cef653932d4d56d11ca6c71b43e348b74","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xd8994603c0b7b5d114f9b6f5b3f2107105b5628b3775c00c93c3f4d62c865477","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x19901f755e4ac85bb505dc760b2fe63d8fae4551c55ad9df64dd2d0ded91bae7","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4a1ad60d8cdb1d41d5579e9ec8e4ef6fd09fa09083651375adca99d3fd0654b1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x29a278a2498ca40d386c1432843e32c04ee996c99e7e16db7e09bfb60f83f47e","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x5eeacd8e0f093e7366cf7c8b2a4cda1b01eb15b11e9a0c535dd2bd63c481e56f","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x448880888c693ba05c3835d28069089aedfd432d65e4629fb507c31688319aa9","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0xc06a70738cd7361b53b45f38f19aa3024116988d834b621d34f04c8690dcb604","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xb5da406fba558d217d50b8aa30fbdd269c94d8483518105e70dc2e1d0aca7f6d","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0xb9fa98480a7aa23de1f50432fb50fb182260574ed88a7970ba5a1e1b3f0e3a40","index":"0x2"}},{"since":"0x0","previous_output":{"tx_hash":"0x42e0946009d6d9672fc51945d28a98ccc1fc50d55e889c4dd402f0678ca37d8a","index":"0x1"}},{"since":"0x0","previous_output":{"tx_hash":"0x3336dc8bb6d6614da094743e83ea37f97fcc2e2a3353e0bbac3b8b10073479cc","index":"0x1"}},{"since":"0x0","previous_output":{"tx_hash":"0x752744e9e18a792d94ba7fca113d3b72ea739862d9138ba7098f9395879f839e","index":"0x1"}}],"outputs":[{"capacity":"0x174876e800","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x04517799ca4818d59e5fe160a0d12e2a697c18f2d804517799ca4818d59e5fe160a0d12e2a697c18f2d8"},"type":{"code_hash":"0x4ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c","hash_type":"type","args":"0x"}},{"capacity":"0x3ebab81de","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x0554366bcd1e73baf55449377bd23123274803236e0554366bcd1e73baf55449377bd23123274803236e"},"type":{"code_hash":"0x4ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c","hash_type":"type","args":"0x"}}],"outputs_data":["0x","0x"],"witnesses":["0x7d000000100000007d0000007d000000690000001df032fac1bebdf861a39d5c289a89774f81fd3e21e0eb29eca1e00795b4f6ff127c3212833a8385282cc846fe1bb3e1b57d15d20cac769a579a8539b2599c41005b2b8e879ee80f8312b6e47f6f58c5c82d52961bd7ca3c2cb389179acc5504910000000000000005","0x","0x","0x","0x64617300000000290000000c000000240000001400000077697468647261775f66726f6d5f77616c6c65740100000000","0x64617368000000bd0300001400000015000000ed010000e102000001d801000038000000580000007800000098000000b8000000d8000000f80000001801000038010000580100007801000098010000b80100001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4e8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5ff40000001c000000400000006400000088000000ac000000d0000000448880888c693ba05c3835d28069089aedfd432d65e4629fb507c31688319aa900000000ccf7216a4af8aad8d8872e9660f56edc21a9712016551fcb679051ec6b5e6ee600000000000000000000000000000000000000000000000000000000000000000000000000000000322d876ac47f8d901d39d6ece8d64cf7d91adbeef76e72532581207f415bb266000000004a1ad60d8cdb1d41d5579e9ec8e4ef6fd09fa09083651375adca99d3fd0654b10000000029a278a2498ca40d386c1432843e32c04ee996c99e7e16db7e09bfb60f83f47e00000000dc0000001c0000003c0000005c0000007c0000009c000000bc000000c9fc9f3dc050f8bf11019842a2426f48420f79da511dd169ee243f455e9f84ed991bcf61b6d7a26e6c27bda87d5468313d99ef0cd37113eee9e16c2680fa4532ebb79383a2947f36a095b434dd4f7c670dec6c2a53d925fb5c5f949104e59a6fe3ce977f83fd46cdacc6bccef4d8b045ecf245b39c29558cad2c1405220b6914f8f6b58d548231bc6fe19c1a1ceafa3a429f54c21a458b211097ebe564b146157ab1b06d51c579d528395d7f472582bf1d3dce45ba96c2bff2c19e30f0d90281"]}`
	var skipGroups []int
	digestList, err := txbuilder.GenerateDigestListFromTx(cli, txJson, skipGroups)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Digest:", digestList)
}
