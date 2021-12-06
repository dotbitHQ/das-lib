package example

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"testing"
)

func TestCapacity2Str(t *testing.T) {
	fmt.Println(common.Capacity2Str(210000000000))

	fmt.Println(common.GetAccountLength("unit.bit"))
}
