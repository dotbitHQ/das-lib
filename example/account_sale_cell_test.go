package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/shopspring/decimal"
	"testing"
)

func TestFee(t *testing.T) {
	dec := decimal.NewFromFloat(float64(1000) * 1e8 * 3 / 1e4)
	fmt.Println(dec.String())
	//21000000
}

func TestAccountSale(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x79b120e88f20a7b3f89b623ac89f9c4ffeb39ec60202a097b40173164755ebd2"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		respList, err := witness.IncomeCellDataBuilderListFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("respList:", len(respList))
		for _, v := range respList {
			list := v.Records()
			for _, r := range list {
				fmt.Println(r.Capacity, common.Bytes2Hex(r.BelongTo.Args().RawData()))
			}
		}
		sale, err := witness.AccountSaleCellDataBuilderFromTx(res.Transaction, common.DataTypeOld)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(sale.BuyerInviterProfitRate)
	}
}

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
		fmt.Println(builder.Account)
		fmt.Println(builder.Price)
		fmt.Println(builder.Version)
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
		fmt.Println(builder.Account, builder.Description)
		fmt.Println(builder.Price)
		fmt.Println(builder.StartedAt)
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

func TestOfferCellDataBuilderData(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x4b595c222cf5a2ee1244118ea830ea15613a4cb7fb00e8fba91f8963d092ccad"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else {
		builder, err := witness.OfferCellDataBuilderFromTx(res.Transaction, common.DataTypeOld)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(builder.Version)
	}
}
