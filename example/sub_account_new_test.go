package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
	"time"
)

func TestSubAccountMintSign(t *testing.T) {
	sams := witness.SubAccountMintSign{
		Version:            witness.SubAccountMintSignVersion1,
		Signature:          []byte{},
		ExpiredAt:          uint64(time.Now().Unix()),
		AccountListSmtRoot: []byte{},
	}
	dataBys := sams.GenSubAccountMintSignBytes()

	var sanb witness.SubAccountNewBuilder
	res, _ := sanb.ConvertSubAccountMintSignFromBytes(dataBys)
	fmt.Println(res.Version, res.ExpiredAt, res.Signature, res.AccountListSmtRoot)
}

func TestSubAccountNew(t *testing.T) {
	san := witness.SubAccountNew{
		Version:   0,
		Signature: nil,
		SignRole:  nil,
		NewRoot:   nil,
		Proof:     nil,
		Action:    "",
		SubAccountData: &witness.SubAccountData{
			Lock: &types.Script{
				CodeHash: types.Hash{},
				HashType: "",
				Args:     nil,
			},
			AccountId:            common.Bytes2Hex(common.GetAccountIdByAccount("aaa.bit")),
			AccountCharSet:       nil,
			Suffix:               "",
			RegisteredAt:         0,
			ExpiredAt:            0,
			Status:               0,
			Records:              nil,
			Nonce:                0,
			EnableSubAccount:     0,
			RenewSubAccountPrice: 0,
		},
		EditKey:        "",
		EditValue:      nil,
		EditLockArgs:   nil,
		EditRecords:    nil,
		RenewExpiredAt: 0,
		PrevRoot:       nil,
		CurrentRoot:    nil,
	}
	dataBys, err := san.GenSubAccountNewBytes()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dataBys)

	var sanb witness.SubAccountNewBuilder
	subAcc, err := sanb.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)

	subAcc.Version = witness.SubAccountNewVersion2
	dataBys, err = subAcc.GenSubAccountNewBytes()
	fmt.Println(dataBys)
	subAcc, err = sanb.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)
}

func TestSubAccountNewMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xa2178d7bd194fcd9f9d7533081ee51a0ba76e4028448052a02473a59958a50c7"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		var sanb witness.SubAccountNewBuilder
		resMap, err := sanb.SubAccountNewMapFromTx(res.Transaction)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range resMap {
			fmt.Println(k, v.SubAccountData.AccountId, v.EditKey, v.EditRecords, v.EditLockArgs, v.RenewExpiredAt)
		}
	}

}
