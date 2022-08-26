package example

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestProposalCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreMainNet()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0x107a56fdb804a6b160d4a1876d1793ef05d2ce486fb640898a92d0edc2b2da2e"
	if res, err := dc.Client().GetTransaction(context.Background(), types.HexToHash(hash)); err != nil {
		t.Fatal(err)
	} else if res.TxStatus.Status != types.TransactionStatusCommitted {
		t.Fatal(res.TxStatus.Status)
	} else {
		builder, err := witness.ProposalCellDataBuilderFromTx(res.Transaction, common.DataTypeNew)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(common.Bytes2Hex(builder.ProposalCellData.ProposerLock().Args().RawData()))
	}
}
