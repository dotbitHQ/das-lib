package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"testing"
)

func TestAccountCellRefundLock(t *testing.T) {
	s := molecule.ScriptDefault()
	fmt.Println(s.IsEmpty())

	res := molecule.MoleculeScript2CkbScript(&s)
	fmt.Println(len(res.Args))
	fmt.Println(common.Bytes2Hex(res.Args), res.HashType, res.CodeHash.String())
	temNewBuilder := molecule.NewAccountCellDataBuilder().Build()
	fmt.Println(temNewBuilder.RefundLock().HasExtraFields())
}
