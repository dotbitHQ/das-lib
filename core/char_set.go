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

	var indexMap = make(map[int]string)
	for i, v := range list {
		var tmp common.AccountCharSet
		tmp.Char = v
		if _, ok := common.CharSetTypeEmojiMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeEmoji
		} else if _, ok = common.CharSetTypeDigitMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeDigit
		} else if _, ok = common.CharSetTypeHanSMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeHanS
		} else if _, ok = common.CharSetTypeHanTMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeHanT
		} else if _, ok = common.CharSetTypeJaMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeJa
		} else if _, ok = common.CharSetTypeKoMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeKo
		} else if _, ok = common.CharSetTypeThMap[v]; ok {
			tmp.CharSetName = common.AccountCharTypeTh
		} else {
			tmp.CharSetName = 99
			indexMap[i] = v
		}

		res = append(res, tmp)
	}
	if len(indexMap) > 0 {
		var checkRes = map[common.AccountCharType]bool{
			common.AccountCharTypeEn: true,
			common.AccountCharTypeVi: true,
			common.AccountCharTypeRu: true,
			common.AccountCharTypeTr: true,
		}

		for _, v := range indexMap {
			if _, ok := common.CharSetTypeEnMap[v]; !ok {
				checkRes[common.AccountCharTypeEn] = false
			}
			if _, ok := common.CharSetTypeViMap[v]; !ok {
				checkRes[common.AccountCharTypeVi] = false
			}
			if _, ok := common.CharSetTypeRuMap[v]; !ok {
				checkRes[common.AccountCharTypeRu] = false
			}
			if _, ok := common.CharSetTypeTrMap[v]; !ok {
				checkRes[common.AccountCharTypeTr] = false
			}
		}
		resCharSetType := common.AccountCharType(99)
		if checkRes[common.AccountCharTypeEn] {
			resCharSetType = common.AccountCharTypeEn
		} else if checkRes[common.AccountCharTypeVi] {
			resCharSetType = common.AccountCharTypeVi
		} else if checkRes[common.AccountCharTypeRu] {
			resCharSetType = common.AccountCharTypeRu
		} else if checkRes[common.AccountCharTypeTr] {
			resCharSetType = common.AccountCharTypeTr
		}
		for k, _ := range indexMap {
			res[k].CharSetName = resCharSetType
		}

	}
	return res, nil
}
