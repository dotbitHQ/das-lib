package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"testing"
)

func TestDidEntityBuilder(t *testing.T) {
	d := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  1,
			Source: witness.SourceTypeOutputs,
		},
		ItemId: witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: []witness.Record{{
			Key:   "60",
			Type:  "address",
			Label: "eth addr",
			Value: "0x123",
			TTL:   300,
		}}},
	}
	bys, err := d.ObjToBys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d.Hash())
	fmt.Println(common.Bytes2Hex(bys))
	//0x0d71a28582ad26f7668b4f2de015312b11ff4399
	//0x4449448300000010000000660000006f0000000000000052000000080000004a0000000800000042000000180000002300000029000000350000003e00000007000000616464726573730200000036300800000065746820616464720500000030783132332c0100000101000000000000000d71a28582ad26f7668b00000000000000000000
}

func TestDidEntity(t *testing.T) {
	witnessStr := "0x4449448300000010000000660000006f0000000000000052000000080000004a0000000800000042000000180000002300000029000000350000003e00000007000000616464726573730200000036300800000065746820616464720500000030783132332c0100000101000000000000000d71a28582ad26f7668b4f2de015312b11ff4399"
	var d witness.DidEntity

	err := d.BysToObj(common.Hex2Bytes(witnessStr))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(d.Hash(), d.Target.Source, d.Target.Index)

	fmt.Println(d.ItemId, d.DidCellWitnessDataV0)
}

//func TestDidCellDataBuilder(t *testing.T) {
//	acc := "test.bit"
//	expireAt := uint64(1713758999) //uint64(time.Now().Unix())
//
//	// DidEntity witness data hash
//	witnessHash, err := blake2b.Blake160([]byte(acc))
//	if err != nil {
//		t.Fatal(err)
//	}
//	d := witness.DidCellData{
//		ItemId:      witness.ItemIdDidCellDataV0,
//		Account:     acc,
//		ExpireAt:    expireAt,
//		WitnessHash: common.Bytes2Hex(witnessHash),
//	}
//	bys, err := d.ObjToBys()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(common.Bytes2Hex(bys))
//	fmt.Println(d.WitnessHash)
//	//0x000000003800000010000000240000002c000000b28072bd0201e6feeb4c0000000000000000000017e325660000000008000000746573742e626974
//	//0xb28072bd0201e6feeb4cd96a6879d6422f2218cd
//}

//func TestDidCellData(t *testing.T) {
//	var d witness.DidCellData
//	witnessStr := "0x000000003800000010000000240000002c000000b28072bd0201e6feeb4cd96a6879d6422f2218cd17e325660000000008000000746573742e626974"
//	if err := d.BysToObj(common.Hex2Bytes(witnessStr)); err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(toolib.JsonString(&d))
//}
