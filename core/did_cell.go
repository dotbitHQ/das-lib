package core

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/txbuilder"
	"github.com/dotbitHQ/das-lib/witness"
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
		} else if p.AccountCellOutPoint.TxHash.String() != "" {
			// account cell -> did cell
		} else {
			return nil, fmt.Errorf("DidCellOutPoint and AccountCellOutPoint nil")
		}
	case common.DidCellActionRenew:
		if p.DidCellOutPoint.TxHash.String() == "" {
			// renew by account cell + balance cell
			return nil, fmt.Errorf("DidCellOutPoint is nil")
		}
		// renew by account cell + did cell + balance cell
	default:
		return nil, fmt.Errorf("unsupport did cell action[%s]", p.Action)
	}
	return nil, fmt.Errorf("invalid params")
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

	// witness
	txDidEntity, err := witness.TxToDidEntity(didCellTx.Transaction)
	if err != nil {
		return nil, fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
	}
	didEntity, err := txDidEntity.GetDidEntity(witness.SourceTypeOutputs, uint64(p.DidCellOutPoint.Index))
	if err != nil {
		return nil, fmt.Errorf("txDidEntity.GetDidEntity err: %s", err.Error())
	}
	newDidEntity := witness.DidEntity{
		Target: witness.CellMeta{
			Index:  0,
			Source: witness.SourceTypeInputs,
		},
		ItemId:               witness.ItemIdWitnessDataDidCellV0,
		DidCellWitnessDataV0: didEntity.DidCellWitnessDataV0,
	}
	newWitness, err := newDidEntity.ObjToBys()
	if err != nil {
		return nil, fmt.Errorf("newDidEntity.ObjToBys err: %s", err.Error())
	}
	txParams.Witnesses = append(txParams.Witnesses, newWitness)

	// cell deps
	txParams.CellDeps = append(txParams.CellDeps, timeCell.ToCellDep())

	return &txParams, nil
}

func (d *DasCore) BuildDidCellTxForEditRecords(p DidCellTxParams) (*txbuilder.BuildTransactionParams, error) {
	var txParams txbuilder.BuildTransactionParams
	return &txParams, nil
}
