package witness

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type SporeData struct {
	ContentType []byte
	Content     []byte
	ClusterId   []byte
}

const (
	DidCellCellDepsFalgTestnet string = "0x2066676e9c6cc0d7218b5fbbf721258999f91eb7fbfc43a4ae080a45b54efb27"
	DidCellCellDepsFalgMainnet string = "0x2ffaa212ed7e00cf595b42765d5e6b8908b18d444ace76113cb707247033ec99"
	ClusterIdTestnet           string = "0x38ab2c230a9f44b4ed7ebb4f7f15a7c9ecf79b3d723a2caf4a8e1b621f61dd71"
	ClusterIdMainnet           string = "0xcff856f49d7a01d48c6a167b5f1bf974d31c375548eea3cf63145a233929f938"
	DidCellDataLVVersion       uint8  = 1
	DidCellDataLVFlag          uint8  = 0
)

func GetDidCellRecycleCellDeps(net common.DasNetType) *types.CellDep {
	outPoint := types.OutPoint{
		TxHash: types.HexToHash(DidCellCellDepsFalgMainnet),
		Index:  0,
	}
	if net != common.DasNetTypeMainNet {
		outPoint = types.OutPoint{
			TxHash: types.HexToHash(DidCellCellDepsFalgTestnet),
			Index:  0,
		}
	}

	return &types.CellDep{
		OutPoint: &outPoint,
		DepType:  types.DepTypeCode,
	}
}

func GetClusterId(net common.DasNetType) []byte {
	id := ClusterIdMainnet
	if net != common.DasNetTypeMainNet {
		id = ClusterIdTestnet
	}
	return common.Hex2Bytes(id)
}

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
	var didCellDataLV DidCellDataLV
	if err := didCellDataLV.BysToObj(s.Content); err != nil {
		return nil, fmt.Errorf("didCellDataLV.BysToObj err: %s", err.Error())
	}
	return &didCellDataLV, nil
}

func (d *DidCellDataLV) ObjToBys() ([]byte, error) {
	var data []byte
	flagBys := molecule.GoU8ToMoleculeU8(d.Flag)
	data = append(data, flagBys.RawData()...)

	versionBys := molecule.GoU8ToMoleculeU8(d.Version)
	data = append(data, versionBys.RawData()...)

	data = append(data, d.WitnessHash...)

	expireAtBys := molecule.GoU64ToMoleculeU64(d.ExpireAt)
	data = append(data, expireAtBys.RawData()...)

	accountBys := []byte(d.Account)
	data = append(data, accountBys...)

	return data, nil
}
func (d *DidCellDataLV) BysToObj(bys []byte) error {
	if len(bys) < 1+1+20+8 {
		return fmt.Errorf("did cell data len invalid")
	}

	d.Flag, _ = molecule.Bytes2GoU8(bys[0:1])
	d.Version, _ = molecule.Bytes2GoU8(bys[1:2])
	d.WitnessHash = bys[2:22]
	d.ExpireAt, _ = molecule.Bytes2GoU64(bys[22:30])
	d.Account = string(bys[30:])

	return nil
}

type DidCellDataLV struct {
	Flag        uint8
	Version     uint8
	WitnessHash []byte
	ExpireAt    uint64
	Account     string
}

//////////////////////////////////////

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
func (d *DidEntity) HashBys() []byte {
	return d.hash
}

type TxDidEntityWitness struct {
	CellDeps map[uint64]DidEntity
	Inputs   map[uint64]DidEntity
	Outputs  map[uint64]DidEntity
}

func GetDidEntityFromTx(tx *types.Transaction) (TxDidEntityWitness, error) {
	var res TxDidEntityWitness
	res.CellDeps = make(map[uint64]DidEntity)
	res.Inputs = make(map[uint64]DidEntity)
	res.Outputs = make(map[uint64]DidEntity)

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
			res.CellDeps[d.Target.Index] = d
		case SourceTypeInputs:
			res.Inputs[d.Target.Index] = d
		case SourceTypeOutputs:
			res.Outputs[d.Target.Index] = d
		}
	}
	return res, nil
}
