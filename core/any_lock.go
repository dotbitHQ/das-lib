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
	AnyLockNameNoStr    AnyLockName = "nostr"
)

func (d *DasCore) GetAnyLockCellDep(anyLockName AnyLockName) (*types.CellDep, error) {
	switch anyLockName {
	case AnyLockNameNoStr:
		argsNoStrLockArgs := "0xfad8cb75eb0bb01718e2336002064568bc05887af107f74ed5dd501829e192f8"
		if d.net != common.DasNetTypeMainNet {
			argsNoStrLockArgs = "0x8dc56c6f35f0c535e23ded1629b1f20535477a1b43e59f14617d11e32c50e0aa"
		}
		searchKey := indexer.SearchKey{
			Script: &types.Script{
				CodeHash: types.HexToHash("0x00000000000000000000000000000000000000000000000000545950455f4944"),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(argsNoStrLockArgs),
			},
			ScriptType: indexer.ScriptTypeType,
		}
		res, err := d.client.GetCells(context.Background(), &searchKey, indexer.SearchOrderDesc, 1, "")
		if err != nil {
			return nil, fmt.Errorf("GetCells err: %s", err.Error())
		}
		log.Info("GetAnyLockCellDep:", len(res.Objects))
		if len(res.Objects) == 0 {
			return nil, fmt.Errorf("GetCells is nil")
		}
		return &types.CellDep{
			OutPoint: res.Objects[0].OutPoint,
			DepType:  types.DepTypeCode,
		}, nil
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
		log.Info("GetAnyLockCellDep:", len(res.Objects))
		if len(res.Objects) == 0 {
			return nil, fmt.Errorf("GetCells is nil")
		}
		return &types.CellDep{
			OutPoint: res.Objects[0].OutPoint,
			DepType:  types.DepTypeCode,
		}, nil
	case AnyLockNameJoyID:
		txHash := "0x06d4b4cc802115633b0ed89fac859504b1c08c93869ee9748b1d17c1d0e149ae"
		if d.net != common.DasNetTypeMainNet {
			txHash = "0x759f281588c96979764cb21c196478cf8e13ea81fede7f4ba26d1ff29dbc6a81"
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
