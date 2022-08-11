package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
)

func (d *DasCore) GetAccountCharSetList(account string) ([]common.AccountCharSet, error) {
	var res []common.AccountCharSet

	list, _, err := common.GetDotBitAccountLength(account)
	if err != nil {
		return nil, fmt.Errorf("GetDotBitAccountLength err: %s", err.Error())
	}
	for _, v := range list {
		var tmp common.AccountCharSet
		tmp.Char = v
		if _, ok := common.CharSetTypeEmojiMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeEmoji
		} else if _, ok = common.CharSetTypeDigitMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeDigit
		} else if _, ok = common.CharSetTypeEnMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeEn
		} else if _, ok = common.CharSetTypeHanSMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeHanS
		} else if _, ok = common.CharSetTypeHanTMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeHanT
		} else if _, ok = common.CharSetTypeJaMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeJa
		} else if _, ok = common.CharSetTypeKoMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeKo
		} else if _, ok = common.CharSetTypeRuMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeRu
		} else if _, ok = common.CharSetTypeTrMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeTr
		} else if _, ok = common.CharSetTypeThMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeTh
		} else if _, ok = common.CharSetTypeViMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeVi
		} else {
			tmp.CharSetName = 99
		}
		res = append(res, tmp)
	}
	return res, nil
}
