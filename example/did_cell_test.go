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
		Args:     common.Hex2Bytes("0x01"),
	}
	fmt.Println(dc.GetDidCellOccupiedCapacity(&anyLock, "20240509.bit"))
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

	oldPK, err := schnorr.ParsePubKey(addrTR.ScriptAddress())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("oldPK:", hex.EncodeToString(oldPK.SerializeCompressed()))

	//addressPubKey2, err := btcutil.NewAddressPubKey(oldPK.SerializeCompressed(), &net)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//pkHash2 := addressPubKey2.AddressPubKeyHash().Hash160()[:]
	//fmt.Println("pkHash2:", hex.EncodeToString(pkHash2))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	addrNormal, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId:    common.DasAlgorithmIdBitcoin,
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2TR, //common.DasSubAlgorithmIdBitcoinP2PKH,
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

func TestUniSatP2PKH(t *testing.T) {
	//pk := "025a946b0635ba7540a5dfe1f7a6656bda6e0f17e64e9f0384d962a33d053aee2f"
	//addr := "mk8b5rG8Rpt1Gc61B8YjFk1czZJEjPDSV8"
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
		DasSubAlgorithmId: common.DasSubAlgorithmIdBitcoinP2WPKH, //common.DasSubAlgorithmIdBitcoinP2PKH,
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
