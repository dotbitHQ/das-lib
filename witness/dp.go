package witness

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type DPOrderInfo struct {
	OrderId string   `json:"order_id"`
	Action  DPAction `json:"action"`
}

type DPAction string

const (
	DPActionDefault         DPAction = ""
	DPActionMint            DPAction = "mint"
	DPActionBurn            DPAction = "burn"
	DPActionTransfer        DPAction = "transfer"
	DPActionTransferDeposit DPAction = "transfer_deposit"
	DPActionTransferRefund  DPAction = "transfer_refund"
	DPActionTransferTLDID   DPAction = "transfer_tldid"
	DPActionTransferSLDID   DPAction = "transfer_sldid"
	DPActionTransferAuction DPAction = "transfer_auction"
	DPActionTransferCoupon  DPAction = "transfer_coupon"
)

func DPOrderInfoFromTx(tx *types.Transaction) (DPOrderInfo, error) {
	var dpOrderInfo DPOrderInfo
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeDPOrderInfo:
			var e error
			dpOrderInfo, e = ConvertDPOrderInfoWitness(dataBys)
			if e != nil {
				return false, fmt.Errorf("ConvertDPOrderInfoWitness err: %s", e.Error())
			}
		}
		return true, nil
	})
	if err != nil {
		return dpOrderInfo, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	return dpOrderInfo, fmt.Errorf("not exist account cell")
}

func ConvertDPOrderInfoWitness(dataBys []byte) (DPOrderInfo, error) {
	var dpOrderInfo DPOrderInfo
	orderInfo, err := molecule.OrderInfoFromSlice(dataBys, true)
	if err != nil {
		return dpOrderInfo, fmt.Errorf("OrderInfoFromSlice err: %s", err.Error())
	}
	bys := orderInfo.Memo().RawData()
	if err := json.Unmarshal(bys, &dpOrderInfo); err != nil {
		return dpOrderInfo, fmt.Errorf("json.Unmarshal err: %s", err.Error())
	}
	return dpOrderInfo, nil
}

func GenDPOrderInfoWitness(info DPOrderInfo) (witness []byte, data []byte, err error) {
	bys, err := json.Marshal(&info)
	if err != nil {
		return nil, nil, fmt.Errorf("json.Marshal err: %s", err.Error())
	}
	moleculeData := molecule.NewOrderInfoBuilder().Memo(molecule.GoBytes2MoleculeBytes(bys)).Build()
	data = moleculeData.AsSlice()
	witness = GenDasDataWitnessWithByte(common.ActionDataTypeDPOrderInfo, data)
	return witness, data, nil
}

type DPData struct {
	Value uint64
}

func ConvertBysToDPData(data []byte) (DPData, error) {
	var res DPData
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)

	dataLen, _ = molecule.Bytes2GoU32(data[index : index+indexLen])
	value := data[index+indexLen : index+indexLen+dataLen]
	res.Value, _ = molecule.Bytes2GoU64(value)

	index = index + indexLen + dataLen
	return res, nil
}

func ConvertDPDataToBys(data DPData) ([]byte, error) {
	var dataBys []byte
	valueBys := molecule.GoU64ToMoleculeU64(data.Value)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(valueBys.RawData())))...)
	dataBys = append(dataBys, valueBys.RawData()...)
	return dataBys, nil
}
