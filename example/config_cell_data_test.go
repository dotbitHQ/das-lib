package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestConfigCellDataBuilderByTypeArgsList(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPreservedAccount19)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(builder.ConfigCellPreservedAccountMap)
	//builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsMain, common.ConfigCellTypeArgsPrice)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//status, err := builder.Status()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("status:", status)
	//reg, renew, err := builder.AccountPrice(4)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println("reg, renew:", reg, renew)
}
