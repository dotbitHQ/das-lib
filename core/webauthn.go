package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
)

func (d *DasCore) GetKeyListCell(args []byte) (*indexer.LiveCell, error) {
	keyListCell, err := GetDasContractInfo(common.DasKeyListCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	searchKey := indexer.SearchKey{
		Script:     keyListCell.ToScript(args),
		ScriptType: indexer.ScriptTypeType,
		ArgsLen:    0,
		Filter:     nil,
	}
	keyListCells, err := d.client.GetCells(d.ctx, &searchKey, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if subLen := len(keyListCells.Objects); subLen != 1 {
		return nil, nil
	}

	return keyListCells.Objects[0], nil
}
