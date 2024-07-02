package witness

//
//type ItemIdDidCellData uint32
//
//const (
//	ItemIdDidCellDataV0 ItemIdDidCellData = 0
//)
//
//type DidCellData struct {
//	ItemId      ItemIdDidCellData `json:"item_id"`
//	Account     string            `json:"account"`
//	ExpireAt    uint64            `json:"expire_at"`
//	WitnessHash string            `json:"witness_hash"`
//}
//
//func (d *DidCellData) AccountId() string {
//	return common.Bytes2Hex(common.GetAccountIdByAccount(d.Account))
//}
//
//func (d *DidCellData) BysToObj(bys []byte) error {
//	didCellData, err := molecule.DidCellDataFromSlice(bys, true)
//	if err != nil {
//		return fmt.Errorf("molecule.DidCellDataFromSlice err: %s", err.Error())
//	}
//
//	didCellDataUnion := didCellData.ToUnion()
//	itemId := didCellDataUnion.ItemID()
//	switch ItemIdDidCellData(itemId) {
//	case ItemIdDidCellDataV0:
//		didCellDataV0 := didCellDataUnion.IntoDidCellDataV0()
//		acc := string(didCellDataV0.Account().RawData())
//		expireAt, err := molecule.Bytes2GoU64(didCellDataV0.ExpireAt().RawData())
//		if err != nil {
//			return fmt.Errorf("molecule.Bytes2GoU64 err: %s", err.Error())
//		}
//		witnessHash := common.Bytes2Hex(didCellDataV0.WitnessHash().RawData())
//
//		d.Account = acc
//		d.ExpireAt = expireAt
//		d.WitnessHash = witnessHash
//	default:
//		return fmt.Errorf("unsupport DidCellDataUnion ItemID[%d]", itemId)
//	}
//
//	return nil
//}
//
//func (d *DidCellData) ObjToBys() ([]byte, error) {
//	switch d.ItemId {
//	case ItemIdDidCellDataV0:
//		accBys := molecule.GoString2MoleculeBytes(d.Account)
//		expireAt := molecule.GoU64ToMoleculeU64(d.ExpireAt)
//		witnessHashBys := common.Hex2Bytes(d.WitnessHash)
//
//		witnessHash, err := molecule.GoBytes2MoleculeByte20(witnessHashBys)
//		if err != nil {
//			return nil, fmt.Errorf("molecule.GoBytes2MoleculeByte20 err: %s", err.Error())
//		}
//
//		didCellDataV0Builder := molecule.NewDidCellDataV0Builder()
//		didCellDataV0 := didCellDataV0Builder.Account(accBys).ExpireAt(expireAt).WitnessHash(witnessHash).Build()
//
//		didCellDataUnion := molecule.DidCellDataUnionFromDidCellDataV0(didCellDataV0)
//		didCellDataBuilder := molecule.NewDidCellDataBuilder()
//		didCellData := didCellDataBuilder.Set(didCellDataUnion).Build()
//
//		return didCellData.AsSlice(), nil
//	default:
//		return nil, fmt.Errorf("unsupport DidCellData ItemID[%d]", d.ItemId)
//	}
//}

// ===================================

//
//func BysToDidCellData(bys []byte) (*SporeData, *DidCellData, error) {
//	var sporeData SporeData
//	if err := sporeData.BysToObj(bys); err != nil {
//		log.Error("sporeData.BysToObj err: %s", err.Error())
//		var didCellData DidCellData
//		if err := didCellData.BysToObj(bys); err != nil {
//			return nil, nil, fmt.Errorf("both SporeData and DidCellData fail")
//		}
//		return nil, &didCellData, nil
//	}
//	return &sporeData, nil, nil
//}
//
//func GetAccountAndExpireFromDidCellData(bys []byte) (string, uint64, error) {
//	account := ""
//	expiredAt := uint64(0)
//
//	sporeData, didCellData, err := BysToDidCellData(bys)
//	if err != nil {
//		return "", 0, fmt.Errorf("BysToDidCellData err: %s", err.Error())
//	} else if sporeData != nil {
//		didCellDataLV, err := sporeData.ContentToDidCellDataLV()
//		if err != nil {
//			return "", 0, fmt.Errorf("ContentToDidCellDataLV err: %s", err.Error())
//		}
//		account = didCellDataLV.Account
//		expiredAt = didCellDataLV.ExpireAt
//	} else if didCellData != nil {
//		account = didCellData.Account
//		expiredAt = didCellData.ExpireAt
//	}
//	return account, expiredAt, nil
//}

// ===================================

//func (d *DidEntity) ToInputsDidEntity(index uint64) DidEntity {
//	inputsDidEntity := DidEntity{
//		Target: CellMeta{
//			Index:  index,
//			Source: SourceTypeInputs,
//		},
//		ItemId:               d.ItemId,
//		DidCellWitnessDataV0: d.DidCellWitnessDataV0,
//	}
//	return inputsDidEntity
//}

// =======================

//func TxToOneDidEntity(tx *types.Transaction, source SourceType) (DidEntity, error) {
//	inputsSize := len(tx.Inputs)
//	witnessesSize := len(tx.Witnesses)
//	for i := inputsSize; i < witnessesSize; i++ {
//		dataBys := tx.Witnesses[i]
//		if string(dataBys[:3]) != common.WitnessDID {
//			continue
//		}
//		var d DidEntity
//		err := d.BysToObj(dataBys)
//		if err != nil {
//			continue
//		}
//
//		if d.Target.Source == source {
//			return d, nil
//		}
//	}
//	return DidEntity{}, fmt.Errorf("not exist did entity")
//}

//func (o *TxDidEntity) GetDidEntity(source SourceType, index uint64) (*DidEntity, error) {
//	var didEntityList []DidEntity
//	switch source {
//	case SourceTypeCellDeps:
//		didEntityList = o.CellDeps
//	case SourceTypeInputs:
//		didEntityList = o.Inputs
//	case SourceTypeOutputs:
//		didEntityList = o.Outputs
//	default:
//		return nil, fmt.Errorf("unsupport source type[%d]", source)
//	}
//	for i, v := range didEntityList {
//		if v.Target.Index == index {
//			return &didEntityList[i], nil
//		}
//	}
//	return nil, fmt.Errorf("not exist did entity in source[%d] index[%d]", source, index)
//}

//
//type TxDidEntity struct {
//	CellDeps []DidEntity
//	Inputs   []DidEntity
//	Outputs  []DidEntity
//}
//
//func TxToDidEntity(tx *types.Transaction) (TxDidEntity, error) {
//	var res TxDidEntity
//	inputsSize := len(tx.Inputs)
//	witnessesSize := len(tx.Witnesses)
//	for i := inputsSize; i < witnessesSize; i++ {
//		dataBys := tx.Witnesses[i]
//		if len(dataBys) < 3 {
//			continue
//		}
//		if string(dataBys[:3]) != common.WitnessDID {
//			continue
//		}
//		var d DidEntity
//		err := d.BysToObj(dataBys)
//		if err != nil {
//			continue
//		}
//
//		switch d.Target.Source {
//		case SourceTypeCellDeps:
//			res.CellDeps = append(res.CellDeps, d)
//		case SourceTypeInputs:
//			res.Inputs = append(res.Inputs, d)
//		case SourceTypeOutputs:
//			res.Outputs = append(res.Outputs, d)
//		}
//	}
//	return res, nil
//}
