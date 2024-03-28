package example

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/txscript"
	"github.com/dotbitHQ/das-lib/bitcoin"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"testing"
)

func TestCreateBTCWallet(t *testing.T) {
	err := bitcoin.CreateBTCWallet(bitcoin.BtcAddressTypeP2SHP2WPKH, true)
	if err != nil {
		t.Fatal(err)
	}
	//WIF: L2vKWmpxVFsRCQPxnhvjsLiYB3hTSV85fAm1Jo6CcAJkvgKqjxoh
	//PubKey: 147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM
	//PubHash 2222b81757f47ebe58881573f64fb8c5f59ba533
	//PriKey: aa13ee7c615ef80c9063bf6875fb894b3936c9551d73bfe0361a4682ae7efe8f

	//WIF: L3t7wxUjYs5A11kajfdQy2w1CnTKCbSxYFMMgstuYX7QraQt7nwb
	//ScriptAddr: 35Y6PCZk4zuP1GJkjrqqR7PpvgWbiMVuvx
	//PubHash d6c09590c8515eaaae150871b19a11cb44c54771
	//pkScript: 76a914d6c09590c8515eaaae150871b19a11cb44c5477188ac
	//pkScriptHash: 2a307b6ee071be7d8f484f1f0c06369742e46919
	//PriKey: c6c8a6bf98b562089e93e5f5270ea4468f3a442a88cccfcc74692bad458c32d3

	//WIF: KwVZNWG6fyqSh1uhVM25iNgNL89wxdbZcr3M5dnTtqdq4T4ZQfBt
	//PubKey: bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf
	//PubHash 39f04d79ada815dea22fc8aff1acb173981e8603
	//PriKey: 082720675b373fbaa6c24fb099867dfbbdeba98ab3c7c83c9ecb2ea26b5fa97d

	//WIF: KyMDvdf11J1CydwBNuMQ6uYVJXbV93j2FCi5ts2XZbVRPm7PeVvZ
	//ScriptAddr: 3A3basSqtJZPdA9mKCC1KtQkgXjKSSJnWc
	//PubHash 3a6274d504078fd35d21aff131eb22c7b1af13ef
	//pkScript: 00143a6274d504078fd35d21aff131eb22c7b1af13ef
	//pkScriptHash: 5ba56c93f710da685871a01afd2e47da5ca069b2
	//PriKey: 3f8a2671be95d5301e0bd7239a87ed9bb357e71545b3e8efbe89dfb1e932fdce
}

func TestDecodeAddr(t *testing.T) {
	addrStr := "3NJvsiRTQxhy2UmJCkhah32Vgr98zVsJoH"
	p := bitcoin.GetBTCMainNetParams()
	addr, err := btcutil.DecodeAddress(addrStr, &p)
	if err != nil {
		t.Fatal(err)
	}
	// 0x78d766bb7b4351b5faef9ea0bf476a8338b4caf9
	// 0xa7eb334e76b533ac82f9338cb626ae555ca21611
	fmt.Println(addr.String(), addr.EncodeAddress(), common.Bytes2Hex(addr.ScriptAddress()))
	//pa := bitcoin.GetBTCMainNetParams()
	//addressWPH, err := btcutil.NewAddressWitnessPubKeyHash(addr.ScriptAddress(), &pa)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	////
	//pkScript, err := txscript.PayToAddrScript(addressWPH)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//scriptAddr, err := btcutil.NewAddressScriptHash(pkScript, &pa)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("ScriptAddr:", scriptAddr.EncodeAddress())
}

func TestAddressDecodeP2SHP2WPKH(t *testing.T) {
	addr := "3A3basSqtJZPdA9mKCC1KtQkgXjKSSJnWc"
	payload, version, err := base58.CheckDecode(addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("payload:", hex.EncodeToString(payload))
	fmt.Println("version:", version)
	// payload: 5ba56c93f710da685871a01afd2e47da5ca069b2
	// version: 5

	//
	wifStr := "KyMDvdf11J1CydwBNuMQ6uYVJXbV93j2FCi5ts2XZbVRPm7PeVvZ"
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		t.Fatal(err)
	}
	pk := hex.EncodeToString(wif.SerializePubKey())
	fmt.Println("pk:", pk)
	// pk: 039f5988954fa9538ebb3c3ca630ed9ae71fddae7b1afa64dc39bd31011c463f5d

	//
	params := bitcoin.GetBTCMainNetParams()
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &params)
	if err != nil {
		t.Fatal(err)
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
	fmt.Println("pkHash:", hex.EncodeToString(pkHash), len(pkHash))
	// pkHash: 3a6274d504078fd35d21aff131eb22c7b1af13ef
	pkScript := append([]byte{txscript.OP_0, txscript.OP_DATA_20}, pkHash...)
	fmt.Println("pkScript：", hex.EncodeToString(pkScript))
	// pkScript: 00143a6274d504078fd35d21aff131eb22c7b1af13ef

	pkScriptHash := btcutil.Hash160(pkScript)
	fmt.Println("pkScriptHash:", hex.EncodeToString(pkScriptHash))
	// pkScriptHash: 5ba56c93f710da685871a01afd2e47da5ca069b2

	base58Addr := base58.CheckEncode(pkScriptHash, version)
	fmt.Println("base58Addr:", base58Addr)
	// base58Addr: 3A3basSqtJZPdA9mKCC1KtQkgXjKSSJnWc
}

func TestAddrDecodeP2WPKH(t *testing.T) {
	addr := "bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf"
	prefix, data, version, err := bech32.DecodeGeneric(addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("prefix:", prefix)
	fmt.Println("version:", version)
	fmt.Println("data:", hex.EncodeToString(data))
	//prefix: bc
	//version: 0
	//data: 00070718041a1e0d0d15000a1d1d08110f1902171f030b05110e0e0c011d011003

	// 5bit -> 8bit
	payload, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("payload:", hex.EncodeToString(payload), len(payload))
	//payload: 39f04d79ada815dea22fc8aff1acb173981e8603 20

	//
	wifStr := "KwVZNWG6fyqSh1uhVM25iNgNL89wxdbZcr3M5dnTtqdq4T4ZQfBt"
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		t.Fatal(err)
	}
	pk := hex.EncodeToString(wif.SerializePubKey())
	fmt.Println("pk:", pk)
	// pk: 03d1d583fe9ee37c30553e8b5b078684052e3eeccfad19212f750aa75fc550853d
	//
	params := bitcoin.GetBTCMainNetParams()
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &params)
	if err != nil {
		t.Fatal(err)
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
	fmt.Println("pkHash:", hex.EncodeToString(pkHash))
	// pkHash: 39f04d79ada815dea22fc8aff1acb173981e8603
	// 8bit -> 5bit
	converted, err := bech32.ConvertBits(pkHash, 8, 5, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("converted:", hex.EncodeToString(converted))
	// converted: 070718041a1e0d0d15000a1d1d08110f1902171f030b05110e0e0c011d011003
	// version+converted
	bech32Addr, err := bech32.Encode(prefix, append([]byte{0x00}, converted...))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("bech32Addr:", bech32Addr)
	// prefix + the separator 1 + data + the 6-byte checksum
	// bech32Addr: bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf
}

func TestAddrDecodeP2PKH(t *testing.T) {
	addr := "147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM"
	payload, version, err := base58.CheckDecode(addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("payload:", hex.EncodeToString(payload))
	fmt.Println("version:", version)
	// payload: 2222b81757f47ebe58881573f64fb8c5f59ba533
	// version: 0

	//
	wifStr := "L2vKWmpxVFsRCQPxnhvjsLiYB3hTSV85fAm1Jo6CcAJkvgKqjxoh"
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		t.Fatal(err)
	}
	pk := hex.EncodeToString(wif.SerializePubKey())
	fmt.Println("pk:", pk)
	// pk: 03d25ce14e816acb15f6c1a4032edd17236974086bb14f74e72a688343d695eb40

	//
	params := bitcoin.GetBTCMainNetParams()
	addressPubKey, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), &params)
	if err != nil {
		t.Fatal(err)
	}
	pkHash := addressPubKey.AddressPubKeyHash().Hash160()[:]
	fmt.Println("pkHash:", hex.EncodeToString(pkHash))
	// pkHash: 2222b81757f47ebe58881573f64fb8c5f59ba533
	base58Addr := base58.CheckEncode(pkHash, version)
	fmt.Println("base58Addr:", base58Addr)
	// base58Addr: 147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM
}

// address format payload
func TestAddressFormatPayload(t *testing.T) {
	fmt.Println(common.ChainTypeBitcoin.ToString())
	fmt.Println(common.ChainTypeBitcoin.ToDasAlgorithmId(true))
	fmt.Println(common.DasAlgorithmIdBitcoin.ToChainType())
	fmt.Println(common.DasAlgorithmIdBitcoin.ToCoinType())
	fmt.Println(common.DasAlgorithmIdBitcoin.ToSoScriptType())
	fmt.Println(common.FormatCoinTypeToDasChainType(common.CoinTypeBTC))
	fmt.Println(common.FormatDasChainTypeToCoinType(common.ChainTypeBitcoin))
	fmt.Println(common.FormatAddressByCoinType(string(common.CoinTypeBTC), "147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM"))

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	res, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeBitcoin,
		AddressNormal: "bc1q88cy67dd4q2aag30ezhlrt93wwvpapsruefmrf", //"147VZrBkaWy5zJhpuGAa7EZ9B9YBLu8MuM",
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.DasSubAlgorithmId, res.ChainType, res.AddressHex, res.Payload())

	res1, err := daf.HexToNormal(res)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res1.ChainType, res1.AddressNormal)

	lockScrip, _, err := daf.HexToScript(res)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(lockScrip.CodeHash.String(), hex.EncodeToString(lockScrip.Args))

	args, err := daf.HexToArgs(res, res)
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
}
