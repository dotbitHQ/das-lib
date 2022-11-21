package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
	"time"
)

func TestSubAccountMintSign(t *testing.T) {
	var sab witness.SubAccountBuilderNew
	dataBys := sab.GenSubAccountMintSignBytes(witness.SubAccountMintSign{
		Version:            witness.SubAccountMintSignVersion1,
		Signature:          []byte{},
		ExpiredTimestamp:   uint32(time.Now().Unix()),
		AccountListSmtRoot: []byte{},
	})
	res, _ := sab.ConvertSubAccountMintSignFromBytes(dataBys)
	fmt.Println(res.Version, res.ExpiredTimestamp, res.Signature, res.AccountListSmtRoot)
}

func TestSubAccountNew(t *testing.T) {
	var sab witness.SubAccountBuilderNew
	dataBys, err := sab.GenSubAccountNewBytes(witness.SubAccountNew{
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
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dataBys)
	subAcc, err := sab.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)

	subAcc.Version = witness.SubAccountNewVersion2
	dataBys, err = sab.GenSubAccountNewBytes(*subAcc)
	fmt.Println(dataBys)
	subAcc, err = sab.ConvertSubAccountNewFromBytes(dataBys)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(subAcc.SubAccountData.AccountId, subAcc.Version)
}
