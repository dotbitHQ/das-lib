package witness

import (
	"bytes"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type AccountApprovalAction string

const (
	AccountApprovalActionTransfer AccountApprovalAction = "transfer"
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
	AccountApproval       AccountApproval
	Records               []Record
	RecordsHashBys        []byte
	AccountCellDataV1     *molecule.AccountCellDataV1
	AccountCellDataV2     *molecule.AccountCellDataV2
	AccountCellDataV3     *molecule.AccountCellDataV3
	AccountCellDataV4     *molecule.AccountCellDataV4
	AccountCellData       *molecule.AccountCellData
	DataEntityOpt         *molecule.DataEntityOpt
	AccountChars          *molecule.AccountChars
	RefundLock            *molecule.Script
}

type AccountApproval struct {
	Action AccountApprovalAction `json:"action"`
	Params AccountApprovalParams `json:"params"`
}

type AccountApprovalParams struct {
	Transfer AccountApprovalParamsTransfer `json:"transfer"`
}

type AccountApprovalParamsTransfer struct {
	PlatformLock     *types.Script `json:"platform_lock"`
	ProtectedUntil   uint64        `json:"protected_until"`
	SealedUntil      uint64        `json:"sealed_until"`
	DelayCountRemain uint8         `json:"delay_count_remain"`
	ToLock           *types.Script `json:"to_lock"`
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
	AccountApproval       AccountApproval
	IsUpgradeDidCell      bool
	RefundScript          *molecule.Script
}

func (a *AccountCellDataBuilder) GetRefundLock() *types.Script {
	if a.RefundLock == nil {
		return nil
	}
	refundLock := molecule.MoleculeScript2CkbScript(a.RefundLock)
	if len(refundLock.Args) > 0 {
		return refundLock
	}
	return nil
}

func AccountApprovalFromSlice(bs []byte) (*AccountApproval, error) {
	res := &AccountApproval{}
	defaultApproval := molecule.AccountApprovalDefault()
	if len(bs) == 0 || bytes.Compare(bs, defaultApproval.AsSlice()) == 0 {
		return res, nil
	}

	accountApproval, err := molecule.AccountApprovalFromSlice(bs, true)
	if err != nil {
		return nil, err
	}
	action := AccountApprovalAction(accountApproval.Action().RawData())
	res.Action = action

	switch action {
	case AccountApprovalActionTransfer:
		accountApprovalTransfer, err := molecule.AccountApprovalTransferFromSlice(accountApproval.Params().RawData(), true)
		if err != nil {
			return nil, err
		}

		var platformHashType types.ScriptHashType
		platformHashTypeBs := accountApprovalTransfer.PlatformLock().HashType().AsSlice()
		switch platformHashTypeBs[0] {
		case 0:
			platformHashType = types.HashTypeData
		case 1:
			platformHashType = types.HashTypeType
		case 2:
			platformHashType = types.HashTypeData1
		}

		var toLockHashType types.ScriptHashType
		toLockHashTypeBs := accountApprovalTransfer.ToLock().HashType().AsSlice()
		switch toLockHashTypeBs[0] {
		case 0:
			toLockHashType = types.HashTypeData
		case 1:
			toLockHashType = types.HashTypeType
		case 2:
			toLockHashType = types.HashTypeData1
		}

		transferParams := AccountApprovalParamsTransfer{
			PlatformLock: &types.Script{
				CodeHash: types.BytesToHash(accountApprovalTransfer.PlatformLock().CodeHash().RawData()),
				HashType: platformHashType,
				Args:     accountApprovalTransfer.PlatformLock().Args().RawData(),
			},
			ToLock: &types.Script{
				CodeHash: types.BytesToHash(accountApprovalTransfer.ToLock().CodeHash().RawData()),
				HashType: toLockHashType,
				Args:     accountApprovalTransfer.ToLock().Args().RawData(),
			},
		}
		transferParams.ProtectedUntil, _ = molecule.Bytes2GoU64(accountApprovalTransfer.ProtectedUntil().RawData())
		transferParams.SealedUntil, _ = molecule.Bytes2GoU64(accountApprovalTransfer.SealedUntil().RawData())
		transferParams.DelayCountRemain, _ = molecule.Bytes2GoU8(accountApprovalTransfer.DelayCountRemain().RawData())
		res.Params.Transfer = transferParams
	default:
		return nil, fmt.Errorf("action: %s no exist", action)
	}
	return res, nil
}

func (approval *AccountApproval) GenToMolecule() (*molecule.AccountApproval, error) {
	if approval.Action == "" {
		approvalDefault := molecule.AccountApprovalDefault()
		return &approvalDefault, nil
	}

	var res molecule.AccountApproval

	builder := molecule.NewAccountApprovalBuilder()
	builder.Action(molecule.GoBytes2MoleculeBytes([]byte(approval.Action)))

	switch approval.Action {
	case AccountApprovalActionTransfer:
		transferBuilder := molecule.NewAccountApprovalTransferBuilder()
		transfer := approval.Params.Transfer

		platformHashBuilder := molecule.NewHashBuilder()
		platformHash := [32]molecule.Byte{}
		for idx, v := range transfer.PlatformLock.CodeHash.Bytes() {
			platformHash[idx] = molecule.NewByte(v)
		}
		platformHashBuilder.Set(platformHash)

		platformHashType, err := transfer.PlatformLock.HashType.Serialize()
		if err != nil {
			return nil, err
		}
		platformLockBuilder := molecule.NewScriptBuilder()
		platformLockBuilder.CodeHash(platformHashBuilder.Build())
		platformLockBuilder.HashType(molecule.NewByte(platformHashType[0]))
		platformLockBuilder.Args(molecule.GoBytes2MoleculeBytes(transfer.PlatformLock.Args))

		toLockHashBuilder := molecule.NewHashBuilder()
		toLockHash := [32]molecule.Byte{}
		for idx, v := range transfer.PlatformLock.CodeHash.Bytes() {
			toLockHash[idx] = molecule.NewByte(v)
		}
		toLockHashBuilder.Set(toLockHash)
		toLockHashType, err := transfer.ToLock.HashType.Serialize()
		if err != nil {
			return nil, err
		}
		toLockBuilder := molecule.NewScriptBuilder()
		toLockBuilder.CodeHash(toLockHashBuilder.Build())
		toLockBuilder.HashType(molecule.NewByte(toLockHashType[0]))
		toLockBuilder.Args(molecule.GoBytes2MoleculeBytes(transfer.ToLock.Args))

		transferBuilder.PlatformLock(platformLockBuilder.Build())
		transferBuilder.ToLock(toLockBuilder.Build())
		transferBuilder.ProtectedUntil(molecule.GoU64ToMoleculeU64(transfer.ProtectedUntil))
		transferBuilder.SealedUntil(molecule.GoU64ToMoleculeU64(transfer.SealedUntil))
		transferBuilder.DelayCountRemain(molecule.GoU8ToMoleculeU8(transfer.DelayCountRemain))
		accountTransfer := transferBuilder.Build()
		builder.Params(molecule.GoBytes2MoleculeBytes(accountTransfer.AsSlice()))
		res = builder.Build()
	default:
		res = molecule.AccountApprovalDefault()
	}
	return &res, nil
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

func AccountCellDataBuilderFromTxByName(tx *types.Transaction, dataType common.DataType, acc string) (*AccountCellDataBuilder, error) {
	builderMap, err := AccountCellDataBuilderMapFromTx(tx, dataType)
	if err != nil {
		return nil, err
	}
	builder, ok := builderMap[acc]
	if !ok {
		return nil, fmt.Errorf("builderMap not exist account: %s", acc)
	}
	return builder, nil
}

func AccountCellDataBuilderMapFromTx(tx *types.Transaction, dataType common.DataType) (map[string]*AccountCellDataBuilder, error) {
	var respMap = make(map[string]*AccountCellDataBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
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
				if err = resp.ConvertToAccountCellDataV3(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellDataV3 err: %s", err.Error())
				}
			case common.GoDataEntityVersion4:
				if err = resp.ConvertToAccountCellDataV4(dataEntity.Entity().RawData()); err != nil {
					return false, fmt.Errorf("ConvertToAccountCellDataV4 err: %s", err.Error())
				}
			case common.GoDataEntityVersion5:
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

func (a *AccountCellDataBuilder) ConvertToAccountCellDataV3(slice []byte) error {
	accountData, err := molecule.AccountCellDataV3FromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("AccountCellDataFromSlice err: %s", err.Error())
	}
	a.AccountCellDataV3 = accountData

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

func (a *AccountCellDataBuilder) ConvertToAccountCellDataV4(slice []byte) error {
	accountData, err := molecule.AccountCellDataV4FromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("AccountCellDataFromSlice err: %s", err.Error())
	}
	a.AccountCellDataV4 = accountData

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
	accountApproval, err := AccountApprovalFromSlice(accountData.Approval().AsSlice())
	if err != nil {
		return err
	}
	a.AccountApproval = *accountApproval
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
	accountApproval, err := AccountApprovalFromSlice(accountData.Approval().AsSlice())
	if err != nil {
		return err
	}
	a.AccountApproval = *accountApproval
	a.RefundLock = accountData.RefundLock()
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
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellDataV3.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion3).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	case common.GoDataEntityVersion4:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion4).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	default:
		oldAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(a.AccountCellData.AsSlice())
		oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
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
			Approval(molecule.AccountApprovalDefault()).
			Build()
		newBuilder = *temNewBuilder
	case common.GoDataEntityVersion3:
		temNewBuilder := molecule.NewAccountCellDataBuilder()
		temNewBuilder.Records(*a.AccountCellDataV3.Records()).Id(*a.AccountCellDataV3.Id()).
			Status(*a.AccountCellDataV3.Status()).Account(*a.AccountCellDataV3.Account()).
			RegisteredAt(*a.AccountCellDataV3.RegisteredAt()).
			LastTransferAccountAt(*a.AccountCellDataV3.LastTransferAccountAt()).
			LastEditRecordsAt(*a.AccountCellDataV3.LastEditRecordsAt()).
			LastEditManagerAt(*a.AccountCellDataV3.LastEditManagerAt()).
			EnableSubAccount(*a.AccountCellDataV3.EnableSubAccount()).
			RenewSubAccountPrice(*a.AccountCellDataV3.RenewSubAccountPrice()).
			Approval(molecule.AccountApprovalDefault()).
			Build()
		newBuilder = *temNewBuilder
	case common.GoDataEntityVersion4:
		temNewBuilder := molecule.NewAccountCellDataBuilder()
		temNewBuilder.Records(*a.AccountCellDataV4.Records()).Id(*a.AccountCellDataV4.Id()).
			Status(*a.AccountCellDataV4.Status()).Account(*a.AccountCellDataV4.Account()).
			RegisteredAt(*a.AccountCellDataV4.RegisteredAt()).
			LastTransferAccountAt(*a.AccountCellDataV4.LastTransferAccountAt()).
			LastEditRecordsAt(*a.AccountCellDataV4.LastEditRecordsAt()).
			LastEditManagerAt(*a.AccountCellDataV4.LastEditManagerAt()).
			EnableSubAccount(*a.AccountCellDataV4.EnableSubAccount()).
			RenewSubAccountPrice(*a.AccountCellDataV4.RenewSubAccountPrice()).
			Approval(*a.AccountCellDataV4.Approval()).
			//RefundLock(molecule.ScriptDefault()).
			Build()
		newBuilder = *temNewBuilder
	case common.GoDataEntityVersion5:
		newBuilder = a.AccountCellData.AsBuilder()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionRenewAccount, common.DasActionConfigSubAccountCustomScript, common.DasActionConfigSubAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionTransferAccount:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		if p.IsUpgradeDidCell {
			newBuilder.Status(molecule.GoU8ToMoleculeU8(common.AccountStatusOnUpgrade))
		}
		newBuilder.Records(molecule.RecordsDefault())

		lastTransferAccountAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(p.LastTransferAccountAt)).Build()
		newBuilder.LastTransferAccountAt(lastTransferAccountAt)
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionStartAccountSale, common.DasActionCancelAccountSale, common.DasActionAcceptOffer,
		common.DasActionLockAccountForCrossChain, common.DasActionForceRecoverAccountStatus,
		common.DasActionAccountCellUpgrade:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()
		newBuilder.Status(molecule.GoU8ToMoleculeU8(p.Status))

		if p.Action == common.DasActionForceRecoverAccountStatus {
			newBuilder.Approval(molecule.AccountApprovalDefault())
		}
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionPropose, common.DasActionDeclareReverseRecord,
		common.DasActionRedeclareReverseRecord, common.DasActionEditSubAccount, common.DasActionUpdateSubAccount:
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
				Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
				Status(molecule.GoU8ToMoleculeU8(p.Status)).
				Records(molecule.RecordsDefault()).
				LastTransferAccountAt(molecule.Uint64Default()).
				LastEditRecordsAt(molecule.Uint64Default()).
				LastEditManagerAt(molecule.Uint64Default()).
				RegisteredAt(molecule.GoU64ToMoleculeU64(p.RegisterAt)).
				Id(*accountId).
				Account(*p.AccountChars).
				RefundLock(*p.RefundScript)
			if p.InitialRecords != nil {
				newAccountCellDataBuilder.Records(*p.InitialRecords)
			}
			newAccountCellData := newAccountCellDataBuilder.Build()
			newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

			newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
				Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
				Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
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
	case common.DasActionCreateApproval, common.DasActionDelayApproval,
		common.DasActionRevokeApproval, common.DasActionFulfillApproval:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		switch p.Action {
		case common.DasActionCreateApproval:
			newBuilder.Status(molecule.GoU8ToMoleculeU8(common.AccountStatusOnApproval))
			accountApproval, err := p.AccountApproval.GenToMolecule()
			if err != nil {
				return nil, nil, err
			}
			newBuilder.Approval(*accountApproval)
		case common.DasActionDelayApproval:
			a.AccountApproval.Params.Transfer.DelayCountRemain--
			a.AccountApproval.Params.Transfer.SealedUntil = p.AccountApproval.Params.Transfer.SealedUntil
			accountApproval, err := a.AccountApproval.GenToMolecule()
			if err != nil {
				return nil, nil, err
			}
			newBuilder.Approval(*accountApproval)
		case common.DasActionRevokeApproval:
			newBuilder.Status(molecule.GoU8ToMoleculeU8(common.AccountStatusNormal))
			newBuilder.Approval(molecule.AccountApprovalDefault())
		case common.DasActionFulfillApproval:
			newBuilder.Status(molecule.GoU8ToMoleculeU8(common.AccountStatusNormal))
			newBuilder.Records(molecule.RecordsDefault())
			newBuilder.Approval(molecule.AccountApprovalDefault())
		}

		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()

		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
	case common.DasActionBidExpiredAccountAuction:
		oldDataEntityOpt := a.getOldDataEntityOpt(p)
		newBuilder := a.getNewAccountCellDataBuilder()

		//records
		newBuilder.Records(molecule.RecordsDefault())
		//last_edit_records_at

		defaultTime := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(0)).Build()
		newBuilder.LastEditRecordsAt(defaultTime)
		//last_transfer_account_at
		newBuilder.LastTransferAccountAt(defaultTime)
		//last_edit_manager_at
		newBuilder.LastEditManagerAt(defaultTime)

		//registered_at
		registerdAt := molecule.NewUint64Builder().Set(molecule.GoTimeUnixToMoleculeBytes(int64(p.RegisterAt))).Build()
		newBuilder.RegisteredAt(registerdAt)

		//default record
		if len(p.Records) == 0 {
			newBuilder.Records(molecule.RecordsDefault())
		} else {
			records := ConvertToCellRecords(p.Records)
			newBuilder.Records(*records)
		}
		newBuilder.Status(molecule.GoU8ToMoleculeU8(common.AccountStatusNormal))
		newAccountCellData := newBuilder.Build()
		newAccountCellDataBytes := molecule.GoBytes2MoleculeBytes(newAccountCellData.AsSlice())

		newDataEntity := molecule.NewDataEntityBuilder().Entity(newAccountCellDataBytes).
			Version(DataEntityVersion).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeAccountCell, &tmp)
		return witness, common.Blake2b(newAccountCellData.AsSlice()), nil
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
