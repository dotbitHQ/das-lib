package core

import (
	"context"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type AnyLockName string

const (
	AnyLockNameOmniLock AnyLockName = "omni-lock"
	AnyLockNameJoyID    AnyLockName = "joyid"
)

func (d *DasCore) GetAnyLockCellDep(anyLockName AnyLockName) (*types.CellDep, error) {
	switch anyLockName {
	case AnyLockNameOmniLock:
		argsOmniLockArgs := "0x855508fe0f0ca25b935b070452ecaee48f6c9f1d66cd15f046616b99e948236a"
		if d.net != common.DasNetTypeMainNet {
			argsOmniLockArgs = "0x761f51fc9cd6a504c32c6ae64b3746594d1af27629b427c5ccf6c9a725a89144"
		}
		searchKey := indexer.SearchKey{
			Script: &types.Script{
				CodeHash: types.HexToHash("0x00000000000000000000000000000000000000000000000000545950455f4944"),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(argsOmniLockArgs),
			},
			ScriptType: indexer.ScriptTypeType,
		}
		res, err := d.client.GetCells(context.Background(), &searchKey, indexer.SearchOrderDesc, 1, "")
		if err != nil {
			return nil, fmt.Errorf("GetCells err: %s", err.Error())
		}
		log.Info("GetAnyLockOutpoint:", len(res.Objects))
		if len(res.Objects) == 0 {
			return nil, fmt.Errorf("GetCells is nil")
		}
		return &types.CellDep{
			OutPoint: res.Objects[0].OutPoint,
			DepType:  types.DepTypeCode,
		}, nil
	case AnyLockNameJoyID:
		txHash := "0xf05188e5f3a6767fc4687faf45ba5f1a6e25d3ada6129dae8722cb282f262493"
		if d.net != common.DasNetTypeMainNet {
			txHash = "0x4dcf3f3b09efac8995d6cbee87c5345e812d310094651e0c3d9a730f32dc9263"
		}
		return &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash(txHash),
				Index:  0,
			},
			DepType: types.DepTypeDepGroup,
		}, nil
	default:
		return nil, fmt.Errorf("unsupport")
	}
}
