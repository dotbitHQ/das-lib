package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type AccountCellDataBuilder struct {
	Index                 uint32
	Version               uint32
	AccountId             string
	NextAccountId         string
	Account               string
	Status                uint8
	RegisteredAt          uint64
	UpdatedAt             uint64
	LastTransferAccountAt uint64
	LastEditManagerAt     uint64
	LastEditRecordsAt     uint64
	ExpiredAt             uint64
	EnableSubAccount      uint8
	RenewSubAccountPrice  uint64
	Records               []Record
	RecordsHashBys        []byte
	AccountCellDataV1     *molecule.AccountCellDataV1
	AccountCellDataV2     *molecule.AccountCellDataV2
	AccountCellData       *molecule.AccountCellData
	DataEntityOpt         *molecule.DataEntityOpt
	AccountChars          *molecule.AccountChars
}

type AccountCellParam struct {
	OldIndex              uint32
	NewIndex              uint32
	Status                uint8
	Action                string
	AccountId             string
	RegisterAt            uint64
	SubAction             string
	AccountChars          *molecule.AccountChars
	LastEditRecordsAt     int64
	LastEditManagerAt     int64
	LastTransferAccountAt int64
	InitialRecords        *molecule.Records
	Records               []Record
	EnableSubAccount      uint8
	RenewSubAccountPrice  uint64
	IsCustomScript        bool
	IsClearRecords        bool
}

func AccountCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*AccountCellDataBuilder, error) {
	respMap, err := AccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist account cell")
}

func AccountCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*AccountCellDataBuilder, error) {
	var respMap = make(map[string]*AccountCellDataBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeAccountCell:
			var resp AccountCellDataBuilder
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				if err == ErrDataEntityOptIsNil {
					return true, nil
				}
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.DataEntityOpt = dataEntityOpt

			version, err := molecule.Bytes2GoU32(dataEntity.Version().RawData())
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			resp.Version = version

			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err: %s", err.Error())
			}
			resp.Index = index
			if dataType == common.DataTypeNew {
				expiredAt, err := common.GetAccountCellExpiredAtFromOutputData(tx.OutputsData[index])
				if err != nil {
					return false, fmt.Errorf("GetAccountCellExpiredAtFromOutputData err: %s", err.Error())
				}
				resp.ExpiredAt = expiredAt
				nextAccountId, err := common.GetAccountCellNextAccountIdFromOutputData(tx.OutputsData[index])
				if err != nil {
					return false, fmt.Errorf("GetAccountCellNextAccountIdFromOutputData err: %s", err.Error())
				}
				resp.NextAccountId = common.Bytes2Hex(nextAccountId)
			}

			switch version {
			case common.GoDataEntityVersion1:
				if err = resp.ConvertToAccountCellDataV1(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellDataV1 err: %s", err.Error())
				}
			case common.GoDataEntityVersion2:
				if err = resp.ConvertToAccountCellDataV2(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellDataV2 err: %s", err.Error())
				}
			case common.GoDataEntityVersion3:
				if err = resp.ConvertToAccountCellData(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellData err: %s", err.Error())
				}
			default:
				if err = resp.ConvertToAccountCellData(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellData err: %s", err.Error())
				}
			}
			respMap[resp.Account] = &resp
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist account cell")
	}
	return respMap, nil
}

func (a *AccountCellDataBuilder) ConvertToAccountCellDataV1(slice []byte) error {
	accountData, err := molecule.AccountCellDataV1FromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("AccountCellDataV1FromSlice err: %s", err.Error())
	}
	a.AccountCellDataV1 = accountData

	a.AccountChars = accountData.Account()
	a.Account = common.AccountCharsToAccount(accountData.Account())
	a.AccountId = common.Bytes2Hex(accountData.Id().RawData())
	a.Status, _ = molecule.Bytes2GoU8(accountData.Status().RawData())
	a.RegisteredAt, _ = molecule.Bytes2GoU64(accountData.RegisteredAt().RawData())
	a.UpdatedAt, _ = molecule.Bytes2GoU64(accountData.UpdatedAt().RawData())
	a.Records = ConvertToRecords(accountData.Records())
	a.RecordsHashBys = common.Blake2b(accountData.Records().AsSlice())
	return nil
}

func (a *AccountCellDataBuilder) ConvertToAccountCellDataV2(slice []byte) error {
	accountData, err := molecule.AccountCellDataV2FromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("AccountCellDataV2FromSlice err: %s", err.Error())
	}
	a.AccountCellDataV2 = accountData

	a.AccountChars = accountData.Account()
	a.Account = common.AccountCharsToAccount(accountData.Account())
	a.AccountId = common.Bytes2Hex(accountData.Id().RawData())
	a.Status, _ = molecule.Bytes2GoU8(accountData.Status().RawData())
	a.RegisteredAt, _ = molecule.Bytes2GoU64(accountData.RegisteredAt().RawData())
	a.LastTransferAccountAt, _ = molecule.Bytes2GoU64(accountData.LastTransferAccountAt().RawData())
	a.LastEditManagerAt, _ = molecule.Bytes2GoU64(accountData.LastEditManagerAt().RawData())
	a.LastEditRecordsAt, _ = molecule.Bytes2GoU64(accountData.LastEditRecordsAt().RawData())
	a.Records = ConvertToRecords(accountData.Records())
	a.RecordsHashBys = common.Blake2b(accountData.Records().AsSlice())
	return nil
}

func (a *AccountCellDataBuilder) ConvertToAccountCellData(slice []byte) error {
	accountData, err := molecule.AccountCellDataFromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("AccountCellDataFromSlice err: %s", err.Error())
	}
	a.AccountCellData = accountData

	a.AccountChars = accountData.Account()
	a.Account = common.AccountCharsToAccount(accountData.Account())
	a.AccountId = common.Bytes2Hex(accountData.Id().RawData())
	a.Status, _ = molecule.Bytes2GoU8(accountData.Status().RawData())
	a.RegisteredAt, _ = molecule.Bytes2GoU64(accountData.RegisteredAt().RawData())
	a.LastTransferAccountAt, _ = molecule.Bytes2GoU64(accountData.LastTransferAccountAt().RawData())
	a.LastEditManagerAt, _ = molecule.Bytes2GoU64(accountData.LastEditManagerAt().RawData())
	a.LastEditRecordsAt, _ = molecule.Bytes2GoU64(accountData.LastEditRecordsAt().RawData())
	a.Records = ConvertToRecords(accountData.Records())
	a.RecordsHashBys = common.Blake2b(accountData.Records().AsSlice())
	a.EnableSubAccount, _ = molecule.Bytes2GoU8(accountData.EnableSubAccount().RawData())
	a.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(accountData.RenewSubAccountPrice().RawData())
	return nil
}

func AccountIdCellDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*AccountCellDataBuilder, error) {
	respMap, err := AccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*AccountCellDataBuilder)
	for k, v := range respMap {
		k1 := v.AccountId
		retMap[k1] = respMap[k]
	}
	return retMap, nil
}

func (a *AccountCellDataBuilder) getOldDataEntityOpt(p *AccountCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	switch a.Version {
	case common.GoDataEntityVersion2:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellDataV2.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion2).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	case common.GoDataEntityVersion3:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	}
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}

func (a *AccountCellDataBuilder) getNewAccountCellDataBuilder() *molecule.AccountCellDataBuilder {
	var newBuilder molecule.AccountCellDataBuilder
	switch a.Version {
	case common.GoDataEntityVersion2:
		temNewBuilder := molecule.NewAccountCellDataBuilder()
		temNewBuilder.Records(*a.AccountCellDataV2.Records()).Id(*a.AccountCellDataV2.Id()).
			Status(*a.AccountCellDataV2.Status()).Account(*a.AccountCellDataV2.Account()).
			RegisteredAt(*a.AccountCellDataV2.RegisteredAt()).
			LastTransferAccountAt(*a.AccountCellDataV2.LastTransferAccountAt()).
			LastEditRecordsAt(*a.AccountCellDataV2.LastEditRecordsAt()).
			LastEditManagerAt(*a.AccountCellDataV2.LastEditManagerAt()).
			EnableSubAccount(molecule.Uint8Default()).
			RenewSubAccountPrice(molecule.Uint64Default()).
			Build()
		newBuilder = *temNewBuilder
	default:
		newBuilder = a.AccountCellData.AsBuilder()
	}
	return &newBuilder
}

func (a *AccountCellDataBuilder) GenWitness(p *AccountCellParam) ([]byte, []byte, error) {

	switch p.Action {
	case common.DasActionCollectSubAccountProfit:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, nil, nil
	case common.DasActionCreateSubAccount:
		if p.IsCustomScript {
			oldDataEntityOpt := a.getOldDataEntityOpt(p)
			tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, nil, nil
		}
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionRenewAccount, common.DasActionConfigSubAccountCustomScript:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionEditRecords:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		lastEditRecordsAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastEditRecordsAt)).Build()
		newBuilder.LastEditRecordsAt(lastEditRecordsAt)
		if len(p.Records) == 0 {
			newBuilder.Records(molecule.RecordsDefault())
		} else {
			records := ConvertToCellRecords(p.Records)
			newBuilder.Records(*records)
		}
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionEditManager:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		lastEditManagerAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastEditManagerAt)).Build()
		newBuilder.LastEditManagerAt(lastEditManagerAt)
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionTransferAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		newBuilder.Records(molecule.RecordsDefault())
		lastTransferAccountAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastTransferAccountAt)).Build()
		newBuilder.LastTransferAccountAt(lastTransferAccountAt)
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionBuyAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status))
		newBuilder.Records(molecule.RecordsDefault())

		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionStartAccountSale, common.DasActionCancelAccountSale, common.DasActionAcceptOffer,
		common.DasActionLockAccountForCrossChain, common.DasActionForceRecoverAccountStatus:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status))

		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionUnlockAccountForCrossChain:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status))
		if p.IsClearRecords {
			newBuilder.Records(molecule.RecordsDefault())
		}

		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionPropose, common.DasActionDeclareReverseRecord, common.DasActionRedeclareReverseRecord, common.DasActionEditSubAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		tmp := molecule.NewDataBuilder().Dep(*oldDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, nil, nil
	case common.DasActionConfirmProposal:
		if p.SubAction == "exist" {
			oldDataEntityOpt := a.getOldDataEntityOpt(p)

			newBuilder := a.getNewAccountCellDataBuilder()
			newAccountCellData := newBuilder.Build()
			newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
				Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

			tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
		} else if p.SubAction == "new" {
			accountId, err := molecule.AccountIdFromSlice(common.Hex2Bytes(p.AccountId), true)
			if err != nil {
				return nil, nil, fmt.Errorf("AccountIdFromSlice err: %s", err.Error())
			}
			newAccountCellDataBuilder := molecule.NewAccountCellDataBuilder().
				Status(molecule.GoU8ToMoleculeU8(uint8(0))).
				Records(molecule.RecordsDefault()).
				LastTransferAccountAt(molecule.Uint64Default()).
				LastEditRecordsAt(molecule.Uint64Default()).
				LastEditManagerAt(molecule.Uint64Default()).
				RegisteredAt(molecule.GoU64ToMoleculeU64(p.RegisterAt)).
				Id(*accountId).
				Account(*p.AccountChars)
			if p.InitialRecords != nil {
				newAccountCellDataBuilder.Records(*p.InitialRecords)
			}
			newAccountCellData := newAccountCellDataBuilder.Build()
			newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
				Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
			tmp := molecule.NewDataBuilder().New(newDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
		} else {
			return nil, nil, fmt.Errorf("not exist sub action [%s]", p.SubAction)
		}
	case common.DasActionEnableSubAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder().
			EnableSubAccount(molecule.GoU8ToMoleculeU8(p.EnableSubAccount)).
			RenewSubAccountPrice(molecule.GoU64ToMoleculeU64(p.RenewSubAccountPrice))
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionRecycleExpiredAccount:
		if p.SubAction == "previous" {
			oldDataEntityOpt := a.getOldDataEntityOpt(p)
			newBuilder := a.getNewAccountCellDataBuilder()
			newAccountCellData := newBuilder.Build()
			newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
				Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
			newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

			tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
		} else if p.SubAction == "current" {
			oldDataEntityOpt := a.getOldDataEntityOpt(p)

			tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).Build()
			witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
			return witness, nil, nil
		} else {
			return nil, nil, fmt.Errorf("not exist sub action [%s]", p.SubAction)
		}
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}

type Record struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

func ConvertToRecords(records *molecule.Records) []Record {
	var cellRecords []Record
	for index, lenRecords := uint(0), records.Len(); index < lenRecords; index++ {
		record := records.Get(index)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		cellRecords = append(cellRecords, Record{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   ttl,
		})
	}
	return cellRecords
}

func ConvertToCellRecords(cellRecords []Record) *molecule.Records {
	recordsBuilder := molecule.NewRecordsBuilder()
	for _, v := range cellRecords {
		recordBuilder := molecule.NewRecordBuilder().
			RecordKey(molecule.GoString2MoleculeBytes(v.Key)).
			RecordType(molecule.GoString2MoleculeBytes(v.Type)).
			RecordLabel(molecule.GoString2MoleculeBytes(v.Label)).
			RecordValue(molecule.GoString2MoleculeBytes(v.Value)).
			RecordTtl(molecule.GoU32ToMoleculeU32(v.TTL))
		recordsBuilder.Push(recordBuilder.Build())
	}
	records := recordsBuilder.Build()
	return &records
}
