package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestGetAccountCharType(t *testing.T) {
	var list []common.AccountCharSet
	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEmoji,
		Char:        "",
	})
	res := common.GetAccountCharType(list)
	fmt.Println(res)

	list = append(list, common.AccountCharSet{
		CharSetName: common.AccountCharTypeEn,
		Char:        "",
	})
	res = common.GetAccountCharType(list)
	fmt.Println(res)
}
