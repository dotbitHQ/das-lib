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
	witnessBytes := common.Hex2Bytes("0x6461730a00000004000000010000000600000075706461746541000000c64feb5409fac9c58403122e0a728e1abda434f1dfa04d96ef21ee064256dc2712131cb4645bfce6ca90c742b6afddff2a452d4dc2640f04b31c59ee13c45de30101000000031400000015a33588908cf8edb27d1abe3852bf287abd38912c0100004c4ff951012907cc3e6c5285f85cd8ea20df33ba648c3568f390a5bd9db74939d050aea591000000000000000000000000000000000000000000000000000000000000000051fa81b516d332984ccdb06460b1547e59ccda8537efa77df9524580bf9295bad4cb22936d358a3903559e45106ece5ce4a5cd9f98b6f178f5fe881ee64dfd575c0150223aff0597b15b8085d77236f687408897e069e42024e3a9c0e5bb4cdba679ce50b321bd30ea49713346fc196cfbb70e4154287f428e84d2b2701c923821e832105091a9cfac6718bd04997222ae4079fc3f2d174a8687b4edd5c955ef4998035554504a4fb68c26f7a64cf63e4161955f5acf02325d2b05b21b3a723e5e875edac1695087bf4d5ae8b67c7837a13ca0afb25ccd4e6b017576b74005d0d2631d1c90c5de040000000c0000000c00000032303234303532352e62697420000000410d788baa315dc86d20404c060093a97943cac4b230a8a3587f7d4f377eb4800c00000032303234303533302e626974")
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
