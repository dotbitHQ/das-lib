package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"testing"
)

func TestNormalToHex(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}

	res, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeCkbSingle,
		AddressNormal: "ckt1qyq639uzneswun3lkrj2hej4f8ky5hw6ltts0ycjj6",
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.ChainType, res.AddressHex, res.IsMulti)
	fmt.Println("=======================")

	res, err = daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeEth,
		AddressNormal: "0x15a33588908cF8Edb27D1AbE3852Bf287Abd3891",
		Is712:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.AddressHex, res.IsMulti)
	fmt.Println("=======================")

	res, err = daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeTron,
		AddressNormal: "TQoLh9evwUmZKxpD1uhFttsZk3EBs8BksV", //41A2AC25BF43680C05ABE82C7B1BCC1A779CFF8D5D
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.AddressHex, res.IsMulti)
	fmt.Println("=======================")

	res, err = daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeMixin,
		AddressNormal: "0xe1090ce82474cbe0b196d1e62ec349ec05a61076c68d14129265370ca7e051c4",
		Is712:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.AddressHex, res.IsMulti)
	fmt.Println("=======================")

	res, err = daf.NormalToHex(core.DasAddressNormal{
		ChainType:     100,
		AddressNormal: "0xe1090ce82474cbe0b196d1e62ec349ec05a61076c68d14129265370ca7e051c4",
		Is712:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.AddressHex, res.IsMulti)
	fmt.Println("=======================")

}

func TestHexToNormal(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	//res, err := daf.HexToNormal(core.DasAddressHex{
	//	DasAlgorithmId: common.DasAlgorithmIdCkb,
	//	AddressHex:     "0xa897829e60ee4e3fb0e4abe65549ec4a5ddafad7",
	//	IsMulti:        false,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	//fmt.Println("=======================")

	res, err := daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdCkbMulti,
		AddressHex:     "0xa897829e60ee4e3fb0e4abe65549ec4a5ddafad7",
		IsMulti:        true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("=======================")

	res, err = daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdEth712,
		AddressHex:     "0x15a33588908cF8Edb27D1AbE3852Bf287Abd3891",
		IsMulti:        false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("=======================")

	res, err = daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdTron,
		AddressHex:     "41a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d",
		IsMulti:        false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("=======================")

	res, err = daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdEd25519,
		AddressHex:     "0xe1090ce82474cbe0b196d1e62ec349ec05a61076c68d14129265370ca7e051c4",
		IsMulti:        false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("=======================")
}

func TestArgsToHex(t *testing.T) {

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	args := "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	ownerHex, managerHex, err := daf.ArgsToHex(common.Hex2Bytes(args))
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(ownerHex, managerHex)
	}
	fmt.Println("=======================")

	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x01a897829e60ee4e3fb0e4abe65549ec4a5ddafad701a897829e60ee4e3fb0e4abe65549ec4a5ddafad7"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")
	//
	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x01a897829e60ee4e3fb0e4abe65549ec4a5ddafad70315a33588908cf8edb27d1abe3852bf287abd3891"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")
	//
	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x0315a33588908cf8edb27d1abe3852bf287abd38910315a33588908cf8edb27d1abe3852bf287abd3891"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")
	//
	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x0515a33588908cf8edb27d1abe3852bf287abd38910515a33588908cf8edb27d1abe3852bf287abd3891"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")
	//
	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d04a2ac25bf43680c05abe82c7b1bcc1a779cff8d5d"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")
	//
	//ownerHex, managerHex, err = daf.ArgsToHex(common.Hex2Bytes("0x06e1090ce82474cbe0b196d1e62ec349ec05a61076c68d14129265370ca7e051c406e1090ce82474cbe0b196d1e62ec349ec05a61076c68d14129265370ca7e051c4"))
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	fmt.Println(ownerHex, managerHex)
	//}
	//fmt.Println("=======================")

}
