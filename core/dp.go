package core

import (
	"context"
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

func (d *DasCore) GetDpCells(p *ParamGetDpCells) ([]*indexer.LiveCell, uint64, error) {
	if d.client == nil {
		return nil, 0, fmt.Errorf("client is nil")
	}
	if p == nil {
		return nil, 0, fmt.Errorf("param is nil")
	}
	dpContract, err := GetDasContractInfo(common.DasContractNameDpCellType)
	if err != nil {
		return nil, 0, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
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
	total := uint64(0)
	hasCache := false
	lastCursor := ""

	ok := false
	for {
		liveCells, err := d.client.GetCells(context.Background(), searchKey, p.SearchOrder, indexer.SearchLimit, lastCursor)
		if err != nil {
			return nil, 0, err
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
				return nil, 0, err
			}
			amount, err := molecule.Bytes2GoU64(liveCell.OutputData[idx : idx+int(l)])
			if err != nil {
				return nil, 0, err
			}
			idx += int(l)

			total += amount
			if p.AmountNeed > 0 && total >= p.AmountNeed {
				ok = true
				break
			}
		}

		if ok {
			break
		}
	}

	if p.AmountNeed > 0 && total < p.AmountNeed {
		if hasCache {
			return cells, total, ErrRejectedOutPoint
		}
		return cells, total, ErrInsufficientFunds
	}
	return cells, total, nil
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
