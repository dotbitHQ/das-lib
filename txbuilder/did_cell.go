package txbuilder

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/sjatsh/uint128"
)

var (
	maxRenewYears = 20
)

type DidCellTxParams struct {
	DasCore             *core.DasCore
	DasCache            *dascache.DasCache
	Action              common.DidCellAction
	DidCellOutPoint     *types.OutPoint
	AccountCellOutPoint *types.OutPoint

	EditRecords   []witness.Record
	EditOwnerLock *types.Script

	RenewYears       int
	NormalCellScript *types.Script
}

func (d DidCellTxParams) GetNormalCellCapacity() uint64 {
	if d.NormalCellScript == nil {
		return common.MinCellOccupiedCkb
	}
	cellOutput := types.CellOutput{
		Capacity: 0,
		Lock:     d.NormalCellScript,
		Type:     nil,
	}
	return cellOutput.OccupiedCapacity(nil) * common.OneCkb
}

func BuildDidCellTx(p DidCellTxParams) (*BuildTransactionParams, error) {
	switch p.Action {
	case common.DidCellActionRecycle:
		// did cell -> nil
		return BuildDidCellTxForRecycle(p)
	case common.DidCellActionEditRecords:
		if p.DidCellOutPoint != nil {
			// did cell -> did cell
			return BuildDidCellTxForEditRecords(p)
		} else if p.AccountCellOutPoint != nil {
			// account cell -> account cell
			return BuildAccountCellTxForEditRecords(p)
		} else {
			return nil, fmt.Errorf("DidCellOutPoint and AccountCellOutPoint nil")
		}
	case common.DidCellActionEditOwner:
		if p.DidCellOutPoint != nil {
			// did cell -> did cell
			return BuildDidCellTxForEditOwner(p)
		} else if p.AccountCellOutPoint != nil {
			// check EditOwnerLock is das lock
			contractDispatch, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
			if err != nil {
				return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
			}
			if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
				// account cell -> account cell
				return BuildAccountCellTxForEditOwner(p)
			}
			// account cell -> did cell
			return BuildDidCellTxForEditOwnerFromAccountCell(p)
		} else {
			return nil, fmt.Errorf("DidCellOutPoint and AccountCellOutPoint nil")
		}
	case common.DidCellActionRenew:
		if p.DidCellOutPoint == nil {
			// renew by account cell + balance cell
			return BuildAccountCellTxForRenew(p)
		}
		// renew by account cell + did cell + balance cell
		return BuildDidCellTxForRenew(p)
	case common.DidCellActionUpgrade:
		// account cell -> did cell
		return BuildDidCellTxForUpgrade(p)
	default:
		return nil, fmt.Errorf("unsupport did cell action[%s]", p.Action)
	}
}

func BuildDidCellTxForRecycle(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.DidCellOutPoint == nil {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var sporeData witness.SporeData
	if err := sporeData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("sporeData.BysToObj err: %s", err.Error())
	}

	didCellDataLV, err := sporeData.ContentToDidCellDataLV()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ContentToDidCellDataLV err: %s", err.Error())
	}
	builderConfigCell, err := p.DasCore.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		fmt.Printf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	expirationGracePeriod, err := builderConfigCell.ExpirationGracePeriod()
	if err != nil {
		fmt.Printf("ExpirationGracePeriod err: %s", err.Error())
		return nil, fmt.Errorf("ExpirationGracePeriod err: %s", err.Error())
	}
	if int64(didCellDataLV.ExpireAt+uint64(expirationGracePeriod)) > timeCell.Timestamp() {
		return nil, fmt.Errorf("this expiration time cannot be recycled")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.DidCellOutPoint,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     didCellOutputs.Lock,
		Type:     nil,
	})

	// outputs data
	txParams.OutputsData = append(txParams.OutputsData, []byte{})

	// cell deps
	joyIDCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}
	omniLockCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameOmniLock)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}

	configCellAcc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	recycleCellDep := witness.GetDidCellRecycleCellDeps(p.DasCore.NetType())
	txParams.CellDeps = append(txParams.CellDeps,
		timeCell.ToCellDep(),
		joyIDCellDep,
		omniLockCellDep,
		configCellAcc.ToCellDep(),
		recycleCellDep,
	)

	return &txParams, nil
}

func BuildDidCellTxForEditRecords(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.DidCellOutPoint == nil {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var sporeData witness.SporeData
	if err := sporeData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("sporeData.BysToObj err: %s", err.Error())
	}
	didCellDataLV, err := sporeData.ContentToDidCellDataLV()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ContentToDidCellDataLV err: %s", err.Error())
	}
	if int64(didCellDataLV.ExpireAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.DidCellOutPoint,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     didCellOutputs.Lock,
		Type:     didCellOutputs.Type,
	})

	// witness
	// outputs witness
	outputsDidEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  0,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: p.EditRecords},
	}
	outputsWitness, err := outputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("outputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, outputsWitness)

	// outputs data
	didCellDataLV.WitnessHash = outputsDidEntity.HashBys()
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellDataLV.ObjToBys err: %s", err.Error())
	}
	sporeData.Content = contentBys
	outputsData, err := sporeData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ObjToBys err: %s", err.Error())
	}
	txParams.OutputsData = append(txParams.OutputsData, outputsData)

	// cell deps
	joyIDCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}
	omniLockCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameOmniLock)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		timeCell.ToCellDep(),
		joyIDCellDep,
		omniLockCellDep,
	)

	return &txParams, nil
}

func BuildAccountCellTxForEditRecords(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}
	if accountCellBuilder.Status != common.AccountStatusNormal {
		return nil, fmt.Errorf("account status is not normal")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionEditRecords, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)

	//
	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:          0,
		NewIndex:          0,
		Action:            common.DasActionEditRecords,
		LastEditRecordsAt: timeCell.Timestamp(),
		Records:           p.EditRecords,
	})
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	txParams.OutputsData = append(txParams.OutputsData, accData)

	// cell deps
	heightCell, err := p.DasCore.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}

	configCellAcc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	configCellRecord, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsRecordNamespace)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		heightCell.ToCellDep(),
		timeCell.ToCellDep(),
		configCellAcc.ToCellDep(),
		configCellRecord.ToCellDep(),
	)

	return &txParams, nil
}

func BuildDidCellTxForEditOwner(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.DidCellOutPoint == nil {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var sporeData witness.SporeData
	if err := sporeData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("sporeData.BysToObj err: %s", err.Error())
	}
	didCellDataLV, err := sporeData.ContentToDidCellDataLV()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ContentToDidCellDataLV err: %s", err.Error())
	}

	if int64(didCellDataLV.ExpireAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.DidCellOutPoint,
	})

	// outputs
	// check for lock capacity
	needCapacity := uint64(0)
	oldCapacity, err := p.DasCore.GetDidCellOccupiedCapacity(didCellOutputs.Lock, didCellDataLV.Account)
	if err != nil {
		return nil, fmt.Errorf("GetDidCellOccupiedCapacity old err: %s", err.Error())
	}
	newCapacity, err := p.DasCore.GetDidCellOccupiedCapacity(p.EditOwnerLock, didCellDataLV.Account)
	if err != nil {
		return nil, fmt.Errorf("GetDidCellOccupiedCapacity new err: %s", err.Error())
	}
	if oldCapacity < newCapacity && didCellOutputs.Capacity < newCapacity {
		if newCapacity-didCellOutputs.Capacity >= common.OneCkb {
			needCapacity = newCapacity - didCellOutputs.Capacity
		}
	}

	didCellCapacity := didCellOutputs.Capacity
	change := uint64(0)
	if needCapacity > 0 {
		log.Info("BuildDidCellTxForEditOwner needCapacity:", needCapacity, oldCapacity, needCapacity, didCellOutputs.Capacity)
		didCellCapacity = newCapacity
		// get needCapacity
		changeNormalCell, normalCkbLiveCell, err := p.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
			DasCache:          p.DasCache,
			LockScript:        p.NormalCellScript,
			CapacityNeed:      needCapacity,
			CapacityForChange: p.GetNormalCellCapacity(),
			SearchOrder:       indexer.SearchOrderDesc,
		})
		if err != nil {
			return nil, fmt.Errorf("GetBalanceCellWithLock err: %s", err.Error())
		}
		change = changeNormalCell
		for _, v := range normalCkbLiveCell {
			txParams.Inputs = append(txParams.Inputs, &types.CellInput{
				Since:          0,
				PreviousOutput: v.OutPoint,
			})
		}
	}
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellCapacity,
		Lock:     p.EditOwnerLock,
		Type:     didCellOutputs.Type,
	})

	// outputs witness
	outputsDidEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  0,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: nil},
	}
	outputsWitness, err := outputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("outputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.LatestWitness = append(txParams.LatestWitness, outputsWitness)

	didCellDataLV.WitnessHash = outputsDidEntity.HashBys()
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellDataLV.ObjToBys err: %s", err.Error())
	}
	sporeData.Content = contentBys
	outputsData, err := sporeData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ObjToBys err: %s", err.Error())
	}

	txParams.OutputsData = append(txParams.OutputsData, outputsData)

	// change
	if change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     p.NormalCellScript,
			Type:     nil,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	// cell deps
	joyIDCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}
	omniLockCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameOmniLock)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		timeCell.ToCellDep(),
		joyIDCellDep,
		omniLockCellDep,
	)

	return &txParams, nil
}

func BuildDidCellTxForEditOwnerFromAccountCell(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	contractDispatch, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is das lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})

	// witness

	// witness account cell
	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}
	if accountCellBuilder.Status != common.AccountStatusNormal {
		return nil, fmt.Errorf("account status is not normal")
	}

	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:              0,
		NewIndex:              0,
		Action:                common.DasActionTransferAccount,
		LastTransferAccountAt: timeCell.Timestamp(),
		IsUpgradeDidCell:      true,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	txParams.OutputsData = append(txParams.OutputsData, accData)

	// outputs did cell
	contractDidCell, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	indexDidCell := uint64(1)
	didEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  indexDidCell,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: nil},
	}
	didCellWitness, err := didEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didEntity.ObjToBys err: %s", err.Error())
	}
	txParams.LatestWitness = append(txParams.LatestWitness, didCellWitness)

	didCellArgs, err := common.GetDidCellTypeArgs(txParams.Inputs[0], indexDidCell)
	if err != nil {
		return nil, fmt.Errorf("common.GetDidCellTypeArgs err: %s", err.Error())
	}
	didCell := types.CellOutput{
		Capacity: 0,
		Lock:     p.EditOwnerLock,
		Type:     contractDidCell.ToScript(didCellArgs),
	}

	didCellDataLV := witness.DidCellDataLV{
		Flag:        witness.DidCellDataLVFlag,
		Version:     witness.DidCellDataLVVersion,
		WitnessHash: didEntity.HashBys(),
		ExpireAt:    accountCellBuilder.ExpiredAt,
		Account:     accountCellBuilder.Account,
	}
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("contentBys.ObjToBys err: %s", err.Error())
	}
	sporeData := witness.SporeData{
		ContentType: []byte{},
		Content:     contentBys,
		ClusterId:   witness.GetClusterId(p.DasCore.NetType()),
	}
	didCellDataBys, err := sporeData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ObjToBys err: %s", err.Error())
	}

	didCellCapacity, err := p.DasCore.GetDidCellOccupiedCapacity(didCell.Lock, didCellDataLV.Account)
	if err != nil {
		return nil, fmt.Errorf("GetDidCellOccupiedCapacity err: %s", err.Error())
	}
	didCell.Capacity = didCellCapacity
	txParams.Outputs = append(txParams.Outputs, &didCell)
	txParams.OutputsData = append(txParams.OutputsData, didCellDataBys)

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionTransferAccount, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)
	// account cell witness
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// change

	change, normalCkbLiveCell, err := p.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
		DasCache:          p.DasCache,
		LockScript:        p.NormalCellScript,
		CapacityNeed:      didCellCapacity,
		CapacityForChange: p.GetNormalCellCapacity(),
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCellWithLock err: %s", err.Error())
	}

	// inputs normal cell
	var changeLock, changeType *types.Script
	for i, v := range normalCkbLiveCell {
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: normalCkbLiveCell[i].OutPoint,
		})
	}

	if change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     changeLock,
			Type:     changeType,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	// cell deps
	heightCell, err := p.DasCore.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}

	configCellAcc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		heightCell.ToCellDep(),
		timeCell.ToCellDep(),
		configCellAcc.ToCellDep(),
	)

	return &txParams, nil
}

func BuildAccountCellTxForEditOwner(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}

	// check das lock
	contractDispatch, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is not das lock")
	}
	ownerHex, managerHex, err := p.DasCore.Daf().ArgsToHex(p.EditOwnerLock.Args)
	if err != nil {
		return nil, fmt.Errorf("ArgsToHex err: %s", err.Error())
	}
	if ownerHex.AddressHex != managerHex.AddressHex {
		return nil, fmt.Errorf("EditOwnerLock invalid")
	}

	//  check old lock
	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	oldOwnerHex, _, err := p.DasCore.Daf().ArgsToHex(accountCellOutput.Lock.Args)
	if err != nil {
		return nil, fmt.Errorf("ArgsToHex err: %s", err.Error())
	}
	if oldOwnerHex.AddressHex == ownerHex.AddressHex {
		return nil, fmt.Errorf("EditOwnerLock same as AccountCellOutPoint lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})

	// witness

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionTransferAccount, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)

	// witness account cell
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}
	if accountCellBuilder.Status != common.AccountStatusNormal {
		return nil, fmt.Errorf("account status is not normal")
	}

	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:              0,
		NewIndex:              0,
		Action:                common.DasActionTransferAccount,
		LastTransferAccountAt: timeCell.Timestamp(),
	})
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     p.EditOwnerLock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	txParams.OutputsData = append(txParams.OutputsData, accData)

	// cell deps
	heightCell, err := p.DasCore.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}

	configCellAcc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		heightCell.ToCellDep(),
		timeCell.ToCellDep(),
		configCellAcc.ToCellDep(),
	)

	return &txParams, nil
}

func BuildAccountCellTxForRenew(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	if p.RenewYears <= 0 || p.RenewYears > maxRenewYears {
		return nil, fmt.Errorf("renew years invalid")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})

	// witness
	priceBuilder, err := p.DasCore.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPrice)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgsList err: %s", err.Error())
	}
	quoteCell, err := p.DasCore.GetQuoteCell()
	if err != nil {
		return nil, fmt.Errorf("GetQuoteCell err: %s", err.Error())
	}
	quote := quoteCell.Quote()

	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	builderConfigCell, err := p.DasCore.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		fmt.Printf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	expirationGracePeriod, err := builderConfigCell.ExpirationGracePeriod()
	if err != nil {
		fmt.Printf("ExpirationGracePeriod err: %s", err.Error())
		return nil, fmt.Errorf("ExpirationGracePeriod err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt+uint64(expirationGracePeriod)) < timeCell.Timestamp() {
		return nil, fmt.Errorf("this expiration time cannot be renewed")
	}

	accountLength := uint8(accountCellBuilder.AccountChars.Len())
	_, renewPrice, err := priceBuilder.AccountPrice(accountLength)
	if err != nil {
		return nil, fmt.Errorf("AccountPrice err: %s", err.Error())
	}
	priceCapacity := uint128.From64(renewPrice).Mul(uint128.From64(common.OneCkb)).Div(uint128.From64(quote)).Big().Uint64()
	priceCapacity = priceCapacity * uint64(p.RenewYears)
	log.Info("BuildAccountCellTxForRenew:", priceCapacity, renewPrice, p.RenewYears, quote)

	newExpiredAt := int64(accountCellBuilder.ExpiredAt) + int64(p.RenewYears)*common.OneYearSec
	byteExpiredAt := molecule.Go64ToBytes(newExpiredAt)

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionRenewAccount, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)

	// witness account cell
	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:              0,
		NewIndex:              0,
		Action:                common.DasActionRenewAccount,
		LastTransferAccountAt: timeCell.Timestamp(),
	})
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	accData1 := accData[:common.ExpireTimeEndIndex-common.ExpireTimeLen]
	accData2 := accData[common.ExpireTimeEndIndex:]
	newAccData := append(accData1, byteExpiredAt...)
	newAccData = append(newAccData, accData2...)
	txParams.OutputsData = append(txParams.OutputsData, newAccData) // change expired_at

	// income cell
	incomeCell, err := GenIncomeCell(p.DasCore, p.NormalCellScript, priceCapacity, 1)
	if err != nil {
		return nil, fmt.Errorf("GenIncomeCell err: %s", err.Error())
	}
	txParams.Outputs = append(txParams.Outputs, incomeCell.Cell)
	txParams.OutputsData = append(txParams.OutputsData, incomeCell.Data)
	txParams.Witnesses = append(txParams.Witnesses, incomeCell.Witness)

	// change
	change, normalCkbLiveCell, err := p.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
		DasCache:          p.DasCache,
		LockScript:        p.NormalCellScript,
		CapacityNeed:      incomeCell.Cell.Capacity,
		CapacityForChange: p.GetNormalCellCapacity(),
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCellWithLock err: %s", err.Error())
	}

	// inputs normal cell
	var changeLock, changeType *types.Script
	for i, v := range normalCkbLiveCell {
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: normalCkbLiveCell[i].OutPoint,
		})
	}

	if change > 0 {
		splitCkb := 2000 * common.OneCkb
		changeList, err := core.SplitOutputCell2(change, splitCkb, 10, changeLock, changeType, indexer.SearchOrderAsc)
		if err != nil {
			return nil, fmt.Errorf("SplitOutputCell2 err: %s", err.Error())
		}
		for i := 0; i < len(changeList); i++ {
			txParams.Outputs = append(txParams.Outputs, changeList[i])
			txParams.OutputsData = append(txParams.OutputsData, []byte{})
		}
	}

	// cell deps
	dasLockContract, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	accContract, err := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	incomeContract, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	heightCell, err := p.DasCore.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}
	accountConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	priceConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsPrice)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	incomeConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsIncome)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	txParams.CellDeps = append(txParams.CellDeps,
		dasLockContract.ToCellDep(),
		accContract.ToCellDep(),
		incomeContract.ToCellDep(),
		timeCell.ToCellDep(),
		heightCell.ToCellDep(),
		quoteCell.ToCellDep(),
		accountConfig.ToCellDep(),
		priceConfig.ToCellDep(),
		incomeConfig.ToCellDep(),
	)

	return &txParams, nil
}

func BuildDidCellTxForRenew(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}

	if p.DidCellOutPoint == nil {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	if p.RenewYears <= 0 || p.RenewYears > maxRenewYears {
		return nil, fmt.Errorf("renew years invalid")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.DidCellOutPoint,
	})

	// witness
	priceBuilder, err := p.DasCore.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsPrice)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgsList err: %s", err.Error())
	}
	quoteCell, err := p.DasCore.GetQuoteCell()
	if err != nil {
		return nil, fmt.Errorf("GetQuoteCell err: %s", err.Error())
	}
	quote := quoteCell.Quote()

	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	builderConfigCell, err := p.DasCore.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsAccount)
	if err != nil {
		fmt.Printf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	expirationGracePeriod, err := builderConfigCell.ExpirationGracePeriod()
	if err != nil {
		fmt.Printf("ExpirationGracePeriod err: %s", err.Error())
		return nil, fmt.Errorf("ExpirationGracePeriod err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt+uint64(expirationGracePeriod)) < timeCell.Timestamp() {
		return nil, fmt.Errorf("this expiration time cannot be renewed")
	}

	accountLength := uint8(accountCellBuilder.AccountChars.Len())
	_, renewPrice, err := priceBuilder.AccountPrice(accountLength)
	if err != nil {
		return nil, fmt.Errorf("AccountPrice err: %s", err.Error())
	}
	priceCapacity := uint128.From64(renewPrice).Mul(uint128.From64(common.OneCkb)).Div(uint128.From64(quote)).Big().Uint64()
	priceCapacity = priceCapacity * uint64(p.RenewYears)
	log.Info("BuildAccountCellTxForRenew:", priceCapacity, renewPrice, p.RenewYears, quote)

	newExpiredAt := int64(accountCellBuilder.ExpiredAt) + int64(p.RenewYears)*common.OneYearSec
	byteExpiredAt := molecule.Go64ToBytes(newExpiredAt)

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionRenewAccount, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}

	// witness account cell
	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:              0,
		NewIndex:              0,
		Action:                common.DasActionRenewAccount,
		LastTransferAccountAt: timeCell.Timestamp(),
	})

	txParams.Witnesses = append(txParams.Witnesses, actionWitness)
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// outputs account cell
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	accData1 := accData[:common.ExpireTimeEndIndex-common.ExpireTimeLen]
	accData2 := accData[common.ExpireTimeEndIndex:]
	newAccData := append(accData1, byteExpiredAt...)
	newAccData = append(newAccData, accData2...)
	txParams.OutputsData = append(txParams.OutputsData, newAccData) // change expired_at

	// outputs did cell
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     didCellOutputs.Lock,
		Type:     didCellOutputs.Type,
	})

	var sporeData witness.SporeData
	if err := sporeData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("sporeData.BysToObj err: %s", err.Error())
	}
	didCellDataLV, err := sporeData.ContentToDidCellDataLV()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ContentToDidCellDataLV err: %s", err.Error())
	}

	didCellDataLV.ExpireAt = uint64(newExpiredAt)
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellDataLV.ObjToBys err: %s", err.Error())
	}
	sporeData.Content = contentBys
	outputsData, err := sporeData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellData.ObjToBys err: %s", err.Error())
	}

	txParams.OutputsData = append(txParams.OutputsData, outputsData)

	// income cell
	incomeCell, err := GenIncomeCell(p.DasCore, p.NormalCellScript, priceCapacity, 2)
	if err != nil {
		return nil, fmt.Errorf("GenIncomeCell err: %s", err.Error())
	}
	txParams.Outputs = append(txParams.Outputs, incomeCell.Cell)
	txParams.OutputsData = append(txParams.OutputsData, incomeCell.Data)
	txParams.Witnesses = append(txParams.Witnesses, incomeCell.Witness)

	// check balance
	// change
	change, normalCkbLiveCell, err := p.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
		DasCache:          p.DasCache,
		LockScript:        p.NormalCellScript,
		CapacityNeed:      incomeCell.Cell.Capacity,
		CapacityForChange: p.GetNormalCellCapacity(),
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCellWithLock err: %s", err.Error())
	}

	// inputs normal cell
	var changeLock, changeType *types.Script
	for i, v := range normalCkbLiveCell {
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: normalCkbLiveCell[i].OutPoint,
		})
	}

	if change > 0 {
		splitCkb := 2000 * common.OneCkb
		changeList, err := core.SplitOutputCell2(change, splitCkb, 10, changeLock, changeType, indexer.SearchOrderAsc)
		if err != nil {
			return nil, fmt.Errorf("SplitOutputCell2 err: %s", err.Error())
		}
		for i := 0; i < len(changeList); i++ {
			txParams.Outputs = append(txParams.Outputs, changeList[i])
			txParams.OutputsData = append(txParams.OutputsData, []byte{})
		}
	}

	// cell deps
	dasLockContract, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	accContract, err := core.GetDasContractInfo(common.DasContractNameAccountCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	incomeContract, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	heightCell, err := p.DasCore.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}
	accountConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	priceConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsPrice)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	incomeConfig, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsIncome)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	joyIDCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameJoyID)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}
	omniLockCellDep, err := p.DasCore.GetAnyLockCellDep(core.AnyLockNameOmniLock)
	if err != nil {
		return nil, fmt.Errorf("GetAnyLockCellDep err: %s", err.Error())
	}
	txParams.CellDeps = append(txParams.CellDeps,
		dasLockContract.ToCellDep(),
		accContract.ToCellDep(),
		incomeContract.ToCellDep(),
		timeCell.ToCellDep(),
		heightCell.ToCellDep(),
		quoteCell.ToCellDep(),
		accountConfig.ToCellDep(),
		priceConfig.ToCellDep(),
		incomeConfig.ToCellDep(),
		joyIDCellDep,
		omniLockCellDep,
	)

	return &txParams, nil
}

func BuildDidCellTxForUpgrade(p DidCellTxParams) (*BuildTransactionParams, error) {
	var txParams BuildTransactionParams

	// check
	if p.AccountCellOutPoint == nil {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	contractDispatch, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is das lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: p.AccountCellOutPoint,
	})

	// witness

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionAccountCellUpgrade, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}

	// witness account cell
	accountCellTx, err := p.DasCore.Client().GetTransaction(context.Background(), p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	accountId := common.Bytes2Hex(accountCellOutputsData[32:52])
	accountCellBuilderMap, err := witness.AccountIdCellDataBuilderFromTx(accountCellTx.Transaction, common.DataTypeNew)
	if err != nil {
		return nil, fmt.Errorf("AccountIdCellDataBuilderFromTx err: %s", err.Error())
	}
	accountCellBuilder, ok := accountCellBuilderMap[accountId]
	if !ok {
		return nil, fmt.Errorf("accountCellBuilderMap not exist accountId: %s", accountId)
	}

	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex: 0,
		NewIndex: 0,
		Action:   common.DasActionAccountCellUpgrade,
		Status:   common.AccountStatusOnUpgrade,
	})

	// check expire at
	timeCell, err := p.DasCore.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	if int64(accountCellBuilder.ExpiredAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	txParams.OutputsData = append(txParams.OutputsData, accData)

	// outputs did cell
	contractDidCell, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	indexDidCell := uint64(1)
	didEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  indexDidCell,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: accountCellBuilder.Records},
	}
	didCellWitness, err := didEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didEntity.ObjToBys err: %s", err.Error())
	}
	txParams.LatestWitness = append(txParams.LatestWitness, didCellWitness)
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	didCellArgs, err := common.GetDidCellTypeArgs(txParams.Inputs[0], indexDidCell)
	if err != nil {
		return nil, fmt.Errorf("common.GetDidCellTypeArgs err: %s", err.Error())
	}
	didCell := types.CellOutput{
		Capacity: 0,
		Lock:     p.EditOwnerLock,
		Type:     contractDidCell.ToScript(didCellArgs),
	}

	didCellDataLV := witness.DidCellDataLV{
		Flag:        witness.DidCellDataLVFlag,
		Version:     witness.DidCellDataLVVersion,
		WitnessHash: didEntity.HashBys(),
		ExpireAt:    accountCellBuilder.ExpiredAt,
		Account:     accountCellBuilder.Account,
	}
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("contentBys.ObjToBys err: %s", err.Error())
	}
	sporeData := witness.SporeData{
		ContentType: []byte{},
		Content:     contentBys,
		ClusterId:   witness.GetClusterId(p.DasCore.NetType()),
	}
	didCellDataBys, err := sporeData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("sporeData.ObjToBys err: %s", err.Error())
	}

	didCellCapacity, err := p.DasCore.GetDidCellOccupiedCapacity(didCell.Lock, didCellDataLV.Account)
	if err != nil {
		return nil, fmt.Errorf("GetDidCellOccupiedCapacity err: %s", err.Error())
	}
	didCell.Capacity = didCellCapacity
	txParams.Outputs = append(txParams.Outputs, &didCell)
	txParams.OutputsData = append(txParams.OutputsData, didCellDataBys)

	// change
	change, normalCkbLiveCell, err := p.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
		DasCache:          p.DasCache,
		LockScript:        p.NormalCellScript,
		CapacityNeed:      didCellCapacity,
		CapacityForChange: p.GetNormalCellCapacity(),
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCellWithLock err: %s", err.Error())
	}

	// inputs normal cell
	var changeLock, changeType *types.Script
	for i, v := range normalCkbLiveCell {
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: normalCkbLiveCell[i].OutPoint,
		})
	}

	if change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     changeLock,
			Type:     changeType,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	// cell deps
	configCellAcc, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		timeCell.ToCellDep(),
		configCellAcc.ToCellDep(),
	)

	return &txParams, nil
}

// ==============================

type OutputsIncomeCell struct {
	Cell    *types.CellOutput
	Data    []byte
	Witness []byte
}

func GenIncomeCell(dc *core.DasCore, serverScript *types.Script, capacity uint64, index uint32) (*OutputsIncomeCell, error) {
	var res OutputsIncomeCell
	builder, err := dc.ConfigCellDataBuilderByTypeArgsList(common.ConfigCellTypeArgsIncome)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgsList err: %s", err.Error())
	}
	incomeCellBaseCapacity, err := builder.IncomeBasicCapacity()
	if err != nil {
		return nil, fmt.Errorf("IncomeBasicCapacity err: %s", err.Error())
	}
	log.Info("IncomeBasicCapacity:", incomeCellBaseCapacity, capacity)

	incomeCellCapacity := capacity
	creator := molecule.ScriptDefault()
	var lockList []*types.Script
	var incomeCapacities []uint64
	if capacity < incomeCellBaseCapacity {
		incomeCellCapacity = incomeCellBaseCapacity
		creator = molecule.CkbScript2MoleculeScript(serverScript)
		lockList = append(lockList, serverScript)
		diff := incomeCellBaseCapacity - capacity
		incomeCapacities = append(incomeCapacities, diff)
	}
	asContract, err := core.GetDasContractInfo(common.DasContractNameAlwaysSuccess)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	incomeContract, err := core.GetDasContractInfo(common.DasContractNameIncomeCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	res.Cell = &types.CellOutput{
		Capacity: incomeCellCapacity,
		Lock:     asContract.ToScript(nil),
		Type:     incomeContract.ToScript(nil),
	}

	dasLock := dc.GetDasLock()
	lockList = append(lockList, dasLock)
	incomeCapacities = append(incomeCapacities, capacity)

	res.Witness, res.Data, err = witness.CreateIncomeCellWitness(&witness.NewIncomeCellParam{
		Creator:     &creator,
		BelongTos:   lockList,
		Capacities:  incomeCapacities,
		OutputIndex: index,
	})
	return &res, nil
}
