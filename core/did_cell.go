package core

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

var (
	ErrorNotExistDidCell = errors.New("not exist did cell")
)

func (d *DasCore) TxToDidCellAction(tx *types.Transaction) (common.DidCellAction, error) {
	res, err := witness.TxToDidEntity(tx)
	if err != nil {
		return "", fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
	}
	if len(res.Inputs) == 0 && len(res.Outputs) == 0 {
		return "", ErrorNotExistDidCell
	}

	// recycle
	if len(res.Inputs) > 0 && len(res.Outputs) == 0 {
		return common.DidCellActionRecycle, nil
	}

	// upgrade from account cell
	if len(res.Inputs) == 0 && len(res.Outputs) > 0 {
		actionDataBuilder, err := witness.ActionDataBuilderFromTx(tx)
		if err != nil {
			return "", fmt.Errorf("witness.ActionDataBuilderFromTx err: %s", err.Error())
		}
		switch actionDataBuilder.Action {
		case common.DasActionTransferAccount:
			return common.DidCellActionEditOwner, nil
		case common.DasActionEditRecords:
			return common.DidCellActionEditRecords, nil
		case common.DasActionRenewAccount:
			return common.DidCellActionRenew, nil
		case common.DasActionAccountCellUpgrade:
			return common.DidCellActionUpgrade, nil
		case common.DasActionConfirmProposal:
			return common.DidCellActionRegister, nil
		case common.DasActionBidExpiredAccountAuction:
			return common.DidCellActionAuction, nil
		default:
			return "", fmt.Errorf("unsupport das action[%s]", actionDataBuilder.Action)
		}
	}

	// not upgrade
	if len(res.Inputs) != 1 || len(res.Outputs) != 1 {
		return "", fmt.Errorf("unsupport did cell action")
	}

	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(tx.OutputsData[res.Outputs[0].Target.Index]); err != nil {
		return "", fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
	}

	inputsIndex := res.Inputs[0].Target.Index
	previousOutputHash := tx.Inputs[inputsIndex].PreviousOutput.TxHash
	previousOutputIndex := tx.Inputs[inputsIndex].PreviousOutput.Index

	previousTx, err := d.client.GetTransaction(d.ctx, previousOutputHash)
	if err != nil {
		return "", fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	previousOutputsData := previousTx.Transaction.OutputsData[previousOutputIndex]

	var previousDidCellData witness.DidCellData
	if err := previousDidCellData.BysToObj(previousOutputsData); err != nil {
		return "", fmt.Errorf("previousDidCellData.BysToObj err: %s", err.Error())
	}
	if didCellData.ExpireAt != previousDidCellData.ExpireAt {
		return common.DidCellActionRenew, nil
	}

	lockHash, err := tx.Outputs[res.Outputs[0].Target.Index].Lock.Hash()
	if err != nil {
		return "", fmt.Errorf("lock hash err: %s", err.Error())
	}
	previousLockHash, err := previousTx.Transaction.Outputs[previousOutputIndex].Lock.Hash()
	if err != nil {
		return "", fmt.Errorf("previous lock hash err: %s", err.Error())
	}
	if lockHash.Hex() != previousLockHash.Hex() {
		return common.DidCellActionEditOwner, nil
	}

	if res.Outputs[0].Hash() != res.Inputs[0].Hash() {
		return common.DidCellActionEditRecords, nil
	}

	return "", ErrorNotExistDidCell
}

// =========================

type DidCellTxParams struct {
	Action              common.DidCellAction
	DidCellOutPoint     types.OutPoint
	AccountCellOutPoint types.OutPoint

	EditRecords   []witness.Record
	EditOwnerLock *types.Script

	NormalCkbLiveCell []*indexer.LiveCell
}

func (d *DasCore) BuildDidCellTx(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	switch p.Action {
	case common.DidCellActionRecycle:
		// did cell -> nil
		return d.BuildDidCellTxForRecycle(p)
	case common.DidCellActionEditRecords:
		// did cell -> did cell
		return d.BuildDidCellTxForEditRecords(p)
	case common.DidCellActionEditOwner:
		if p.DidCellOutPoint.TxHash.String() != "" {
			// did cell -> did cell
			return d.BuildDidCellTxForEditOwner(p)
		} else if p.AccountCellOutPoint.TxHash.String() != "" {
			// check EditOwnerLock is das lock
			contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
			if err != nil {
				return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
			}
			if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
				// account cell -> account cell
				return d.BuildAccountCellTxForEditOwner(p)
			}
			// account cell -> did cell
			return d.BuildDidCellTxForEditOwnerFromAccountCell(p)
		} else {
			return nil, fmt.Errorf("DidCellOutPoint and AccountCellOutPoint nil")
		}
	case common.DidCellActionRenew:
		if p.DidCellOutPoint.TxHash.String() == "" {
			// todo renew by account cell + balance cell
			return nil, fmt.Errorf("DidCellOutPoint is nil")
		}
		// todo renew by account cell + did cell + balance cell
		return d.BuildDidCellTxForRenew(p)
	case common.DidCellActionUpgrade:
		// account cell -> did cell
		return d.BuildDidCellTxForUpgrade(p)
	default:
		return nil, fmt.Errorf("unsupport did cell action[%s]", p.Action)
	}
}

func (d *DasCore) BuildDidCellTxForRecycle(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.DidCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := d.client.GetTransaction(d.ctx, p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
	}
	if int64(didCellData.ExpireAt+3*30*24*60*60) > timeCell.Timestamp() {
		return nil, fmt.Errorf("this expiration time cannot be recycled")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.DidCellOutPoint,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     didCellOutputs.Lock,
		Type:     nil,
	})

	// outputs data
	txParams.OutputsData = append(txParams.OutputsData, []byte{})

	// witness
	txDidEntity, err := witness.TxToDidEntity(didCellTx.Transaction)
	if err != nil {
		return nil, fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
	}
	didEntity, err := txDidEntity.GetDidEntity(witness.SourceTypeOutputs, uint64(p.DidCellOutPoint.Index))
	if err != nil {
		return nil, fmt.Errorf("txDidEntity.GetDidEntity err: %s", err.Error())
	}
	inputsDidEntity := didEntity.ToInputsDidEntity(0)
	inputsWitness, err := inputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("inputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, inputsWitness)

	// cell deps
	txParams.CellDeps = append(txParams.CellDeps, timeCell.ToCellDep())

	return &txParams, nil
}

func (d *DasCore) BuildDidCellTxForEditRecords(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.DidCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := d.client.GetTransaction(d.ctx, p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
	}
	if int64(didCellData.ExpireAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.DidCellOutPoint,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     didCellOutputs.Lock,
		Type:     didCellOutputs.Type,
	})

	// witness
	txDidEntity, err := witness.TxToDidEntity(didCellTx.Transaction)
	if err != nil {
		return nil, fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
	}
	didEntity, err := txDidEntity.GetDidEntity(witness.SourceTypeOutputs, uint64(p.DidCellOutPoint.Index))
	if err != nil {
		return nil, fmt.Errorf("txDidEntity.GetDidEntity err: %s", err.Error())
	}

	// inputs witness
	inputsDidEntity := didEntity.ToInputsDidEntity(0)
	inputsWitness, err := inputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("inputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, inputsWitness)

	// outputs witness
	outputsDidEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  0,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               didEntity.ItemId,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: p.EditRecords},
	}
	outputsWitness, err := outputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("outputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, outputsWitness)

	// outputs data
	didCellData.WitnessHash = outputsDidEntity.Hash()
	outputsData, err := didCellData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellData.ObjToBys err: %s", err.Error())
	}
	txParams.OutputsData = append(txParams.OutputsData, outputsData)

	// cell deps
	txParams.CellDeps = append(txParams.CellDeps, timeCell.ToCellDep())

	return &txParams, nil
}

func (d *DasCore) BuildDidCellTxForEditOwner(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.DidCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("DidCellOutPoint is nil")
	}
	didCellTx, err := d.client.GetTransaction(d.ctx, p.DidCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	didCellOutputs := didCellTx.Transaction.Outputs[p.DidCellOutPoint.Index]

	// check did cell type
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDidCell.IsSameTypeId(didCellOutputs.Type.CodeHash) {
		return nil, fmt.Errorf("DidCellOutPoint is invalid: %s-%d", p.DidCellOutPoint.TxHash.String(), p.DidCellOutPoint.Index)
	}

	// check expire at
	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	var didCellData witness.DidCellData
	if err := didCellData.BysToObj(didCellTx.Transaction.OutputsData[p.DidCellOutPoint.Index]); err != nil {
		return nil, fmt.Errorf("didCellData.BysToObj err: %s", err.Error())
	}
	if int64(didCellData.ExpireAt) < timeCell.Timestamp() {
		return nil, fmt.Errorf("expired and unavailable")
	}

	// todo check owner lock

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.DidCellOutPoint,
	})

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: didCellOutputs.Capacity,
		Lock:     p.EditOwnerLock,
		Type:     didCellOutputs.Type,
	})

	// witness
	txDidEntity, err := witness.TxToDidEntity(didCellTx.Transaction)
	if err != nil {
		return nil, fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
	}
	didEntity, err := txDidEntity.GetDidEntity(witness.SourceTypeOutputs, uint64(p.DidCellOutPoint.Index))
	if err != nil {
		return nil, fmt.Errorf("txDidEntity.GetDidEntity err: %s", err.Error())
	}

	// inputs witness
	inputsDidEntity := didEntity.ToInputsDidEntity(0)
	inputsWitness, err := inputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("inputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, inputsWitness)

	// outputs witness
	outputsDidEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  0,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               didEntity.ItemId,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: nil},
	}
	outputsWitness, err := outputsDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("outputsDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, outputsWitness)

	// outputs data
	didCellData.WitnessHash = outputsDidEntity.Hash()
	outputsData, err := didCellData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellData.ObjToBys err: %s", err.Error())
	}
	txParams.OutputsData = append(txParams.OutputsData, outputsData)

	// cell deps
	txParams.CellDeps = append(txParams.CellDeps, timeCell.ToCellDep())

	return &txParams, nil
}

func (d *DasCore) BuildDidCellTxForEditOwnerFromAccountCell(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.AccountCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is das lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.AccountCellOutPoint,
	})

	// inputs normal ckb cell
	var capacityTotal uint64
	var changeLock, changeType *types.Script
	for i, v := range p.NormalCkbLiveCell {
		capacityTotal += v.Output.Capacity
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: p.NormalCkbLiveCell[i].OutPoint,
		})
	}

	// witness

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionTransferAccount, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)

	// witness account cell
	accountCellTx, err := d.client.GetTransaction(d.ctx, p.AccountCellOutPoint.TxHash)
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

	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}

	accWitness, accData, err := accountCellBuilder.GenWitness(&witness.AccountCellParam{
		OldIndex:              0,
		NewIndex:              0,
		Action:                common.DasActionTransferAccount,
		LastTransferAccountAt: timeCell.Timestamp(),
		IsUpgradeDidCell:      true,
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

	// outputs did cell
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	didEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  1,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: accountCellBuilder.Records},
	}
	didCellWitness, err := didEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, didCellWitness)

	didCell := types.CellOutput{
		Capacity: 0,
		Lock:     p.EditOwnerLock,
		Type:     contractDidCell.ToScript(nil),
	}
	didCellData := witness.DidCellData{
		ItemId:      witness.ItemIdDidCellDataV0,
		Account:     accountCellBuilder.Account,
		ExpireAt:    accountCellBuilder.ExpiredAt,
		WitnessHash: didEntity.Hash(),
	}
	didCellDataBys, err := didCellData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellData.ObjToBys err: %s", err.Error())
	}

	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
	didCell.Capacity = didCellCapacity
	txParams.Outputs = append(txParams.Outputs, &didCell)
	txParams.OutputsData = append(txParams.OutputsData, didCellDataBys)

	// change
	if change := capacityTotal - didCellCapacity; change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     changeLock,
			Type:     changeType,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	// cell deps
	heightCell, err := d.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}

	configCellAcc, err := GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
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

func (d *DasCore) BuildAccountCellTxForEditOwner(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.AccountCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}

	// check das lock
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if !contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is not das lock")
	}
	ownerHex, managerHex, err := d.Daf().ArgsToHex(p.EditOwnerLock.Args)
	if err != nil {
		return nil, fmt.Errorf("ArgsToHex err: %s", err.Error())
	}
	if ownerHex.AddressHex != managerHex.AddressHex {
		return nil, fmt.Errorf("EditOwnerLock invalid")
	}

	//  check old lock
	accountCellTx, err := d.client.GetTransaction(d.ctx, p.AccountCellOutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	accountCellOutput := accountCellTx.Transaction.Outputs[p.AccountCellOutPoint.Index]
	accountCellOutputsData := accountCellTx.Transaction.OutputsData[p.AccountCellOutPoint.Index]
	oldOwnerHex, _, err := d.Daf().ArgsToHex(accountCellOutput.Lock.Args)
	if err != nil {
		return nil, fmt.Errorf("ArgsToHex err: %s", err.Error())
	}
	if oldOwnerHex.AddressHex == ownerHex.AddressHex {
		return nil, fmt.Errorf("EditOwnerLock same as AccountCellOutPoint lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.AccountCellOutPoint,
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

	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
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
	heightCell, err := d.GetHeightCell()
	if err != nil {
		return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	}

	configCellAcc, err := GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
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

func (d *DasCore) BuildDidCellTxForRenew(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	return &txParams, nil
}

func (d *DasCore) BuildDidCellTxForUpgrade(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams

	// check
	if p.AccountCellOutPoint.TxHash.String() == "" {
		return nil, fmt.Errorf("AccountCellOutPoint is nil")
	}
	contractDispatch, err := GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	if contractDispatch.IsSameTypeId(p.EditOwnerLock.CodeHash) {
		return nil, fmt.Errorf("EditOwnerLock is das lock")
	}

	// inputs
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		Since:          0,
		PreviousOutput: &p.AccountCellOutPoint,
	})

	// inputs normal ckb cell
	var capacityTotal uint64
	var changeLock, changeType *types.Script
	for i, v := range p.NormalCkbLiveCell {
		capacityTotal += v.Output.Capacity
		changeLock = v.Output.Lock
		changeType = v.Output.Type
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			Since:          0,
			PreviousOutput: p.NormalCkbLiveCell[i].OutPoint,
		})
	}

	// witness

	// witness action
	actionWitness, err := witness.GenActionDataWitness(common.DasActionAccountCellUpgrade, nil)
	if err != nil {
		return nil, fmt.Errorf("GenActionDataWitness err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, actionWitness)

	// witness account cell
	accountCellTx, err := d.client.GetTransaction(d.ctx, p.AccountCellOutPoint.TxHash)
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
	txParams.Witnesses = append(txParams.Witnesses, accWitness)

	// outputs
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: accountCellOutput.Capacity,
		Lock:     accountCellOutput.Lock,
		Type:     accountCellOutput.Type,
	})
	accData = append(accData, accountCellOutputsData[32:]...)
	txParams.OutputsData = append(txParams.OutputsData, accData)

	// outputs did cell
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	didEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  1,
			Source: witness.SourceTypeOutputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: accountCellBuilder.Records},
	}
	didCellWitness, err := didEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, didCellWitness)

	didCell := types.CellOutput{
		Capacity: 0,
		Lock:     p.EditOwnerLock,
		Type:     contractDidCell.ToScript(nil),
	}
	didCellData := witness.DidCellData{
		ItemId:      witness.ItemIdDidCellDataV0,
		Account:     accountCellBuilder.Account,
		ExpireAt:    accountCellBuilder.ExpiredAt,
		WitnessHash: didEntity.Hash(),
	}
	didCellDataBys, err := didCellData.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("didCellData.ObjToBys err: %s", err.Error())
	}

	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
	didCell.Capacity = didCellCapacity
	txParams.Outputs = append(txParams.Outputs, &didCell)
	txParams.OutputsData = append(txParams.OutputsData, didCellDataBys)

	// change
	if change := capacityTotal - didCellCapacity; change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     changeLock,
			Type:     changeType,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	// cell deps
	timeCell, err := d.GetTimeCell()
	if err != nil {
		return nil, fmt.Errorf("GetTimeCell err: %s", err.Error())
	}
	//
	//heightCell, err := d.GetHeightCell()
	//if err != nil {
	//	return nil, fmt.Errorf("GetHeightCell err: %s", err.Error())
	//}

	configCellAcc, err := GetDasConfigCellInfo(common.ConfigCellTypeArgsAccount)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}

	txParams.CellDeps = append(txParams.CellDeps,
		//heightCell.ToCellDep(),
		timeCell.ToCellDep(),
		configCellAcc.ToCellDep(),
	)

	return &txParams, nil
}
