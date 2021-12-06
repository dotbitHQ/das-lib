package example

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestAccountSaleCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x114dcdb52147d5886b4fa62757dff30aa3144800d6b2583018b5c7a793ce61ff"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.AccountSaleCellDataBuilderFromTx(res.Transaction, common.DataTypeOld)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Account())
		fmt.Println(builder.Price())
	}
}

func TestAccountSaleCellGenWitnessData(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x60b009156855325ad296f7a03f27d0274dead1704c7c64fc7dd8ae2527b4f35c"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.AccountSaleCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Account(), builder.Description())
		fmt.Println(builder.Price())
		fmt.Println(builder.StartedAt())
	}

	var accountSale witness.AccountSaleCellDataBuilder
	_, bys, err := accountSale.GenWitness(&witness.AccountSaleCellParam{
		Price:       43200000000,
		Description: "12qw",
		Account:     "tang000002.bit",
		StartAt:     1632882826,
		Action:      common.DasActionStartAccountSale,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(common.Bytes2Hex(bys))
}
