package witness

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ItemIdDidCellData uint32

const (
	ItemIdDidCellDataV0 ItemIdDidCellData = 0
)

type DidCellData struct {
	ItemId      ItemIdDidCellData `json:"item_id"`
	Account     string            `json:"account"`
	ExpireAt    uint64            `json:"expire_at"`
	WitnessHash string            `json:"witness_hash"`
}

func (d *DidCellData) AccountId() string {
	return common.Bytes2Hex(common.GetAccountIdByAccount(d.Account))
}

func (d *DidCellData) BysToObj(bys []byte) error {
	didCellData, err := molecule.DidCellDataFromSlice(bys, true)
	if err != nil {
		return fmt.Errorf("molecule.DidCellDataFromSlice err: %s", err.Error())
	}

	didCellDataUnion := didCellData.ToUnion()
	itemId := didCellDataUnion.ItemID()
	switch ItemIdDidCellData(itemId) {
	case ItemIdDidCellDataV0:
		didCellDataV0 := didCellDataUnion.IntoDidCellDataV0()
		acc := string(didCellDataV0.Account().RawData())
		expireAt, err := molecule.Bytes2GoU64(didCellDataV0.ExpireAt().RawData())
		if err != nil {
			return fmt.Errorf("molecule.Bytes2GoU64 err: %s", err.Error())
		}
		witnessHash := common.Bytes2Hex(didCellDataV0.WitnessHash().RawData())

		d.Account = acc
		d.ExpireAt = expireAt
		d.WitnessHash = witnessHash
	default:
		return fmt.Errorf("unsupport DidCellDataUnion ItemID[%d]", itemId)
	}

	return nil
}

func (d *DidCellData) ObjToBys() ([]byte, error) {
	switch d.ItemId {
	case ItemIdDidCellDataV0:
		accBys := molecule.GoString2MoleculeBytes(d.Account)
		expireAt := molecule.GoU64ToMoleculeU64(d.ExpireAt)
		witnessHashBys := common.Hex2Bytes(d.WitnessHash)

		witnessHash, err := molecule.GoBytes2MoleculeByte20(witnessHashBys)
		if err != nil {
			return nil, fmt.Errorf("molecule.GoBytes2MoleculeByte20 err: %s", err.Error())
		}

		didCellDataV0Builder := molecule.NewDidCellDataV0Builder()
		didCellDataV0 := didCellDataV0Builder.Account(accBys).ExpireAt(expireAt).WitnessHash(witnessHash).Build()

		didCellDataUnion := molecule.DidCellDataUnionFromDidCellDataV0(didCellDataV0)
		didCellDataBuilder := molecule.NewDidCellDataBuilder()
		didCellData := didCellDataBuilder.Set(didCellDataUnion).Build()

		return didCellData.AsSlice(), nil
	default:
		return nil, fmt.Errorf("unsupport DidCellData ItemID[%d]", d.ItemId)
	}
}

// ===================================

type SporeData struct {
	ContentType []byte
	Content     []byte
	ClusterId   []byte
}

const (
	ClusterId string = "0xcdb443dd0f9d98f530fd8945b86f3ea946f56ee4d015882beb757571bbd529f1"
)

func (s *SporeData) ObjToBys() ([]byte, error) {
	sporeDataBuilder := molecule.NewSporeDataBuilder()

	if s.ContentType == nil {
		s.ContentType = make([]byte, 0)
	}
	contentType := molecule.GoBytes2MoleculeBytes(s.ContentType)
	sporeDataBuilder.ContentType(contentType)

	clusterIdBuilder := molecule.NewBytesOptBuilder()
	clusterId := clusterIdBuilder.Set(molecule.GoBytes2MoleculeBytes(s.ClusterId)).Build()
	sporeDataBuilder.ClusterId(clusterId)

	content := molecule.GoBytes2MoleculeBytes(s.Content)
	sporeDataBuilder.Content(content)

	sporeData := sporeDataBuilder.Build()
	return sporeData.AsSlice(), nil
}

func (s *SporeData) BysToObj(bys []byte) error {
	sd, err := molecule.SporeDataFromSlice(bys, true)
	if err != nil {
		return fmt.Errorf("molecule.SporeDataFromSlice err: %s", err.Error())
	}

	s.ContentType = sd.ContentType().RawData()
	s.Content = sd.Content().RawData()
	clusterIdBys, err := sd.ClusterId().IntoBytes()
	if err != nil {
		return fmt.Errorf("ClusterId().IntoBytes err: %s", err.Error())
	}
	s.ClusterId = clusterIdBys.RawData()

	return nil
}

func (s *SporeData) ContentToDidCellDataLV() (*DidCellDataLV, error) {
	var contentBys [][]byte
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)
	for index == uint32(len(s.Content)) {
		dataLen, _ = molecule.Bytes2GoU32(s.Content[index : index+indexLen])
		content := s.Content[index+indexLen : index+indexLen+dataLen]
		index = index + indexLen + dataLen
		contentBys = append(contentBys, content)
	}

	var didCellDataLV DidCellDataLV
	if err := didCellDataLV.BysToObj(s.Content); err != nil {
		return nil, fmt.Errorf("didCellDataLV.BysToObj err: %s", err.Error())
	}
	return &didCellDataLV, nil
}

func (d *DidCellDataLV) ObjToBys() ([]byte, error) {
	var data []byte
	flagBys := molecule.GoU8ToMoleculeU8(d.Flag)
	data = append(data, molecule.GoU32ToBytes(uint32(len(flagBys.RawData())))...)
	data = append(data, flagBys.RawData()...)

	versionBys := molecule.GoU8ToMoleculeU8(d.Version)
	data = append(data, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
	data = append(data, versionBys.RawData()...)

	//witnessHash, err := molecule.GoBytes2MoleculeByte20(d.WitnessHash)
	//if err != nil {
	//	return nil, fmt.Errorf("molecule.GoBytes2MoleculeByte20 err: %s", err.Error())
	//}
	//data = append(data, molecule.GoU32ToBytes(uint32(len(witnessHash.RawData())))...)
	//data = append(data, witnessHash.RawData()...)

	data = append(data, molecule.GoU32ToBytes(uint32(len(d.WitnessHash)))...)
	data = append(data, d.WitnessHash...)

	expireAtBys := molecule.GoU64ToMoleculeU64(d.ExpireAt)
	data = append(data, molecule.GoU32ToBytes(uint32(len(expireAtBys.RawData())))...)
	data = append(data, expireAtBys.RawData()...)

	//accountBys := molecule.GoString2MoleculeBytes(d.Account)
	//data = append(data, molecule.GoU32ToBytes(uint32(len(accountBys.RawData())))...)
	//data = append(data, accountBys.RawData()...)

	accountBys := []byte(d.Account)
	data = append(data, molecule.GoU32ToBytes(uint32(len(accountBys)))...)
	data = append(data, accountBys...)

	return data, nil
}
func (d *DidCellDataLV) BysToObj(bys []byte) error {
	var data [][]byte

	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)
	for index < uint32(len(bys)) {
		dataLen, _ = molecule.Bytes2GoU32(bys[index : index+indexLen])
		dataBys := bys[index+indexLen : index+indexLen+dataLen]
		data = append(data, dataBys)
		index = index + indexLen + dataLen
	}
	if len(data) != 5 {
		return fmt.Errorf("data len err")
	}

	d.Flag, _ = molecule.Bytes2GoU8(data[0])
	d.Version, _ = molecule.Bytes2GoU8(data[1])
	d.WitnessHash = data[2]
	d.ExpireAt, _ = molecule.Bytes2GoU64(data[3])
	d.Account = string(data[4])

	return nil
}

type DidCellDataLV struct {
	Flag        uint8
	Version     uint8
	WitnessHash []byte
	ExpireAt    uint64
	Account     string
}

func BysToDidCellData(bys []byte) (*SporeData, *DidCellData, error) {
	var sporeData SporeData
	if err := sporeData.BysToObj(bys); err != nil {
		log.Error("sporeData.BysToObj err: %s", err.Error())
		var didCellData DidCellData
		if err := didCellData.BysToObj(bys); err != nil {
			return nil, nil, fmt.Errorf("both SporeData and DidCellData fail")
		}
		return nil, &didCellData, nil
	}
	return &sporeData, nil, nil
}

// ===================================

type ItemIdWitnessData uint32
type SourceType byte

const (
	ItemIdWitnessDataDidCellV0 ItemIdWitnessData = 0

	SourceTypeInputs   SourceType = 0
	SourceTypeOutputs  SourceType = 1
	SourceTypeCellDeps SourceType = 2
)

type DidCellWitnessDataV0 struct {
	Records []Record
}
type CellMeta struct {
	Index  uint64     `json:"index"`
	Source SourceType `json:"source"`
}
type DidEntity struct {
	hash        []byte
	witnessData *molecule.WitnessData

	Target               CellMeta              `json:"target"`
	ItemId               ItemIdWitnessData     `json:"item_id"`
	DidCellWitnessDataV0 *DidCellWitnessDataV0 `json:"witness_data_v_0"`
}

var (
	ErrorNotDidEntityWitness = errors.New("not did entity witness")
)

func (d *DidEntity) BysToObj(bys []byte) error {
	if string(bys[:3]) != common.WitnessDID {
		return ErrorNotDidEntityWitness
	}

	didEntity, err := molecule.DidEntityFromSlice(bys[3:], true)
	if err != nil {
		return fmt.Errorf("molecule.DidEntityFromSlice err: %s", err.Error())
	}

	cellMeta, err := didEntity.Target().IntoCellMeta()
	if err != nil {
		return fmt.Errorf("IntoCellMeta err: %s", err.Error())
	}
	index, err := molecule.Bytes2GoU64(cellMeta.Index().RawData())
	if err != nil {
		return fmt.Errorf("molecule.Bytes2GoU64 err: %s", err.Error())
	}

	d.hash = didEntity.Hash().AsSlice()
	d.witnessData = didEntity.Data()

	d.Target.Index = index
	d.Target.Source = SourceType(cellMeta.Source()[0])
	d.ItemId = ItemIdWitnessData(didEntity.Data().ItemID())

	witnessDataUnion := didEntity.Data().ToUnion()
	switch d.ItemId {
	case ItemIdWitnessDataDidCellV0:
		var data DidCellWitnessDataV0
		didCellWitnessDataV0 := witnessDataUnion.IntoDidCellWitnessDataV0()
		data.Records = ConvertToRecords(didCellWitnessDataV0.Records())
		d.DidCellWitnessDataV0 = &data
	default:
		return fmt.Errorf("unsupport WitnessDataUnion ItemID[%d]", d.ItemId)
	}

	return nil
}

func (d *DidEntity) ObjToBys() ([]byte, error) {
	var witnessBys []byte
	switch d.ItemId {
	case ItemIdWitnessDataDidCellV0:
		source := molecule.NewByte(byte(d.Target.Source))
		index := molecule.GoU64ToMoleculeU64(d.Target.Index)
		cellMetaBuilder := molecule.NewCellMetaBuilder()
		cellMeta := cellMetaBuilder.Index(index).Source(source).Build()
		cellMetaOptBuilder := molecule.NewCellMetaOptBuilder()
		cellMetaOpt := cellMetaOptBuilder.Set(cellMeta).Build()

		records := ConvertToCellRecords(d.DidCellWitnessDataV0.Records)
		didCellWitnessDataV0Builder := molecule.NewDidCellWitnessDataV0Builder()
		didCellWitnessDataV0 := didCellWitnessDataV0Builder.Records(*records).Build()

		witnessDataUnion := molecule.WitnessDataUnionFromDidCellWitnessDataV0(didCellWitnessDataV0)
		witnessDataBuilder := molecule.NewWitnessDataBuilder()
		witnessData := witnessDataBuilder.Set(witnessDataUnion).Build()
		d.witnessData = &witnessData

		//
		hash, err := blake2b.Blake160(witnessData.AsSlice())
		if err != nil {
			return nil, fmt.Errorf("blake2b.Blake160 err: %s", err.Error())
		}
		d.hash = hash

		dataHash, err := molecule.GoBytes2MoleculeByte20(hash)
		if err != nil {
			return nil, fmt.Errorf("molecule.GoBytes2MoleculeByte20 err: %s", err.Error())
		}
		byte20OptBuilder := molecule.NewByte20OptBuilder()
		byte20Opt := byte20OptBuilder.Set(dataHash).Build()
		//
		didEntity := molecule.DidEntityDefault()
		didEntityBuilder := didEntity.AsBuilder()
		didEntity = didEntityBuilder.Target(cellMetaOpt).Data(witnessData).Hash(byte20Opt).Build()

		witnessBys = didEntity.AsSlice()
	default:
		return nil, fmt.Errorf("unsupport WitnessData ItemID[%d]", d.ItemId)
	}

	return append([]byte(common.WitnessDID), witnessBys...), nil
}

func (d *DidEntity) Hash() string {
	return common.Bytes2Hex(d.hash)
}

func (d *DidEntity) ToInputsDidEntity(index uint64) DidEntity {
	inputsDidEntity := DidEntity{
		Target: CellMeta{
			Index:  index,
			Source: SourceTypeInputs,
		},
		ItemId:               d.ItemId,
		DidCellWitnessDataV0: d.DidCellWitnessDataV0,
	}
	return inputsDidEntity
}

// =======================

func TxToOneDidEntity(tx *types.Transaction, source SourceType) (DidEntity, error) {
	inputsSize := len(tx.Inputs)
	witnessesSize := len(tx.Witnesses)
	for i := inputsSize; i < witnessesSize; i++ {
		dataBys := tx.Witnesses[i]
		if string(dataBys[:3]) != common.WitnessDID {
			continue
		}
		var d DidEntity
		err := d.BysToObj(dataBys)
		if err != nil {
			continue
		}

		if d.Target.Source == source {
			return d, nil
		}
	}
	return DidEntity{}, fmt.Errorf("not exist did entity")
}

type TxDidEntity struct {
	CellDeps []DidEntity
	Inputs   []DidEntity
	Outputs  []DidEntity
}

func TxToDidEntity(tx *types.Transaction) (TxDidEntity, error) {
	var res TxDidEntity
	inputsSize := len(tx.Inputs)
	witnessesSize := len(tx.Witnesses)
	for i := inputsSize; i < witnessesSize; i++ {
		dataBys := tx.Witnesses[i]
		if len(dataBys) < 3 {
			continue
		}
		if string(dataBys[:3]) != common.WitnessDID {
			continue
		}
		var d DidEntity
		err := d.BysToObj(dataBys)
		if err != nil {
			continue
		}

		switch d.Target.Source {
		case SourceTypeCellDeps:
			res.CellDeps = append(res.CellDeps, d)
		case SourceTypeInputs:
			res.Inputs = append(res.Inputs, d)
		case SourceTypeOutputs:
			res.Outputs = append(res.Outputs, d)
		}
	}
	return res, nil
}

func (o *TxDidEntity) GetDidEntity(source SourceType, index uint64) (*DidEntity, error) {
	var didEntityList []DidEntity
	switch source {
	case SourceTypeCellDeps:
		didEntityList = o.CellDeps
	case SourceTypeInputs:
		didEntityList = o.Inputs
	case SourceTypeOutputs:
		didEntityList = o.Outputs
	default:
		return nil, fmt.Errorf("unsupport source type[%d]", source)
	}
	for i, v := range didEntityList {
		if v.Target.Index == index {
			return &didEntityList[i], nil
		}
	}
	return nil, fmt.Errorf("not exist did entity in source[%d] index[%d]", source, index)
}
