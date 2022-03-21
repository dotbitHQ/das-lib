package common

import (
	"github.com/DeAccountSystems/das-lib/molecule"
	"strings"
)

type AccountCharType uint32

const (
	AccountCharTypeEmoji  AccountCharType = 0
	AccountCharTypeNumber AccountCharType = 1
	AccountCharTypeEn     AccountCharType = 2
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
	if strings.HasSuffix(account, DasAccountSuffix) {
		account = strings.TrimSuffix(account, DasAccountSuffix)
	}

	var list []AccountCharSet
	for _, v := range account {
		char := string(v)
		charSetName := AccountCharTypeEmoji
		if strings.Contains("0123456789", char) {
			charSetName = AccountCharTypeNumber
		} else if strings.Contains("abcdefghijklmnopqrstuvwxyz", char) {
			charSetName = AccountCharTypeEn
		}
		list = append(list, AccountCharSet{
			CharSetName: charSetName,
			Char:        char,
		})
	}
	return list, nil
}
