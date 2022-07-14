package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"testing"
)

func TestCapacity2Str(t *testing.T) {
	fmt.Println(common.Capacity2Str(210000000000))

	fmt.Println(common.GetAccountLength("unit.bit"))
}

func TestConvertAddressKey(t *testing.T) {
	fmt.Println(common.ConvertRecordsAddressCoinType("address.0"))
	fmt.Println(common.ConvertRecordsAddressKey("address.btc"))
}
