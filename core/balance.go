package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/nervosnetwork/ckb-sdk-go/collector"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

var (
	ErrRejectedOutPoint  = errors.New("RejectedOutPoint")
	ErrInsufficientFunds = errors.New("InsufficientFunds")
	ErrNotEnoughChange   = errors.New("NotEnoughChange")
)

func GetSatisfiedLimitLiveCell(client rpc.Client, dasCache *dascache.DasCache, searchKey *indexer.SearchKey, needLimit uint64, order indexer.SearchOrder) ([]*indexer.LiveCell, error) {
	co := collector.NewLiveCellCollector(client, searchKey, order, indexer.SearchLimit, "")
	co.TypeScript = searchKey.Filter.Script
	if searchKey.ScriptType == indexer.ScriptTypeType {
		co.TypeScript = searchKey.Script
	}

	iterator, err := co.Iterator()
	if err != nil {
		return nil, fmt.Errorf("iterator err:%s", err.Error())
	}
	var cells []*indexer.LiveCell
	foundLimit := uint64(0)
	for iterator.HasNext() {
		liveCell, err := iterator.CurrentItem()
		if err != nil {
			return nil, fmt.Errorf("CurrentItem err:%s", err.Error())
		}
		if dasCache != nil && dasCache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
		} else {
			cells = append(cells, liveCell)
			foundLimit = foundLimit + 1
			if foundLimit >= needLimit {
				break
			}
		}
		if err = iterator.Next(); err != nil {
			return nil, fmt.Errorf("next err:%s", err.Error())
		}
	}
	return cells, nil
}

func GetSatisfiedCapacityLiveCellWithOrder(client rpc.Client, dasCache *dascache.DasCache, dasLockScript, dasTypeScript *types.Script, capacityNeed, capacityForChange uint64, order indexer.SearchOrder) ([]*indexer.LiveCell, uint64, error) {
	if client == nil {
		return nil, 0, fmt.Errorf("client is nil")
	}
	searchKey := &indexer.SearchKey{
		Script:     dasLockScript,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script:             dasTypeScript,
			OutputDataLenRange: &[2]uint64{0, 1},
		},
	}
	co := collector.NewLiveCellCollector(client, searchKey, order, indexer.SearchLimit, "")
	co.TypeScript = searchKey.Filter.Script
	iterator, err := co.Iterator()
	if err != nil {
		return nil, 0, fmt.Errorf("iterator err:%s", err.Error())
	}
	var cells []*indexer.LiveCell
	total := uint64(0)
	hasCache := false
	for iterator.HasNext() {
		liveCell, err := iterator.CurrentItem()
		if err != nil {
			return nil, 0, fmt.Errorf("CurrentItem err:%s", err.Error())
		}
		if capacityNeed > 0 && dasCache != nil && dasCache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
			hasCache = true
		} else {
			cells = append(cells, liveCell)
			total += liveCell.Output.Capacity
			if capacityNeed > 0 && (total == capacityNeed || total >= capacityNeed+capacityForChange) { // limit 为转账金额+手续费
				break
			}
		}
		if err = iterator.Next(); err != nil {
			return nil, 0, fmt.Errorf("next err:%s", err.Error())
		}
	}
	if capacityNeed > 0 && total != capacityNeed {
		if total < capacityNeed {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrInsufficientFunds
			}
		} else if total < capacityNeed+capacityForChange {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrNotEnoughChange
			}
		}
	}
	return cells, total, nil
}

func GetSatisfiedCapacityLiveCell(client rpc.Client, dasCache *dascache.DasCache, dasLockScript, dasTypeScript *types.Script, capacityNeed, capacityForChange uint64) ([]*indexer.LiveCell, uint64, error) {
	return GetSatisfiedCapacityLiveCellWithOrder(client, dasCache, dasLockScript, dasTypeScript, capacityNeed, capacityForChange, indexer.SearchOrderDesc)
}

func SplitOutputCell(total, base, limit uint64, lockScript, typeScript *types.Script) ([]*types.CellOutput, error) {
	log.Info("total: ", total, "base: ", base, "limit: ", limit)
	formatCell := &types.CellOutput{
		Capacity: base,
		Lock:     lockScript,
		Type:     typeScript,
	}
	realBase := formatCell.OccupiedCapacity(nil) * 1e8
	if total < realBase || base < realBase {
		return nil, fmt.Errorf("total(%d) or base(%d) should not less than real base(%d)", total, base, realBase)
	}
	cellLen := total / base
	left := total % base
	var cellList []*types.CellOutput
	if cellLen <= limit {
		limit = 0
	}
	var baseLen, leftCapacity uint64
	if limit == 0 {
		if left < realBase && cellLen > 0 {
			baseLen = cellLen - 1
			leftCapacity = base + left
			log.Info("left: ", left, "realBase: ", realBase)
		} else if left >= realBase {
			baseLen = cellLen
			leftCapacity = left
			log.Info("left: ", left, "realBase: ", realBase)
		} else {
			return nil, fmt.Errorf("total(%d) should not less than base(%d)", total, base)
		}
	} else {
		baseLen = limit
		leftCapacity = (cellLen-limit)*base + left
	}
	log.Info("baseLen: ", baseLen, "leftCapacity: ", leftCapacity)
	for i := uint64(0); i < baseLen; i++ {
		cellList = append(cellList, formatCell)
	}
	tmp := &types.CellOutput{
		Capacity: leftCapacity,
		Lock:     lockScript,
		Type:     typeScript,
	}
	cellList = append(cellList, tmp)

	return cellList, nil
}

type ParamGetBalanceCells struct {
	DasCache          *dascache.DasCache
	LockScript        *types.Script
	CapacityNeed      uint64
	CapacityForChange uint64
	SearchOrder       indexer.SearchOrder
}

func (d *DasCore) GetBalanceCells(p *ParamGetBalanceCells) ([]*indexer.LiveCell, uint64, error) {
	if d.client == nil {
		return nil, 0, fmt.Errorf("client is nil")
	}
	if p == nil {
		return nil, 0, fmt.Errorf("param is nil")
	}
	balanceContract, err := GetDasContractInfo(common.DasContractNameBalanceCellType)
	if err != nil {
		return nil, 0, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	searchKey := &indexer.SearchKey{
		Script:     p.LockScript,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script:             nil,
			OutputDataLenRange: &[2]uint64{0, 1},
		},
	}
	var cells []*indexer.LiveCell
	total := uint64(0)
	hasCache := false
	lastCursor := ""

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
			if liveCell.Output.Type != nil && !balanceContract.IsSameTypeId(liveCell.Output.Type.CodeHash) {
				continue
			}
			if p.CapacityNeed > 0 && p.DasCache != nil && p.DasCache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
				hasCache = true
				continue
			}
			cells = append(cells, liveCell)
			total += liveCell.Output.Capacity
			if p.CapacityNeed > 0 {
				if total == p.CapacityNeed || total >= p.CapacityNeed+p.CapacityForChange {
					break
				}
			}
		}
	}

	if p.CapacityNeed > 0 {
		if total < p.CapacityNeed {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrInsufficientFunds
			}
		} else if total < p.CapacityNeed+p.CapacityForChange {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrNotEnoughChange
			}
		}
	}
	return cells, total, nil
}

func SplitOutputCell2(total, base, limit uint64, lockScript, typeScript *types.Script) ([]*types.CellOutput, error) {
	log.Info("SplitOutputCell2:", "total: ", total, "base: ", base, "limit: ", limit)
	formatCell := &types.CellOutput{
		Capacity: base,
		Lock:     lockScript,
		Type:     typeScript,
	}
	realBase := formatCell.OccupiedCapacity(nil) * 1e8
	if total < realBase || base < realBase {
		return nil, fmt.Errorf("total(%d) or base(%d) should not less than real base(%d)", total, base, realBase)
	}
	log.Info("realBase:", realBase)

	var cellList []*types.CellOutput
	splitTotal := uint64(0)
	for i := uint64(0); i < limit && splitTotal+2*base < total; i++ {
		cellList = append(cellList, formatCell)
		splitTotal += base
	}
	cellList = append(cellList, &types.CellOutput{
		Capacity: total - splitTotal,
		Lock:     lockScript,
		Type:     typeScript,
	})

	return cellList, nil
}
