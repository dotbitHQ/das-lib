package example

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"testing"
	"time"
)

func TestDidEntityBuilder(t *testing.T) {
	source := molecule.NewByte(1)
	index := molecule.GoU64ToMoleculeU64(1)
	cellMetaBuilder := molecule.NewCellMetaBuilder()
	cellMeta := cellMetaBuilder.Index(index).Source(source).Build()
	cellMetaOptBuilder := molecule.NewCellMetaOptBuilder()
	cellMetaOpt := cellMetaOptBuilder.Set(cellMeta).Build()
	//
	recordList := []witness.Record{
		{
			Key:   "60",
			Type:  "address",
			Label: "eth",
			Value: "0x123",
			TTL:   300,
		},
	}
	records := witness.ConvertToCellRecords(recordList)
	didCellWitnessDataV0Builder := molecule.NewDidCellWitnessDataV0Builder()
	didCellWitnessDataV0 := didCellWitnessDataV0Builder.Records(*records).Build()
	witnessDataUnion := molecule.WitnessDataUnionFromDidCellWitnessDataV0(didCellWitnessDataV0)
	witnessDataBuilder := molecule.NewWitnessDataBuilder()
	witnessData := witnessDataBuilder.Set(witnessDataUnion).Build()
	//
	hash, err := blake2b.Blake160(witnessData.AsSlice())
	if err != nil {
		t.Fatal(err)
	}
	dataHash, err := molecule.GoBytes2MoleculeByte20(hash)
	if err != nil {
		t.Fatal(err)
	}
	byte20OptBuilder := molecule.NewByte20OptBuilder()
	byte20Opt := byte20OptBuilder.Set(dataHash).Build()
	//
	didEntity := molecule.DidEntityDefault()
	didEntityBuilder := didEntity.AsBuilder()
	didEntity = didEntityBuilder.Target(cellMetaOpt).Data(witnessData).Hash(byte20Opt).Build()

	witnessBys := didEntity.AsSlice()
	fmt.Println(common.Bytes2Hex(witnessBys))
	// 0x7e00000010000000610000006a000000000000004d0000000800000045000000080000003d00000018000000230000002900000030000000390000000700000061646472657373020000003630030000006574680500000030783132332c01000001010000000000000085f350985d8e6bbc2f5600000000000000000000
}

func TestDidEntity(t *testing.T) {
	witnessStr := "0x7e00000010000000610000006a000000000000004d0000000800000045000000080000003d00000018000000230000002900000030000000390000000700000061646472657373020000003630030000006574680500000030783132332c01000001010000000000000085f350985d8e6bbc2f5600000000000000000000"
	didEntity, err := molecule.DidEntityFromSlice(common.Hex2Bytes(witnessStr), true)
	if err != nil {
		t.Fatal(err)
	}
	cellMeta, err := didEntity.Target().IntoCellMeta()
	if err != nil {
		t.Fatal(err)
	}

	index, err := molecule.Bytes2GoU64(cellMeta.Index().RawData())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("index:", index)
	source := cellMeta.Source()[0]
	fmt.Println(source)

	didCellWitnessDataV0 := didEntity.Data().ToUnion().IntoDidCellWitnessDataV0()
	recordList := witness.ConvertToRecords(didCellWitnessDataV0.Records())
	fmt.Println(recordList)
}

func TestDidCellDataBuilder(t *testing.T) {
	acc := "test.bit"
	accBys := molecule.GoString2MoleculeBytes(acc)
	fmt.Println("acc:", acc, common.Bytes2Hex(accBys.AsSlice()))
	expireAt := uint64(time.Now().Unix())
	expireAtM := molecule.GoU64ToMoleculeU64(expireAt)
	fmt.Println("expireAt:", expireAt, common.Bytes2Hex(expireAtM.AsSlice()))
	hash, err := blake2b.Blake160([]byte(acc))
	if err != nil {
		t.Fatal(err)
	}

	witnessHash, err := molecule.GoBytes2MoleculeByte20(hash)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("witnessHash:", common.Bytes2Hex(witnessHash.RawData()), len(witnessHash.RawData()))

	didCellDataV0 := molecule.DidCellDataV0Default()
	didCellDataV0Builder := didCellDataV0.AsBuilder()
	didCellDataV0 = didCellDataV0Builder.Account(accBys).ExpireAt(expireAtM).WitnessHash(witnessHash).Build()
	didCellDataUnion := molecule.DidCellDataUnionFromDidCellDataV0(didCellDataV0)
	didCellData := molecule.DidCellDataDefault()
	didCellDataBuilder := didCellData.AsBuilder()
	didCellData = didCellDataBuilder.Set(didCellDataUnion).Build()

	witnessBys := didCellData.AsSlice()
	fmt.Println(common.Bytes2Hex(witnessBys))
	// 0x000000003800000010000000240000002c000000
	// b28072bd0201e6feeb4c00000000000000000000
	// c31a226600000000
	// 08000000746573742e626974
}

func TestDidCellData(t *testing.T) {
	witnessStr := "0x000000003800000010000000240000002c000000b28072bd0201e6feeb4c00000000000000000000c31a22660000000008000000746573742e626974"
	didCellData, err := molecule.DidCellDataFromSlice(common.Hex2Bytes(witnessStr), true)
	if err != nil {
		t.Fatal(err)
	}
	didCellDataUnion := didCellData.ToUnion()
	fmt.Println("itmeId:", didCellDataUnion.ItemID())
	didCellDataV0 := didCellDataUnion.IntoDidCellDataV0()
	acc := string(didCellDataV0.Account().RawData())
	expireAt, err := molecule.Bytes2GoU64(didCellDataV0.ExpireAt().RawData())
	if err != nil {
		t.Fatal(err)
	}
	witnessHash := common.Bytes2Hex(didCellDataV0.WitnessHash().RawData())
	fmt.Println(acc, expireAt, witnessHash)
}
