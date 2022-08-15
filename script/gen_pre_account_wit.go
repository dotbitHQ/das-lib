package main

import (
	"encoding/json"
	"flag"
	"fmt"
	dasCommon "github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

var ownerArgs = flag.String("o", "", "owner args, full format, include 0x")
var inviterArgs = flag.String("i", "0x053a6cab3323833f53754db4202f5741756c436ede053a6cab3323833f53754db4202f5741756c436ede", "inviter args, full format, include 0x")
var channelArgs = flag.String("c", "0x053a6cab3323833f53754db4202f5741756c436ede053a6cab3323833f53754db4202f5741756c436ede", "channel args, full format, include 0x")
var accountName = flag.String("a", "", "account's name, exclude .bit")
var inviterAccountName = flag.String("n", "", "inviter account's name, exclude .bit")
var registerTime = flag.Int64("t", 0, "timestamp in time cell")
var ckbPrice = flag.Uint64("q", 0, "ckb price in quota cell")
var dispatchTypeId = flag.String("l", "", "das lock's type id")
var charListJson = flag.String("j", "", "char set list in json")

func main() {
	flag.Parse()
	inviterScript := &types.Script{
		CodeHash: types.HexToHash(*dispatchTypeId),
		HashType: types.HashTypeType,
		Args:     common.FromHex(*inviterArgs),
	}

	channelScript := &types.Script{
		CodeHash: types.HexToHash(*dispatchTypeId),
		HashType: types.HashTypeType,
		Args:     common.FromHex(*channelArgs),
	}
	var content AccountCharStrList
	if *charListJson != "" {
		if err := json.Unmarshal([]byte(*charListJson), &content); err != nil {
			//fmt.Println("json unmarshal failed")
		}
	} else {
		name := *accountName
		if strings.HasSuffix(*accountName, ".bit") {
			name = strings.TrimSuffix(*accountName, ".bit")
		}
		content.AccountCharStr = AccountToCharSet(name)
	}
	//fmt.Println(content.AccountCharStr)
	nameLength := len(content.AccountCharStr) - 4
	//fmt.Println(nameLength)
	priceListConfigBuilder := molecule.NewPriceConfigListBuilder()
	priceList := []uint64{0, 1024000000, 1024000000, 660000000, 160000000, 5000000, 5000000, 5000000, 5000000}
	for nameLen := 1; nameLen < 9; nameLen++ {
		tmp := molecule.NewPriceConfigBuilder().
			Length(molecule.GoU8ToMoleculeU8(uint8(nameLen))).
			New(molecule.GoU64ToMoleculeU64(priceList[nameLen])).
			Renew(molecule.GoU64ToMoleculeU64(priceList[nameLen])).
			Build()
		priceListConfigBuilder.Push(tmp)
	}
	priceListConfig := priceListConfigBuilder.Build()

	var inviterAccountId []byte
	if *inviterAccountName == "" {
		inviterAccountId = common.FromHex("0x0000000000000000000000000000000000000000")
	} else {
		inviterAccountId = dasCommon.GetAccountIdByAccount(*inviterAccountName)
	}

	var preBuilder witness.PreAccountCellDataBuilder
	preWitness, preData, err := preBuilder.GenWitness(&witness.PreAccountCellParam{
		NewIndex:        0,
		Action:          "pre_register",
		CreatedAt:       *registerTime,
		InvitedDiscount: 500,
		Quote:           *ckbPrice,
		InviterScript:   inviterScript,
		ChannelScript:   channelScript,
		InviterId:       inviterAccountId,
		OwnerLockArgs:   common.FromHex(*ownerArgs),
		RefundLock:      inviterScript,
		Price:           *priceListConfig.Get(uint(nameLength - 1)),
		AccountChars:    AccountCharSetListToMoleculeAccountChars(content.AccountCharStr),
	})
	if err == nil {
		fmt.Println("0x" + common.Bytes2Hex(preData), "0x" + common.Bytes2Hex(preWitness))
	} else {
		fmt.Println(err.Error())
	}
}

type AccountCharStrList struct {
	AccountCharStr []dasCommon.AccountCharSet `json:"account_char_str"`
}
func AccountCharSetListToMoleculeAccountChars(list []dasCommon.AccountCharSet) molecule.AccountChars {
	accountChars := molecule.NewAccountCharsBuilder()
	for _, item := range list {
		if item.Char == "." {
			break
		}
		accountChar := molecule.NewAccountCharBuilder().
			CharSetName(molecule.GoU32ToMoleculeU32(uint32(item.CharSetName))).
			Bytes(molecule.GoBytes2MoleculeBytes([]byte(item.Char))).Build()
		accountChars.Push(accountChar)
	}
	return accountChars.Build()
}
func AccountToCharSet(account string) (accountChars []dasCommon.AccountCharSet) {
	chars := []rune(account)
	for _, v := range chars {
		char := string(v)
		charSetName := dasCommon.AccountCharTypeEmoji
		if strings.Contains("qwertyuiopasdfghjklzxcvbnm", char) {
			charSetName = dasCommon.AccountCharTypeEn
		} else if strings.Contains("1234567890", char) {
			charSetName = dasCommon.AccountCharTypeDigit
		}
		accountChars = append(accountChars, dasCommon.AccountCharSet{
			CharSetName: charSetName,
			Char:        char,
		})
	}
	return
}