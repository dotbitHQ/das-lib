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
	//"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestGetDidCellOccupiedCapacity(t *testing.T) {
	didCell := types.CellOutput{
		Capacity: 0,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes("0x01"),
		},
		Type: &types.Script{
			CodeHash: types.HexToHash("0x0b1f412fbae26853ff7d082d422c2bdd9e2ff94ee8aaec11240a5b34cc6e890f"),
			HashType: types.HashTypeType,
			Args:     nil,
		},
	}

	defaultWitnessHash := molecule.Byte20Default()
	didCellData := witness.DidCellData{
		ItemId:      witness.ItemIdDidCellDataV0,
		Account:     "20240509.bit",
		ExpireAt:    0,
		WitnessHash: common.Bytes2Hex(defaultWitnessHash.RawData()),
	}
	didCellDataBys, err := didCellData.ObjToBys()
	if err != nil {
		t.Fatal(err)
	}

	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
	fmt.Println(didCellCapacity)
}

func TestGetDidCellOccupiedCapacity2(t *testing.T) {
	dc, _ := getNewDasCoreTestnet2()

	anyLock := types.Script{
		CodeHash: types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"),
		HashType: types.HashTypeType,
		Args:     common.Hex2Bytes("0x045ef634a3ddc0b2cf9a6804c6a3cc3251ea5c8e4400"),
	}
	fmt.Println(dc.GetDidCellOccupiedCapacity(&anyLock, "20230616.bit"))
}

func TestTxToDidCellAction(t *testing.T) {
	dc, _ := getNewDasCoreTestnet2()
	res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash("0x4b5cb65d2203d00d755133797feced8c0e43292cb60cb2b0b4ebcab0ac917024"))
	action, _ := dc.TxToDidCellAction(res.Transaction)
	fmt.Println(action)
}

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
		ClusterId:   common.Hex2Bytes(witness.ClusterId),
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

func TestBysToDidCellData(t *testing.T) {
	s, _, err := witness.BysToDidCellData(common.Hex2Bytes("0x66000000100000001400000042000000000000002a0000000001a7d4860aaf1dc83daedf75d6022811d2c2ae250b1b666d660000000032303233303631362e62697420000000cdb443dd0f9d98f530fd8945b86f3ea946f56ee4d015882beb757571bbd529f1"))
	if err != nil {
		t.Fatal(err)
	}
	c, err := s.ContentToDidCellDataLV()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ExpireAt:", c.ExpireAt)
}

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
	dc, _ := getNewDasCoreTestnet2()
	res, err := dc.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.OutPoint.TxHash.String())
}

func TestRenewTx(t *testing.T) {
	dc, _ := getNewDasCoreTestnet2()
	h := types.HexToHash("0xa4907d7f014d83426ed83fb9540537a1cc2e22c0ef8c893b0e619d77213e127a")
	tx, _ := dc.Client().GetTransaction(context.Background(), h)

	var oldDidCellOutpoint string
	didCellAction, _ := dc.TxToDidCellAction(tx.Transaction)
	fmt.Println(didCellAction)

	txDidEntity, _ := witness.TxToDidEntity(tx.Transaction)

	oldDidCellOutpoint = common.OutPointStruct2String(tx.Transaction.Inputs[txDidEntity.Inputs[0].Target.Index].PreviousOutput)
	fmt.Println(oldDidCellOutpoint)
	//didCellInfo.Outpoint = common.OutPoint2String(req.Tx.Hash.Hex(), uint(txDidEntity.Outputs[0].Target.Index))
	//didCellInfo.ExpiredAt = accountInfo.ExpiredAt
	//didCellInfo.BlockNumber = accountInfo.BlockNumber
}

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
