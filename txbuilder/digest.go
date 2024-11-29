package txbuilder

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sort"
)

type SignData struct {
	SignType common.DasAlgorithmId `json:"sign_type"`
	SignMsg  string                `json:"sign_msg"`
}

func (d *DasTxBuilder) GenerateMultiSignDigest(group []int, firstN uint8, signatures [][]byte, sortArgsList [][]byte) ([]byte, error) {
	if len(group) == 0 {
		return nil, fmt.Errorf("group is nil")
	}

	wa := GenerateMultiSignWitnessArgs(firstN, signatures, sortArgsList)
	data, err := wa.Serialize()
	if err != nil {
		return nil, err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := d.Transaction.ComputeHash()
	if err != nil {
		return nil, err
	}
	message := append(hash.Bytes(), length...)
	message = append(message, data...)

	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
		}
	}

	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}

	message, err = blake2b.Blake256(message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (d *DasTxBuilder) GenerateDigestListFromTx(skipGroups []int) ([]SignData, error) {
	skipGroups = append(skipGroups, d.ServerSignGroup...)
	groups, err := d.getGroupsFromTx()
	if err != nil {
		return nil, fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	log.Info("groups:", len(groups), groups)
	var digestList []SignData
	for _, group := range groups {
		if digest, err := d.generateDigestByGroup(group, skipGroups); err != nil {
			return nil, fmt.Errorf("generateDigestByGroup err: %s", err.Error())
		} else {
			digestList = append(digestList, digest)
		}
	}
	return digestList, nil
}

func (d *DasTxBuilder) getGroupsFromTx() ([][]int, error) {
	//input
	//code1-1
	//code1-2
	//code2-1
	var tmpMapForGroup = make(map[string][]int)
	var sortList []string
	for i, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
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
	//sortList = [code1,code2]
	//tmpMapForGroup = [
	//	code1=>[0,1]
	//	code2=>[2]
	//]
	sort.Strings(sortList)
	var list [][]int
	for _, v := range sortList {
		item, _ := tmpMapForGroup[v]
		list = append(list, item)
	}
	//list = [[0,1], [2]]
	return list, nil
}

func (d *DasTxBuilder) generateDigestByGroup(group []int, skipGroups []int) (SignData, error) {
	var signData = SignData{}
	if group == nil || len(group) < 1 {
		return signData, fmt.Errorf("invalid param")
	}
	// check AlgorithmId
	item, err := d.getInputCell(d.Transaction.Inputs[group[0]].PreviousOutput)
	if err != nil {
		return signData, fmt.Errorf("getInputCell err: %s", err.Error())
	}
	has712, action := false, ""

	//didType, err := core.GetDasContractInfo(common.DasContractNameDidCellType)
	//if err != nil {
	//	return signData, fmt.Errorf("core.GetDasContractInfo err: %s", err.Error())
	//}
	//else if dasLock.IsSameTypeId(item.Cell.Output.Lock.CodeHash) && item.Cell.Output.Type != nil && didType.IsSameTypeId(item.Cell.Output.Type.CodeHash) {
	//	log.Info("generateDigestByGroup did cell with daslock")
	//	daf := core.DasAddressFormat{DasNetType: d.dasCore.NetType()}
	//	ownerHex, _, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
	//	ownerAlgorithmId := ownerHex.DasAlgorithmId
	//
	//	signData.SignType = ownerAlgorithmId
	//	if signData.SignType == common.DasAlgorithmIdEth712 {
	//		has712 = true
	//	}
	//}

	dasLock, err := core.GetDasContractInfo(common.DasContractNameDispatchCellType)
	if err != nil {
		return signData, fmt.Errorf("core.GetDasContractInfo err: %s", err.Error())
	} else if dasLock.IsSameTypeId(item.Cell.Output.Lock.CodeHash) {
		daf := core.DasAddressFormat{DasNetType: d.dasCore.NetType()}
		ownerHex, managerHex, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
		ownerAlgorithmId, managerAlgorithmId := ownerHex.DasAlgorithmId, managerHex.DasAlgorithmId

		signData.SignType = ownerAlgorithmId
		actionDataBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction)
		if err != nil {
			return signData, fmt.Errorf("ActionDataBuilderFromTx err: %s", err.Error())
		}
		if actionDataBuilder.ParamsStr == common.ParamManager {
			signData.SignType = managerAlgorithmId
		}

		actionBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction)
		//actionBuilder.Params
		if err != nil {
			if err != witness.ErrNotExistActionData {
				return signData, fmt.Errorf("witness.ActionDataBuilderFromTx err: %s", err.Error())
			}
		} else {
			action = actionBuilder.Action
			switch actionBuilder.Action {
			case common.DasActionEditRecords:
				signData.SignType = managerAlgorithmId
			case common.DasActionEnableSubAccount, common.DasActionCreateSubAccount,
				common.DasActionConfigSubAccountCustomScript, common.DasActionConfigSubAccount:
				if signData.SignType == common.DasAlgorithmIdEth712 {
					signData.SignType = common.DasAlgorithmIdEth
				}
			case common.DasActionRevokeApproval:
				signData.SignType = common.DasAlgorithmIdEth
			}
			// 712
			switch actionBuilder.Action {
			case common.DasActionEditManager, common.DasActionEditRecords,
				common.DasActionTransferAccount, common.DasActionTransfer,
				common.DasActionWithdrawFromWallet, common.DasActionStartAccountSale,
				common.DasActionEditAccountSale, common.DasActionCancelAccountSale,
				common.DasActionBuyAccount, common.DasActionDeclareReverseRecord,
				common.DasActionRedeclareReverseRecord, common.DasActionRetractReverseRecord,
				common.DasActionMakeOffer, common.DasActionEditOffer, common.DasActionCancelOffer,
				common.DasActionAcceptOffer, common.DasActionLockAccountForCrossChain,
				common.DasActionCreateApproval, common.DasActionDelayApproval, common.DasActionFulfillApproval,
				common.DasActionMintDP, common.DasActionTransferDP, common.DasActionBurnDP, common.DasActionBidExpiredAccountAuction:
				has712 = true
			}
		}
		// gen digest
		log.Warn("generateDigestByGroup:", len(group), group, action, has712, actionDataBuilder.ParamsStr)
	} else if item.Cell.Output.Lock.CodeHash.Hex() == transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH &&
		d.equalArgs(common.Bytes2Hex(item.Cell.Output.Lock.Args), d.serverArgs) {
		signData.SignType = common.DasAlgorithmIdCkb
	} else {
		signData.SignType = common.DasAlgorithmIdAnyLock
	}

	emptyWitnessArg := types.WitnessArgs{
		Lock:       make([]byte, 65),
		InputType:  nil,
		OutputType: nil,
	}
	if signData.SignType == common.DasAlgorithmIdDogeChain {
		emptyWitnessArg.Lock = make([]byte, 66)
	} else if signData.SignType == common.DasAlgorithmIdEth712 && has712 {
		emptyWitnessArg.Lock = make([]byte, 105)
	} else if signData.SignType == common.DasAlgorithmIdWebauthn {
		emptyWitnessArg.Lock = make([]byte, 800)
	}
	data, err := emptyWitnessArg.Serialize()
	if err != nil {
		return signData, err
	}

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := d.Transaction.ComputeHash()
	if err != nil {
		return signData, err
	}
	//fmt.Println("tx_hash:", hash.Hex())

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	//fmt.Println("init witness:", common.Bytes2Hex(message))
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
			//fmt.Println("add group other witness:", common.Bytes2Hex(message))
		}
	}
	//fmt.Println("add group other witness:", common.Bytes2Hex(message))
	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}
	//fmt.Println("add other witness:", common.Bytes2Hex(message))

	message, err = blake2b.Blake256(message)
	if err != nil {
		return signData, err
	}

	signData.SignMsg = common.Bytes2Hex(message)
	if signData.SignType == common.DasAlgorithmIdWebauthn {
		signData.SignMsg = signData.SignMsg[2:]
	}
	//03 04 07 sign string
	if signData.SignType == common.DasAlgorithmIdEth || signData.SignType == common.DasAlgorithmIdDogeChain || signData.SignType == common.DasAlgorithmIdTron || signData.SignType == common.DasAlgorithmIdWebauthn {
		signData.SignMsg = common.DotBitPrefix + hex.EncodeToString(message)
	}
	log.Info("digest:", signData.SignType, signData.SignMsg)

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

func (d *DasTxBuilder) getInputCell(o *types.OutPoint) (*types.CellWithStatus, error) {
	if o == nil {
		return nil, fmt.Errorf("OutPoint is nil")
	}
	key := fmt.Sprintf("%s-%d", o.TxHash.Hex(), o.Index)
	if item, ok := d.MapInputsCell[key]; ok {
		if item.Cell != nil && item.Cell.Output != nil && item.Cell.Output.Lock != nil {
			return item, nil
		}
	}
	if cell, err := d.dasCore.Client().GetLiveCell(d.ctx, o, true); err != nil {
		return nil, fmt.Errorf("GetLiveCell err: %s", err.Error())
	} else if cell.Cell.Output.Lock != nil {
		d.MapInputsCell[key] = cell
		return cell, nil
	} else {
		log.Warn("GetLiveCell:", key, cell.Status)
		if !d.notCheckInputs {
			return nil, fmt.Errorf("cell [%s] not live", key)
		} else {
			d.MapInputsCell[key] = cell
			return cell, nil
		}
	}
}
