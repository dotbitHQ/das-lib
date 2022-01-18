package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestIncomeCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x551cad2c0567c5397df134880e224b5b98ddefb21c0f5ba87d5ac08cf88d23f7"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		respList, err := witness.IncomeCellDataBuilderListFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(len(respList))
		for _, v := range respList {
			list := v.Records()
			for _, r := range list {
				fmt.Println(r.Capacity, common.Bytes2Hex(r.BelongTo.Args().RawData()))
			}
		}
	}
}
