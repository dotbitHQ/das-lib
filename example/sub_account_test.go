package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/smt"
	"github.com/dotbitHQ/das-lib/witness"
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
		var sab witness.SubAccountNewBuilder
		builderMaps, err := sab.SubAccountNewMapFromTx(res.Transaction) //witness.SubAccountBuilderMapFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}

		for _, builder := range builderMaps {
			fmt.Println("--------------------------------------------")
			fmt.Println(common.Bytes2Hex(builder.PrevRoot))
			fmt.Println(common.Bytes2Hex(builder.CurrentRoot))
			fmt.Println(builder.Version)
			fmt.Println(builder.Account)
			fmt.Println(builder.SubAccountData)
			fmt.Println(builder.SubAccountData.ExpiredAt, builder.CurrentSubAccountData.ExpiredAt)
			fmt.Println(len(builder.SubAccountData.Records), len(builder.CurrentSubAccountData.Records))
			fmt.Println(common.Bytes2Hex(builder.SubAccountData.Lock.Args), common.Bytes2Hex(builder.CurrentSubAccountData.Lock.Args))
			fmt.Println(builder.SubAccountData.ToH256())
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
	var sab witness.SubAccountNewBuilder
	for _, hash := range hashList {
		if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
			t.Fatal(err)
		} else {
			builderMaps, err := sab.SubAccountNewMapFromTx(res.Transaction) //witness.SubAccountBuilderMapFromTx(res.Transaction)
			if err != nil {
				t.Fatal(err)
			}

			for _, builder := range builderMaps {
				fmt.Println(strings.Repeat("-", 100))
				fmt.Println(fmt.Sprintf("%-20s %s", "current root", common.Bytes2Hex(builder.CurrentRoot)))

				key := smt.AccountIdToSmtH256(builder.SubAccountData.AccountId)
				value := builder.CurrentSubAccountData.ToH256()
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
	var sab witness.SubAccountNewBuilder
	for _, v := range hashList {
		res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash(v))
		builderMaps, _ := sab.SubAccountNewMapFromTx(res.Transaction) //witness.SubAccountBuilderMapFromTx(res.Transaction)
		for k, v := range builderMaps {
			fmt.Println(k, v.CurrentSubAccountData.ToH256())
			key := smt.AccountIdToSmtH256(v.SubAccountData.AccountId)
			value := v.CurrentSubAccountData.ToH256()
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

	res2 := witness.ConvertSubAccountCellOutputData(common.Hex2Bytes("0xee27cf64f3dd463d7a3264fde553c3c2247c02ec4938d7b4961ab2096fe5b1be00f92fc6010000000060e8000700000001f15f519ecb226cd763b2bcbcab093e63f89100c07ac0caebc032c788b187ec993d7b681cb328ffec1323"))
	fmt.Println(res2.DasProfit, res2.OwnerProfit)
	//7520000000，29680000000
	//7620000000，30080000000
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
	res, _ := dc.Client().GetTransaction(context.Background(), types.HexToHash("0x3467032e97ca67cb2df10e4a922d4c9fc540bab1199e583338966cda2964816c"))
	fmt.Println(witness.ConvertCustomScriptConfigByTx(res.Transaction))
}

func TestArgsAndConfigHash(t *testing.T) {
	var de witness.SubAccountCellDataDetail
	fmt.Println(de.ArgsAndConfigHash())
}
