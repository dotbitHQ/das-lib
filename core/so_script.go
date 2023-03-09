package core

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"time"
)

type SoScript struct {
	Name         common.SoScriptType
	OutPoint     types.OutPoint
	SoScriptArgs string
}

func (d *DasCore) InitDasSoScript() error {
	mapSoScript := EnvMainNet.MapSoScript
	if d.net == common.DasNetTypeTestnet2 {
		mapSoScript = EnvTestnet2.MapSoScript
	} else if d.net == common.DasNetTypeTestnet3 {
		mapSoScript = EnvTestnet3.MapSoScript
	}
	//log.Info("mapSoScript:", mapSoScript)
	for k, v := range mapSoScript {
		//log.Info("InitDasSoScript:", k)
		DasSoScriptMap.Store(k, &SoScript{Name: k, SoScriptArgs: v})
	}
	return d.asyncDasSoScript()
}

func (d *DasCore) RunAsyncDasSoScript(t time.Duration) {
	contractTicker := time.NewTicker(t) // update SO
	d.wg.Add(1)
	go func() {
		for {
			select {
			case <-contractTicker.C:
				log.Info("asyncDasSoScript begin ...")
				if err := d.asyncDasSoScript(); err != nil {
					log.Error("asyncDasConfigCell err:", err.Error())
				}
				log.Info("asyncDasSoScript end ...")
			case <-d.ctx.Done():
				d.wg.Done()
				return
			}
		}
	}()
}

func (d *DasCore) asyncDasSoScript() error {
	DasSoScriptMap.Range(func(key, value interface{}) bool {
		item, ok := value.(*SoScript)
		if !ok {
			return true
		}
		if item.SoScriptArgs == "" {
			log.Warn("asyncDasSoScriptByTypeId so script args is nil:", key)
			return true
		}
		searchKey := &indexer.SearchKey{
			Script: &types.Script{
				CodeHash: types.HexToHash(d.dasContractCodeHash),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(item.SoScriptArgs),
			},
			ScriptType: indexer.ScriptTypeType,
		}
		now := time.Now()
		res, err := d.client.GetCells(d.ctx, searchKey, indexer.SearchOrderDesc, 1, "")
		if err != nil {
			log.Error("GetCells err:", key, err.Error())
			return true
		}
		if len(res.Objects) == 0 {
			log.Warn("asyncDasSoScriptByTypeId:", key, len(res.Objects))
		}
		if len(res.Objects) > 0 {
			item.OutPoint.Index = res.Objects[0].OutPoint.Index
			item.OutPoint.TxHash = res.Objects[0].OutPoint.TxHash
			typeId := common.ScriptToTypeId(searchKey.Script)
			log.Info("asyncDasSoScriptByTypeId:", key, item.OutPoint.TxHash, item.OutPoint.Index, typeId, time.Since(now).Seconds())
		}
		return true
	})
	return nil
}

func GetDasSoScript(soScriptName common.SoScriptType) (*SoScript, error) {
	if value, ok := DasSoScriptMap.Load(soScriptName); ok {
		if item, okSo := value.(*SoScript); okSo {
			return item, nil
		}
	}
	return nil, fmt.Errorf("not exist so script: [%s]", soScriptName)
}

func (d *SoScript) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: &d.OutPoint,
		DepType:  types.DepTypeCode,
	}
}
