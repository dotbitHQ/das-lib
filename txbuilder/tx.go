package txbuilder

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"strings"
)

func (d *DasTxBuilder) newTx() error {
	systemScriptCell, err := utils.NewSystemScripts(d.dasCore.Client())
	if err != nil {
		return err
	}
	baseTx := transaction.NewSecp256k1SingleSigTx(systemScriptCell)
	d.Transaction = baseTx
	return nil
}

func (d *DasTxBuilder) equalArgs(src, dst string) bool {
	if common.Has0xPrefix(src) {
		src = src[2:]
	}
	if common.Has0xPrefix(dst) {
		dst = dst[2:]
	}
	return src == dst
}
func (d *DasTxBuilder) addInputsForTx(inputs []*types.CellInput) error {
	if len(inputs) == 0 {
		return fmt.Errorf("inputs is nil")
	}
	startIndex := len(d.Transaction.Inputs)
	_, _, err := transaction.AddInputsForTransaction(d.Transaction, inputs)
	if err != nil {
		return fmt.Errorf("AddInputsForTransaction err: %s", err.Error())
	}

	var cellDepList []*types.CellDep
	for i, v := range inputs {
		if v == nil {
			return fmt.Errorf("input is nil")
		}
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return fmt.Errorf("getInputCell err: %s", err.Error())
		}

		if item.Cell.Output.Type != nil {
			if contractName, ok := core.DasContractByTypeIdMap[item.Cell.Output.Type.CodeHash.Hex()]; ok {
				dasContract, err := core.GetDasContractInfo(contractName)
				if err != nil {
					return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
				}
				cellDepList = append(cellDepList, dasContract.ToCellDep())
			}
		}

		if item.Cell.Output.Lock != nil &&
			item.Cell.Output.Lock.CodeHash.Hex() == transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH &&
			d.equalArgs(common.Bytes2Hex(item.Cell.Output.Lock.Args), d.serverArgs) {
			d.ServerSignGroup = append(d.ServerSignGroup, startIndex+i)
		}
		if item.Cell.Output.Lock != nil {
			if contractName, ok := core.DasContractByTypeIdMap[item.Cell.Output.Lock.CodeHash.Hex()]; ok {
				if dasContract, err := core.GetDasContractInfo(contractName); err != nil {
					return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
				} else {
					cellDepList = append(cellDepList, &types.CellDep{OutPoint: dasContract.OutPoint, DepType: types.DepTypeCode})
					if contractName == common.DasContractNameDispatchCellType {
						log.Info("addInputsForTx:", v.PreviousOutput.TxHash, v.PreviousOutput.Index)
						daf := core.DasAddressFormat{DasNetType: d.dasCore.NetType()}
						ownerHex, managerHex, _ := daf.ArgsToHex(item.Cell.Output.Lock.Args)
						oID, mID := ownerHex.DasAlgorithmId, managerHex.DasAlgorithmId
						oSo, _ := core.GetDasSoScript(oID.ToSoScriptType())
						mSo, _ := core.GetDasSoScript(mID.ToSoScriptType())
						cellDepList = append(cellDepList, oSo.ToCellDep())
						cellDepList = append(cellDepList, mSo.ToCellDep())
					}
				}
			}

		}
	}
	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *DasTxBuilder) addOutputsForTx(outputs []*types.CellOutput, outputsData [][]byte) error {
	lo := len(outputs)
	lod := len(outputsData)
	if lo == 0 || lod == 0 || lo != lod {
		return fmt.Errorf("outputs[%d] or outputDatas[%d]", lo, lod)
	}
	var cellDepList []*types.CellDep
	for i := 0; i < lo; i++ {
		output := outputs[i]
		outputData := outputsData[i]
		d.Transaction.Outputs = append(d.Transaction.Outputs, output)
		d.Transaction.OutputsData = append(d.Transaction.OutputsData, outputData)

		if output.Type == nil {
			continue
		}
		contractName, ok := core.DasContractByTypeIdMap[output.Type.CodeHash.Hex()]
		if !ok {
			continue
		}
		dasContract, err := core.GetDasContractInfo(contractName)
		if err != nil {
			return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
		}
		cellDepList = append(cellDepList, dasContract.ToCellDep())
	}

	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *DasTxBuilder) addWebauthnInfo() error {
	var cellDepList []*types.CellDep
	//Remove duplicate keylist witness
	keyListMap := make(map[string]bool, 0)
	for _, v := range d.Transaction.Inputs {
		if v == nil {
			return fmt.Errorf("input is nil")
		}
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return fmt.Errorf("getInputCell err: %s", err.Error())
		}
		if item == nil || item.Status == "unknown" || item.Cell == nil {
			log.Warn("addWebauthnInfo unknown:", v.PreviousOutput.TxHash, v.PreviousOutput.Index)
			continue
		}
		if args := item.Cell.Output.Lock.Args; len(args) > 0 {
			actionDataBuilder, err := witness.ActionDataBuilderFromTx(d.Transaction)
			if err != nil {
				log.Warn("witness.ActionDataBuilderFromTx err:", err.Error())
				return nil
			}
			log.Info("args: ", item.Cell.Output.Lock.Args)
			ownerHex, managerHex, err := d.dasCore.Daf().ArgsToHex(item.Cell.Output.Lock.Args)
			if err != nil && !strings.Contains(err.Error(), "len(args) error") {
				return fmt.Errorf("ArgsToHex err: %s", err.Error())
			}
			log.Info("actionDataBuilder.Params: ", actionDataBuilder.Params)
			//Obtain the role of owner or manager for the current signature verification through the action witness parameter
			if len(actionDataBuilder.Params) == 0 {
				continue
			}
			verifyRole := ownerHex
			if len(actionDataBuilder.Params[0]) > 0 {
				if actionDataBuilder.Params[0][0] == 0 {
					verifyRole = ownerHex
				} else {
					verifyRole = managerHex
				}
			}
			log.Info("verifyRole :", verifyRole)
			lockArgs, err := d.dasCore.Daf().HexToArgs(verifyRole, verifyRole)
			if err != nil {
				return fmt.Errorf("HexToArgs err: %s", err.Error())
			}
			log.Info("lockArgs[0] ", lockArgs[0])
			if common.ChainType(lockArgs[0]) == common.ChainTypeWebauthn {
				keyListCfgCell, err := core.GetDasContractInfo(common.DasKeyListCellType)
				if err != nil {
					return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
				}
				cellDepList = append(cellDepList, keyListCfgCell.ToCellDep())

				//exclude create and update keylist tx (balance cell type is nil)
				if item.Cell.Output.Type == nil || !keyListCfgCell.IsSameTypeId(item.Cell.Output.Type.CodeHash) {
					//select args=owner owner or  args=manager manager keylistCell
					log.Info("is webauthn")
					cell, err := d.dasCore.GetKeyListCell(lockArgs)
					if err != nil {
						log.Warn("dasCore.GetKeyListCell err: ", err.Error())
						continue
					}
					if cell != nil {
						if _, ok := keyListMap[cell.OutPoint.TxHash.Hex()]; ok {
							continue
						}
						cellDepList = append(cellDepList, &types.CellDep{
							OutPoint: cell.OutPoint,
							DepType:  types.DepTypeCode,
						})

						keyListConfigTx, err := d.dasCore.Client().GetTransaction(d.ctx, cell.OutPoint.TxHash)
						if err != nil {
							return err
						}
						webAuthnKeyListConfigBuilder, err := witness.WebAuthnKeyListDataBuilderFromTx(keyListConfigTx.Transaction, common.DataTypeNew)
						if err != nil {
							return err
						}
						webAuthnKeyListConfigBuilder.DataEntityOpt.AsSlice()
						tmp := webAuthnKeyListConfigBuilder.DeviceKeyListCellData.AsSlice()
						keyListWitness := witness.GenDasDataWitnessWithByte(common.ActionDataTypeKeyListCfgCellData, tmp)
						d.otherWitnesses = append(d.otherWitnesses, keyListWitness)
						keyListMap[cell.OutPoint.TxHash.Hex()] = true
						log.Info("add key list cell :  ", cell.OutPoint.TxHash.Hex())
					} else {
						log.Warn("dasCore.GetKeyListCell cell is nil")
					}
				}
			}
		}
	}
	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *DasTxBuilder) checkTxWitnesses() error {
	if len(d.Transaction.Witnesses) == 0 {
		return fmt.Errorf("witness is nil")
	}
	lenI := len(d.Transaction.Inputs)
	lenW := len(d.Transaction.Witnesses)
	if lenW < lenI {
		return fmt.Errorf("len witness[%d]<len inputs[%d]", lenW, lenI)
	} else if lenW > lenI {
		for i := lenI; i < lenW; i++ {
			if _, err := witness.ActionDataBuilderFromWitness(d.Transaction.Witnesses[i]); err == nil {
				return nil
			}
		}
	}
	return fmt.Errorf("action data check fail")
}

func (d *DasTxBuilder) addCellDepListIntoMapCellDep(cellDepList []*types.CellDep) {
	for i, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		d.mapCellDep[k] = cellDepList[i]
	}
}

func (d *DasTxBuilder) addMapCellDepWitnessForBaseTx(cellDepList []*types.CellDep) error {
	configCellMain, err := core.GetDasConfigCellInfo(common.ConfigCellTypeArgsMain)
	if err != nil {
		return fmt.Errorf("GetDasConfigCellInfo err: %s", err.Error())
	}
	cellDepList = append(cellDepList, configCellMain.ToCellDep())

	contractEip712Lib, err := core.GetDasContractInfo(common.DASContractNameEip712LibCellType)
	if err != nil {
		log.Warn("core.GetDasContractInfo 712 err: ", err.Error())
		//return fmt.Errorf("GetDasContractInfo err: %s", err.Error())
	} else {
		cellDepList = append(cellDepList, contractEip712Lib.ToCellDep())
	}

	soEth, err := core.GetDasSoScript(common.SoScriptTypeEth)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeEth err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soEth.ToCellDep())
	}
	soTron, err := core.GetDasSoScript(common.SoScriptTypeTron)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeTron err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soTron.ToCellDep())
	}
	soEd25519, err := core.GetDasSoScript(common.SoScriptTypeEd25519)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeEd25519 err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soEd25519.ToCellDep())
	}
	soCkbMulti, err := core.GetDasSoScript(common.SoScriptTypeCkbMulti)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeCkbMulti err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soCkbMulti.ToCellDep())
	}
	soCkbSingle, err := core.GetDasSoScript(common.SoScriptTypeCkbSingle)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeCkbSingle err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soCkbSingle.ToCellDep())
	}
	soDoge, err := core.GetDasSoScript(common.SoScriptTypeDogeCoin)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeDogeCoin err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soDoge.ToCellDep())
	}
	soBtc, err := core.GetDasSoScript(common.SoScriptBitcoin)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptBitcoin err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, soBtc.ToCellDep())
	}
	webauthn, err := core.GetDasSoScript(common.SoScriptWebauthn)
	if err != nil {
		log.Warn("GetDasSoScript SoScriptTypeWebauthn err: ", err.Error())
	} else {
		cellDepList = append(cellDepList, webauthn.ToCellDep())
	}

	tmpMap := make(map[string]bool)
	var tmpCellDeps []*types.CellDep
	for _, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		if _, ok := tmpMap[k]; ok {
			continue
		}
		tmpMap[k] = true
		tmpCellDeps = append(tmpCellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
		if _, ok := core.DasConfigCellByTxHashMap.Load(v.OutPoint.TxHash.Hex()); !ok {
			continue
		}
		if res, err := d.dasCore.Client().GetTransaction(d.ctx, v.OutPoint.TxHash); err != nil {
			return fmt.Errorf("GetTransaction err: %s [%s]", err.Error(), k)
		} else {
			d.Transaction.Witnesses = append(d.Transaction.Witnesses, res.Transaction.Witnesses[len(res.Transaction.Witnesses)-1])
		}
	}
	if len(tmpCellDeps) > 0 {
		d.Transaction.CellDeps = append(tmpCellDeps, d.Transaction.CellDeps...)
	}

	for k, v := range d.mapCellDep {
		if _, ok := tmpMap[k]; ok {
			continue
		}
		d.Transaction.CellDeps = append(d.Transaction.CellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
		if _, ok := core.DasConfigCellByTxHashMap.Load(v.OutPoint.TxHash.Hex()); !ok {
			continue
		}
		if res, err := d.dasCore.Client().GetTransaction(d.ctx, v.OutPoint.TxHash); err != nil {
			return fmt.Errorf("GetTransaction err: %s [%s]", err.Error(), k)
		} else {
			d.Transaction.Witnesses = append(d.Transaction.Witnesses, res.Transaction.Witnesses[len(res.Transaction.Witnesses)-1])
		}
	}
	if len(d.otherWitnesses) > 0 {
		d.Transaction.Witnesses = append(d.Transaction.Witnesses, d.otherWitnesses...)
	}
	return nil
}

func (d *DasTxBuilder) SendTransactionWithCheck(needCheck bool) (*types.Hash, error) {
	if needCheck {
		err := d.checkTxBeforeSend()
		if err != nil {
			return nil, fmt.Errorf("checkTxBeforeSend err: %s", err.Error())
		}
	}

	err := d.serverSignTx()
	if err != nil {
		return nil, fmt.Errorf("remoteSignTx err: %s", err.Error())
	}

	log.Info("before sent: ", d.TxString())
	txHash, err := d.dasCore.Client().SendTransactionNoneValidation(d.ctx, d.Transaction)
	if err != nil {
		return nil, fmt.Errorf("SendTransaction err: %v", err)
	}
	log.Info("SendTransaction success:", txHash.Hex())
	return txHash, nil
}

func (d *DasTxBuilder) SendTransaction() (*types.Hash, error) {
	return d.SendTransactionWithCheck(true)
}

func (d *DasTxBuilder) checkTxBeforeSend() error {
	// check total num of inputs and outputs
	if len(d.Transaction.Inputs)+len(d.Transaction.Outputs) > 9000 {
		return fmt.Errorf("checkTxBeforeSend, failed len of inputs: %d, ouputs: %d", len(d.Transaction.Inputs), len(d.Transaction.Outputs))
	}
	// check tx fee < 1 CKB
	totalCapacityFromInputs, err := d.getCapacityFromInputs()
	if err != nil {
		return err
	}
	totalCapacityFromOutputs := d.Transaction.OutputsCapacity()
	txFee := totalCapacityFromInputs - totalCapacityFromOutputs
	log.Info("checkTxBeforeSend:", totalCapacityFromInputs, totalCapacityFromOutputs, txFee)
	if totalCapacityFromInputs <= totalCapacityFromOutputs || txFee >= common.OneCkb {
		return fmt.Errorf("checkTxBeforeSend failed, txFee: %d totalCapacityFromInputs: %d totalCapacityFromOutputs: %d", txFee, totalCapacityFromInputs, totalCapacityFromOutputs)
	}

	// check witness format
	err = d.checkTxWitnesses()
	if err != nil {
		log.Warn("checkTxWitnesses:", err.Error())
		//return err
	}
	// check the occupied capacity
	for i, cell := range d.Transaction.Outputs {
		occupied := cell.OccupiedCapacity(d.Transaction.OutputsData[i])
		if cell.Capacity < occupied {
			return fmt.Errorf("checkTxBeforeSend occupied capacity failed, occupied: %d capacity: %d index: %d", occupied, cell.Capacity, i)
		}
	}
	log.Info("check success before sent")
	return nil
}

func (d *DasTxBuilder) getCapacityFromInputs() (uint64, error) {
	total := uint64(0)
	for _, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return 0, fmt.Errorf("getInputCell err: %s", err.Error())
		}
		total += item.Cell.Output.Capacity
	}
	return total, nil
}
