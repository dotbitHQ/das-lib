package core

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

//var (
//	ErrorNotExistDidCell = errors.New("not exist did cell")
//)

type TxDidCellMap struct {
	Inputs  map[string]DidCellInfo
	Outputs map[string]DidCellInfo
}

type DidCellInfo struct {
	Index       uint64
	OutPoint    *types.OutPoint
	Lock        *types.Script
	OutputsData []byte
}

func (d *DidCellInfo) GetDataInfo() (*witness.SporeData, *witness.DidCellDataLV, error) {
	var sporeData witness.SporeData
	err := sporeData.BysToObj(d.OutputsData)
	if err != nil {
		return nil, nil, fmt.Errorf("sporeData.BysToObj err: %s", err.Error())
	}
	didCellDataLV, err := sporeData.ContentToDidCellDataLV()
	if err != nil {
		return nil, nil, fmt.Errorf("sporeData.ContentToDidCellDataLV err: %s", err.Error())
	}
	return &sporeData, didCellDataLV, nil
}

func (d *DidCellInfo) GetLockAddress(netType common.DasNetType) (string, error) {
	mode := address.Mainnet
	if netType != common.DasNetTypeMainNet {
		mode = address.Testnet
	}
	addr, err := address.ConvertScriptToAddress(mode, d.Lock)
	if err != nil {
		return "", fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
	}
	return addr, nil
}

func (d *DasCore) TxToDidCellEntityAndAction(tx *types.Transaction) (common.DidCellAction, TxDidCellMap, error) {
	var res TxDidCellMap
	res.Inputs = make(map[string]DidCellInfo)
	res.Outputs = make(map[string]DidCellInfo)

	didCellType, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return "", res, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	for i, v := range tx.Outputs {
		if v.Type == nil {
			continue
		}
		if !didCellType.IsSameTypeId(v.Type.CodeHash) {
			continue
		}
		if len(v.Type.Args) == 0 {
			continue
		}
		res.Outputs[common.Bytes2Hex(v.Type.Args)] = DidCellInfo{
			Index: uint64(i),
			OutPoint: &types.OutPoint{
				TxHash: tx.Hash,
				Index:  uint(i),
			},
			Lock:        v.Lock,
			OutputsData: tx.OutputsData[i],
		}
	}

	isDidCellTx := false
	if len(res.Outputs) == 0 {
		for _, v := range tx.CellDeps {
			switch v.OutPoint.TxHash.String() {
			case witness.DidCellCellDepsFalgTestnet,
				witness.DidCellCellDepsFalgMainnet:
				isDidCellTx = true
				break
			}
		}
	} else {
		isDidCellTx = true
	}

	if isDidCellTx {
		for i, v := range tx.Inputs {
			txRes, err := d.client.GetTransaction(context.Background(), v.PreviousOutput.TxHash)
			if err != nil {
				return "", res, fmt.Errorf("GetTransaction err: %s", err.Error())
			}
			cell := txRes.Transaction.Outputs[v.PreviousOutput.Index]
			if cell.Type == nil {
				continue
			}
			if !didCellType.IsSameTypeId(cell.Type.CodeHash) {
				continue
			}
			if len(cell.Type.Args) == 0 {
				continue
			}
			res.Inputs[common.Bytes2Hex(cell.Type.Args)] = DidCellInfo{
				Index:       uint64(i),
				OutPoint:    v.PreviousOutput,
				Lock:        cell.Lock,
				OutputsData: txRes.Transaction.OutputsData[v.PreviousOutput.Index],
			}
		}
	}

	if len(res.Inputs) == 0 && len(res.Outputs) > 0 {
		return common.DidCellActionUpgrade, res, nil
	}
	if len(res.Inputs) > 0 && len(res.Outputs) == 0 {
		return common.DidCellActionRecycle, res, nil
	}
	if len(res.Inputs) > 0 && len(res.Outputs) > 0 {
		return common.DidCellActionUpdate, res, nil
	}

	return "", res, nil
}

//func (d *DasCore) TxToDidCellAction(tx *types.Transaction) (common.DidCellAction, error) {
//	res, err := witness.TxToDidEntity(tx)
//	if err != nil {
//		return "", fmt.Errorf("witness.TxToDidEntity err: %s", err.Error())
//	}
//	if len(res.Inputs) == 0 && len(res.Outputs) == 0 {
//		return "", ErrorNotExistDidCell
//	}
//
//	// recycle
//	if len(res.Inputs) > 0 && len(res.Outputs) == 0 {
//		return common.DidCellActionRecycle, nil
//	}
//
//	// upgrade from account cell
//	if len(res.Inputs) == 0 && len(res.Outputs) > 0 {
//		actionDataBuilder, err := witness.ActionDataBuilderFromTx(tx)
//		if err != nil {
//			return "", fmt.Errorf("witness.ActionDataBuilderFromTx err: %s", err.Error())
//		}
//		switch actionDataBuilder.Action {
//		case common.DasActionTransferAccount:
//			return common.DidCellActionEditOwner, nil
//		case common.DasActionEditRecords:
//			return common.DidCellActionEditRecords, nil
//		case common.DasActionRenewAccount:
//			return common.DidCellActionRenew, nil
//		case common.DasActionAccountCellUpgrade:
//			return common.DidCellActionUpgrade, nil
//		case common.DasActionConfirmProposal:
//			return common.DidCellActionRegister, nil
//		case common.DasActionBidExpiredAccountAuction:
//			return common.DidCellActionAuction, nil
//		default:
//			return "", fmt.Errorf("unsupport das action[%s]", actionDataBuilder.Action)
//		}
//	}
//
//	// not upgrade
//	if len(res.Inputs) != 1 || len(res.Outputs) != 1 {
//		return "", fmt.Errorf("unsupport did cell action")
//	}
//
//	_, expireAt, err := witness.GetAccountAndExpireFromDidCellData(tx.OutputsData[res.Outputs[0].Target.Index])
//	if err != nil {
//		return "", fmt.Errorf("witness.GetAccountAndExpireFromDidCellData err: %s", err.Error())
//	}
//
//	inputsIndex := res.Inputs[0].Target.Index
//	previousOutputHash := tx.Inputs[inputsIndex].PreviousOutput.TxHash
//	previousOutputIndex := tx.Inputs[inputsIndex].PreviousOutput.Index
//
//	previousTx, err := d.client.GetTransaction(d.ctx, previousOutputHash)
//	if err != nil {
//		return "", fmt.Errorf("GetTransaction err: %s", err.Error())
//	}
//	previousOutputsData := previousTx.Transaction.OutputsData[previousOutputIndex]
//
//	_, preExpireAt, err := witness.GetAccountAndExpireFromDidCellData(previousOutputsData)
//	if err != nil {
//		return "", fmt.Errorf("witness.GetAccountAndExpireFromDidCellData pre err: %s", err.Error())
//	}
//
//	if expireAt != preExpireAt {
//		return common.DidCellActionRenew, nil
//	}
//
//	lockHash, err := tx.Outputs[res.Outputs[0].Target.Index].Lock.Hash()
//	if err != nil {
//		return "", fmt.Errorf("lock hash err: %s", err.Error())
//	}
//	previousLockHash, err := previousTx.Transaction.Outputs[previousOutputIndex].Lock.Hash()
//	if err != nil {
//		return "", fmt.Errorf("previous lock hash err: %s", err.Error())
//	}
//	if lockHash.Hex() != previousLockHash.Hex() {
//		return common.DidCellActionEditOwner, nil
//	}
//
//	if res.Outputs[0].Hash() != res.Inputs[0].Hash() {
//		return common.DidCellActionEditRecords, nil
//	}
//
//	return "", ErrorNotExistDidCell
//}

func (d *DasCore) GetDidCellOccupiedCapacity(lock *types.Script, account string) (uint64, error) {
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return 0, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	// type args: blake2b.Blake256(first_input + output_index)
	defaultArgs := make([]byte, 32)

	didCell := types.CellOutput{
		Capacity: 0,
		Lock:     lock,
		Type:     contractDidCell.ToScript(defaultArgs),
	}

	defaultWitnessHash := molecule.Byte20Default()

	didCellDataLV := witness.DidCellDataLV{
		Flag:        witness.DidCellDataLVFlag,
		Version:     witness.DidCellDataLVVersion,
		WitnessHash: defaultWitnessHash.RawData(),
		ExpireAt:    0,
		Account:     account,
	}
	contentBys, err := didCellDataLV.ObjToBys()
	if err != nil {
		return 0, fmt.Errorf("didCellDataLV.ObjToBys() err: %s", err.Error())
	}
	sporeData := witness.SporeData{
		ContentType: []byte{},
		Content:     contentBys,
		ClusterId:   witness.GetClusterId(d.net),
	}
	didCellDataBys, err := sporeData.ObjToBys()
	if err != nil {
		return 0, fmt.Errorf("sporeData.ObjToBys() err: %s", err.Error())
	}

	didCellCapacity := didCell.OccupiedCapacity(didCellDataBys)
	didCellCapacity = didCellCapacity*common.OneCkb + common.OneCkb

	switch lock.CodeHash.Hex() {
	case common.AnyLockCodeHashOfMainnetNoStrLock,
		common.AnyLockCodeHashOfTestnetNoStrLock:
		didCellCapacity += common.OneCkb
	}

	return didCellCapacity, nil
}

type GenDidCellParam struct {
	DidCellLock *types.Script
	Account     string
	ExpireAt    uint64
}

func (d *DasCore) GenDidCellList(input0 *types.CellInput, indexDidCellFrom uint64, didCellParamList []GenDidCellParam) (didCellList []*types.CellOutput, outputsDataList [][]byte, witnessList [][]byte, err error) {
	contractDidCell, err := GetDasContractInfo(common.DasContractNameDidCellType)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	for i, v := range didCellParamList {
		didEntity := witness.DidEntity{
			Target: witness.CellMeta{
				Index:  indexDidCellFrom,
				Source: witness.SourceTypeOutputs,
			},
			ItemId:               witness.ItemIdWitnessDataDidCellV0,
			DidCellWitnessDataV0: &witness.DidCellWitnessDataV0{Records: nil},
		}
		didCellWitness, err := didEntity.ObjToBys()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("didEntity.ObjToBys err: %s", err.Error())
		}
		witnessList = append(witnessList, didCellWitness)
		didCellTypeArgs, err := common.GetDidCellTypeArgs(input0, indexDidCellFrom)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("common.GetDidCellTypeArgs err: %s", err.Error())
		}
		didCell := &types.CellOutput{
			Capacity: 0,
			Lock:     didCellParamList[i].DidCellLock,
			Type:     contractDidCell.ToScript(didCellTypeArgs),
		}

		didCellDataLV := witness.DidCellDataLV{
			Flag:        witness.DidCellDataLVFlag,
			Version:     witness.DidCellDataLVVersion,
			WitnessHash: didEntity.HashBys(),
			ExpireAt:    v.ExpireAt,
			Account:     v.Account,
		}
		contentBys, err := didCellDataLV.ObjToBys()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("contentBys.ObjToBys err: %s", err.Error())
		}
		sporeData := witness.SporeData{
			ContentType: []byte{},
			Content:     contentBys,
			ClusterId:   witness.GetClusterId(d.NetType()),
		}
		didCellDataBys, err := sporeData.ObjToBys()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("sporeData.ObjToBys err: %s", err.Error())
		}

		didCellCapacity, err := d.GetDidCellOccupiedCapacity(didCell.Lock, didCellDataLV.Account)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("GetDidCellOccupiedCapacity err: %s", err.Error())
		}
		didCell.Capacity = didCellCapacity
		didCellList = append(didCellList, didCell)
		outputsDataList = append(outputsDataList, didCellDataBys)

		indexDidCellFrom++
	}
	return
}
