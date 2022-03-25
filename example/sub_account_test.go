package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/dascache"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sync"
	"testing"
)

func TestSubAccountCellFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xd2ed490f6cec9543291b3b730d0f38a2e46258c8848c6ec7ac12a6f9fa0ffd7f"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("SMTRoot")
		for k, v := range res.Transaction.Outputs {
			contract, err := core.GetDasContractInfo(common.DASContractNameSubAccountCellType)
			if err != nil {
				t.Fatal(err)
			}
			if contract.IsSameTypeId(v.Type.CodeHash) {
				fmt.Println(common.Bytes2Hex(res.Transaction.OutputsData[k]))
			}
		}
	}
}

func TestSubAccountBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x3cefd87b4c0102e3679ea456ac3766df6028296ba7e2d51185ccc5a29399ec49"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.SubAccountBuilderFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(common.Bytes2Hex(builder.PrevRoot))
		fmt.Println(common.Bytes2Hex(builder.CurrentRoot))
		fmt.Println(builder.Version)
		fmt.Println(builder.Account)
		fmt.Println(builder.SubAccount)
	}
}

func TestNewSubAccountWitness(t *testing.T) {
	p := witness.SubAccountParam{
		Signature:      nil,
		PrevRoot:       nil,
		CurrentRoot:    nil,
		Proof:          nil,
		SubAccount:     nil,
		EditKey:        "",
		EditLockScript: nil,
		RenewExpiredAt: 0,
		EditRecords:    nil,
	}
	bys, err := p.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
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
	bys := common.Hex2Bytes("0x646173080000000000000001000000020100000003010000000404000000010000002401000024010000300000008f000000a3000000f2000000fe000000060100000e0100000f010000130100001b0100001c0100005f0000001000000030000000310000008bb0413701cdd2e3a661cc8914e6790e16d619ce674930671e695807274bd14c012a00000005c9f53b1d85356b60453f867610888d89a0b667ad0515a33588908cf8edb27d1abe3852bf287abd3891338e9410a195ddf7fedccd99834ea6c5b6e5449c4f00000010000000250000003a000000150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061150000000c00000010000000020000000100000061080000002e6161612e62697401000000000000000200000000000000000400000000000000000000000000000000000000000000000000000000")
	res, err := witness.SubAccountBuilderFromBytes(bys[common.WitnessDasTableTypeEndIndex:])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
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
		PrevRoot:       []byte{2},
		CurrentRoot:    []byte{3},
		Proof:          []byte{4},
		SubAccount:     &subAccount,
		EditKey:        "",
		EditLockScript: nil,
		EditRecords:    nil,
		RenewExpiredAt: 0,
	}
	bys, err := param.NewSubAccountWitness()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
}
