package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestGetAccountCharType(t *testing.T) {
	var res = make(map[common.AccountCharType]struct{})
	var list []common.AccountCharSet
	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEmoji,
		Char:        "",
	})
	common.GetAccountCharType(res, list)
	fmt.Println(res)

	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEn,
		Char:        "",
	})
	common.GetAccountCharType(res, list)
	fmt.Println(res)
}

func TestGetAccountCharTypeExclude(t *testing.T) {
	var res = make(map[common.AccountCharType]struct{})
	var list = []common.AccountCharSet{
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeEn, Char: "."},
		{CharSetName: common.AccountCharTypeEn, Char: "b"},
		{CharSetName: common.AccountCharTypeEn, Char: "i"},
		{CharSetName: common.AccountCharTypeEn, Char: "t"},
	}
	common.GetAccountCharTypeExclude(res, list)
	fmt.Println(res)

	var subRes = make(map[common.AccountCharType]struct{})
	var subList = []common.AccountCharSet{
		{CharSetName: common.AccountCharTypeEn, Char: "h"},
		{CharSetName: common.AccountCharTypeEn, Char: "e"},
		{CharSetName: common.AccountCharTypeEn, Char: "l"},
		{CharSetName: common.AccountCharTypeEn, Char: "l"},
		{CharSetName: common.AccountCharTypeEn, Char: "o"},
		{CharSetName: common.AccountCharTypeEn, Char: "."},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeJa, Char: "ぁ"},
		{CharSetName: common.AccountCharTypeEn, Char: "."},
		{CharSetName: common.AccountCharTypeEn, Char: "b"},
		{CharSetName: common.AccountCharTypeEn, Char: "i"},
		{CharSetName: common.AccountCharTypeEn, Char: "t"},
	}
	common.GetAccountCharTypeExclude(subRes, subList)
	fmt.Println(subRes)
}
