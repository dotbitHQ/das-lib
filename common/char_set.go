package common

import (
	"fmt"
	"github.com/Andrew-M-C/go.emoji/official"
	"github.com/clipperhouse/uax29/graphemes"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/ethereum/go-ethereum/common/math"
	"strings"
)

type AccountCharType uint32

const (
	AccountCharTypeEmoji AccountCharType = 0
	AccountCharTypeDigit AccountCharType = 1
	AccountCharTypeEn    AccountCharType = 2  // English
	AccountCharTypeHanS  AccountCharType = 3  // Chinese Simplified
	AccountCharTypeHanT  AccountCharType = 4  // Chinese Traditional
	AccountCharTypeJa    AccountCharType = 5  // Japanese
	AccountCharTypeKo    AccountCharType = 6  // Korean
	AccountCharTypeRu    AccountCharType = 7  // Russian
	AccountCharTypeTr    AccountCharType = 8  // Turkish
	AccountCharTypeTh    AccountCharType = 9  // Thai
	AccountCharTypeVi    AccountCharType = 10 // Vietnamese
)

var CharSetTypeEmojiMap = make(map[string]struct{})
var CharSetTypeDigitMap = make(map[string]struct{})
var CharSetTypeEnMap = make(map[string]struct{})
var CharSetTypeHanSMap = make(map[string]struct{})
var CharSetTypeHanTMap = make(map[string]struct{})
var CharSetTypeJaMap = make(map[string]struct{})
var CharSetTypeKoMap = make(map[string]struct{})
var CharSetTypeViMap = make(map[string]struct{})
var CharSetTypeRuMap = make(map[string]struct{})
var CharSetTypeThMap = make(map[string]struct{})
var CharSetTypeTrMap = make(map[string]struct{})

var AccountCharTypeMap = map[AccountCharType]map[string]struct{}{
	AccountCharTypeEmoji: CharSetTypeEmojiMap,
	AccountCharTypeDigit: CharSetTypeDigitMap,
	AccountCharTypeEn:    CharSetTypeEnMap,
	AccountCharTypeHanS:  CharSetTypeHanSMap,
	AccountCharTypeHanT:  CharSetTypeHanTMap,
	AccountCharTypeJa:    CharSetTypeJaMap,
	AccountCharTypeKo:    CharSetTypeKoMap,
	AccountCharTypeRu:    CharSetTypeRuMap,
	AccountCharTypeTr:    CharSetTypeTrMap,
	AccountCharTypeTh:    CharSetTypeThMap,
	AccountCharTypeVi:    CharSetTypeViMap,
}

var AccountCharTypeNameMap = map[string]AccountCharType{
	"Emoji":  AccountCharTypeEmoji,
	"Digit":  AccountCharTypeDigit,
	"En":     AccountCharTypeEn,
	"ZhHans": AccountCharTypeHanS,
	"ZhHant": AccountCharTypeHanT,
	"Ja":     AccountCharTypeJa,
	"Ko":     AccountCharTypeKo,
	"Ru":     AccountCharTypeRu,
	"Tr":     AccountCharTypeTr,
	"Th":     AccountCharTypeTh,
	"Vi":     AccountCharTypeVi,
}

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
		if v == "" {
			continue
		}
		CharSetTypeEmojiMap[v] = struct{}{}
	}
}

func InitDigitMap(numbers []string) {
	for _, v := range numbers {
		if v == "" {
			continue
		}
		CharSetTypeDigitMap[v] = struct{}{}
	}
}

func InitEnMap(ens []string) {
	for _, v := range ens {
		if v == "" {
			continue
		}
		CharSetTypeEnMap[v] = struct{}{}
	}
}

func InitHanSMap(hanSs []string) {
	for _, v := range hanSs {
		if v == "" {
			continue
		}
		CharSetTypeHanSMap[v] = struct{}{}
	}
}

func InitHanTMap(hanTs []string) {
	for _, v := range hanTs {
		if v == "" {
			continue
		}
		CharSetTypeHanTMap[v] = struct{}{}
	}
}

func InitJaMap(jas []string) {
	for _, v := range jas {
		if v == "" {
			continue
		}
		CharSetTypeJaMap[v] = struct{}{}
	}
}

func InitKoMap(kos []string) {
	for _, v := range kos {
		if v == "" {
			continue
		}
		CharSetTypeKoMap[v] = struct{}{}
	}
}

func InitRuMap(rus []string) {
	for _, v := range rus {
		if v == "" {
			continue
		}
		CharSetTypeRuMap[v] = struct{}{}
	}
}

func InitTrMap(trs []string) {
	for _, v := range trs {
		if v == "" {
			continue
		}
		CharSetTypeTrMap[v] = struct{}{}
	}
}

func InitThMap(ths []string) {
	for _, v := range ths {
		if v == "" {
			continue
		}
		CharSetTypeThMap[v] = struct{}{}
	}
}

func InitViMap(vis []string) {
	for _, v := range vis {
		if v == "" {
			continue
		}
		CharSetTypeViMap[v] = struct{}{}
	}
}

// GetAccountCharType 'res' for sub-account multi AccountCharType
func GetAccountCharType(res map[AccountCharType]struct{}, list []AccountCharSet) {
	if res == nil {
		return
	}
	for _, v := range list {
		res[v.CharSetName] = struct{}{}
	}
}

func CheckAccountCharTypeDiff(list []AccountCharSet) bool {
	var res = make(map[AccountCharType]struct{})
	for _, v := range list {
		if v.CharSetName == AccountCharTypeEmoji || v.CharSetName == AccountCharTypeDigit {
			continue
		}
		if v.Char == "." {
			break
		}
		res[v.CharSetName] = struct{}{}
	}
	if len(res) > 1 {
		return true
	}
	return false
}

func CheckAccountCharSetList(list []AccountCharSet) (account string, err error) {
	for i, v := range list {
		if v.Char == "" {
			err = fmt.Errorf("char[%d] is nil", i)
			return
		}
		switch v.CharSetName {
		case AccountCharTypeEmoji:
			if _, ok := CharSetTypeEmojiMap[v.Char]; !ok {
				err = fmt.Errorf("emoji char[%d] is nil", i)
				return
			}
		case AccountCharTypeDigit:
			if _, ok := CharSetTypeDigitMap[v.Char]; !ok {
				err = fmt.Errorf("digit char[%d] is nil", i)
				return
			}
		case AccountCharTypeEn:
			if _, ok := CharSetTypeEnMap[v.Char]; !ok {
				err = fmt.Errorf("en char[%d] is nil", i)
				return
			}
		case AccountCharTypeJa:
			if _, ok := CharSetTypeJaMap[v.Char]; !ok {
				err = fmt.Errorf("ja char[%d] is nil", i)
				return
			}
		case AccountCharTypeRu:
			if _, ok := CharSetTypeRuMap[v.Char]; !ok {
				err = fmt.Errorf("ru char[%d] is nil", i)
				return
			}
		case AccountCharTypeTr:
			if _, ok := CharSetTypeTrMap[v.Char]; !ok {
				err = fmt.Errorf("tr char[%d] is nil", i)
				return
			}
		case AccountCharTypeVi:
			if _, ok := CharSetTypeViMap[v.Char]; !ok {
				err = fmt.Errorf("vi char[%d] is nil", i)
				return
			}
		case AccountCharTypeTh:
			if _, ok := CharSetTypeThMap[v.Char]; !ok {
				err = fmt.Errorf("th char[%d] is nil", i)
				return
			}
		case AccountCharTypeKo:
			if _, ok := CharSetTypeKoMap[v.Char]; !ok {
				err = fmt.Errorf("ko char[%d] is nil", i)
				return
			}
		default:
			err = fmt.Errorf("char type [%d] is invalid", v.CharSetName)
			return
		}
		account += v.Char
	}
	return
}

// deprecated
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
		} else if _, ok = CharSetTypeJaMap[char]; ok {
			charSetName = AccountCharTypeJa
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

// deprecated
func GetAccountLength(account string) uint8 {
	account = strings.TrimSuffix(account, DasAccountSuffix)
	nextIndex := 0
	accLen := uint8(0)
	for i, _ := range account {
		if i < nextIndex {
			continue
		}
		match, length := official.AllSequences.HasEmojiPrefix(account[i:])
		if match {
			nextIndex = i + length
		}
		accLen++
	}
	return accLen
}

func GetDotBitAccountLength(account string) ([]string, int, error) {
	if !strings.HasSuffix(account, DasAccountSuffix) {
		return nil, 0, fmt.Errorf("account [%s] invalid", account)
	}
	index := strings.Index(account, ".")
	if index > -1 {
		account = account[:index]
	}

	var res []string
	segments := graphemes.NewSegmenter([]byte(account))
	for segments.Next() {
		res = append(res, segments.Text())
	}

	if err := segments.Err(); err != nil {
		return res, len(res), fmt.Errorf("segments.Err: %s", err.Error())
	}

	return res, len(res), nil
}

func AccountCharTypeToUint64(accountCharType AccountCharType) uint64 {
	return math.BigPow(2, int64(accountCharType)).Uint64()
}

func Uint64ToAccountCharType(num uint64) map[AccountCharType]struct{} {
	var charMap = make(map[AccountCharType]struct{})
	for i := 0; num > 0; {
		lsb := int(num % 2)
		if lsb == 1 {
			charMap[AccountCharType(i)] = struct{}{}
		}
		num /= 2
		i += 1
	}
	return charMap
}

func ConvertAccountCharsToCharsetNum(accountChars *molecule.AccountChars) uint64 {
	charsetList := ConvertToAccountCharSets(accountChars)
	var charsetMap = make(map[AccountCharType]struct{})
	GetAccountCharType(charsetMap, charsetList)
	var charsetNum uint64
	for c, _ := range charsetMap {
		numTmp := AccountCharTypeToUint64(c)
		charsetNum += numTmp
	}
	return charsetNum
}
