package core

import (
	"context"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/go-redis/redis"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"golang.org/x/sync/syncmap"
	"sync"
)

var (
	log                      = logger.NewLogger("das-core", mylog.LevelDebug)
	DasContractMap           syncmap.Map                               // map[contact name]{contract info}
	DasContractByTypeIdMap   = make(map[string]common.DasContractName) // map[contract type id]{contract name}
	DasConfigCellMap         syncmap.Map                               // map[ConfigCellTypeArgs]config cell info
	DasConfigCellByTxHashMap syncmap.Map                               // map[tx hash]{true}
	DasSoScriptMap           syncmap.Map                               // map[so script type]
)

type DasCore struct {
	client              rpc.Client
	ctx                 context.Context
	wg                  *sync.WaitGroup
	dasContractCodeHash string // contract code hash
	dasContractArgs     string // contract owner args
	thqCodeHash         string // time,height,quote cell code hash
	net                 common.DasNetType
	daf                 *DasAddressFormat
	red                 *redis.Client
}

func NewDasCore(ctx context.Context, wg *sync.WaitGroup, opts ...DasCoreOption) *DasCore {
	var dc DasCore
	dc.ctx = ctx
	dc.wg = wg
	for _, opt := range opts {
		opt(&dc)
	}
	return &dc
}

func (d *DasCore) Client() rpc.Client {
	return d.client
}

func (d *DasCore) NetType() common.DasNetType {
	return d.net
}

func (d *DasCore) Daf() *DasAddressFormat {
	return d.daf
}

func (d *DasCore) GetDasLock() *types.Script {
	switch d.net {
	case common.DasNetTypeMainNet:
		return common.GetNormalLockScriptByMultiSig(EnvMainNet.ContractArgs)
	case common.DasNetTypeTestnet2:
		return common.GetNormalLockScriptByMultiSig(EnvTestnet2.ContractArgs)
	case common.DasNetTypeTestnet3:
		return common.GetNormalLockScript(EnvTestnet3.ContractArgs)
	default:
		return nil
	}
}

func SetLogLevel(level int) {
	log = logger.NewLogger("das-core", level)
}
