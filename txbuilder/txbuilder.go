package txbuilder

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/sign"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sort"
)

var log = mylog.NewLogger("txbuilder", mylog.LevelDebug)

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

	if err := d.addMapCellDepWitnessForBaseTx(p.CellDeps); err != nil {
		return fmt.Errorf("addMapCellDepWitnessForBaseTx err: %s", err.Error())
	}

	for _, v := range p.HeadDeps {
		d.Transaction.HeaderDeps = append(d.Transaction.HeaderDeps, v)
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

func generateDigestByGroup(cli rpc.Client, group []int, skipGroups []int, Tx *types.Transaction) (SignData, error) {
	var signData = SignData{}
	if group == nil || len(group) < 1 {
		return signData, fmt.Errorf("invalid param")
	}

	digest := ""
	data, err := transaction.EmptyWitnessArg.Serialize()
	if err != nil {
		return signData, err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := Tx.ComputeHash()
	if err != nil {
		return signData, err
	}

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = Tx.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
		}
	}
	// hash witnesses which do not in any input group
	for _, wit := range Tx.Witnesses[len(Tx.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}

	message, err = blake2b.Blake256(message)
	if err != nil {
		return signData, err
	}
	digest = common.Bytes2Hex(message)
	item, err := getInputCell(cli, Tx.Inputs[group[0]].PreviousOutput)
	if err != nil {
		return signData, fmt.Errorf("getInputCell err: %s", err.Error())
	}

	daf := core.DasAddressFormat{DasNetType: common.DasNetTypeTestnet2}
	ownerHex, managerHex, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
	ownerAlgorithmId, managerAlgorithmId := ownerHex.DasAlgorithmId, managerHex.DasAlgorithmId

	signData.SignMsg = digest
	signData.SignType = ownerAlgorithmId

	actionBuilder, err := witness.ActionDataBuilderFromTx(Tx)
	if err != nil {
		if err != witness.ErrNotExistActionData {
			return signData, fmt.Errorf("witness.ActionDataBuilderFromTx err: %s", err.Error())
		}
	} else {
		switch actionBuilder.Action {
		case common.DasActionEditRecords:
			signData.SignType = managerAlgorithmId
		case common.DasActionEnableSubAccount, common.DasActionCreateSubAccount,
			common.DasActionConfigSubAccountCustomScript:
			if signData.SignType == common.DasAlgorithmIdEth712 {
				signData.SignType = common.DasAlgorithmIdEth
			}
		}
	}

	if signData.SignType == common.DasAlgorithmIdTron {
		signData.SignMsg += "04" // fix tron sign
	}

	// skip useless signature
	if len(skipGroups) != 0 {
		skip := false
		for i := range group {
			for j := range skipGroups {
				if group[i] == skipGroups[j] {
					skip = true
					break
				}
			}
			if skip {
				break
			}
		}
		if skip {
			signData.SignMsg = ""
		}
	}
	return signData, nil
}
func getInputCell(cli rpc.Client, o *types.OutPoint) (*types.CellWithStatus, error) {
	if o == nil {
		return nil, fmt.Errorf("OutPoint is nil")
	}
	if cell, err := cli.GetLiveCell(context.Background(), o, true); err != nil {
		return nil, fmt.Errorf("GetLiveCell err: %s", err.Error())
	} else if cell.Cell.Output.Lock != nil {
		return cell, nil
	} else {
		return cell, nil
	}
}

func getGroupsFromTx(cli rpc.Client, Tx *types.Transaction) ([][]int, error) {
	var tmpMapForGroup = make(map[string][]int)
	var sortList []string
	for i, v := range Tx.Inputs {
		item, err := getInputCell(cli, v.PreviousOutput)
		if err != nil {
			return nil, fmt.Errorf("getInputCell err: %s", err.Error())
		}

		cellHash, err := item.Cell.Output.Lock.Hash()
		if err != nil {
			return nil, fmt.Errorf("inputs lock to hash err: %s", err.Error())
		}
		indexList, okTmp := tmpMapForGroup[cellHash.String()]
		if !okTmp {
			sortList = append(sortList, cellHash.String())
		}
		indexList = append(indexList, i)
		tmpMapForGroup[cellHash.String()] = indexList
	}
	sort.Strings(sortList)
	var list [][]int
	for _, v := range sortList {
		item, _ := tmpMapForGroup[v]
		list = append(list, item)
	}
	return list, nil
}

func GenerateDigestListFromTx(cli rpc.Client, txJson string, skipGroups []int) ([]SignData, error) {
	Tx, err := rpc.TransactionFromString(txJson)
	if err != nil {
		return nil, err

	}
	groups, err := getGroupsFromTx(cli, Tx)
	if err != nil {
		return nil, err
	}
	var digestList []SignData
	for _, group := range groups {
		if digest, err := generateDigestByGroup(cli, group, skipGroups, Tx); err != nil {
			return nil, err
		} else {
			digestList = append(digestList, digest)
		}
	}
	return digestList, nil
}
