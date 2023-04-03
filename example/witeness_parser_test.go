package example

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/smt"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"strings"
	"testing"
)

func TestParserWitnessData(t *testing.T) {
	witnessBytes := common.Hex2Bytes("0x646173080000000400000002000000040000006564697414000000ffffffffffffffffffffffffffffffffffffffff010000000008000000ffffffffffffffff200000007de0c01714f7437d92fe5043aef2fa5fcc4d31e32e3db351d46a97d46dfffd4f470000004c4ffd51fdf54a73a606a92e4cf109469e472e18bf2b88f3134a45df0d635db3052bf05c8318dcfcebb59f9ed433b39fa14bf18e894ca61b03a58a264c5e09447d3e5ff5004f025801000058010000300000008f000000a300000024010000320100003a0100004201000043010000470100004f010000500100005f000000100000003000000031000000ebd2ca43797df1eae21f5a0d20a09a3851beab063ca06d7b86a1e1e8ef9c7698012a000000020000000000000000000000000000000000001111020000000000000000000000000000000000001111b5e8f063fe55e67fe87310ddde7d05fea4dbe28281000000180000002d00000042000000570000006c000000150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c00000010000000010000000100000030150000000c000000100000000100000001000000300a0000002e78787878782e6269745af6086000000000ffffffffffffffff000400000000000000000000000000000000000000000a000000657870697265645f617408000000ffffffffffffffff")
	b, err := json.Marshal(witness.ParserWitnessData(witnessBytes))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestParseFromTx(t *testing.T) {
	dc, err := getNewDasCoreTestnet2()
	if err != nil {
		t.Fatal(err)
	}

	str := ``

	list := strings.Split(str, "\n")
	mapR := make(map[string]string)

	tree := smt.NewSparseMerkleTree(nil)

	for _, v := range list {
		outpoint := common.String2OutPointStruct(v)
		res, err := dc.Client().GetTransaction(context.Background(), outpoint.TxHash)
		if err != nil {
			t.Fatal(err)
		}
		var rList []*witness.ReverseSmtRecord
		if err := witness.ParseFromTx(res.Transaction, common.ActionDataTypeReverseSmt, &rList); err != nil {
			t.Fatal(err)
		}
		fmt.Println("rList:", len(rList))
		for _, r := range rList {
			key := fmt.Sprintf("%d-%s", r.SignType, common.Bytes2Hex(r.Address))
			mapR[key] = v
			fmt.Println(key)
			//
			k, _ := blake2b.Blake256(r.Address)
			valBs := make([]byte, 0)
			nonce := molecule.GoU32ToMoleculeU32(r.PrevNonce + 1)
			valBs = append(valBs, nonce.RawData()...)
			valBs = append(valBs, []byte(r.NextAccount)...)
			fmt.Println("valBs:", common.Bytes2Hex(valBs), r.NextAccount)

			value, _ := blake2b.Blake256(valBs)

			_ = tree.Update(k, value)
		}
	}
	root, err := tree.Root()
	fmt.Println("root:", common.Bytes2Hex(root))
	for k, v := range mapR {
		fmt.Println(k, v)
	}
}
