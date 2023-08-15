package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type SubAccountNewBuilder struct{}

// === SubAccountMintSign ===

type SubAccountMintSignVersion = uint32
type SubAccountVersion = uint32

const (
	SubAccountMintSignVersion1 SubAccountMintSignVersion = 1

	SubAccountVersion1      SubAccountVersion = 1
	SubAccountVersion2      SubAccountVersion = 2
	SubAccountVersionLatest                   = SubAccountVersion2
)

type SubAccountMintSign struct {
	versionBys   []byte
	expiredAtBys []byte

	Version            SubAccountMintSignVersion
	Signature          []byte
	SignRole           []byte
	ExpiredAt          uint64
	AccountListSmtRoot []byte
}

func (s *SubAccountNewBuilder) ConvertSubAccountMintSignFromBytes(dataBys []byte) (*SubAccountMintSign, error) {
	var res SubAccountMintSign
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.versionBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.Version, _ = molecule.Bytes2GoU32(res.versionBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.Signature = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.SignRole = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.expiredAtBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.ExpiredAt, _ = molecule.Bytes2GoU64(res.expiredAtBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.AccountListSmtRoot = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	return &res, nil
}

func (s *SubAccountMintSign) GenSubAccountMintSignBytes() (dataBys []byte) {
	versionBys := molecule.GoU32ToMoleculeU32(s.Version)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
	dataBys = append(dataBys, versionBys.RawData()...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.Signature)))...)
	dataBys = append(dataBys, s.Signature...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.SignRole)))...)
	dataBys = append(dataBys, s.SignRole...)

	expiredAtBys := molecule.GoU64ToMoleculeU64(s.ExpiredAt)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(expiredAtBys.RawData())))...)
	dataBys = append(dataBys, expiredAtBys.RawData()...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.AccountListSmtRoot)))...)
	dataBys = append(dataBys, s.AccountListSmtRoot...)

	return
}

func (s *SubAccountMintSign) GenWitness() []byte {
	return GenDasDataWitnessWithByte(common.ActionDataTypeSubAccountMintSign, s.GenSubAccountMintSignBytes())
}

func (s *SubAccountMintSign) GenWitnessWithAction(action common.ActionDataType) []byte {
	return GenDasDataWitnessWithByte(action, s.GenSubAccountMintSignBytes())
}

// === SubAccountNew ===

type SubAccountNewVersion = uint32

const (
	SubAccountNewVersion1      SubAccountNewVersion = 1
	SubAccountNewVersion2      SubAccountNewVersion = 2
	SubAccountNewVersion3      SubAccountNewVersion = 3
	SubAccountNewVersionLatest                      = SubAccountNewVersion3
)

type SubAccountNew struct {
	// v2
	Index                int
	Version              SubAccountNewVersion
	Action               string
	actionBys            []byte
	versionBys           []byte
	Signature            []byte
	SignRole             []byte
	SignExpiredAt        uint64
	signExpiredAtBys     []byte
	NewRoot              []byte
	Proof                []byte
	OldSubAccountVersion SubAccountVersion
	NewSubAccountVersion SubAccountVersion
	SubAccountData       *SubAccountData
	subAccountDataBys    []byte
	EditKey              common.EditKey
	editKeyBys           []byte
	EditValue            []byte
	//
	EditLockArgs          []byte
	EditRecords           []Record
	CurrentSubAccountData *SubAccountData
	Account               string
	// v1
	PrevRoot    []byte
	CurrentRoot []byte
}

func (s *SubAccountNew) genSubAccountNewBytesV1() (dataBys []byte, err error) {
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.Signature)))...)
	dataBys = append(dataBys, s.Signature...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.SignRole)))...)
	dataBys = append(dataBys, s.SignRole...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.PrevRoot)))...)
	dataBys = append(dataBys, s.PrevRoot...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.CurrentRoot)))...)
	dataBys = append(dataBys, s.CurrentRoot...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.Proof)))...)
	dataBys = append(dataBys, s.Proof...)

	versionBys := molecule.GoU32ToMoleculeU32(s.Version)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
	dataBys = append(dataBys, versionBys.RawData()...)

	if s.SubAccountData == nil {
		return nil, fmt.Errorf("SubAccountData is nil")
	}
	subAccountData, err := s.SubAccountData.ConvertToMoleculeSubAccount(SubAccountVersion1)
	if err != nil {
		return nil, fmt.Errorf("ConvertToMoleculeSubAccount err: %s", err.Error())
	}
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(subAccountData)))...)
	dataBys = append(dataBys, subAccountData...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len([]byte(s.EditKey))))...)
	dataBys = append(dataBys, s.EditKey...)

	var editValue []byte
	if len(s.EditValue) > 0 {
		editValue = s.EditValue
	}
	switch s.EditKey {
	case common.EditKeyOwner, common.EditKeyManager:
		editValue = s.EditLockArgs
	case common.EditKeyRecords:
		records := ConvertToCellRecords(s.EditRecords)
		editValue = records.AsSlice()
	}

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(editValue)))...)
	dataBys = append(dataBys, editValue...)
	return
}

func (s *SubAccountNew) genSubAccountNewBytes() (dataBys []byte, err error) {
	versionBys := molecule.GoU32ToMoleculeU32(s.Version)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
	dataBys = append(dataBys, versionBys.RawData()...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len([]byte(s.Action))))...)
	dataBys = append(dataBys, s.Action...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.Signature)))...)
	dataBys = append(dataBys, s.Signature...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.SignRole)))...)
	dataBys = append(dataBys, s.SignRole...)

	signExpiredAtBys := molecule.GoU64ToMoleculeU64(s.SignExpiredAt)
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(signExpiredAtBys.RawData())))...)
	dataBys = append(dataBys, signExpiredAtBys.RawData()...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.NewRoot)))...)
	dataBys = append(dataBys, s.NewRoot...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(s.Proof)))...)
	dataBys = append(dataBys, s.Proof...)

	subAccountVersion := SubAccountVersion1
	if s.Version >= SubAccountNewVersion3 {
		oldSubAccountVersionBys := molecule.GoU32ToMoleculeU32(s.OldSubAccountVersion)
		dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(oldSubAccountVersionBys.RawData())))...)
		dataBys = append(dataBys, oldSubAccountVersionBys.RawData()...)

		newSubAccountVersionBys := molecule.GoU32ToMoleculeU32(s.NewSubAccountVersion)
		dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(newSubAccountVersionBys.RawData())))...)
		dataBys = append(dataBys, newSubAccountVersionBys.RawData()...)

		subAccountVersion = s.OldSubAccountVersion
		log.Infof("genSubAccountNewBytesV3 oldSubAccountVersion: %d newSubAccountVersion: %d", s.OldSubAccountVersion, s.NewSubAccountVersion)
	}
	if s.SubAccountData == nil {
		return nil, fmt.Errorf("SubAccountData is nil")
	}

	subAccountData, err := s.SubAccountData.ConvertToMoleculeSubAccount(subAccountVersion)
	if err != nil {
		return nil, fmt.Errorf("ConvertToMoleculeSubAccount err: %s", err.Error())
	}
	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(subAccountData)))...)
	dataBys = append(dataBys, subAccountData...)

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len([]byte(s.EditKey))))...)
	dataBys = append(dataBys, s.EditKey...)

	var editValue []byte
	if len(s.EditValue) > 0 {
		editValue = s.EditValue
	}
	switch s.EditKey {
	case common.EditKeyOwner, common.EditKeyManager:
		editValue = s.EditLockArgs
	case common.EditKeyRecords:
		records := ConvertToCellRecords(s.EditRecords)
		editValue = records.AsSlice()
	}

	if s.Action == common.SubActionCreate && len(s.EditValue) > 0 {
		editValue = s.EditValue
	}

	dataBys = append(dataBys, molecule.GoU32ToBytes(uint32(len(editValue)))...)
	dataBys = append(dataBys, editValue...)
	return
}

func (s *SubAccountNew) GenSubAccountNewBytes() (dataBys []byte, err error) {
	if s.Version == 0 || s.Version == SubAccountNewVersion1 {
		return s.genSubAccountNewBytesV1()
	}
	return s.genSubAccountNewBytes()
}

func (s *SubAccountNew) GenWitness() ([]byte, error) {
	dataBys, err := s.GenSubAccountNewBytes()
	if err != nil {
		return nil, fmt.Errorf("GenSubAccountNewBytes err: %s", err.Error())
	}
	witness := GenDasDataWitnessWithByte(common.ActionDataTypeSubAccount, dataBys)
	return witness, nil
}

func (s *SubAccountNewBuilder) convertSubAccountNewFromBytesV1(dataBys []byte) (*SubAccountNew, error) {
	var res SubAccountNew
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.Signature = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.SignRole = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.PrevRoot = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.CurrentRoot = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.Proof = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.versionBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.Version, _ = molecule.Bytes2GoU32(res.versionBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.subAccountDataBys = dataBys[index+indexLen : index+indexLen+dataLen]
	subAccount, err := s.ConvertSubAccountDataFromBytes(SubAccountVersion1, res.subAccountDataBys)
	if err != nil {
		return nil, fmt.Errorf("ConvertSubAccountDataFromBytes err: %s", err.Error())
	}
	res.SubAccountData = subAccount
	res.Account = subAccount.Account()
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.editKeyBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.EditKey = string(res.editKeyBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.EditValue = dataBys[index+indexLen : index+indexLen+dataLen]
	if err := s.convertCurrentSubAccountData(&res); err != nil {
		return nil, err
	}
	index = index + indexLen + dataLen

	return &res, nil
}

func (s *SubAccountNewBuilder) convertSubAccountNewFromBytes(dataBys []byte) (*SubAccountNew, error) {
	var res SubAccountNew
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.versionBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.Version, _ = molecule.Bytes2GoU32(res.versionBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.actionBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.Action = string(res.actionBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.Signature = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.SignRole = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.signExpiredAtBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.SignExpiredAt, _ = molecule.Bytes2GoU64(res.signExpiredAtBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.NewRoot = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.Proof = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.OldSubAccountVersion, _ = molecule.Bytes2GoU32(dataBys[index+indexLen : index+indexLen+dataLen])
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.NewSubAccountVersion, _ = molecule.Bytes2GoU32(dataBys[index+indexLen : index+indexLen+dataLen])
	index = index + indexLen + dataLen

	subAccountVersion := SubAccountVersion1
	if res.Version >= SubAccountNewVersion3 {
		subAccountVersion = res.OldSubAccountVersion
	}
	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.subAccountDataBys = dataBys[index+indexLen : index+indexLen+dataLen]
	subAccount, err := s.ConvertSubAccountDataFromBytes(subAccountVersion, res.subAccountDataBys)
	if err != nil {
		return nil, fmt.Errorf("ConvertSubAccountDataFromBytes err: %s", err.Error())
	}
	res.SubAccountData = subAccount
	res.Account = subAccount.Account()
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.editKeyBys = dataBys[index+indexLen : index+indexLen+dataLen]
	res.EditKey = string(res.editKeyBys)
	index = index + indexLen + dataLen

	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	res.EditValue = dataBys[index+indexLen : index+indexLen+dataLen]
	index = index + indexLen + dataLen

	if err := s.convertCurrentSubAccountData(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubAccountNewBuilder) ConvertSubAccountNewFromBytes(dataBys []byte) (*SubAccountNew, error) {
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)
	dataLen, _ = molecule.Bytes2GoU32(dataBys[index : index+indexLen])
	if dataLen == 4 {
		return s.convertSubAccountNewFromBytes(dataBys)
	} else {
		return s.convertSubAccountNewFromBytesV1(dataBys)
	}
}

func (s *SubAccountNewBuilder) SubAccountNewMapFromTx(tx *types.Transaction) (map[string]*SubAccountNew, error) {
	var respMap = make(map[string]*SubAccountNew)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeSubAccount:
			subAccountNew, err := s.ConvertSubAccountNewFromBytes(dataBys)
			if err != nil {
				return false, err
			}
			subAccountNew.Index = index
			respMap[subAccountNew.SubAccountData.AccountId] = subAccountNew
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist sub account")
	}
	return respMap, nil
}

// === EditValue ===
func (s *SubAccountNewBuilder) convertCurrentSubAccountData(p *SubAccountNew) error {
	if p.Action == common.SubActionRecycle {
		p.CurrentSubAccountData = &SubAccountData{}
		return nil
	}
	currentSubAccountData := *p.SubAccountData
	p.CurrentSubAccountData = &currentSubAccountData

	if (p.Action == "" && p.EditKey != "") || (p.Action != "" && p.Action != common.SubActionCreate) {
		p.CurrentSubAccountData.Nonce++
	}
	switch p.EditKey {
	case common.EditKeyOwner:
		p.CurrentSubAccountData.Lock = &types.Script{
			CodeHash: p.CurrentSubAccountData.Lock.CodeHash,
			HashType: p.CurrentSubAccountData.Lock.HashType,
			Args:     p.EditValue,
		}
		p.EditLockArgs = p.EditValue
		p.CurrentSubAccountData.Records = nil
	case common.EditKeyManager:
		p.CurrentSubAccountData.Lock = &types.Script{
			CodeHash: p.CurrentSubAccountData.Lock.CodeHash,
			HashType: p.CurrentSubAccountData.Lock.HashType,
			Args:     p.EditValue,
		}
		p.EditLockArgs = p.EditValue
	case common.EditKeyRecords:
		records, _ := molecule.RecordsFromSlice(p.EditValue, true)
		p.EditRecords = ConvertToRecords(records)
		p.CurrentSubAccountData.Records = p.EditRecords
	}

	switch p.Action {
	case common.SubActionRenew:
		expiredAt, _ := molecule.Bytes2GoU64(p.EditValue[:8])
		p.CurrentSubAccountData.ExpiredAt = expiredAt
	case common.SubActionCreateApproval:
		p.CurrentSubAccountData.Status = common.AccountStatusOnApproval
		accountApproval, err := AccountApprovalFromSlice(p.EditValue)
		if err != nil {
			return fmt.Errorf("AccountApprovalFromSlice err: %s", err.Error())
		}
		p.CurrentSubAccountData.AccountApproval = *accountApproval
	case common.DasActionDelayApproval:
		accountApproval, err := AccountApprovalFromSlice(p.EditValue)
		if err != nil {
			return fmt.Errorf("AccountApprovalFromSlice err: %s", err.Error())
		}
		p.CurrentSubAccountData.AccountApproval = *accountApproval
	case common.SubActionRevokeApproval:
		p.CurrentSubAccountData.Status = common.AccountStatusNormal
	case common.SubActionFullfillApproval:
		p.CurrentSubAccountData.Status = common.AccountStatusNormal
		p.CurrentSubAccountData.Lock = p.CurrentSubAccountData.AccountApproval.Params.Transfer.ToLock
	}
	return nil
}

// === SubAccountData ===
type SubAccountData struct {
	Lock                 *types.Script           `json:"lock"`
	AccountId            string                  `json:"account_id"`
	AccountCharSet       []common.AccountCharSet `json:"account_char_set"`
	Suffix               string                  `json:"suffix"`
	RegisteredAt         uint64                  `json:"registered_at"`
	ExpiredAt            uint64                  `json:"expired_at"`
	Status               uint8                   `json:"status"`
	Records              []Record                `json:"records"`
	Nonce                uint64                  `json:"nonce"`
	EnableSubAccount     uint8                   `json:"enable_sub_account"`
	RenewSubAccountPrice uint64                  `json:"renew_sub_account_price"`
	AccountApproval      AccountApproval         `json:"account_approval"`
}

func (s *SubAccountNewBuilder) ConvertSubAccountDataFromBytes(version SubAccountVersion, dataBys []byte) (*SubAccountData, error) {
	var tmp SubAccountData
	switch version {
	case SubAccountVersion1:
		subAccount, err := molecule.SubAccountV1FromSlice(dataBys, true)
		if err != nil {
			return nil, fmt.Errorf("SubAccountV1FromSlice err: %s", err.Error())
		}
		tmp.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
		tmp.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
		tmp.AccountCharSet = common.ConvertToAccountCharSets(subAccount.Account())
		tmp.Suffix = string(subAccount.Suffix().RawData())
		tmp.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
		tmp.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
		tmp.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
		tmp.Records = ConvertToRecords(subAccount.Records())
		tmp.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
		tmp.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
		tmp.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())
	case SubAccountVersion2:
		subAccount, err := molecule.SubAccountFromSlice(dataBys, true)
		if err != nil {
			return nil, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
		}
		tmp.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
		tmp.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
		tmp.AccountCharSet = common.ConvertToAccountCharSets(subAccount.Account())
		tmp.Suffix = string(subAccount.Suffix().RawData())
		tmp.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
		tmp.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
		tmp.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
		tmp.Records = ConvertToRecords(subAccount.Records())
		tmp.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
		tmp.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
		tmp.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())
		accountApproval, err := AccountApprovalFromSlice(subAccount.Approval().AsSlice())
		if err != nil {
			return nil, err
		}
		tmp.AccountApproval = *accountApproval
	default:
		return nil, fmt.Errorf("unknown subAccount version: %d", version)
	}
	return &tmp, nil
}

func (s *SubAccountData) ConvertToMoleculeSubAccount(version SubAccountVersion) ([]byte, error) {
	if s.Lock == nil {
		return nil, fmt.Errorf("lock is nil")
	}
	lock := molecule.CkbScript2MoleculeScript(s.Lock)
	accountChars := common.ConvertToAccountChars(s.AccountCharSet)
	accountId, err := molecule.AccountIdFromSlice(common.Hex2Bytes(s.AccountId), true)
	if err != nil {
		return nil, fmt.Errorf("AccountIdFromSlice err: %s", err.Error())
	}
	suffix := molecule.GoBytes2MoleculeBytes([]byte(s.Suffix))
	registeredAt := molecule.GoU64ToMoleculeU64(s.RegisteredAt)
	expiredAt := molecule.GoU64ToMoleculeU64(s.ExpiredAt)
	status := molecule.GoU8ToMoleculeU8(s.Status)
	records := ConvertToCellRecords(s.Records)
	nonce := molecule.GoU64ToMoleculeU64(s.Nonce)
	enableSubAccount := molecule.GoU8ToMoleculeU8(s.EnableSubAccount)
	renewSubAccountPrice := molecule.GoU64ToMoleculeU64(s.RenewSubAccountPrice)
	accountApproval, err := s.AccountApproval.GenToMolecule()
	if err != nil {
		return nil, err
	}

	var res []byte
	switch version {
	case SubAccountVersion1:
		moleculeSubAccount := molecule.NewSubAccountV1Builder().
			Lock(lock).
			Id(*accountId).
			Account(*accountChars).
			Suffix(suffix).
			RegisteredAt(registeredAt).
			ExpiredAt(expiredAt).
			Status(status).
			Records(*records).
			Nonce(nonce).
			EnableSubAccount(enableSubAccount).
			RenewSubAccountPrice(renewSubAccountPrice).
			Build()
		res = moleculeSubAccount.AsSlice()
	case SubAccountVersion2:
		moleculeSubAccount := molecule.NewSubAccountBuilder().
			Lock(lock).
			Id(*accountId).
			Account(*accountChars).
			Suffix(suffix).
			RegisteredAt(registeredAt).
			ExpiredAt(expiredAt).
			Status(status).
			Records(*records).
			Nonce(nonce).
			EnableSubAccount(enableSubAccount).
			RenewSubAccountPrice(renewSubAccountPrice).
			Approval(*accountApproval).
			Build()
		res = moleculeSubAccount.AsSlice()
	default:
		return nil, fmt.Errorf("unknown subAccount version: %d", version)
	}
	return res, nil
}

func (s *SubAccountData) Account() string {
	var account string
	for _, v := range s.AccountCharSet {
		account += v.Char
	}
	return account + s.Suffix
}

func (s *SubAccountData) ToH256() ([]byte, error) {
	if s.AccountId == "" { // for recycle sub_account
		return make([]byte, 32), nil
	}
	moleculeSubAccount, err := s.ConvertToMoleculeSubAccount(SubAccountVersionLatest)
	if err != nil {
		log.Error("ToH256 ConvertToMoleculeSubAccount err:", err.Error())
		return nil, err
	}
	res, err := blake2b.Blake256(moleculeSubAccount)
	if err != nil {
		log.Error("ToH256 blake2b.Blake256 err:", err.Error())
		return nil, err
	}
	return res, nil
}
