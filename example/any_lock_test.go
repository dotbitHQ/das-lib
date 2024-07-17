package example

import (
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"testing"
)

// address format payload
func TestAddressFormatPayload2(t *testing.T) {
	fmt.Println(common.ChainTypeBitcoin.ToString())
	fmt.Println(common.ChainTypeBitcoin.ToDasAlgorithmId(true))
	fmt.Println(common.DasAlgorithmIdBitcoin.ToChainType())
	fmt.Println(common.DasAlgorithmIdBitcoin.ToCoinType())
	fmt.Println(common.DasAlgorithmIdBitcoin.ToSoScriptType())
	fmt.Println(common.FormatCoinTypeToDasChainType(common.CoinTypeBTC))
	fmt.Println(common.FormatDasChainTypeToCoinType(common.ChainTypeBitcoin))
	fmt.Println(common.FormatAddressByCoinType(string(common.CoinTypeBTC), "147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM"))

	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	//daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	daf := dc.Daf()
	res, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeBitcoin,
		AddressNormal: "tb1qumrp5k2es0d0hy5z6044zr2305pyzc978qz0ju", //"bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf", //"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM",
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	res2, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeBitcoin,
		AddressNormal: "mk8b5rG8Rpt1Gc61B8YjFk1czZJEjPDSV8", //"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM",
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.DasSubAlgorithmId, res.ChainType, res.AddressHex, res.Payload())

	res1, err := daf.HexToNormal(res2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res1:", res1.ChainType, res1.AddressNormal)

	lockScrip, _, err := daf.HexToScript(res2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(lockScrip.Args))

	args, err := daf.HexToArgs(res, res2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(args))

	owner, manager, err := daf.ArgsToNormal(args)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(owner.ChainType, owner.AddressNormal, manager.ChainType, manager.AddressNormal)

	oHex, mHex, err := daf.ScriptToHex(lockScrip)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(oHex.ChainType, oHex.DasAlgorithmId, oHex.DasSubAlgorithmId, oHex.AddressHex, oHex.Payload())
	fmt.Println(mHex.ChainType, mHex.DasAlgorithmId, mHex.DasSubAlgorithmId, mHex.AddressHex, mHex.Payload())

	cta := core.ChainTypeAddress{
		Type: "blockchain",
		KeyInfo: core.KeyInfo{
			CoinType: common.CoinTypeBTC,
			ChainId:  "",
			Key:      "bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf", //"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM",
		},
	}
	hexAddr, err := cta.FormatChainTypeAddress(common.DasNetTypeMainNet, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hexAddr.ChainType, hexAddr.DasAlgorithmId, hexAddr.DasSubAlgorithmId, hexAddr.AddressHex, hexAddr.Payload())
}

func TestFormatAnyLock(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrHex, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeCkb,
		AddressNormal: "ckt1qp4wtmsvhzrm9h66ngvpxuc4hx7u2klg65nr0vk7qcjqjt2lpjga2qgqng57g6ce9y0lj3rul6fjkgl3jsk47p5u96hss9",
		//AddressNormal: "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjzk3ntzys3nuwmvnar2lrs54l9pat6wy3qq5glj65",
		Is712: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrHex.AddressHex, addrHex.ChainType, addrHex.DasAlgorithmId, addrHex.DasSubAlgorithmId)
	anyLockHex, err := addrHex.FormatAnyLock()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(anyLockHex.AddressHex, anyLockHex.ChainType, anyLockHex.DasAlgorithmId, anyLockHex.DasSubAlgorithmId)
}

func TestSigVerify(t *testing.T) {
	msg := "From .bit: 7992045aec9c90e39f48addb28ccf3f8e07893c5a6ab8625fa51513f11638062"
	sig := "0x7df51a5d516cd0595bf4a202277931599b1c07720da24c81e25c368f5d5371fb56cca7c8a330e6a0612d24e70bd4189b33b05321b5adcc61967cf385ca4b62ca010100"
	addr := "5ef634a3ddc0b2cf9a6804c6a3cc3251ea5c8e44"
	verifyRes, _, err := http_api.VerifySignature(common.DasAlgorithmIdBitcoin, msg, sig, addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(verifyRes)
}

func TestBtcReverse(t *testing.T) {
	str := `{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0x55542796de898b0d792fa607175c5b100fc704352f8693990719975380a009fd","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x37e84424afadb18d5abea00f9abcfa1821ccd644a9b657feb5d960b1ce6687f8","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xbc0cbfe61010776302d3a0c8bef47d14529f73550f7122e441e5db32e28193a2","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x7ccc73c799ef509840323149adb28161e76caca8fcf2816070eb414b631076e4","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x77cdb8d076e3780ef46c42e8f473e9ec2ea1d9521e1cf8ee0db9efb01671d341","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf249a946f1302c34d63d437eaf345ce77b96c91f142cef3c356ec16f0ecc3f34","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x33ddf8335bb61dd61570c54582afc5d82c4ff45fc353037529423f5dee743430","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x97cf6fb6d0500d677f6a4989b90216e0adf5ddf4869b58b484c600781e86c983","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x06c4da6db0f5f9df3df09a17b026379871b82035a2229985f9a03c807af0d29b","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x17077b8116677d152718a27d5af0d0c4b12c5767ea25a0fdb61ea365daf507fa","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x41ee6bafb7a7a0ad65232c49a1bb3daa85f476041c505acfce8a2cea73442a3f","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xe3938457b467dde9a31a1d987fd8862aa9d678b8e3c9d9ea8ecb7ae568e65a0b","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x5394d678301851ac563fb512bc2bb99a4bd6ff38fddd6e5ecaf41607062e8140","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xe982cc62a53629312fd5132c6d06a3e44ee440c3acd29fdbe78ed22198ac96c1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x53d9a95eca84b8f32dc84495c454298a4b28957c71aff781e2dde1225086e388","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0x447a67b0b1728478729c56aea07c2380a2b682a0073f04f67429d66e6d966db4","index":"0x0"},"dep_type":"code"}],"header_deps":["0x465dce92e57c6fdd9992508213cc5c59729260e9c32c63b1425240e77f0097aa"],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0x226f89ed59fbbac38acd6f3744d4817d013dc33a92ad5f7e617c71a698d12daf","index":"0x0"}},{"since":"0x0","previous_output":{"tx_hash":"0x226f89ed59fbbac38acd6f3744d4817d013dc33a92ad5f7e617c71a698d12daf","index":"0x1"}}],"outputs":[{"capacity":"0x4a817c800","lock":{"code_hash":"0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f","hash_type":"type","args":"0x"},"type":{"code_hash":"0x8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5f","hash_type":"type","args":"0x"}},{"capacity":"0xe8d4630614","lock":{"code_hash":"0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8","hash_type":"type","args":"0xda44ed9db97056a06e471d3a1b6a1b82219e7232"},"type":null}],"outputs_data":["0x68743c8ed1b3f67ad393b619a870d210acbd86ebb9f6bac536f8c15c581cca56","0x"],"witnesses":["0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","0x5500000010000000550000005500000041000000a5efd14b9cb5d1428ff1267bfe7d4a77ff5ea5d89ff1d50363227185c539285424d843bbfcdc5cfcab365b81f54a76419dd85d8cf3e2180348f31b4d6cec944501","0x646173000000002f0000000c0000002a0000001a0000007570646174655f726576657273655f7265636f72645f726f6f740100000000","0x6461730a000000040000000100000006000000757064617465430000007df51a5d516cd0595bf4a202277931599b1c07720da24c81e25c368f5d5371fb56cca7c8a330e6a0612d24e70bd4189b33b05321b5adcc61967cf385ca4b62ca0101000100000009140000005ef634a3ddc0b2cf9a6804c6a3cc3251ea5c8e44ea0000004c4ffa504f3528a830511d93952465424186d2b55f46c0967335bc7d97ddf0aad97e5abc51fb2273a351c4867f6604b17d3285f073f9afa1c95b8b53deb506a3cc824871e16eda6c8e75bd81e2d9d60d1ce2df4256186cffe598ed3ce3477719a0ae32dcf80150a4c5137c3d4b212e126910433aa3c00adf46a427b5426adf9fefcdf1c7ab2b13500425c9884c96ee9ba0e75725de6cc6d134d752f25090e8d8098a2e56deb01a56502bd6035a9f621b4f3b07c340ebbea00e65f71389bd426991a2d1f0b86f98ca62506264846d3d53b85ffbd9a736e81aa6e26318158090bee751d729ee8b97a04f4800000000000000002000000068743c8ed1b3f67ad393b619a870d210acbd86ebb9f6bac536f8c15c581cca560c00000032303234303533302e626974","0x646173680000009d040000140000001500000035020000790300000120020000400000006000000080000000a0000000c0000000e00000000001000020010000400100006001000080010000a0010000c0010000e0010000000200001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4e8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5f9986d68bbf798e21238f8e5f58178354a8aeb7cc3f38e2abcb683e6dbb08f7375988ce37f185904477f120742b191a0730da0d5de9418a8bdf644e6bb3bd8c124401000024000000480000006c00000090000000b4000000d8000000fc000000200100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002401000024000000440000006400000084000000a4000000c4000000e400000004010000c9fc9f3dc050f8bf11019842a2426f48420f79da511dd169ee243f455e9f84ed991bcf61b6d7a26e6c27bda87d5468313d99ef0cd37113eee9e16c2680fa4532ebb79383a2947f36a095b434dd4f7c670dec6c2a53d925fb5c5f949104e59a6f6d0f4c38ae82383c619b9752ed8140019aa49128e39d48b271239a668c40a174f8f6b58d548231bc6fe19c1a1ceafa3a429f54c21a458b211097ebe564b146157ab1b06d51c579d528395d7f472582bf1d3dce45ba96c2bff2c19e30f0d90281b2d54e4da02130a9f7a9067ced1996180c0f2b122a6399090649a1050a66b2d82b8d30fdc9419104531fc1f2c5019c7ca061d438d534281fe3128dbd4acba5d9","0x646173700000002800000010000000180000002000000000c817a80400000000e1f505000000001027000000000000","0x646173740000002400000061336d521b8c43e3b38686c3923f05051a1e0416ff556907b37a6ee06ce84246"]}`
	tx, err := rpc.TransactionFromString(str)
	if err != nil {
		t.Fatal(err)
	}
	txReverseSmtRecord := make([]*witness.ReverseSmtRecord, 0)
	if err := witness.ParseFromTx(tx, common.ActionDataTypeReverseSmt, &txReverseSmtRecord); err != nil {
		t.Fatal(err)
	}
	for _, v := range txReverseSmtRecord {
		fmt.Println(v.GetP2SHP2WPKH(common.DasNetTypeTestnet2))
		fmt.Println(v.GetP2TR(common.DasNetTypeTestnet2))
	}
	//0x24cf4a6e349f84dcc3e122245c5abca848d7c1280e271556bb26bdbc4c29490f
	//0x7df51a5d516cd0595bf4a202277931599b1c07720da24c81e25c368f5d5371fb56cca7c8a330e6a0612d24e70bd4189b33b05321b5adcc61967cf385ca4b62ca01
	//0xdb1fd4ab021e6cab45c5386fc1444c5c2a927cea57189a36082f04470055b060
	//0x7df51a5d516cd0595bf4a202277931599b1c07720da24c81e25c368f5d5371fb56cca7c8a330e6a0612d24e70bd4189b33b05321b5adcc61967cf385ca4b62ca01
}
