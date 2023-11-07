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
	log.Info("GetDpCells:", common.Bytes2Hex(p.LockScript.Args))
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	searchKey := &indexer.SearchKey{
		Script:     p.LockScript,
		ScriptType: indexer.ScriptTypeLock,
		//Filter: &indexer.CellsFilter{
		//	OutputDataLenRange: &[2]uint64{12, 13},
		//},
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
			log.Info("GetDpCells:", liveCell.OutPoint.TxHash.Hex())
			if liveCell.Output.Type != nil && !dpContract.IsSameTypeId(liveCell.Output.Type.CodeHash) {
				continue
			}
			if p.AmountNeed > 0 && p.DasCache != nil && p.DasCache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
				hasCache = true
				continue
			}
			cells = append(cells, liveCell)

			dpData, err := witness.ConvertBysToDPData(liveCell.OutputData)
			if err != nil {
				return nil, 0, 0, err
			}

			totalAmount += dpData.Value
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
	log.Info("GetDpCells:", p.AmountNeed, totalAmount)
	if p.AmountNeed > 0 && totalAmount < p.AmountNeed {
		if hasCache {
			return cells, totalAmount, totalCapacity, ErrRejectedOutPoint
		}
		log.Info("GetDpCells:", p.AmountNeed, totalAmount)
		return cells, totalAmount, totalCapacity, ErrInsufficientFunds
	}
	return cells, totalAmount, totalCapacity, nil
}

func (d *DasCore) GetDPointTransferWhitelist() (map[string]*types.Script, error) {
	builder, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return builder.GetDPointTransferWhitelist()
}

func (d *DasCore) GetDPointCapacityRecycleWhitelist() (map[string]*types.Script, error) {
	builder, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return builder.GetDPointCapacityRecycleWhitelist()
}

func (d *DasCore) GetDPBaseCapacity() (uint64, uint64, error) {
	builder, err := d.ConfigCellDataBuilderByTypeArgs(common.ConfigCellTypeArgsDPoint)
	if err != nil {
		return 0, 0, fmt.Errorf("ConfigCellDataBuilderByTypeArgs err: %s", err.Error())
	}
	return builder.GetDPBaseCapacity()
}

type ParamSplitDPCell struct {
	FromLock           *types.Script
	ToLock             *types.Script
	DPLiveCell         []*indexer.LiveCell
	DPLiveCellCapacity uint64
	DPTotalAmount      uint64
	DPTransferAmount   uint64
	DPSplitCount       int
	DPSplitAmount      uint64
	NormalCellLock     *types.Script
}

func (d *DasCore) SplitDPCell(p *ParamSplitDPCell) ([]*types.CellOutput, [][]byte, uint64, error) {
	basicCapacity, preparedFeeCapacity, err := d.GetDPBaseCapacity()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("GetDPBaseCapacity err: %s", err.Error())
	}
	dpBaseCapacity := basicCapacity + preparedFeeCapacity
	//
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("GetDasContractInfo err: %s", err)
	}
	var outputs []*types.CellOutput
	var outputsData [][]byte
	// transfer
	outputs = append(outputs, &types.CellOutput{
		Capacity: dpBaseCapacity,
		Lock:     p.ToLock,
		Type:     dpContract.ToScript(nil),
	})
	oData, err := witness.ConvertDPDataToBys(witness.DPData{Value: p.DPTransferAmount})
	if err != nil {
		return nil, nil, 0, fmt.Errorf("ConvertDPDataToBys err: %s", err.Error())
	}
	outputsData = append(outputsData, oData)
	// split
	dpBalanceAmount := p.DPTotalAmount - p.DPTransferAmount
	for i := 1; i < p.DPSplitCount; i++ {
		if dpBalanceAmount > p.DPSplitAmount*2 {
			outputs = append(outputs, &types.CellOutput{
				Capacity: dpBaseCapacity,
				Lock:     p.FromLock,
				Type:     dpContract.ToScript(nil),
			})
			dpBalanceAmount -= p.DPSplitAmount
			oData, err = witness.ConvertDPDataToBys(witness.DPData{Value: p.DPSplitAmount})
			if err != nil {
				return nil, nil, 0, fmt.Errorf("ConvertDPDataToBys err: %s", err.Error())
			}
			outputsData = append(outputsData, oData)
		}
	}
	if dpBalanceAmount > 0 {
		outputs = append(outputs, &types.CellOutput{
			Capacity: dpBaseCapacity,
			Lock:     p.FromLock,
			Type:     dpContract.ToScript(nil),
		})
		oData, err = witness.ConvertDPDataToBys(witness.DPData{Value: dpBalanceAmount})
		if err != nil {
			return nil, nil, 0, fmt.Errorf("ConvertDPDataToBys err: %s", err.Error())
		}
		outputsData = append(outputsData, oData)
	}

	// capacity
	normalCellCapacity := uint64(0)
	outputsCapacity := uint64(len(outputs)) * dpBaseCapacity
	if p.DPLiveCellCapacity > outputsCapacity {
		outputs = append(outputs, &types.CellOutput{
			Capacity: p.DPLiveCellCapacity - outputsCapacity,
			Lock:     p.NormalCellLock,
			Type:     nil,
		})
		outputsData = append(outputsData, []byte{})
	} else if p.DPLiveCellCapacity < outputsCapacity {
		diff := outputsCapacity - p.DPLiveCellCapacity
		if diff < preparedFeeCapacity {
			outputs[len(outputs)-1].Capacity -= diff
		} else {
			normalCellCapacity = outputsCapacity - p.DPLiveCellCapacity
		}
	}
	return outputs, outputsData, normalCellCapacity, nil
}

//

type TxDPInfo struct {
	AlgId    common.DasAlgorithmId    `json:"alg_id"`
	SubAlgId common.DasSubAlgorithmId `json:"sub_alg_id"`
	Payload  string                   `json:"payload"`
	AmountDP uint64                   `json:"amount_dp"`
}

func (d *DasCore) GetOutputsDPInfo(tx *types.Transaction) (map[string]TxDPInfo, error) {
	var res = make(map[string]TxDPInfo)
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
		//fmt.Println("Args", common.Bytes2Hex(v.Lock.Args))
		ownerScript, _, err := d.daf.ScriptToHex(v.Lock)
		if err != nil {
			return res, fmt.Errorf("ScriptToHex err: %s", err.Error())
		}
		//fmt.Println("AddressPayload", common.Bytes2Hex(ownerScript.AddressPayload))
		payload := hex.EncodeToString(ownerScript.AddressPayload)
		dpData, err := witness.ConvertBysToDPData(tx.OutputsData[i])
		if err != nil {
			return res, fmt.Errorf("Bytes2GoU64 err: %s", err.Error())
		}
		if item, ok := res[payload]; !ok {
			res[payload] = TxDPInfo{
				AlgId:    ownerScript.DasAlgorithmId,
				SubAlgId: ownerScript.DasSubAlgorithmId,
				Payload:  payload,
				AmountDP: dpData.Value,
			}
		} else {
			item.AmountDP += dpData.Value
			res[payload] = item
		}
	}
	return res, nil
}

func (d *DasCore) GetInputsDPInfo(tx *types.Transaction) (map[string]TxDPInfo, error) {
	var res = make(map[string]TxDPInfo)
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return res, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}
	var mapTx = make(map[string]*types.Transaction)
	for _, v := range tx.Inputs {
		tmpTx, ok := mapTx[v.PreviousOutput.TxHash.Hex()]
		if !ok {
			txStatus, err := d.client.GetTransaction(d.ctx, v.PreviousOutput.TxHash)
			if err != nil {
				return res, fmt.Errorf("GetTransaction err: %s", err.Error())
			}
			mapTx[v.PreviousOutput.TxHash.Hex()] = txStatus.Transaction
			tmpTx = txStatus.Transaction
		}
		preOutput := tmpTx.Outputs[v.PreviousOutput.Index]
		if preOutput.Type == nil {
			continue
		}
		if !dpContract.IsSameTypeId(preOutput.Type.CodeHash) {
			continue
		}
		ownerScript, _, err := d.daf.ScriptToHex(preOutput.Lock)
		if err != nil {
			return res, fmt.Errorf("ScriptToHex err: %s", err.Error())
		}
		payload := hex.EncodeToString(ownerScript.AddressPayload)
		amountDP, err := molecule.Bytes2GoU64(tmpTx.OutputsData[v.PreviousOutput.Index])
		if err != nil {
			return res, fmt.Errorf("Bytes2GoU64 err: %s", err.Error())
		}
		if item, ok := res[payload]; !ok {
			res[payload] = TxDPInfo{
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
