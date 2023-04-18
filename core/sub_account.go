package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func (d *DasCore) GetSubAccountCell(parentAccountId string) (*indexer.LiveCell, error) {
	contractSubAcc, err := GetDasContractInfo(common.DASContractNameSubAccountCellType)
	if err != nil {
		return nil, fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	}

	searchKey := indexer.SearchKey{
		Script:     contractSubAcc.ToScript(common.Hex2Bytes(parentAccountId)),
		ScriptType: indexer.ScriptTypeType,
		ArgsLen:    0,
		Filter:     nil,
	}
	subAccLiveCells, err := d.client.GetCells(d.ctx, &searchKey, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if subLen := len(subAccLiveCells.Objects); subLen != 1 {
		return nil, fmt.Errorf("sub account cell len: %d", subLen)
	}
	return subAccLiveCells.Objects[0], nil
}

func (d *DasCore) GetCustomScriptLiveCell(data []byte) (*indexer.LiveCell, error) {
	subDataDetail := witness.ConvertSubAccountCellOutputData(data)
	var customScript *types.Script
	switch subDataDetail.Flag {
	case 1:
		customScript = &types.Script{
			CodeHash: types.HexToHash("0x00000000000000000000000000000000000000000000000000545950455f4944"),
			HashType: types.HashTypeType,
			Args:     subDataDetail.CustomScriptArgs,
		}
	}
	if customScript == nil {
		return nil, fmt.Errorf("customScript is nil")
	}
	searchKey := indexer.SearchKey{
		Script:     customScript,
		ScriptType: indexer.ScriptTypeType,
	}
	customScriptCell, err := d.client.GetCells(d.ctx, &searchKey, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	if subLen := len(customScriptCell.Objects); subLen != 1 {
		return nil, fmt.Errorf("sub account outpoint len: %d", subLen)
	}
	log.Info("getCustomScriptLiveCell:", common.OutPointStruct2String(customScriptCell.Objects[0].OutPoint))
	return customScriptCell.Objects[0], nil
}
