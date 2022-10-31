package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"testing"
)

func TestSplitOutputCell2(t *testing.T) {
	total := 10000 * common.OneCkb //10000000 * common.OneCkb
	base := 61 * common.OneCkb

	lockScript := common.GetNormalLockScript("0xc866479211cadf63ad115b9da50a6c16bd3d226d")
	list, err := core.SplitOutputCell2(total, base, 10, lockScript, nil, indexer.SearchOrderDesc)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range list {
		fmt.Println(v.Capacity / common.OneCkb)
	}
}
