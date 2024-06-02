package example

import (
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
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
		AddressNormal: "ckt1qrejnmlar3r452tcg57gvq8patctcgy8acync0hxfnyka35ywafvkqgytmmrfg7aczevlxngqnr28npj2849erjyqqhe2guh",
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
