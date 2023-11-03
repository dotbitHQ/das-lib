package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ParamGetDpCells struct {
	DasCache           *dascache.DasCache
	LockScript         *types.Script
	AmountNeed         uint64
	CurrentBlockNumber uint64
	SearchOrder        indexer.SearchOrder
}

func (d *DasCore) GetDpCells(p *ParamGetDpCells) ([]*indexer.LiveCell, uint64, uint64, error) {
	if d.client == nil {
		return nil, 0, 0, fmt.Errorf("client is nil")
	}
	if p == nil {
		return nil, 0, 0, fmt.Errorf("param is nil")
	}
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	searchKey := &indexer.SearchKey{
		Script:     p.LockScript,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			OutputDataLenRange: &[2]uint64{12, 13},
		},
	}
	if p.CurrentBlockNumber > 0 {
		searchKey.Filter.BlockRange = &[2]uint64{0, p.CurrentBlockNumber - 20}
	}

	var cells []*indexer.LiveCell
	totalAmount := uint64(0)
	totalCapacity := uint64(0)
	hasCache := false
	lastCursor := ""

	ok := false
	for {
		liveCells, err := d.client.GetCells(context.Background(), searchKey, p.SearchOrder, indexer.SearchLimit, lastCursor)
		if err != nil {
			return nil, 0, 0, err
		}
		log.Info("liveCells:", liveCells.LastCursor, len(liveCells.Objects))
		if len(liveCells.Objects) == 0 || lastCursor == liveCells.LastCursor {
			break
		}
		lastCursor = liveCells.LastCursor

		for _, liveCell := range liveCells.Objects {
			if liveCell.Output.Type != nil && !dpContract.IsSameTypeId(liveCell.Output.Type.CodeHash) {
				continue
			}
			if p.AmountNeed > 0 && p.DasCache != nil && p.DasCache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
				hasCache = true
				continue
			}
			cells = append(cells, liveCell)

			idx := 4
			l, err := molecule.Bytes2GoU32(liveCell.OutputData[:idx])
			if err != nil {
				return nil, 0, 0, err
			}
			amount, err := molecule.Bytes2GoU64(liveCell.OutputData[idx : idx+int(l)])
			if err != nil {
				return nil, 0, 0, err
			}
			idx += int(l)

			totalAmount += amount
			totalCapacity += liveCell.Output.Capacity
			if p.AmountNeed > 0 && totalAmount >= p.AmountNeed {
				ok = true
				break
			}
		}

		if ok {
			break
		}
	}

	if p.AmountNeed > 0 && totalAmount < p.AmountNeed {
		if hasCache {
			return cells, totalAmount, totalCapacity, ErrRejectedOutPoint
		}
		return cells, totalAmount, totalCapacity, ErrInsufficientFunds
	}
	return cells, totalAmount, totalCapacity, nil
}

func (d *DasCore) GetDPointTransferWhitelist() (map[string]*types.Script, error) {
	cell, err := GetDasConfigCellInfo(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	tx, err := d.Client().GetTransaction(d.ctx, cell.OutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	builder, err := witness.ConfigCellDataBuilderByTypeArgs(tx.Transaction, common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return builder.GetDPointTransferWhitelist()
}

func (d *DasCore) GetDPointCapacityRecycleWhitelist() (map[string]*types.Script, error) {
	cell, err := GetDasConfigCellInfo(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	tx, err := d.Client().GetTransaction(d.ctx, cell.OutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("GetTransaction err: %s", err.Error())
	}
	builder, err := witness.ConfigCellDataBuilderByTypeArgs(tx.Transaction, common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return builder.GetDPointCapacityRecycleWhitelist()
}

type ParamSplitDPCell struct {
	FromLock           *types.Script
	ToLock             *types.Script
	DPLiveCell         []*indexer.LiveCell
	DPLiveCellCapacity uint64
	DPTotalAmount      uint64
	DPTransferAmount   uint64
	DPBaseCapacity     uint64
	DPContract         *DasContractInfo
	DPSplitCount       int
	DPSplitAmount      uint64
	NormalCellLock     *types.Script
}

func SplitDPCell(p *ParamSplitDPCell) ([]*types.CellOutput, [][]byte, uint64, error) {
	var outputs []*types.CellOutput
	var outputsData [][]byte
	// transfer
	outputs = append(outputs, &types.CellOutput{
		Capacity: p.DPBaseCapacity,
		Lock:     p.ToLock,
		Type:     p.DPContract.ToScript(nil),
	})
	moleculeData := molecule.GoU64ToMoleculeU64(p.DPTransferAmount)
	outputsData = append(outputsData, moleculeData.RawData())
	// split
	dpBalanceAmount := p.DPTotalAmount - p.DPTransferAmount
	for i := 1; i < p.DPSplitCount; i++ {
		if dpBalanceAmount > p.DPSplitAmount*2 {
			outputs = append(outputs, &types.CellOutput{
				Capacity: p.DPBaseCapacity,
				Lock:     p.FromLock,
				Type:     p.DPContract.ToScript(nil),
			})
			dpBalanceAmount -= p.DPSplitAmount
			moleculeData = molecule.GoU64ToMoleculeU64(p.DPSplitAmount)
			outputsData = append(outputsData, moleculeData.RawData())
		}
	}
	outputs = append(outputs, &types.CellOutput{
		Capacity: p.DPBaseCapacity,
		Lock:     p.FromLock,
		Type:     p.DPContract.ToScript(nil),
	})
	moleculeData = molecule.GoU64ToMoleculeU64(dpBalanceAmount)
	outputsData = append(outputsData, moleculeData.RawData())
	// capacity
	normalCellCapacity := uint64(0)
	outputsCapacity := uint64(len(outputs)) * p.DPBaseCapacity
	if p.DPLiveCellCapacity > outputsCapacity {
		outputs = append(outputs, &types.CellOutput{
			Capacity: p.DPLiveCellCapacity - outputsCapacity,
			Lock:     p.NormalCellLock,
			Type:     nil,
		})
		outputsData = append(outputsData, []byte{})
	} else {
		normalCellCapacity = outputsCapacity - p.DPLiveCellCapacity
	}
	return outputs, outputsData, normalCellCapacity, nil
}

//

type OutputsDPInfo struct {
	AlgId    common.DasAlgorithmId    `json:"alg_id"`
	SubAlgId common.DasSubAlgorithmId `json:"sub_alg_id"`
	Payload  string                   `json:"payload"`
	AmountDP uint64                   `json:"amount_dp"`
}

func (d *DasCore) GetOutputsDPInfo(tx *types.Transaction) (map[string]OutputsDPInfo, error) {
	var res = make(map[string]OutputsDPInfo)
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return res, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	for i, v := range tx.Outputs {
		if v.Type == nil {
			continue
		}
		if !dpContract.IsSameTypeId(v.Type.CodeHash) {
			continue
		}
		ownerScript, _, err := d.daf.ScriptToHex(v.Lock)
		if err != nil {
			return res, fmt.Errorf("ScriptToHex err: %s", err.Error())
		}
		payload := hex.EncodeToString(ownerScript.AddressPayload)
		amountDP, err := molecule.Bytes2GoU64(tx.OutputsData[i])
		if err != nil {
			return res, fmt.Errorf("Bytes2GoU64 err: %s", err.Error())
		}
		if item, ok := res[payload]; !ok {
			res[payload] = OutputsDPInfo{
				AlgId:    ownerScript.DasAlgorithmId,
				SubAlgId: ownerScript.DasSubAlgorithmId,
				Payload:  payload,
				AmountDP: amountDP,
			}
		} else {
			item.AmountDP += amountDP
			res[payload] = item
		}
	}
	return res, nil
}
