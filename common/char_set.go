package common

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/molecule"
	"strings"
)

type AccountCharType uint32

const (
	AccountCharTypeEmoji AccountCharType = 0
	AccountCharTypeDigit AccountCharType = 1
	AccountCharTypeEn    AccountCharType = 2  // English
	AccountCharTypeHanS  AccountCharType = 3  // Simplified
	AccountCharTypeHanT  AccountCharType = 4  // Traditional
	AccountCharTypeJp    AccountCharType = 5  // Japan
	AccountCharTypeKo    AccountCharType = 6  // Korea
	AccountCharTypeVi    AccountCharType = 7  // Vietnam
	AccountCharTypeRu    AccountCharType = 8  // Russian
	AccountCharTypeTh    AccountCharType = 9  // Thailand
	AccountCharTypeTr    AccountCharType = 10 // Turkey
)

var CharSetTypeEmojiMap = make(map[string]struct{})
var CharSetTypeDigitMap = make(map[string]struct{})
var CharSetTypeEnMap = make(map[string]struct{})
var CharSetTypeHanSMap = make(map[string]struct{})
var CharSetTypeHanTMap = make(map[string]struct{})
var CharSetTypeJpMap = make(map[string]struct{})
var CharSetTypeKoMap = make(map[string]struct{})
var CharSetTypeViMap = make(map[string]struct{})
var CharSetTypeRuMap = make(map[string]struct{})
var CharSetTypeThMap = make(map[string]struct{})
var CharSetTypeTrMap = make(map[string]struct{})

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
		} else if _, ok = CharSetTypeDigitMap[char]; ok {
			charSetName = AccountCharTypeDigit
		} else if _, ok = CharSetTypeEnMap[char]; ok {
			charSetName = AccountCharTypeEn
		} else if _, ok = CharSetTypeHanSMap[char]; ok {
			charSetName = AccountCharTypeHanS
		} else if _, ok = CharSetTypeHanTMap[char]; ok {
			charSetName = AccountCharTypeHanT
		} else if _, ok = CharSetTypeJpMap[char]; ok {
			charSetName = AccountCharTypeJp
		} else if _, ok = CharSetTypeKoMap[char]; ok {
			charSetName = AccountCharTypeKo
		} else if _, ok = CharSetTypeViMap[char]; ok {
			charSetName = AccountCharTypeVi
		} else if _, ok = CharSetTypeRuMap[char]; ok {
			charSetName = AccountCharTypeRu
		} else if _, ok = CharSetTypeThMap[char]; ok {
			charSetName = AccountCharTypeTh
		} else if _, ok = CharSetTypeTrMap[char]; ok {
			charSetName = AccountCharTypeTr
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

func InitEmojiMap(emojis []string) {
	for _, v := range emojis {
		CharSetTypeEmojiMap[v] = struct{}{}
	}
}

func InitDigitMap(numbers []string) {
	for _, v := range numbers {
		CharSetTypeDigitMap[v] = struct{}{}
	}
}

func InitEnMap(ens []string) {
	for _, v := range ens {
		CharSetTypeEnMap[v] = struct{}{}
	}
}

func InitHanSMap(hanSs []string) {
	for _, v := range hanSs {
		CharSetTypeHanSMap[v] = struct{}{}
	}
}

func InitHanTMap(hanTs []string) {
	for _, v := range hanTs {
		CharSetTypeHanTMap[v] = struct{}{}
	}
}

func InitJpMap(jps []string) {
	for _, v := range jps {
		CharSetTypeJpMap[v] = struct{}{}
	}
}

func InitKrMap(krs []string) {
	for _, v := range krs {
		CharSetTypeKoMap[v] = struct{}{}
	}
}

func InitVnMap(vns []string) {
	for _, v := range vns {
		CharSetTypeViMap[v] = struct{}{}
	}
}

func InitRuMap(rus []string) {
	for _, v := range rus {
		CharSetTypeRuMap[v] = struct{}{}
	}
}

func InitThMap(ths []string) {
	for _, v := range ths {
		CharSetTypeThMap[v] = struct{}{}
	}
}

func InitTrMap(trs []string) {
	for _, v := range trs {
		CharSetTypeTrMap[v] = struct{}{}
	}
}

func GetAccountCharType(res map[AccountCharType]struct{}, list []AccountCharSet) {
	if res == nil {
		return
	}
	for _, v := range list {
		res[v.CharSetName] = struct{}{}
	}
}

func GetAccountCharTypeExclude(res map[AccountCharType]struct{}, list []AccountCharSet) {
	if res == nil {
		return
	}
	for _, v := range list {
		if v.CharSetName == AccountCharTypeEmoji || v.CharSetName == AccountCharTypeDigit {
			continue
		}
		res[v.CharSetName] = struct{}{}
	}
}
