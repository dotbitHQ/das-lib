package core

import (
	"encoding/binary"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

// quote cell

type QuoteCell struct {
	LiveCell *indexer.LiveCell
}

func (q *QuoteCell) Quote() uint64 {
	return binary.BigEndian.Uint64(q.LiveCell.OutputData[2:])
}

func (q *QuoteCell) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: q.LiveCell.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func (d *DasCore) GetQuoteCell() (*QuoteCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsQuoteCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist quote cell")
	}
	var qc QuoteCell
	qc.LiveCell = res.Objects[0]
	return &qc, nil
}

func (d *DasCore) GetQuoteCellList() ([]*QuoteCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsQuoteCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist quote cell")
	}

	var list []*QuoteCell
	for i, _ := range res.Objects {
		list = append(list, &QuoteCell{LiveCell: res.Objects[i]})
	}
	return list, nil
}

// time cell

type TimeCell struct {
	LiveCell *indexer.LiveCell
}

func (t *TimeCell) Timestamp() int64 {
	return int64(binary.BigEndian.Uint64(t.LiveCell.OutputData[2:]))
}

func (t *TimeCell) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: t.LiveCell.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func (d *DasCore) GetTimeCell() (*TimeCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsTimeCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist time cell")
	}
	var tc TimeCell
	tc.LiveCell = res.Objects[0]
	return &tc, nil
}

func (d *DasCore) GetTimeCellList() ([]*TimeCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsTimeCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist time cell")
	}
	var list []*TimeCell
	for i, _ := range res.Objects {
		list = append(list, &TimeCell{LiveCell: res.Objects[i]})
	}
	return list, nil
}

// height cell

type HeightCell struct {
	LiveCell *indexer.LiveCell
}

func (t *HeightCell) BlockNumber() int64 {
	return int64(binary.BigEndian.Uint64(t.LiveCell.OutputData[2:]))
}

func (t *HeightCell) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: t.LiveCell.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func (d *DasCore) GetHeightCell() (*HeightCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsHeightCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist height cell")
	}
	var hc HeightCell
	hc.LiveCell = res.Objects[0]
	return &hc, nil
}

func (d *DasCore) GetHeightCellList() ([]*HeightCell, error) {
	searchKey := &indexer.SearchKey{
		Script:     common.GetScript(d.thqCodeHash, common.ArgsHeightCell),
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 20, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if len(res.Objects) == 0 {
		return nil, fmt.Errorf("not exist height cell")
	}

	var list []*HeightCell
	for i, _ := range res.Objects {
		list = append(list, &HeightCell{LiveCell: res.Objects[i]})
	}
	return list, nil
}

func (d *DasCore) GetTxQuote(tx *types.Transaction) (uint64, error) {
	if tx == nil {
		return 0, fmt.Errorf("tx is nil")
	}
	var quote uint64
	for _, v := range tx.CellDeps {
		cellDepTx, err := d.client.GetTransaction(d.ctx, v.OutPoint.TxHash)
		if err != nil {
			return 0, fmt.Errorf("GetTransaction CellDeps err: %s", err.Error())
		}
		cell := cellDepTx.Transaction.Outputs[v.OutPoint.Index]
		if cell.Type != nil {
			if common.Bytes2Hex(cell.Type.Args) == common.ArgsQuoteCell && d.thqCodeHash == cell.Type.CodeHash.Hex() {
				quote = binary.BigEndian.Uint64(cellDepTx.Transaction.OutputsData[v.OutPoint.Index][2:])
				break
			}
		}
	}
	return quote, nil
}
