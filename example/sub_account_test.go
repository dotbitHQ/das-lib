package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
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
	fmt.Println(witness.GenActionDataWitnessV2(common.DasActionCreateSubAccount, nil, common.ParamManager))
}
