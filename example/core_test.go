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
	_, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	script, err := core.GetDasSoScript(common.SoScriptBitcoin)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
	cont, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(cont.ContractName, cont.ContractTypeId, cont.OutPoint.TxHash.Hex())
	//script, err = core.GetDasSoScript(common.SoScriptTypeEd25519)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
}

func TestNewDasCore(t *testing.T) {
	_, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}

	//builder, err := dc.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsApply)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//applyMaxWaitingBlockNumber, err := molecule.Bytes2GoU32(builder.ConfigCellApply.ApplyMaxWaitingBlockNumber().RawData())
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res := utils.SinceFromRelativeBlockNumber(uint64(applyMaxWaitingBlockNumber))
	//fmt.Println(res)
	//fmt.Println(hexutil.Uint64(res))

	core.DasContractMap.Range(func(key, value any) bool {
		item, ok := value.(*core.DasContractInfo)
		if !ok {
			return true
		}
		fmt.Println(item.ContractName, item.OutPoint.TxHash)
		return true
	})
	core.DasConfigCellMap.Range(func(key, value any) bool {
		item, ok := value.(*core.DasConfigCellInfo)
		if !ok {
			return true
		}
		fmt.Println(item.Name, item.OutPoint.TxHash)
		return true
	})

	//did-cell-type 0xc4e296d8dc96636603d404cc5762cc865baff234e2e0adf07b4de6ca059312ba
	//apply-register-cell-type 0x4cdae664cd2ba5c7b4ee5f67401191fb0d8ee559cda5befae318fc4f9cfb3563
	//config-cell-type 0x18557adb9154627384d89ddbcab66b5f4d8bf6fc3dffe4cf1c8f0f8700ba90f4
	//account-cell-type 0x03316fb2b69fd159545b566dbc6f76d3dc1af159ebfdaf88ef04e406c5a90747
	//dpoint-cell-type 0xe927000070b6cc8b20f0907511ebc71dfdca874631665bddce1e6cb3edfbcad5
	//
	//CharSetDigit 0xbf51e3db4991d14cd5673a5b921e5ba43b5be7790affe52db0522f98760d5df1
	//ConfigCellApply 0x4c6b3436a4cecb98e30b2ef007a7d572dc8275fdd550196f4015fd03d8abd468
	//ConfigCellIncome 0x7893fe9d4d94444a9dd4a46c3d0d8192d718fda4648e463b6588fc3ac9d6abd1
	//ConfigCellProfitRate 0x1277a5e133e76e7e99e1e9a9cb350f2db2c57e1179c1575fc5b9e0ccd18bc094
	//ConfigCellSubAccount 0xc328cb4517283dff08202af3560e455b8b41dd12fb3216a5d9637334169cc8b6
	//PreservedAccount06 0x9d943d9d250ec9a5334b2221b4c97ebdcede9ecd57b28d76f98a3e103cc01b4d
	//PreservedAccount11 0x2dac14930859ca748fa868fe94fbd74481dc97b6aede33afb44ba867d08a9920
	//PreservedAccount18 0xf54f6bb5ee6f8fa3e36eda6e875b5687b7dae41db9488e67f89ca7d90c4f1198
	//CharSetRu 0x67ba188e1669f3279f14b4847de4f1f3bde470e7c46b0ef57eeefe718cbae6a7
	//ConfigCellRelease 0x8e77aac24337f25cd048800547214be990449a24b61161056159861536a0b6fe
	//ConfigCellReverseRecord 0x10ad9923086ba2acfa259776d52c057e4a2d8029c33bcbdb8420a00a7a1523f9
	//ConfigCellTypeArgsDPoint 0x879d0fab1c733aa3ed05136db45ee6a657e10bd8fe9950886e04c10bf1dfdfa0
	//PreservedAccount04 0xeab9413a8ecd01e1804657ca0e3b064326fefcb1377da6723f7d482f1f9ab82c
	//PreservedAccount08 0x0c7e18918e85ca3d5a4bbcecbcc229f52550dbb96a6bc79a4727e9e2f705df7f
	//PreservedAccount14 0x7ea2c5ee592ad353bf37a39b5f47a3daf170ead4d870f8c31c0a8102aecc7dae
	//CharSetEmoji 0x51b6856350eb28b4b25f6a22aa7d4b33d2d8a12d3a583912e401ba16ab9370d9
	//ConfigCellMain 0xd0d3d7669f1daea227d2ba156b5d56663c94422bcb65dc63fe4ca7c541e6b4dc
	//ConfigCellRecordNamespace 0x88e24f1989758700428135530f6e0d944ca6a78880e2980e5c3695245a7bdcca
	//PreservedAccount00 0x829585f1d4573b3bb62ae7dad51526d27aa012343378e98a5f3bee78ab019f98
	//PreservedAccount01 0x70a620f2979db005c1d3f8ed7b85c68d75524d0cae9af61d7c9be6934b6c6238
	//PreservedAccount19 0xd9396ed8a26bab05b3fe2b1ead6106a0df98e00e843102aeb8d48bfc706e8057
	//CharSetVi 0x30e29931fde86782f380dc452ca820ae3805892405255ceaa89bf8d43b2d030d
	//PreservedAccount02 0x1ed2a65bee97d51b25b4e16f16ff3e0c754bcab76c3ed62c708da96dbb819f43
	//PreservedAccount16 0x57ff26ec98db15a9dcf544a8045186d9ae2d575296d0828b255c6db9318b4285
	//CharSetEn 0x90a7248e6104e5ae3fd89457b015e54e7282ad91b3670a350095fbd35d8dcc13
	//CharSetTr 0x6ef6afa3565f05022f7c455e522ea3ac6be045c105c1274aee96dbf1458c4852
	//ConfigCellAccount 0x6e673f8d3cd37013d26dee4fc05e9564fab99450c4d4a96c2160d6cd9de6d526
	//PreservedAccount03 0x973ba8c4151b46610ca7b6953d09b24013caaba255e200b0593b981f1220669a
	//PreservedAccount10 0x1946bf1e1a63f44a2f1e295e686cb0dbec91437ebe57cd3cd74a9b764def2c2a
	//CharSetJp 0x508cd9c813d1f56bf4822a9fe0995e7cb64ee7098ad6ac11b9960dd9b0a5d388
	//CharSetKo 0x3c5247c838312e0896ed1673bdd7e2ccee890c7888051d0d4a65e710fc68cc6a
	//ConfigCellPrice 0x948bb1a3a4540600c2229dcf76072e793cf6d48b79cda7fc0db4b23fdf9f7c73
	//ConfigCellProposal 0x64d74471a4cb46114fa2e273328966193f241da8134422dee9b55eb8a429b9ff
	//ConfigCellUnavailable 0x5fb5d65e46502407ea01f88dc76739c66b2a2dd228d6bf16d0420c02b2853674
	//PreservedAccount07 0x15ac9c6ac382958c81dbffbc57665182e3fc74b8bcc765d0ebeb367a34cda3ad
	//PreservedAccount13 0x432184ce7d30435988917e5dfe41c8ecafb835815c96958db7f2adb9690e476e
	//CharSetTh 0x21f5823d792142891c61a87a731543cb1cdce34cc73930f1c200bb521be0e354
	//PreservedAccount05 0xe06ef3d0ce77119233a1bd5e393681c65237d9f41e5edb1f6ca9500e597bf693
	//PreservedAccount12 0x85c1e2b13cde6de2218dc8b809d97df0b56d61d73d0a7e7632e6ca046a779c75
	//PreservedAccount17 0x4adf6ca822b8a9e5eede57dbcd9955c3ba9dfbe8d1a970d0dd1b9b9e11cf8da4
	//CharSetHanT 0xa0048cb33e52dbdd28b69ed5cfd34b17614cec3cedc7fb7410c4e2ad6a48cecf
	//ConfigCellSecondaryMarket 0x617b5b0e084749b596cabd9ecb049beb49e2016db1057f961d6d14fcec76decc
	//ConfigCellSubAccountWhiteList 0x29eab3da6b2a3c8f3c627bdb0f8b3c41a07b7073348a032db3bb285d78dd5ddc
	//ConfigCellTypeArgsSystemStatus 0xddfdf51c663bbe8f3b59044dceed499e2fdeee9decc82e4883393bc5e24d4e28
	//ConfigCellTypeArgsSMTNodeWhitelist 0x16b42743ce99c75fe1a649d5fd2f2db83bcf3fec86d6d5bcade8e6352fab18d8
	//PreservedAccount09 0x3a7e49ed7c07fdc90edeaf76bd72e7d846796cf108e241e19b357859add4ca98
	//PreservedAccount15 0xa547ed326117771df7eb2b440321500b9fb9ccdace5b522af5daf9aac16ebdbb
	//CharSetHanS 0x95f887ed322e9cf194ca612696f85929ecf881848e9995cb50beb1a009d715c7

	// contract
	//cont, err := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(cont.ContractName, cont.ContractTypeId, cont.OutPoint.TxHash.Hex())
	// config cell
	//cc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsMain)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(cc.Name, cc.OutPoint.TxHash.Hex(), cc.OutPoint.Index)
	// so script
	//script, err := core.GetDasSoScript(common.SoScriptBitcoin)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(script.Name, script.OutPoint.TxHash.Hex(), script.OutPoint.Index)
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
