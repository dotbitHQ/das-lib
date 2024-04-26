package core

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
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
