package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"

	//"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

//
//func TestGetDidCellOccupiedCapacity(t *testing.T) {
//	didCell := types.CellOutput{
//		Capacity: 0,
//		Lock: &types.Script{
//			CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
//			HashType: types.HashTypeType,
//			Args:     common.Hex2Bytes("0x01"),
//		},
//		Type: &types.Script{
//			CodeHash: types.HexToHash("0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f"),
//			HashType: types.HashTypeType,
//			Args:     nil,
//		},
//	}
//
//	defaultWitnessHash := molecule.Byte20Default()
//	didCellData := witness.DidCellData{
//		ItemId:      witness.ItemIdDidCellDataV0,
//		Account:     "20240509.bit",
//		ExpireAt:    0,
//		WitnessHash: common.Bytes2Hex(defaultWitnessHash.RawData()),
//	}
//	didCellDataBys, err := didCellData.ObjToBys()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
//	fmt.Println(didCellCapacity)
//}

func TestGetDidCellOccupiedCapacity2(t *testing.T) {
	dc, _ := getNewDasCoreTestnet2()

	anyLock := types.Script{
		CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes("0x045ef634a3ddc0b2cf9a6804c6a3cc3251ea5c8e4400"),
	}
	fmt.Println(dc.GetDidCellOccupiedCapacity(&anyLock, "12345.bit"))
}

//func TestTxToDidCellAction(t *testing.T) {
//	dc, _ := getNewDasCoreTestnet2()
//	res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash("0x4b5cb65d2203d00d755133797feced8c0e43292cb60cb2b0b4ebcab0ac917024"))
//	action, _ := dc.TxToDidCellAction(res.Transaction)
//	fmt.Println(action)
//}

func TestUniSatP2TR(t *testing.T) {
	pk := "02c0f888c8490ca3f7f095222b91afb2efbf5d21915506dd3f77e915b845eaaf17"
	addr := "tb1pzl9nkuavvt303hly08u3ug0v55yd3a8x86d8g5jsrllsaell8j5s8gzedg"
	net := bitcoin.GetBTCTestNetParams()
	res, err := btcutil.DecodeAddress(addr, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.EncodeAddress(), hex.EncodeToString(res.ScriptAddress()))

	addressPubKey, err := btcutil.NewAddressPubKey(common.Hex2Bytes(pk), &net)
	if err != nil {
		t.Fatal(err)
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
	fmt.Println("pkHash:", hex.EncodeToString(pkHash[:]), len(pkHash))

	//

	fmt.Println("addressPubKey.PubKey():", hex.EncodeToString(addressPubKey.PubKey().SerializeCompressed()))
	tapKey := txscript.ComputeTaprootKeyNoScript(addressPubKey.PubKey())
	fmt.Println("tapKey:", hex.EncodeToString(tapKey.SerializeCompressed()))
	addrTR, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey),
		&net,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrTR.EncodeAddress(), hex.EncodeToString(addrTR.ScriptAddress()))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrNormal, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdBitcoin,
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2TR,
		AddressHex:        hex.EncodeToString(res.ScriptAddress()),
		AddressPayload:    nil,
		IsMulti:           false,
		ChainType:         common.ChainTypeBitcoin,
		ParsedAddress:     nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrNormal.AddressNormal)

	addrHex, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeBitcoin,
		AddressNormal: addrNormal.AddressNormal,
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrHex.DasAlgorithmId, addrHex.DasSubAlgorithmId, addrHex.AddressHex)
}

func TestUniSatP2SHP2WPKH(t *testing.T) {
	pk := "03878e05aea052e38eded01b290b2e5e9c4f20182d0d3f2bd2932def1200c9aed6"
	addr := "2MtBBWYHheRwj1zLrf5KA6j68XmMkBbtzAS"
	net := bitcoin.GetBTCTestNetParams()
	res, err := btcutil.DecodeAddress(addr, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.EncodeAddress(), hex.EncodeToString(res.ScriptAddress()))

	addressPubKey, err := btcutil.NewAddressPubKey(common.Hex2Bytes(pk), &net)
	if err != nil {
		t.Fatal(err)
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
	fmt.Println("pkHash:", hex.EncodeToString(pkHash[:]))
	addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &net)
	if err != nil {
		t.Fatal(err)
	}

	pkScript, err := txscript.PayToAddrScript(addressWPH)
	if err != nil {
		t.Fatal(err)
	}
	scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(scriptAddr.EncodeAddress(), hex.EncodeToString(scriptAddr.ScriptAddress()))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrNormal, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdBitcoin,
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2SHP2WPKH, //common.DasSubAlgorithmIdBitcoinP2PKH,
		AddressHex:        hex.EncodeToString(res.ScriptAddress()),
		AddressPayload:    nil,
		IsMulti:           false,
		ChainType:         common.ChainTypeBitcoin,
		ParsedAddress:     nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrNormal.AddressNormal)

	addrHex, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeBitcoin,
		AddressNormal: addrNormal.AddressNormal,
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrHex.DasAlgorithmId, addrHex.DasSubAlgorithmId, addrHex.AddressHex)
}

func TestUniSatP2WPKH(t *testing.T) {
	pk := "0262c6eb28bc42cc168f61319dfa54fa64267bc3626ab05094cd1195fdf49a3009"
	addr := "tb1qumrp5k2es0d0hy5z6044zr2305pyzc978qz0ju"

	net := bitcoin.GetBTCTestNetParams()
	res, err := btcutil.DecodeAddress(addr, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.EncodeAddress(), hex.EncodeToString(res.ScriptAddress()))

	addrPK, err := btcutil.NewAddressPubKey(common.Hex2Bytes(pk), &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrPK.EncodeAddress(), hex.EncodeToString(addrPK.AddressPubKeyHash().ScriptAddress()), hex.EncodeToString(addrPK.ScriptAddress()))
	pkHash := addrPK.AddressPubKeyHash().Hash160()[:]
	addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addressWPH.EncodeAddress(), hex.EncodeToString(addressWPH.ScriptAddress()))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrNormal, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdBitcoin,
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2WPKH,
		AddressHex:        hex.EncodeToString(res.ScriptAddress()),
		AddressPayload:    nil,
		IsMulti:           false,
		ChainType:         common.ChainTypeBitcoin,
		ParsedAddress:     nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrNormal.AddressNormal)
}

func TestUniSatP2PKH(t *testing.T) {
	pk := "025a946b0635ba7540a5dfe1f7a6656bda6e0f17e64e9f0384d962a33d053aee2f"
	addr := "mk8b5rG8Rpt1Gc61B8YjFk1czZJEjPDSV8"

	net := bitcoin.GetBTCTestNetParams()
	res, err := btcutil.DecodeAddress(addr, &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.EncodeAddress(), hex.EncodeToString(res.ScriptAddress()))

	addrPK, err := btcutil.NewAddressPubKey(common.Hex2Bytes(pk), &net)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrPK.EncodeAddress(), hex.EncodeToString(addrPK.AddressPubKeyHash().ScriptAddress()), hex.EncodeToString(addrPK.ScriptAddress()))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrNormal, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdBitcoin,
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2PKH,
		AddressHex:        hex.EncodeToString(res.ScriptAddress()),
		AddressPayload:    nil,
		IsMulti:           false,
		ChainType:         common.ChainTypeBitcoin,
		ParsedAddress:     nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addrNormal.AddressNormal)
}

func TestSporeData(t *testing.T) {

	defaultWitnessHash := molecule.Byte20Default()
	dcdLV := witness.DidCellDataLV{
		Flag:        0,
		Version:     0,
		WitnessHash: defaultWitnessHash.RawData(),
		ExpireAt:    1714201479,
		Account:     "test.bit",
	}
	contentBys, err := dcdLV.ObjToBys()
	if err != nil {
		t.Fatal(err)
	}
	sd := witness.SporeData{
		ContentType: []byte{},
		Content:     contentBys,
		ClusterId:   witness.GetClusterId(common.DasNetTypeTestnet2),
	}

	res, err := sd.ObjToBys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", common.Bytes2Hex(res))

	if err = sd.BysToObj(res); err != nil {
		t.Fatal(err)
	}

	dcdLV2, err := sd.ContentToDidCellDataLV()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dcdLV2)
}

//
//func TestBysToDidCellData(t *testing.T) {
//	s, _, err := witness.BysToDidCellData(common.Hex2Bytes("0x66000000100000001400000042000000000000002a0000000001a7d4860aaf1dc83daedf75d6022811d2c2ae250b1b666d660000000032303233303631362e62697420000000cdb443dd0f9d98f530fd8945b86f3ea946f56ee4d015882beb757571bbd529f1"))
//	if err != nil {
//		t.Fatal(err)
//	}
//	c, err := s.ContentToDidCellDataLV()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println("ExpireAt:", c.ExpireAt)
//}

func TestAnyLockCodeHash(t *testing.T) {
	h := "0x65a7ee8deea9f4ca61aee11c4f4a04b349393f8b26672767e10b0ba2fd19badf"
	c, err := getClientTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.GetTransaction(context.Background(), types.HexToHash(h))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res.Transaction.CellDeps {
		fmt.Println(v.OutPoint.TxHash.String())
		tx, err := c.GetTransaction(context.Background(), v.OutPoint.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		typeId := common.ScriptToTypeId(tx.Transaction.Outputs[v.OutPoint.Index].Type)
		fmt.Println("typeId:", typeId.String())
	}

}

func TestGetAnyLockOutpoint(t *testing.T) {
	dc, _ := getNewDasCoreMainNet()
	res, err := dc.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.OutPoint.TxHash.String())
}

//func TestRenewTx(t *testing.T) {
//	dc, _ := getNewDasCoreTestnet2()
//	h := types.HexToHash("0xa4907d7f014d83426ed83fb9540537a1cc2e22c0ef8c893b0e619d77213e127a")
//	tx, _ := dc.Client().GetTransaction(context.Background(), h)
//
//	var oldDidCellOutpoint string
//	didCellAction, _ := dc.TxToDidCellAction(tx.Transaction)
//	fmt.Println(didCellAction)
//
//	txDidEntity, _ := witness.TxToDidEntity(tx.Transaction)
//
//	oldDidCellOutpoint = common.OutPointStruct2String(tx.Transaction.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
//	fmt.Println(oldDidCellOutpoint)
//	//didCellInfo.Outpoint = common.OutPoint2String(req.Tx.Hash.Hex(), uint(txDidEntity.Outputs[0].Target.Index))
//	//didCellInfo.ExpiredAt = accountInfo.ExpiredAt
//	//didCellInfo.BlockNumber = accountInfo.BlockNumber
//}

func TestAddr(t *testing.T) {
	cta := core.ChainTypeAddress{
		Type: "blockchain",
		KeyInfo: core.KeyInfo{
			CoinType: common.CoinTypeCKB,
			ChainId:  "",
			Key:      "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgjqgygrcl4k7pjuafdzmlzwy8ws4dxja7uqqmuv8c5", //"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM",
		},
	}
	hexAddr, err := cta.FormatChainTypeAddress(common.DasNetTypeMainNet, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hexAddr.ChainType, hexAddr.DasAlgorithmId, hexAddr.DasSubAlgorithmId, hexAddr.AddressHex, hexAddr.Payload())
}

func TestAnyLockTypeArgs(t *testing.T) {
	h := "0x86e00aec4ec93fc7b8690be66cfeb5147f29dce1e85e8a277fee00ef471d999e"
	c, _ := getClientTestnet2()
	res, _ := c.GetTransaction(context.Background(), types.HexToHash(h))

	bys, _ := common.GetDidCellTypeArgs(res.Transaction.Inputs[0], 1)

	fmt.Println(common.Bytes2Hex(bys))
}

func TestGetDidEntityFromTx(t *testing.T) {
	//c, _ := getClientMainNet()
	//
	//h := "0xceef5f05fcc18875bc3f99a3b410ea2435a97e52e10769f5e1cb3f9c92e7b1d3"
	//tx, err := c.GetTransaction(context.Background(), types.HexToHash(h))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//res, err := witness.GetDidEntityFromTx(tx.Transaction)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(len(res.Inputs), len(res.Outputs))
	//for k, v := range res.Outputs {
	//	fmt.Println(k, v.Hash())
	//}

	str := `{"version":"0x0","cell_deps":[{"out_point":{"tx_hash":"0xb88e5384d37ff91c963082d54b457b3ed404cb658fddad0956fe23942aaacb4e","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xe014f0bfec433816c36f3e40f5b70fd2618418bf15a8ff8314a959b3ffb198b1","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x266f9651d6c29935529ba559112a73d9524b28e215f32b074c21f9a2a4b58fd4","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x181715c3290858010ffd166f6bd334ed561b9c3cd97a9270544298ef063691d9","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0xb348e0338f22ed6e2f5e1cb467a9d62cdc1a4d684f9c13c8287229391ab0ae0e","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x9cbdf906b183dea8c3474e75a2648c70e54ac58b2d8f984d4db6140d522d423f","index":"0x1"},"dep_type":"code"},{"out_point":{"tx_hash":"0x30f946501bb4980a1c61e66eea3778be4fc8c356e2c28460ecffca1fa4338356","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x0978fc2cb86b3e771f9f2ae7e0ea59324a8185b74752c6904109963e9ebe9248","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x2917ade8b8a5e222a6ae86839ab8bdf1e324fef1f716fd9e27f8122ec0397c7f","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4dcf3f3b09efac8995d6cbee87c5345e812d310094651e0c3d9a730f32dc9263","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0xec18bf0d857c981c3d1f4e17999b9b90c484b303378e94de1a57b0872f5d4602","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x4d2e707d015079afcdb2a5d153f601c1239c01dfe5ca1fea3e20e125d2030291","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8eadc1fc9f8ec8315e09f0b2fc064e52a9184d63467082f03c9c755f6182b12","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x97cf6fb6d0500d677f6a4989b90216e0adf5ddf4869b58b484c600781e86c983","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x838d61207008ccd33fc8349a47f7af70cedea12be84d85427341e7c206696f9d","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x065ee2a9d9d8bf4d920cfd21fbd7648b0fc34d362701a9356f623401545a55da","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x0aab06d601056cbfd6c01312b6a0fc6c8b440f4f1c0bbad78f06209779267c1e","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xb35e2dcb608ff7373b6f51ae1dca751451f7d9033e71d083c51068f5fc2df30b","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x5394d678301851ac563fb512bc2bb99a4bd6ff38fddd6e5ecaf41607062e8140","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x0013ddcb388677e703f1816b53f5f17dc5d29b06ebf70398d1a2db6509d9a5c5","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0x53d9a95eca84b8f32dc84495c454298a4b28957c71aff781e2dde1225086e388","index":"0x0"},"dep_type":"code"},{"out_point":{"tx_hash":"0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37","index":"0x0"},"dep_type":"dep_group"},{"out_point":{"tx_hash":"0x9a730f91df450b47af0bd47be8ec4005f9a2599ac3e3a6c7aab22839569f453b","index":"0x0"},"dep_type":"code"}],"header_deps":[],"inputs":[{"since":"0x0","previous_output":{"tx_hash":"0xba6b198a926f6ab51a901cb3131ae04f644e535daf7644238b65b90c2fe65a51","index":"0x0"}},{"since":"0x0","previous_output":{"tx_hash":"0xb54e8c5a469598b7aea73eac1c9db9bcf13ada629a410cb2b31e2ebcd29ae192","index":"0x0"}},{"since":"0x0","previous_output":{"tx_hash":"0xba6b198a926f6ab51a901cb3131ae04f644e535daf7644238b65b90c2fe65a51","index":"0x2"}},{"since":"0x0","previous_output":{"tx_hash":"0x96caae566012cbd87447998936da81ff9c04a761bd41637a15639eb0d8c60b18","index":"0x2"}},{"since":"0x0","previous_output":{"tx_hash":"0xf59f5cb280442a5981fc412e4f9260a969d2dac7009694525387417d35e2787a","index":"0x2"}}],"outputs":[{"capacity":"0x51957680f","lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"0x0515a33588908cf8edb27d1abe3852bf287abd38910515a33588908cf8edb27d1abe3852bf287abd3891"},"type":{"code_hash":"0x1106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b1","hash_type":"type","args":"0x"}},{"capacity":"0x560ddf9cd","lock":{"code_hash":"0xf329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb","hash_type":"type","args":"0x04e6c61a595983dafb9282d3eb510d517d024160be00"},"type":{"code_hash":"0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f","hash_type":"type","args":"0x4a203e9eb5c4e45cdd0a9238218bfa2809799cd1f4c6a446c2c70bfa583d5a8f"}},{"capacity":"0x147959e212","lock":{"code_hash":"0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f","hash_type":"type","args":"0x"},"type":{"code_hash":"0x08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e","hash_type":"type","args":"0x"}},{"capacity":"0x132e270e6b","lock":{"code_hash":"0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8","hash_type":"type","args":"0xc866479211cadf63ad115b9da50a6c16bd3d226d"},"type":null}],"outputs_data":["0x9b519ed2089390e0adf990c7707777622637a02979879df6a4935de8ee3bffba7ed693e09a67d0767f2cb4e48e9e17e0f8a6a9097edd5a507a5f7160c7a3beb7cee405a06a036d01bdcefa6a0000000032303233313131372e626974","0x66000000100000001400000042000000000000002a0000000001f50a7c43f98400d3a6af9bcd0b8f6f2fda07524abdcefa6a0000000032303233313131372e6269742000000038ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71","0xc47237b93500b6a00871eb6436fcb19e3c246e9da4f59a969b31e14032680d22","0x"],"witnesses":["0x","0x6900000010000000690000006900000055000000550000001000000055000000550000004100000020fedf7721b64d7c459a4a81cb92ea880ff261cfc5d05094090f2b1edf6864f87312a8c94f3c85327de26a256f46d58bb4d0355ea27ecab518504a60a90a78488e","0x55000000100000005500000055000000410000006b03234ce20693827e1d7342fe8dd59c1627f90b33f06e4cfc20ceae9366d7c571581b1360c885fb95f624d3e3ec79d83b21e4c2416d58f674867331f832e53001","0x","0x","0x64617300000000210000000c0000001d0000000d00000072656e65775f6163636f756e7400000000","0x64617301000000ec02000010000000100000007e0100006e01000010000000140000001800000000000000040000005201000052010000300000004400000010010000180100002001000028010000300100003101000035010000360100003e0100007ed693e09a67d0767f2cb4e48e9e17e0f8a6a909cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000032150000000c00000010000000010000000100000030150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000031150000000c00000010000000010000000100000031150000000c00000010000000010000000100000031150000000c000000100000000100000001000000373d345765000000001ef97f6600000000000000000000000000000000000000009904000000000000000000000000140000000c0000001000000000000000000000006e01000010000000140000001800000000000000040000005201000052010000300000004400000010010000180100002001000028010000300100003101000035010000360100003e0100007ed693e09a67d0767f2cb4e48e9e17e0f8a6a909cc00000024000000390000004e00000063000000780000008d000000a2000000b7000000150000000c00000010000000010000000100000032150000000c00000010000000010000000100000030150000000c00000010000000010000000100000032150000000c00000010000000010000000100000033150000000c00000010000000010000000100000031150000000c00000010000000010000000100000031150000000c00000010000000010000000100000031150000000c000000100000000100000001000000373d345765000000001ef97f6600000000000000000000000000000000000000009904000000000000000000000000140000000c000000100000000000000000000000","0x64617306000000d2000000100000001000000010000000c20000001000000014000000180000000200000001000000a6000000a60000000c00000041000000350000001000000030000000310000000000000000000000000000000000000000000000000000000000000000000000000000000065000000080000005d0000000c00000055000000490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000efbf497f752ff7a655a8ec6f3c8f3feaaed6e41012e2597914000000","0x64617364000000a4000000480000004c000000540000005c000000600000006400000068000000700000007800000080000000880000008c0000009000000094000000980000009c000000a00000002a000000000edbcb0400000000e1f5050000000000a776002c0100008813000010270000000000001027000000000000102700000000000010270000000000002c0100002c0100002c0100002c010000f07e2700100e000000e1f505","0x64617369000000440100000c000000180000000c00000008000000f40100002c01000024000000450000006600000087000000a8000000c9000000ea0000000b0100002100000010000000110000001900000001ffffffffffffffffffffffffffffffff210000001000000011000000190000000280c3c9010000000080c3c901000000002100000010000000110000001900000003002d310100000000002d3101000000002100000010000000110000001900000004809698000000000080969800000000002100000010000000110000001900000005404b4c0000000000404b4c00000000002100000010000000110000001900000006404b4c0000000000404b4c00000000002100000010000000110000001900000007404b4c0000000000404b4c00000000002100000010000000110000001900000008404b4c0000000000404b4c0000000000","0x646173670000002400000010000000180000001c00000000c817a8040000003200000000f469b302000000","0x646173680000009d040000140000001500000035020000790300000120020000400000006000000080000000a0000000c0000000e00000000001000020010000400100006001000080010000a0010000c0010000e0010000000200001106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b10fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd14ff58f2c76b4ac26fdf675aa82541e02e4cf896279c6d6982d17b959788b2f0c08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b754467d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e449438124abdf4cbbfd61065e8b64523172bef5eefe27cb769c40acaf036aa89c200000000000000000000000000000000000000000000000000000000000000001a3f02aa89651a18112f0c21d0ae370a86e13f6a060c378184cd859a7bb6520361711416468fa5211ead5f24c6f3efadfbbc332274c5d40e50c6feadcb5f96068bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c4fd085557b4ef857b0577723bbf0a2e94081bbe3114de847cd9db01abaeb4f4e8041560ab6bd812c4523c824f2dcf5843804a099cb2f69fcbd57c8afcef2ed5f9986d68bbf798e21238f8e5f58178354a8aeb7cc3f38e2abcb683e6dbb08f7375988ce37f185904477f120742b191a0730da0d5de9418a8bdf644e6bb3bd8c124401000024000000480000006c00000090000000b4000000d8000000fc000000200100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002401000024000000440000006400000084000000a4000000c4000000e400000004010000c9fc9f3dc050f8bf11019842a2426f48420f79da511dd169ee243f455e9f84ed991bcf61b6d7a26e6c27bda87d5468313d99ef0cd37113eee9e16c2680fa4532ebb79383a2947f36a095b434dd4f7c670dec6c2a53d925fb5c5f949104e59a6f6d0f4c38ae82383c619b9752ed8140019aa49128e39d48b271239a668c40a174f8f6b58d548231bc6fe19c1a1ceafa3a429f54c21a458b211097ebe564b146157ab1b06d51c579d528395d7f472582bf1d3dce45ba96c2bff2c19e30f0d90281b2d54e4da02130a9f7a9067ced1996180c0f2b122a6399090649a1050a66b2d82b8d30fdc9419104531fc1f2c5019c7ca061d438d534281fe3128dbd4acba5d9"]}`
	txres, err := rpc.TransactionFromString(str)
	if err != nil {
		t.Fatal(err)
	}
	res, err := witness.GetDidEntityFromTx(txres)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(res.Inputs), len(res.Outputs))
	for k, v := range res.Outputs {
		fmt.Println(k, v.Hash())
	}
}

func TestTxToDidCellEntityAndAction(t *testing.T) {
	c, _ := getNewDasCoreTestnet2()
	h := "0x92c99efc4e8837fdebd247bc3bc689b4f5bab5afc456569a09dd21aa538af07c"
	tx, err := c.Client().GetTransaction(context.Background(), types.HexToHash(h))
	if err != nil {
		t.Fatal(err)
	}

	action := ""
	builder, err := witness.ActionDataBuilderFromTx(tx.Transaction)
	if err != nil {
		if err != witness.ErrNotExistActionData {
			t.Fatal(err)
		}
		didCellAction, _, err := c.TxToDidCellEntityAndAction(tx.Transaction)
		if err != nil {
			t.Fatal(err)
		}
		action = didCellAction
	} else {

		action = builder.Action
	}
	fmt.Println("parsingBlockData action:", action, h)

	//builder, err := witness.ActionDataBuilderFromTx(tx.Transaction)
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println(builder.Action)
	//}
	//action, res, err := c.TxToDidCellEntityAndAction(tx.Transaction)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(action)
	//for k, v := range res.Inputs {
	//	fmt.Println(k, v.OutPoint.TxHash, v.OutPoint.Index, v.Index)
	//	fmt.Println(common.Bytes2Hex(v.Lock.Args))
	//}
	//for k, v := range res.Outputs {
	//	fmt.Println(k, v.OutPoint.TxHash, v.OutPoint.Index, v.Index)
	//	fmt.Println(common.Bytes2Hex(v.Lock.Args))
	//}
}

func TestData(t *testing.T) {
	str := "0x66000000100000001400000042000000000000002a0000000001f50a7c43f98400d3a6af9bcd0b8f6f2fda07524abdcefa6a0000000032303233313131372e6269742000000038ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71"
	str = "0x66000000100000001400000042000000000000002a0000000001f50a7c43f98400d3a6af9bcd0b8f6f2fda07524abd6738670000000032303233313131372e6269742000000038ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71"
	bys := common.Hex2Bytes(str)
	var s witness.SporeData
	if err := s.BysToObj(bys); err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(s.ContentType))
	fmt.Println(common.Bytes2Hex(s.ClusterId))
	fmt.Println(common.Bytes2Hex(s.Content))
	d, err := s.ContentToDidCellDataLV()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d.Flag, d.Version, common.Bytes2Hex(d.WitnessHash), d.Account, d.ExpireAt)
	//0x0
	//0x38ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71
	//0x0001f50a7c43f98400d3a6af9bcd0b8f6f2fda07524abdcefa6a0000000032303233313131372e626974
	//0 1 0xf50a7c43f98400d3a6af9bcd0b8f6f2fda07524a 20231117.bit 1794821821
}

func TestCa(t *testing.T) {
	//co := types.CellOutput{
	//	Capacity: 0,
	//	Lock: &types.Script{
	//		CodeHash: types.HexToHash("0x493510d54e815611a643af97b5ac93bfbb45ddc2aae0f2dceffaf3408b4fcfcd"),
	//		HashType: "type",
	//		Args:     common.Hex2Bytes("0x4b000000100000003000000031000000f329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb011600000004591d65e5f82c02d139868889a91543ab327245cb000400000000000000000000004f3b75eb00"),
	//	},
	//	Type: &types.Script{
	//		CodeHash: types.HexToHash("0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f"),
	//		HashType: "type",
	//		Args:     common.Hex2Bytes("0x541407d01c47f578f27ac994dc3f97e1b3154bb1909e813caf6edeeb069d80be"),
	//	},
	//}
	//res := co.OccupiedCapacity(common.Hex2Bytes("0x66000000100000001400000042000000000000002a0000000001a7d4860aaf1dc83daedf75d6022811d2c2ae250b404f94670000000032303234303132362e6269742000000038ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71"))
	//
	//fmt.Println(res)
	//
	//s := common.Hex2Bytes("0x4b000000100000003000000031000000f329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb011600000004591d65e5f82c02d139868889a91543ab327245cb000400000000000000000000004f3b75eb00")
	//fmt.Println(len(s))
	//
	//p.NormalCellScript.OccupiedCapacity() * common.OneCkb,

	s := types.Script{
		CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: "type",
		Args:     common.Hex2Bytes("0x6d91285768e7c96f1cea0173c8167ada2cfeabe8"),
	}
	fmt.Println(s.OccupiedCapacity() * common.OneCkb)
}
