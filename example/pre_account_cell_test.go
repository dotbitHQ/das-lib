package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestPreAccountCellDataBuilderMapFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xaa03df07b0dd48ba8e746b1bf7650ef9bb0f01c00df4cd8c0820d7cc01854207"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builderMap, err := witness.PreAccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range builderMap {
			//v.PreAccountCellData.InvitedDiscount()
			fmt.Println(k)
			fmt.Println(v.OwnerLockArgsStr())
			//fmt.Println(v.ChannelLock())
			//fmt.Println(v.InviterId())
			//fmt.Println(v.PreAccountCellData.InviterLock().AsSlice())
			//d := molecule.ScriptDefault()
			//fmt.Println(d.AsSlice())
			//fmt.Println(v.InviterLock())
			//s := molecule.ScriptDefault()
			//fmt.Println(common.Bytes2Hex(s.Args().RawData()))
			//fmt.Println(common.Bytes2Hex(s.CodeHash().RawData()))
		}
	}
}

func TestAddressFormat(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x170d5812d4fe148a67786263364791eff69c7c3bf280b24972c5c82cbf8dbcef"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		for i := 1; i < len(res.Transaction.Outputs); i++ {
			ownerHex, _, err := dc.Daf().ArgsToHex(res.Transaction.Outputs[i].Lock.Args)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(ownerHex)
		}
		//preMap, err := witness.PreAccountCellDataBuilderMapFromTx(res.Transaction, common.DataTypeOld)
		//if err != nil {
		//	t.Fatal(err)
		//}
		//for _, v := range preMap {
		//	inviterLock, _ := v.InviterLock()
		//	if inviterLock == nil {
		//		tmp := molecule.ScriptDefault()
		//		inviterLock = &tmp
		//	}
		//	inviterHex, _, err := dc.Daf().ScriptToHex(molecule.MoleculeScript2CkbScript(inviterLock))
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	fmt.Println(inviterHex.ChainType, inviterHex.AddressHex, inviterHex.DasAlgorithmId)
		//}
	}

}
