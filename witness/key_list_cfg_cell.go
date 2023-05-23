package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type WebauthnKey struct {
	MinAlgId uint8  `json:"min_alg_id"`
	SubAlgId uint8  `json:"sub_alg_id"`
	Cid      string `json:"cid"`
	PubKey   string `json:"pubkey"`
}
type WebAuthnKeyListDataBuilder struct {
	WebauthnKeyList []WebauthnKey
	Index           uint32
	Version         uint32

	//WebAuthnKeyListData *molecule.WebAuthnKey
	WebAuthnKeyListData *molecule.DeviceKeyList
	DataEntityOpt       *molecule.DataEntityOpt
}

type WebauchnKeyListCellParam struct {
	AddWebauthnKeyList WebauthnKey //keylist need to be added
	Action             string
	OldIndex           uint32
	NewIndex           uint32
}

func WebAuthnKeyListDataBuilderFromTx(tx *types.Transaction, dataType common.DataType) (*WebAuthnKeyListDataBuilder, error) {
	//var respList = make([]*WebAuthnKeyListDataBuilder, 0)
	var resp *WebAuthnKeyListDataBuilder
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeKeyListCfgCell:
			//var resp WebAuthnKeyListDataBuilder
			dataEntityOpt, dataEntity, err := getDataEntityOpt(dataBys, dataType)
			if err != nil {
				if err == ErrDataEntityOptIsNil {
					//log.Warn("getDataEntityOpt err:", err.Error(), tx.Hash)
					return true, nil
				}
				return false, fmt.Errorf("getDataEntityOpt err: %s", err.Error())
			}
			resp.DataEntityOpt = dataEntityOpt

			version, err := molecule.Bytes2GoU32(dataEntity.Version().RawData())
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			resp.Version = version
			index, err := molecule.Bytes2GoU32(dataEntity.Index().RawData())
			if err != nil {
				return false, fmt.Errorf("get index err: %s", err.Error())
			}
			resp.Index = index

			webauthnKeyListCellData, err := molecule.DeviceKeyListFromSlice(dataEntity.Entity().RawData(), true)
			if err != nil {
				return false, fmt.Errorf("WebauthnKeyListCellDataFromSlice err : %s", err.Error())
			}
			resp.WebAuthnKeyListData = webauthnKeyListCellData
			return true, nil
			//respList = append(respList,&resp)
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if resp == nil {
		return nil, fmt.Errorf("not exist new keylistcfg cell")
	}
	return resp, nil

}

func (w *WebAuthnKeyListDataBuilder) GenWitness(p *WebauchnKeyListCellParam) (witness []byte, accData []byte, err error) {
	switch p.Action {
	case common.DasActionCreateKeyList:

	case common.DasActionUpdateKeyList:
		oldDataEntityOpt := w.getOldDataEntityOpt(p)
		//TODO 是否去重，如果以及存在是否还继续追加。
		oldList := ConvertToWebauthnKeyList(w.WebAuthnKeyListData)
		newList := append(oldList, p.AddWebauthnKeyList)
		webauthnKeyLists, err := ConvertToWebKeyList(newList)
		if err != nil {
			return
		}
		w.WebAuthnKeyListData = webauthnKeyLists
		newBuilder := w.WebAuthnKeyListData.AsBuilder()
		newWebauthnKeyData := newBuilder.Build()
		newWebauthnKeyDataBytes := molecule.GoBytes2MoleculeBytes(newWebauthnKeyData.AsSlice())
		newDataEntity := molecule.NewDataEntityBuilder().Entity(newWebauthnKeyDataBytes).
			Version(molecule.GoU32ToMoleculeU32(w.Version)).Index(molecule.GoU32ToMoleculeU32(p.NewIndex)).Build()
		newDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(newDataEntity).Build()
		tmp := molecule.NewDataBuilder().Old(*oldDataEntityOpt).New(newDataEntityOpt).Build()
		witness := GenDasDataWitness(common.ActionDataTypeKeyListCfgCell, &tmp)
		return witness, common.Blake2b(newWebauthnKeyData.AsSlice()), nil
	}
	return nil, nil, fmt.Errorf("not exist action [%s]", p.Action)
}

func (w *WebAuthnKeyListDataBuilder) getOldDataEntityOpt(p *WebauchnKeyListCellParam) *molecule.DataEntityOpt {
	var oldDataEntity molecule.DataEntity
	oldWebauthnDataBytes := molecule.GoBytes2MoleculeBytes(w.WebAuthnKeyListData.AsSlice())
	oldDataEntity = molecule.NewDataEntityBuilder().Entity(oldWebauthnDataBytes).
		Version(molecule.GoU32ToMoleculeU32(w.Version)).Index(molecule.GoU32ToMoleculeU32(p.OldIndex)).Build()
	oldDataEntityOpt := molecule.NewDataEntityOptBuilder().Set(oldDataEntity).Build()
	return &oldDataEntityOpt
}

func ConvertToWebauthnKeyList(keyLists *molecule.DeviceKeyList) []WebauthnKey {
	var keyList []WebauthnKey
	for index, lenKeyLists := uint(0), keyLists.Len(); index < lenKeyLists; index++ {
		value := keyLists.Get(index)
		alg, _ := molecule.Bytes2GoU8(value.Alg().RawData())
		keyList = append(keyList, WebauthnKey{
			MinAlgId: alg,
			SubAlgId: alg,
			Cid:      string(value.Cid().RawData()),
			PubKey:   string(value.Pubkey().RawData()),
		})
	}
	return keyList
}

func ConvertToWebKeyList(keyLists []WebauthnKey) (*molecule.DeviceKeyList, error) {
	keyListsBuilder := molecule.NewDeviceKeyListBuilder()

	for _, v := range keyLists {
		cid, err := molecule.GoString2MoleculeByte10(v.Cid)
		if err != nil {
			return nil, err
		}
		pubKey, err := molecule.GoString2MoleculeByte10(v.Cid)
		if err != nil {
			return nil, err
		}
		keyListBuilder := molecule.NewDeviceKeyBuilder().
			Alg(molecule.GoU8ToMoleculeU8(v.MinAlgId)).
			Cid(cid).
			Pubkey(pubKey)
		keyListsBuilder.Push(keyListBuilder.Build())
	}
	keylist := keyListsBuilder.Build()
	return &keylist, nil
}
