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
