package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestPreAccountCellDataBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x039f448c0f56892aaf3026f8ef29f3c940d488a04b662fe95ef3a52d3e02af27"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMap, err := witness.PreAccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range builderMap {
			v.PreAccountCellData.InvitedDiscount()
			fmt.Println(k)
			fmt.Println(v.ChannelLock())
			fmt.Println(v.InviterId())
			//fmt.Println(v.PreAccountCellData.InviterLock().AsSlice())
			//d := molecule.ScriptDefault()
			//fmt.Println(d.AsSlice())
			fmt.Println(v.InviterLock())
			//s := molecule.ScriptDefault()
			//fmt.Println(common.Bytes2Hex(s.Args().RawData()))
			//fmt.Println(common.Bytes2Hex(s.CodeHash().RawData()))
		}
	}
}
