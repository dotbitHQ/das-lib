package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"testing"
)

func TestSplitOutputCell2(t *testing.T) {
	total := 122 * common.OneCkb //10000000 * common.OneCkb
	base := 61 * common.OneCkb

	lockScript := common.GetNormalLockScript("0xc866479211cadf63ad115b9da50a6c16bd3d226d")
	list, err := core.SplitOutputCell2(total, base, 1000, lockScript, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range list {
		fmt.Println(v.Capacity / common.OneCkb)
	}
}
