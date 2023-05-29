package witness

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type ActionDataBuilder struct {
	ActionData *molecule.ActionData
	Action     common.DasAction
	Params     [][]byte
	ParamsStr  string
}

func (a *ActionDataBuilder) ActionBuyAccountInviterScript() (*molecule.Script, error) {
	if len(a.Params) != 3 {
		return nil, fmt.Errorf("len params err:[%d]", len(a.Params))
	}
	inviterScript, err := molecule.ScriptFromSlice(a.Params[0], true)
	if err != nil {
		return nil, fmt.Errorf("ScriptFromSlice err: %s", err.Error())
	}
	return inviterScript, nil
}

func (a *ActionDataBuilder) ActionBuyAccountChannelScript() (*molecule.Script, error) {
	if len(a.Params) != 3 {
		return nil, fmt.Errorf("len params err:[%d]", len(a.Params))
	}
	channelScript, err := molecule.ScriptFromSlice(a.Params[1], true)
	if err != nil {
		return nil, fmt.Errorf("ScriptFromSlice err: %s", err.Error())
	}
	return channelScript, nil
}

var ErrNotExistActionData = errors.New("not exist action data")

func ActionDataBuilderFromTx(tx *types.Transaction) (*ActionDataBuilder, error) {
	var resp ActionDataBuilder
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeActionData:
			if err := resp.ConvertToActionData(dataBys); err != nil {
				return false, fmt.Errorf("ConvertToActionData err: %s", err.Error())
			}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp.ActionData == nil {
		return nil, ErrNotExistActionData
	}
	return &resp, nil
}

func ActionDataBuilderFromWitness(wit []byte) (*ActionDataBuilder, error) {
	if len(wit) <= common.WitnessDasTableTypeEndIndex+1 {
		return nil, fmt.Errorf("action data len is invalid")
	} else if string(wit[0:common.WitnessDasCharLen]) != common.WitnessDas {
		return nil, fmt.Errorf("not a das data")
	}
	actionDataType := common.Bytes2Hex(wit[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])
	dataBys := wit[common.WitnessDasTableTypeEndIndex:]
	if actionDataType != common.ActionDataTypeActionData {
		return nil, fmt.Errorf("not a action data")
	}

	var resp ActionDataBuilder
	if err := resp.ConvertToActionData(dataBys); err != nil {
		return nil, fmt.Errorf("ConvertToActionData err: %s", err.Error())
	}
	return &resp, nil
}

func (a *ActionDataBuilder) ConvertToActionData(slice []byte) error {
	actionData, err := molecule.ActionDataFromSlice(slice, true)
	if err != nil {
		return fmt.Errorf("ActionDataFromSlice err: %s", err.Error())
	}
	a.ActionData = actionData
	a.Action = string(actionData.Action().RawData())
	if a.Action == common.DasActionBuyAccount {
		raw := actionData.Params().RawData()

		lenRaw := len(raw)
		inviterLockBytesLen, err := molecule.Bytes2GoU32(raw[:4])
		if err != nil {
			return fmt.Errorf("Bytes2GoU32 err: %s", err.Error())
		}
		inviterLockRaw := raw[:inviterLockBytesLen]
		channelLockRaw := raw[inviterLockBytesLen : lenRaw-1]

		a.Params = append(a.Params, inviterLockRaw)
		a.Params = append(a.Params, channelLockRaw)
		a.Params = append(a.Params, raw[lenRaw-1:lenRaw])
		a.ParamsStr = common.GetMaxHashLenParams(common.Bytes2Hex(inviterLockRaw)) + "," + common.GetMaxHashLenParams(common.Bytes2Hex(channelLockRaw)) + "," + common.Bytes2Hex(raw[lenRaw-1:lenRaw])
	} else if a.Action == common.DasActionLockAccountForCrossChain {
		raw := actionData.Params().RawData()
		if len(raw) == 17 {
			coinType := raw[:8]
			chainId := raw[8:16]
			a.Params = append(a.Params, coinType)
			a.Params = append(a.Params, chainId)
			a.Params = append(a.Params, raw[16:])
			a.ParamsStr = common.GetMaxHashLenParams(common.Bytes2Hex(coinType)) + "," + common.GetMaxHashLenParams(common.Bytes2Hex(chainId)) + "," + common.Bytes2Hex(raw[16:])
		}
	} else {
		a.Params = append(a.Params, actionData.Params().RawData())
		a.ParamsStr = common.Bytes2Hex(actionData.Params().RawData())
	}
	return nil
}

func GenActionDataWitnessV2(action common.DasAction, params []byte, signer string) ([]byte, error) {
	if action == "" {
		return nil, fmt.Errorf("action is nil")
	}
	params = append(params, common.Hex2Bytes(signer)...)

	actionBytes := molecule.GoString2MoleculeBytes(action)
	paramsBytes := molecule.GoBytes2MoleculeBytes(params)
	actionData := molecule.NewActionDataBuilder().Action(actionBytes).Params(paramsBytes).Build()

	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(common.ActionDataTypeActionData)...)
	tmp = append(tmp, actionData.AsSlice()...)
	return tmp, nil
}

func GenActionDataWitness(action common.DasAction, params []byte) ([]byte, error) {
	if action == "" {
		return nil, fmt.Errorf("action is nil")
	}
	if params == nil {
		params = []byte{}
	}
	if action == common.DasActionEditRecords {
		params = append(params, common.Hex2Bytes(common.ParamManager)...)
	} else if action == common.DasActionRenewAccount {
		params = []byte{}
	} else {
		params = append(params, common.Hex2Bytes(common.ParamOwner)...)
	}
	actionBytes := molecule.GoString2MoleculeBytes(action)
	paramsBytes := molecule.GoBytes2MoleculeBytes(params)
	actionData := molecule.NewActionDataBuilder().Action(actionBytes).Params(paramsBytes).Build()

	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(common.ActionDataTypeActionData)...)
	tmp = append(tmp, actionData.AsSlice()...)
	return tmp, nil
}

func GenActionDataWitnessV3(action common.DasAction, params []byte) ([]byte, error) {
	if action == "" {
		return nil, fmt.Errorf("action is nil")
	}
	if params == nil {
		params = []byte{}
	}
	actionBytes := molecule.GoString2MoleculeBytes(action)
	paramsBytes := molecule.GoBytes2MoleculeBytes(params)
	actionData := molecule.NewActionDataBuilder().Action(actionBytes).Params(paramsBytes).Build()

	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(common.ActionDataTypeActionData)...)
	tmp = append(tmp, actionData.AsSlice()...)
	return tmp, nil
}

func GenBuyAccountParams(inviterScript, channelScript *types.Script) []byte {
	iScript := molecule.ScriptDefault()
	if inviterScript != nil {
		iScript = molecule.CkbScript2MoleculeScript(inviterScript)
	}
	paramsInviter := iScript.AsSlice()
	cScript := molecule.ScriptDefault()
	if channelScript != nil {
		cScript = molecule.CkbScript2MoleculeScript(channelScript)
	}
	paramsChannel := cScript.AsSlice()
	return append(paramsInviter, paramsChannel...)
}
