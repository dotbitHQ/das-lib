package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
)

func (d *DasCore) GetReverseRecordSmtCell() (*indexer.LiveCell, error) {
	contractReverseRecord, err := GetDasContractInfo(common.DasContractNameReverseRecordRootCellType)
	if err != nil {
		return nil, fmt.Errorf("GetReverseRecordSmtCell GetDasContractInfo err: %s", err.Error())
	}

	searchKey := indexer.SearchKey{
		Script:     contractReverseRecord.ToScript(nil),
		ScriptType: indexer.ScriptTypeType,
		ArgsLen:    0,
		Filter:     nil,
	}
	reverseRecordLiveCells, err := d.client.GetCells(d.ctx, &searchKey, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, fmt.Errorf("GetReverseRecordSmtCell GetCells err: %s", err.Error())
	}
	if subLen := len(reverseRecordLiveCells.Objects); subLen != 1 {
		return nil, fmt.Errorf("GetReverseRecordSmtCell %s cell len: %d", common.DasContractNameReverseRecordRootCellType, subLen)
	}
	return reverseRecordLiveCells.Objects[0], nil
}
