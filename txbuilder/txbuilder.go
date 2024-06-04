package txbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/dascache"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sync"
)

var log = logger.NewLogger("txbuilder", mylog.LevelDebug)

type DasTxBuilder struct {
	*DasTxBuilderBase                                  // for base
	*DasTxBuilderTransaction                           // for tx
	DasMMJson                                          // for mmjson
	mapCellDep               map[string]*types.CellDep // for memory
	notCheckInputs           bool
	otherWitnesses           [][]byte
}

func NewDasTxBuilderBase(ctx context.Context, dasCore *core.DasCore, handle sign.HandleSignCkbMessage, serverArgs string) *DasTxBuilderBase {
	var base DasTxBuilderBase
	base.ctx = ctx
	base.dasCore = dasCore
	base.handleServerSign = handle
	base.serverArgs = serverArgs
	return &base
}

func NewDasTxBuilderFromBase(base *DasTxBuilderBase, tx *DasTxBuilderTransaction) *DasTxBuilder {
	var b DasTxBuilder
	b.DasTxBuilderBase = base
	b.DasTxBuilderTransaction = tx
	if tx == nil {
		b.DasTxBuilderTransaction = &DasTxBuilderTransaction{}
		b.MapInputsCell = make(map[string]*types.CellWithStatus)
	}
	b.mapCellDep = make(map[string]*types.CellDep)
	return &b
}

type DasTxBuilderBase struct {
	ctx              context.Context
	dasCore          *core.DasCore
	handleServerSign sign.HandleSignCkbMessage
	serverArgs       string
}

type DasTxBuilderTransaction struct {
	Transaction     *types.Transaction               `json:"transaction"`
	MapInputsCell   map[string]*types.CellWithStatus `json:"map_inputs_cell"`
	ServerSignGroup []int                            `json:"server_sign_group"`
}

type DasMMJson struct {
	account   string
	salePrice uint64
	offers    int // cancel offer count
}

type BuildTransactionParams struct {
	CellDeps       []*types.CellDep    `json:"cell_deps"`
	HeadDeps       []types.Hash        `json:"head_deps"`
	Inputs         []*types.CellInput  `json:"inputs"`
	Outputs        []*types.CellOutput `json:"outputs"`
	OutputsData    [][]byte            `json:"outputs_data"`
	Witnesses      [][]byte            `json:"witnesses"`
	OtherWitnesses [][]byte            `json:"other_witnesses"`
	LatestWitness  [][]byte
}

func (d *DasTxBuilder) BuildTransactionWithCheckInputs(p *BuildTransactionParams, notCheckInputs bool) error {
	d.notCheckInputs = notCheckInputs
	err := d.newTx()
	if err != nil {
		return fmt.Errorf("newBaseTx err: %s", err.Error())
	}

	err = d.addInputsForTx(p.Inputs)
	if err != nil {
		return fmt.Errorf("addInputsForBaseTx err: %s", err.Error())
	}

	err = d.addOutputsForTx(p.Outputs, p.OutputsData)
	if err != nil {
		return fmt.Errorf("addOutputsForBaseTx err: %s", err.Error())
	}

	d.Transaction.Witnesses = append(d.Transaction.Witnesses, p.Witnesses...)
	d.otherWitnesses = append(d.otherWitnesses, p.OtherWitnesses...)

	err = d.addWebauthnInfo()
	if err != nil {
		return fmt.Errorf("addWebauthnInfo err: %s", err.Error())
	}

	if err := d.addMapCellDepWitnessForBaseTx(p.CellDeps); err != nil {
		return fmt.Errorf("addMapCellDepWitnessForBaseTx err: %s", err.Error())
	}

	for _, v := range p.HeadDeps {
		d.Transaction.HeaderDeps = append(d.Transaction.HeaderDeps, v)
	}
	for _, v := range p.LatestWitness {
		d.Transaction.Witnesses = append(d.Transaction.Witnesses, v)
	}
	return nil
}

func (d *DasTxBuilder) BuildTransaction(p *BuildTransactionParams) error {
	return d.BuildTransactionWithCheckInputs(p, false)
}

func (d *DasTxBuilder) TxString() string {
	txStr, _ := rpc.TransactionString(d.Transaction)
	return txStr
}

func (d *DasTxBuilder) GetDasTxBuilderTransactionString() string {
	bys, err := json.Marshal(d.DasTxBuilderTransaction)
	if err != nil {
		return ""
	}
	return string(bys)
}

func GenerateDigestListFromTx(cli rpc.Client, txJson string, skipGroups []int) ([]SignData, error) {
	Tx, err := rpc.TransactionFromString(txJson)
	if err != nil {
		return nil, err
	}
	hash, _ := Tx.ComputeHash()
	fmt.Println(hash.Hex())

	var dasTxBuilderTransaction DasTxBuilderTransaction
	var txBuilder DasTxBuilder
	var netType common.DasNetType
	wgServer := sync.WaitGroup{}
	ctxServer := context.Background()
	blockInfo, err := cli.GetBlockchainInfo(ctxServer)
	if err != nil {
		return nil, err
	}

	if blockInfo.Chain == "ckb" {
		netType = common.DasNetTypeMainNet
	} else if blockInfo.Chain == "ckb_testnet" {
		netType = common.DasNetTypeTestnet2
	} else {
		netType = common.DasNetTypeTestnet3
	}

	dasTxBuilderTransaction.Transaction = Tx
	dasTxBuilderTransaction.MapInputsCell = make(map[string]*types.CellWithStatus)

	txBuilder.DasTxBuilderTransaction = &dasTxBuilderTransaction
	ops := []core.DasCoreOption{
		core.WithClient(cli),
		core.WithDasNetType(netType),
	}
	dasCore := core.NewDasCore(ctxServer, &wgServer, ops...)

	txBuilderBase := NewDasTxBuilderBase(ctxServer, dasCore, nil, "")
	txBuilder.DasTxBuilderBase = txBuilderBase
	digestList, err := txBuilder.GenerateDigestListFromTx(skipGroups)
	return digestList, nil
}

type CheckTxFeeParam struct {
	TxParams      *BuildTransactionParams
	DasCache      *dascache.DasCache
	TxFee         uint64
	FeeLock       *types.Script
	TxBuilderBase *DasTxBuilderBase
	DasCore       *core.DasCore
}

func CheckTxFee(checkTxFeeParam *CheckTxFeeParam) (*DasTxBuilder, error) {
	if checkTxFeeParam.FeeLock == nil {
		log.Info("checkTxFeeParam.FeeLock is nil")
		return nil, nil
	}
	if checkTxFeeParam.TxFee >= common.UserCellTxFeeLimit {
		log.Info("buildTx das fee:", checkTxFeeParam.TxFee)
		change, liveBalanceCell, err := checkTxFeeParam.DasCore.GetBalanceCellWithLock(&core.ParamGetBalanceCells{
			LockScript:   checkTxFeeParam.FeeLock,
			CapacityNeed: checkTxFeeParam.TxFee,
			DasCache:     checkTxFeeParam.DasCache,
			SearchOrder:  indexer.SearchOrderDesc,
		})
		if err != nil {
			return nil, fmt.Errorf("GetBalanceCell err %s", err.Error())
		}
		for _, v := range liveBalanceCell {
			checkTxFeeParam.TxParams.Inputs = append(checkTxFeeParam.TxParams.Inputs, &types.CellInput{
				PreviousOutput: v.OutPoint,
			})
		}
		// change balance_cell
		checkTxFeeParam.TxParams.Outputs = append(checkTxFeeParam.TxParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     checkTxFeeParam.FeeLock,
		})

		checkTxFeeParam.TxParams.OutputsData = append(checkTxFeeParam.TxParams.OutputsData, []byte{})
		txBuilder := NewDasTxBuilderFromBase(checkTxFeeParam.TxBuilderBase, nil)
		err = txBuilder.BuildTransaction(checkTxFeeParam.TxParams)
		if err != nil {
			return nil, fmt.Errorf("txBuilder.BuildTransaction err: %s", err.Error())
		}
		log.Info("buildTx: das:", txBuilder.TxString())
		return txBuilder, nil
	}
	return nil, nil
}

func DeepCopyTxParams(src interface{}) (*BuildTransactionParams, error) {
	var params BuildTransactionParams
	jsonString, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err %s", err.Error())
	}
	err = json.Unmarshal(jsonString, &params)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal err %s", err.Error())

	}
	return &params, nil
}

//func TxToTransaction(tx *types.Transaction) *Transaction {
//	var res Transaction
//	res.Version = tx.Version
//	for _, v := range tx.CellDeps {
//		res.CellDeps = append(res.CellDeps, &CellDep{
//			OutPoint: &OutPoint{
//				TxHash: v.OutPoint.TxHash,
//				Index:  v.OutPoint.Index,
//			},
//			DepType: v.DepType,
//		})
//	}
//	res.HeaderDeps = tx.HeaderDeps
//	for _, v := range tx.Inputs {
//		res.Inputs = append(res.Inputs, &CellInput{
//			Since: v.Since,
//			PreviousOutput: &OutPoint{
//				TxHash: v.PreviousOutput.TxHash,
//				Index:  v.PreviousOutput.Index,
//			},
//		})
//	}
//	for _, v := range tx.Outputs {
//		tmp := CellOutput{
//			Capacity: v.Capacity,
//			Lock:     nil,
//			Type:     nil,
//		}
//		if v.Lock != nil {
//			tmp.Lock = &Script{
//				CodeHash: v.Lock.CodeHash,
//				HashType: v.Lock.HashType,
//				Args:     v.Lock.Args,
//			}
//		}
//		if v.Type != nil {
//			tmp.Type = &Script{
//				CodeHash: v.Type.CodeHash,
//				HashType: v.Type.HashType,
//				Args:     v.Type.Args,
//			}
//		}
//		res.Outputs = append(res.Outputs, &tmp)
//	}
//	res.OutputsData = tx.OutputsData
//	res.Witnesses = tx.Witnesses
//
//	return &res
//}
//
//func TransactionToTx(tx *Transaction) *types.Transaction {
//	var res types.Transaction
//	res.Version = tx.Version
//	for _, v := range tx.CellDeps {
//		res.CellDeps = append(res.CellDeps, &types.CellDep{
//			OutPoint: &types.OutPoint{
//				TxHash: v.OutPoint.TxHash,
//				Index:  v.OutPoint.Index,
//			},
//			DepType: v.DepType,
//		})
//	}
//	res.HeaderDeps = tx.HeaderDeps
//	for _, v := range tx.Inputs {
//		res.Inputs = append(res.Inputs, &types.CellInput{
//			Since: v.Since,
//			PreviousOutput: &types.OutPoint{
//				TxHash: v.PreviousOutput.TxHash,
//				Index:  v.PreviousOutput.Index,
//			},
//		})
//	}
//	for _, v := range tx.Outputs {
//		tmp := types.CellOutput{
//			Capacity: v.Capacity,
//			Lock:     nil,
//			Type:     nil,
//		}
//		if v.Lock != nil {
//			tmp.Lock = &types.Script{
//				CodeHash: v.Lock.CodeHash,
//				HashType: v.Lock.HashType,
//				Args:     v.Lock.Args,
//			}
//		}
//		if v.Type != nil {
//			tmp.Type = &types.Script{
//				CodeHash: v.Type.CodeHash,
//				HashType: v.Type.HashType,
//				Args:     v.Type.Args,
//			}
//		}
//		res.Outputs = append(res.Outputs, &tmp)
//	}
//	res.OutputsData = tx.OutputsData
//	res.Witnesses = tx.Witnesses
//
//	return &res
//}
//
//type Transaction struct {
//	Version     uint          `json:"version"`
//	CellDeps    []*CellDep    `json:"cellDeps"`
//	HeaderDeps  []types.Hash  `json:"headerDeps"`
//	Inputs      []*CellInput  `json:"inputs"`
//	Outputs     []*CellOutput `json:"outputs"`
//	OutputsData [][]byte      `json:"outputsData"`
//	Witnesses   [][]byte      `json:"witnesses"`
//}
//
//type CellDep struct {
//	OutPoint *OutPoint     `json:"outPoint"`
//	DepType  types.DepType `json:"depType"`
//}
//
//type OutPoint struct {
//	TxHash types.Hash `json:"txHash"`
//	Index  uint       `json:"index"`
//}
//
//type CellInput struct {
//	Since          uint64    `json:"since"`
//	PreviousOutput *OutPoint `json:"previousOutput"`
//}
//
//type CellOutput struct {
//	Capacity uint64  `json:"capacity"`
//	Lock     *Script `json:"lock"`
//	Type     *Script `json:"type"`
//}
//
//type Script struct {
//	CodeHash types.Hash           `json:"codeHash"`
//	HashType types.ScriptHashType `json:"hashType"`
//	Args     []byte               `json:"args"`
//}
