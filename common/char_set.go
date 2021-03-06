package common

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/molecule"
	"strings"
)

type AccountCharType uint32

const (
	AccountCharTypeEmoji  AccountCharType = 0
	AccountCharTypeNumber AccountCharType = 1
	AccountCharTypeEn     AccountCharType = 2
)

var CharSetTypeEmojiMap = make(map[string]struct{})

const (
	CharSetTypeNumber = "0123456789-"
	CharSetTypeEn     = "abcdefghijklmnopqrstuvwxyz"
)

type AccountCharSet struct {
	CharSetName AccountCharType `json:"char_set_name"`
	Char        string          `json:"char"`
}

func AccountCharsToAccount(accountChars *molecule.AccountChars) string {
	index := uint(0)
	var accountRawBytes []byte
	accountCharsSize := accountChars.ItemCount()
	for ; index < accountCharsSize; index++ {
		char := accountChars.Get(index)
		accountRawBytes = append(accountRawBytes, char.Bytes().RawData()...)
	}
	accountStr := string(accountRawBytes)
	if accountStr != "" && !strings.HasSuffix(accountStr, DasAccountSuffix) {
		accountStr = accountStr + DasAccountSuffix
	}
	return accountStr
}

func AccountToAccountChars(account string) ([]AccountCharSet, error) {
	if index := strings.Index(account, "."); index > 0 {
		account = account[:index]
	}

	chars := []rune(account)
	var list []AccountCharSet
	for _, v := range chars {
		char := string(v)
		var charSetName AccountCharType
		if _, ok := CharSetTypeEmojiMap[char]; ok {
			charSetName = AccountCharTypeEmoji
		} else if strings.Contains(CharSetTypeNumber, char) {
			charSetName = AccountCharTypeNumber
		} else if strings.Contains(CharSetTypeEn, char) {
			charSetName = AccountCharTypeEn
		} else {
			return nil, fmt.Errorf("invilid char type")
		}
		list = append(list, AccountCharSet{
			CharSetName: charSetName,
			Char:        char,
		})
	}
	return list, nil
}

func ConvertToAccountCharSets(accountChars *molecule.AccountChars) []AccountCharSet {
	index := uint(0)
	var accountCharSets []AccountCharSet
	for ; index < accountChars.ItemCount(); index++ {
		char := accountChars.Get(index)
		charSetName, _ := molecule.Bytes2GoU32(char.CharSetName().RawData())
		accountCharSets = append(accountCharSets, AccountCharSet{
			CharSetName: AccountCharType(charSetName),
			Char:        string(char.Bytes().RawData()),
		})
	}
	return accountCharSets
}

func ConvertToAccountChars(accountCharSet []AccountCharSet) *molecule.AccountChars {
	accountCharsBuilder := molecule.NewAccountCharsBuilder()
	for _, item := range accountCharSet {
		if item.Char == "." {
			break
		}
		accountChar := molecule.NewAccountCharBuilder().
			CharSetName(molecule.GoU32ToMoleculeU32(uint32(item.CharSetName))).
			Bytes(molecule.GoBytes2MoleculeBytes([]byte(item.Char))).Build()
		accountCharsBuilder.Push(accountChar)
	}
	accountChars := accountCharsBuilder.Build()
	return &accountChars
}

func InitEmoji(emojis []string) {
	for _, v := range emojis {
		CharSetTypeEmojiMap[v] = struct{}{}
	}
	//fmt.Println(CharSetTypeEmojiMap)
}

func GetAccountCharType(res map[AccountCharType]struct{}, list []AccountCharSet) {
	if res == nil {
		return
	}
	for _, v := range list {
		res[v.CharSetName] = struct{}{}
	}
}
