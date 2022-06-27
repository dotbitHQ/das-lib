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
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xc552a8430e2d5e81d78836979f2e41507954295faab72c435d051f722dc5ccd5"
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
