package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"log"
	"testing"
)

func TestAddrToScript(t *testing.T) {
	// CKB 地址
	ckbAddress := "ckt1qqexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6qggqu4qyfuzauwmj9k6qeenhmyt039rhu5xaqyqw2szy7pw78dezmdqvuemaj9hcj3m72rwsv94j9m"

	// 解析 CKB 地址
	parsedAddress, err := address.Parse(ckbAddress)
	if err != nil {
		log.Fatalf("Failed to parse CKB address: %v", err)
	}
	fmt.Println(common.Bytes2Hex(parsedAddress.Script.Args[:]))
	//parsedAddress.Script.Args
	//// 获取 Lock Script
	//lockScript, err := secp256k1.Script(parsedAddress.ScriptArgs)
	//if err != nil {
	//	log.Fatalf("Failed to get Lock Script: %v", err)
	//}
	//
	//// 打印 Lock Script
	//fmt.Printf("Lock Script: %x\n", lockScript)
}

func TestNormalToHex(t *testing.T) {
	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}

	//webauthn
	res, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeWebauthn,
		AddressNormal: "ckt1qqexmutxu0c2jq9q4msy8cc6fh4q7q02xvr7dc347zw3ks3qka0m6qggqajr5je2ylnz9jsuue986vvt2ld4v7f4hvyqwep6fv4z0e3zegwwvjnaxx940k6k0y6mkresszm",
		Is712:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.DasAlgorithmId, res.DasSubAlgorithmId, res.AddressHex, common.Bytes2Hex(res.AddressPayload), res.IsMulti)
	fmt.Println("=======================")

	res, err = daf.NormalToHex(core.DasAddressNormal{
		ChainType:     common.ChainTypeCkbSingle,
		AddressNormal: "ckt1qyq0wjp2jda08xztr7w2s0gqll4aa8z0nq4s9gnzg5",
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
	fmt.Println("222222 =======================")

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
		DasAlgorithmId: common.DasAlgorithmIdCkbSingle,
		AddressHex:     "0xd437b8e9ca16fce24bf3258760c3567214213c5",
		IsMulti:        true,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("111 =======================")

	res, err = daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdEth712,
		AddressHex:     "0xd437b8e9ca16fce24bf3258760c3567214213c5a",
		IsMulti:        false,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("05算法： ", res.ChainType, res.AddressNormal, res.Is712)
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

	res, err = daf.HexToNormal(core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdWebauthn,
		AddressHex:     "0x643a4b2a27e622ca1ce64a7d318b57db567935bb",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.ChainType, res.AddressNormal, res.Is712)
	fmt.Println("=======================")
}

func TestArgsToHex(t *testing.T) {

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	//webauthn args
	args := "0x08072a022782ef1db916da06733bec8b7c4a3bf286e808072a022782ef1db916da06733bec8b7c4a3bf286e8"

	ownerHex, managerHex, err := daf.ArgsToHex(common.Hex2Bytes(args))
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(ownerHex, managerHex)
		fmt.Println("owner cid1: ", common.Bytes2Hex(ownerHex.AddressPayload[:10]))
		fmt.Println("owner pk1: ", common.Bytes2Hex(ownerHex.AddressPayload[10:]))
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
