package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/smt"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
	"sync"
	"testing"
)

func TestGetSMTRoot(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xd2ed490f6cec9543291b3b730d0f38a2e46258c8848c6ec7ac12a6f9fa0ffd7f"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("SMTRoot")
		contract, err := core.GetDasContractInfo(common.DASContractNameSubAccountCellType)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range res.Transaction.Outputs {
			if v.Type != nil && contract.IsSameTypeId(v.Type.CodeHash) {
				fmt.Println(common.Bytes2Hex(res.Transaction.OutputsData[k]))
			}
		}
	}
}

func TestSubAccountBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x14bbfd4b1e576264bc844e74e7c7e5f39891877e3865621ef40f83d9c3d6cf20"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMaps, err := witness.SubAccountBuilderMapFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}

		for _, builder := range builderMaps {
			fmt.Println("--------------------------------------------")
			fmt.Println(common.Bytes2Hex(builder.PrevRoot))
			fmt.Println(common.Bytes2Hex(builder.CurrentRoot))
			fmt.Println(builder.Version)
			fmt.Println(builder.Account)
			fmt.Println(builder.SubAccount)
			fmt.Println(builder.SubAccount.ExpiredAt, builder.CurrentSubAccount.ExpiredAt)
			fmt.Println(len(builder.SubAccount.Records), len(builder.CurrentSubAccount.Records))
			fmt.Println(common.Bytes2Hex(builder.SubAccount.Lock.Args), common.Bytes2Hex(builder.CurrentSubAccount.Lock.Args))
			fmt.Println(common.Bytes2Hex(builder.SubAccount.ToH256()))
		}
	}
}

func TestGenActionDataWitnessV2(t *testing.T) {
	fmt.Println(witness.GenActionDataWitness(common.DasActionCreateSubAccount, nil))
	fmt.Println(witness.GenActionDataWitnessV2(common.DasActionCreateSubAccount, nil, common.ParamManager))
}

func TestGetRefund(t *testing.T) {
	client, _ := getClientMainNet()
	res, _ := client.GetTransaction(context.Background(), types.HexToHash("0xfd1153209e99d26cfffc1b3583b223548cb8259af09ea9e4dbb2c8391d9dde46"))
	for i, v := range res.Transaction.Outputs {
		fmt.Println(common.Bytes2Hex(v.Lock.Args), i, v.Capacity)
	}
}

func TestPre(t *testing.T) {
	var wg sync.WaitGroup
	ca := dascache.NewDasCache(context.Background(), &wg)
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: types.HexToHash("0x18ab87147e8e81000ab1b9f319a5784d4c7b6c98a9cec97d738a5c11f69e7254"),
			HashType: types.HashTypeType,
			Args:     nil,
		},
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script: &types.Script{
				CodeHash: types.HexToHash("0x303ead37be5eebfcf3504847155538cb623a26f237609df24bd296750c123078"),
				HashType: types.HashTypeType,
				Args:     nil,
			},
			OutputDataLenRange: &[2]uint64{52, 53}, // hash + account id
		},
	}
	ca.AddOutPoint([]string{"0x204f2b221eaef83af915139e6ff5cc4ca1a08e245ec7a513b9c8049f0920b6a9-0"})
	client, _ := getClientMainNet()
	cells, _ := core.GetSatisfiedLimitLiveCell(client, ca, searchKey, 400, indexer.SearchOrderAsc)
	for _, v := range cells {
		fmt.Println(common.OutPointStruct2String(v.OutPoint))
	}
}

func TestSubAccountBuilderFromBytes(t *testing.T) {
	bys := common.Hex2Bytes("0x6461730800000041000000f0dec13ab2e0ae4c41124c8e5fd1952c46a95e30933072138a7a3b3b0355a8fd0069dc85a146001b13953b20bfa415dd678c6e9feeea421ce7d75fb77620dbff01010000000020000000f956c0221862a44289cb7988f61615adc0263e34d2cca448828256bf6e2b5617200000005d3612fbc1d4445fb87ce09534d2a3bbc7c2451f6b56662a180fd33b1a376502680000004c4f9e519e4c31391b4bc4bd6c3aa37b1bd8b99476b5f1bfb2f0974d654f5cfec84f73da81d074be748970d8ed944054b1635d481b7d3322060000000000000000000000005073a3a4e3b147097e6cee17110384d6c9dcd0c97b24776233db805bda3128bbce4f6004000000010000009d0100009d010000300000008f000000a300000024010000360100003e01000046010000470100008c01000094010000950100005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a000000053a6cab3323833f53754db4202f5741756c436ede053a6cab3323833f53754db4202f5741756c436ede5559c11ec7dbcfec80f21dfe05454358df95090a81000000180000002d00000042000000570000006c000000150000000c00000010000000020000000100000061150000000c0000001000000002000000010000006c150000000c00000010000000020000000100000069150000000c00000010000000020000000100000063150000000c000000100000000200000001000000650e0000002e746f706269646465722e62697455524e6200000000d5852f64000000000045000000080000003d00000018000000230000002e00000032000000390000000700000070726f66696c65070000007477697474657200000000030000003132332c0100000100000000000000000000000000000000050000006f776e65722a000000052ce62318bc5bf2eeb34b8f2f12880fdef20e356f052ce62318bc5bf2eeb34b8f2f12880fdef20e356f")
	res, err := witness.SubAccountBuilderFromBytes(bys[common.WitnessDasTableTypeEndIndex:])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(res.SubAccount.Lock.Args), res.Account)
	fmt.Println(string(res.EditKey), common.Bytes2Hex(res.EditValue))
	fmt.Println(common.Bytes2Hex(res.SubAccount.ToH256()))
	res.SubAccount.Lock.Args = res.EditValue
	res.SubAccount.Nonce++
	res.SubAccount.Records = nil
	fmt.Println(common.Bytes2Hex(res.SubAccount.ToH256()))

	fmt.Println(string(common.Hex2Bytes("0x66726f6d206469643a2035633265343434626234343263323663383033316633346134386161346637363763396664633537376630663636663263363332386666636433306336663731")))
	bys2, _ := blake2b.Blake256(common.Hex2Bytes("0x5559c11ec7dbcfec80f21dfe05454358df95090a6f776e6572052ce62318bc5bf2eeb34b8f2f12880fdef20e356f052ce62318bc5bf2eeb34b8f2f12880fdef20e356f0000000000000000"))
	fmt.Println(common.Bytes2Hex(bys2))

}

func TestGenSubAccountBytes(t *testing.T) {
	account := "aaa.bit"

	accountCharSet, err := common.AccountToAccountChars(account)
	if err != nil {
		t.Fatal(err)
	}
	subAccount := witness.SubAccount{
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x8bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c"),
			HashType: types.HashTypeType,
			Args:     common.Hex2Bytes("0x05c9f53b1d85356b60453f867610888d89a0b667ad0515a33588908cf8edb27d1abe3852bf287abd3891"),
		},
		AccountId:            "0x338e9410a195ddf7fedccd99834ea6c5b6e5449c",
		AccountCharSet:       accountCharSet,
		Suffix:               ".aaa.bit",
		RegisteredAt:         1,
		ExpiredAt:            2,
		Status:               0,
		Records:              nil,
		Nonce:                0,
		EnableSubAccount:     0,
		RenewSubAccountPrice: 0,
	}
	param := witness.SubAccountParam{
		Signature:      nil,
		SignRole:       nil,
		PrevRoot:       []byte{2},
		CurrentRoot:    []byte{3},
		Proof:          []byte{4},
		SubAccount:     &subAccount,
		EditKey:        "",
		EditLockArgs:   nil,
		EditRecords:    nil,
		RenewExpiredAt: 0,
	}
	bys, err := param.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
}

func TestTestGenSubAccountBytes2(t *testing.T) {
	var param witness.SubAccountParam
	str := `{"Signature":null,"SignRole":null,"PrevRoot":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=","CurrentRoot":"2rrwuKes5atj5TKCbLBcyCroH8/yw/7b2fvBtOH9Pb8=","Proof":"TE8A","SubAccount":{"lock":{"code_hash":"0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd","hash_type":"type","args":"Bcn1Ox2FNWtgRT+GdhCIjYmgtmetBcn1Ox2FNWtgRT+GdhCIjYmgtmet"},"account_id":"0xdbcaa515cbd79477e17502a6e51dcdccadad8690","account_char_set":[{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"0"},{"char_set_name":1,"char":"1"}],"suffix":".0001.bit","registered_at":1648195840,"expired_at":1679731840,"status":0,"records":null,"nonce":0,"enable_sub_account":0,"renew_sub_account_price":0},"EditKey":"","EditLockScript":null,"EditRecords":null,"RenewExpiredAt":0}`
	_ = json.Unmarshal([]byte(str), &param)
	bys, err := param.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))

}

func TestSMTRootVerify(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hashList := []string{
		"0x23d91172a788a407ecc01c65bc02e6129f109d3a62b21ed7b3cc16daa93c465c",
		"0xff5f6472261175c588607f1a4a3d70829dcd48c6382a8f858ccd2a533c3e344b",
		"0xe7135538f29787308a06a9c5f789b6ecea647ba242631c38d8ee29281a556652",
		"0x14bbfd4b1e576264bc844e74e7c7e5f39891877e3865621ef40f83d9c3d6cf20",
	}
	tree := smt.NewSparseMerkleTree(nil)

	for _, hash := range hashList {
		if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
			t.Fatal(err)
		} else {
			builderMaps, err := witness.SubAccountBuilderMapFromTx(res.Transaction)
			if err != nil {
				t.Fatal(err)
			}

			for _, builder := range builderMaps {
				fmt.Println(strings.Repeat("-", 100))
				fmt.Println(fmt.Sprintf("%-20s %s", "current root", common.Bytes2Hex(builder.CurrentRoot)))

				key := smt.AccountIdToSmtH256(builder.SubAccount.AccountId)
				value := builder.CurrentSubAccount.ToH256()
				_ = tree.Update(key, value)
				root, _ := tree.Root()
				fmt.Println(fmt.Sprintf("%-20s %s", "tree value", common.Bytes2Hex(value)))
				fmt.Println(fmt.Sprintf("%-20s %s", "tree root", common.Bytes2Hex(root)))
			}
		}
	}
}

func TestConvertSubAccountCellOutputData(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	hashList := []string{
		"0x6058c58ed9f6b5565691399da68eca595bcff7a8722cc2097f5d9ccd61905d34",
		"0xa252c9f7601672426df90386b41f51b1488b51e06aa976a753bd4205f588df33",
		"0xb02c7d76f9b59d332bfa089d2a3a4fa66c8815e0ca581bce7329ff1f3ad2333e",
	}
	tree := smt.NewSparseMerkleTree(nil)
	for _, v := range hashList {
		res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash(v))
		builderMaps, _ := witness.SubAccountBuilderMapFromTx(res.Transaction)
		for k, v := range builderMaps {
			fmt.Println(k, common.Bytes2Hex(v.CurrentSubAccount.ToH256()))
			key := smt.AccountIdToSmtH256(v.SubAccount.AccountId)
			value := v.CurrentSubAccount.ToH256()
			_ = tree.Update(key, value)
		}
	}
	root, _ := tree.Root()
	fmt.Println("root:", root.String())
}

func TestBuildSubAccountCellOutputData(t *testing.T) {
	//fmt.Println(len(common.Hex2Bytes("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000")))
	//root := "0x9b9f8b3b4f7e1121a6a48d4b359285b55c6160bb889436d29d5ef4bb58821b9e"
	//profit := uint64(200)
	//
	//detail := witness.SubAccountCellDataDetail{
	//	SmtRoot:    common.Hex2Bytes(root),
	//	DasProfit:  profit,
	//	HashType:   nil,
	//	CustomArgs: nil,
	//}
	//res := witness.BuildSubAccountCellOutputData(detail)
	//
	//detailNew := witness.ConvertSubAccountCellOutputData(res)
	//fmt.Println(detailNew)

	res2 := witness.ConvertSubAccountCellOutputData(common.Hex2Bytes("0x303b6e3d3cf64bcbe3df4c6922c38aafd1f4bedb80f0726abd804293daa239a8008f6f640100000000b8e6790500000001f15f519ecb226cd763b2bcbcab093e63f89100c07ac0caebc032c788b187ec9993a7a331393b6c4b3b77"))
	fmt.Println(res2.DasProfit, res2.OwnerProfit)

	//rootBys, capacity := witness.ConvertSubAccountCellOutputData(common.Hex2Bytes("0x315f75ac13f22bd687004c944455347f386b11dbf49144bdfc50df5e2d5e554a007ddaac00000000"))
	//fmt.Println(common.Bytes2Hex(rootBys), capacity)
}

func TestGetCustomScriptConfig(t *testing.T) {
	//bys, hash := witness.BuildCustomScriptConfig(witness.CustomScriptConfig{
	//	Header:  witness.Script001,
	//	Version: 1,
	//	Body:    map[uint8]witness.CustomScriptPrice{
	//		//1: {1, 2},
	//		//2: {3, 4},
	//	},
	//	MaxLength: 0,
	//})
	//fmt.Println(common.Bytes2Hex(bys), len(hash))
	bys := common.Hex2Bytes("0x7363726970742d30303100000000e20000001c0000003d0000005e0000007f000000a0000000c1000000210000001000000011000000190000000540420f000000000040420f0000000000210000001000000011000000190000000640420f000000000040420f00000000002100000010000000110000001900000001404b4c0000000000404b4c0000000000210000001000000011000000190000000200093d000000000000093d00000000002100000010000000110000001900000003c0c62d0000000000c0c62d0000000000210000001000000011000000190000000480841e000000000080841e0000000000")
	res, err := witness.ConvertCustomScriptConfig(bys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Header, res.Version, res.Body, res.MaxLength)
	wit, hash := witness.BuildCustomScriptConfig(*res)
	fmt.Println(common.Bytes2Hex(wit))
	fmt.Println(common.Bytes2Hex(hash))
	//fmt.Println(common.Bytes2Hex([]byte("script-001")))
	//script-001 0 map[1:{5000000 5000000} 2:{4000000 4000000} 3:{3000000 3000000} 4:{2000000 2000000} 5:{1000000 1000000} 6:{1000000 1000000}] 6
	//script-001 0 map[1:{5000000 5000000} 2:{4000000 4000000} 3:{3000000 3000000} 4:{2000000 2000000} 5:{1000000 1000000} 6:{1000000 1000000}] 6
	//0x7363726970742d30303100000000e20000001c0000003d0000005e0000007f000000a0000000c1000000210000001000000011000000190000000200093d000000000000093d00000000002100000010000000110000001900000003c0c62d0000000000c0c62d0000000000210000001000000011000000190000000480841e000000000080841e0000000000210000001000000011000000190000000540420f000000000040420f0000000000210000001000000011000000190000000640420f000000000040420f00000000002100000010000000110000001900000001404b4c0000000000404b4c0000000000
	//0x7363726970742d30303100000000e20000001c0000003d0000005e0000007f000000a0000000c1000000210000001000000011000000190000000540420f000000000040420f0000000000210000001000000011000000190000000640420f000000000040420f00000000002100000010000000110000001900000001404b4c0000000000404b4c0000000000210000001000000011000000190000000200093d000000000000093d00000000002100000010000000110000001900000003c0c62d0000000000c0c62d0000000000210000001000000011000000190000000480841e000000000080841e0000000000
}

func TestConvertCustomScriptConfigByTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash("0xd9f965cff85edf5b3ef8a2beb039754b79a415d2c0356193c6f5cc07f4e0206a"))
	fmt.Println(witness.ConvertCustomScriptConfigByTx(res.Transaction))
}

func TestArgsAndConfigHash(t *testing.T) {
	var de witness.SubAccountCellDataDetail
	fmt.Println(de.ArgsAndConfigHash())
}
