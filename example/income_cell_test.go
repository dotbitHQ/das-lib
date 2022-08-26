package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

func TestIncomeCellDataBuilderFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}
	hash := "0xfce3cca3c4392e0e65ee99e738d421b8d3d6a4d690202570123209d00e0bdcbc"
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

func TestParserIncomeCell(t *testing.T) {
	witnessByte := common.Hex2Bytes("0x646173060000001f0200001000000010000000100000000f0200001000000014000000180000000100000001000000f3010000f30100000c00000055000000490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000d0e1f9a79ab9361821cbc3b31fccee094cafacd09e0100001400000071000000e4000000410100005d0000000c00000055000000490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000d0e1f9a79ab9361821cbc3b31fccee094cafacd00050d6dc01000000730000000c0000006b0000005f000000100000003000000031000000326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd012a000000053a8ac9ad3efd980ffaad37aec39ba9455aa8bb76053a8ac9ad3efd980ffaad37aec39ba9455aa8bb7600e40b54020000005d0000000c00000055000000490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000d0e1f9a79ab9361821cbc3b31fccee094cafacd000ca9a3b000000005d0000000c00000055000000490000001000000030000000310000009bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce80114000000efbf497f752ff7a655a8ec6f3c8f3feaaed6e41000ca9a3b00000000")
	b, _ := json.Marshal(witness.ParserWitnessData(witnessByte))
	t.Log(string(b))
}
