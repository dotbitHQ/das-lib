package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type SubAccountBuilder struct {
	Signature          []byte
	PrevRoot           []byte
	CurrentRoot        []byte
	Proof              []byte
	Version            uint32
	SubAccount         *SubAccount
	EditKey            []byte
	EditValue          []byte
	MoleculeSubAccount *molecule.SubAccount
	Account            string
}

type SubAccountParam struct {
	Action      string
	Signature   []byte
	PrevRoot    []byte
	CurrentRoot []byte
	Proof       []byte
	Version     uint32
	SubAccount  *SubAccount
	EditKey     []byte
	EditValue   []byte
}

type SubAccount struct {
	Lock                 *types.Script       `json:"lock"`
	AccountId            string              `json:"account_id"`
	AccountCharSet       []*AccountCharSet   `json:"account_char_set"`
	Suffix               string              `json:"suffix"`
	RegisteredAt         uint64              `json:"registered_at"`
	ExpiredAt            uint64              `json:"expired_at"`
	Status               uint8               `json:"status"`
	Records              []*SubAccountRecord `json:"records"`
	Nonce                uint64              `json:"nonce"`
	EnableSubAccount     uint8               `json:"enable_sub_account"`
	RenewSubAccountPrice uint64              `json:"renew_sub_account_price"`
}

type SubAccountRecord struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

type SubAccountEditValue struct {
	Lock      *types.Script       `json:"lock"`
	Records   []*SubAccountRecord `json:"records"`
	ExpiredAt uint64              `json:"expired_at"`
}

func SubAccountDataBuilderFromTx(tx *types.Transaction) (*SubAccountBuilder, error) {
	respMap, err := SubAccountDataBuilderMapFromTx(tx)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist sub account")
}

func SubAccountIdDataBuilderFromTx(tx *types.Transaction) (map[string]*SubAccountBuilder, error) {
	respMap, err := SubAccountDataBuilderMapFromTx(tx)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*SubAccountBuilder)
	for k, v := range respMap {
		k1 := v.SubAccount.AccountId
		retMap[k1] = respMap[k]
	}
	return retMap, nil
}

func SubAccountDataBuilderMapFromTx(tx *types.Transaction) (map[string]*SubAccountBuilder, error) {
	var respMap = make(map[string]*SubAccountBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeSubAccount:
			var resp SubAccountBuilder
			index, length := 0, 4

			signatureLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.Signature = dataBys[length:signatureLen]
			index = length + int(signatureLen)

			prevRootLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.PrevRoot = dataBys[index+length : prevRootLen]
			index = length + int(prevRootLen)

			currentRootLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.CurrentRoot = dataBys[index+length : currentRootLen]
			index = length + int(currentRootLen)

			proofLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.Proof = dataBys[index+length : proofLen]
			index = length + int(proofLen)

			versionLen, err := molecule.Bytes2GoU32(dataBys[index:length])
			if err != nil {
				return false, fmt.Errorf("get version len err: %s", err.Error())
			}
			resp.Version, err = molecule.Bytes2GoU32(dataBys[index+length : versionLen])
			if err != nil {
				return false, fmt.Errorf("get version err: %s", err.Error())
			}
			index = length + int(versionLen)

			subAccountLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			subAccountBys := dataBys[index+length : subAccountLen]
			index = length + int(subAccountLen)

			keyLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.EditKey = dataBys[index+length : keyLen]
			index = length + int(keyLen)

			valueLen, _ := molecule.Bytes2GoU32(dataBys[index:length])
			resp.EditValue = dataBys[index+length : valueLen]

			switch resp.Version {
			case common.GoDataEntityVersion1:
				subAccount, err := molecule.SubAccountFromSlice(subAccountBys, false)
				if err != nil {
					return false, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
				}
				resp.SubAccount.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
				resp.SubAccount.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
				resp.SubAccount.AccountCharSet = ConvertToAccountCharSets(subAccount.Account())
				resp.SubAccount.Suffix = string(subAccount.Suffix().RawData())
				resp.SubAccount.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
				resp.SubAccount.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
				resp.SubAccount.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
				resp.SubAccount.Records = ConvertToSubAccountRecords(subAccount.Records())
				resp.SubAccount.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
				resp.SubAccount.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
				resp.SubAccount.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())
				resp.MoleculeSubAccount = subAccount
				resp.Account = common.AccountCharsToAccount(subAccount.Account())
				respMap[resp.Account] = &resp
			default:
				return false, fmt.Errorf("sub account version: %d", resp.Version)
			}
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

func (s *SubAccountBuilder) ConvertToEditValue() (*SubAccountEditValue, error) {
	var editValue SubAccountEditValue
	editKey := string(s.EditKey)
	switch editKey {
	case common.EditKeyOwner, common.EditKeyManager:
		lock := s.ConvertEditValueToLock()
		editValue.Lock = molecule.MoleculeScript2CkbScript(lock)
	case common.EditKeyRecords:
		records := s.ConvertEditValueToRecords()
		editValue.Records = ConvertToSubAccountRecords(records)
	case common.EditKeyExpiredAt:
		expiredAt := s.ConvertEditValueToExpiredAt()
		editValue.ExpiredAt, _ = molecule.Bytes2GoU64(expiredAt.RawData())
	default:
		return nil, fmt.Errorf("not support edit key[%s]", editKey)
	}
	return &editValue, nil
}

func (s *SubAccountBuilder) ConvertEditValueToLock() *molecule.Script {
	lock, _ := molecule.ScriptFromSlice(s.EditValue, false)
	return lock
}

func (s *SubAccountBuilder) ConvertEditValueToExpiredAt() *molecule.Uint64 {
	expiredAt, _ := molecule.Uint64FromSlice(s.EditValue, false)
	return expiredAt
}

func (s *SubAccountBuilder) ConvertEditValueToRecords() *molecule.Records {
	records, _ := molecule.RecordsFromSlice(s.EditValue, false)
	return records
}

func ConvertToSubAccountRecords(records *molecule.Records) []*SubAccountRecord {
	var subAccountRecords []*SubAccountRecord
	for index, lenRecords := uint(0), records.Len(); index < lenRecords; index++ {
		record := records.Get(index)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		subAccountRecords = append(subAccountRecords, &SubAccountRecord{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   ttl,
		})
	}
	return subAccountRecords
}

func ConvertToRecordsHash(records *molecule.Records) []byte {
	bys, _ := blake2b.Blake256(records.AsSlice())
	return bys
}

type AccountCharType uint32

const (
	AccountCharTypeEmoji  AccountCharType = 0
	AccountCharTypeNumber AccountCharType = 1
	AccountCharTypeEn     AccountCharType = 2
)

type AccountCharSet struct {
	CharSetName AccountCharType `json:"char_set_name"`
	Char        string          `json:"char"`
}

func ConvertToAccountCharSets(accountChars *molecule.AccountChars) []*AccountCharSet {
	index := uint(0)
	var accountCharSets []*AccountCharSet
	for ; index < accountChars.ItemCount(); index++ {
		char := accountChars.Get(index)
		charSetName, _ := molecule.Bytes2GoU32(char.CharSetName().RawData())
		accountCharSets = append(accountCharSets, &AccountCharSet{
			CharSetName: AccountCharType(charSetName),
			Char:        string(char.Bytes().RawData()),
		})
	}
	return accountCharSets
}

/****************************************** Parting Line ******************************************/

func ConvertToScript(script *types.Script) *molecule.Script {
	lock := molecule.CkbScript2MoleculeScript(script)
	return &lock
}

func ConvertToAccountChars(accountCharSet []*AccountCharSet) *molecule.AccountChars {
	accountCharsBuilder := molecule.NewAccountCharsBuilder()
	for _, item := range accountCharSet {
		if item.Char == "." {
			break
		}
		accountChar := molecule.NewAccountCharBuilder().
			CharSetName(molecule.GoU32ToMoleculeU32(uint32(item.CharSetName))).
			Bytes(molecule.GoBytes2MoleculeBytes([]byte(item.Char))).Build()
		accountCharsBuilder.Push(accountChar)
	}
	accountChars := accountCharsBuilder.Build()
	return &accountChars
}

func ConvertToExpiredAt(expiredAt uint64) molecule.Uint64 {
	return molecule.GoU64ToMoleculeU64(expiredAt)
}

func ConvertToRecords(subAccountRecords []*SubAccountRecord) *molecule.Records {
	recordsBuilder := molecule.NewRecordsBuilder()
	for _, v := range subAccountRecords {
		record := molecule.RecordDefault()
		recordBuilder := record.AsBuilder()
		recordBuilder.RecordKey(molecule.GoString2MoleculeBytes(v.Key)).
			RecordType(molecule.GoString2MoleculeBytes(v.Type)).
			RecordLabel(molecule.GoString2MoleculeBytes(v.Label)).
			RecordValue(molecule.GoString2MoleculeBytes(v.Value)).
			RecordTtl(molecule.GoU32ToMoleculeU32(v.TTL))
		recordsBuilder.Push(recordBuilder.Build())
	}
	records := recordsBuilder.Build()
	return &records
}

func ConvertToMoleculeSubAccount(subAccount *SubAccount) *molecule.SubAccount {
	lock := ConvertToScript(subAccount.Lock)
	accountChars := ConvertToAccountChars(subAccount.AccountCharSet)
	accountId, _ := molecule.AccountIdFromSlice(common.Hex2Bytes(subAccount.AccountId), false)
	suffix := molecule.GoBytes2MoleculeBytes([]byte(subAccount.Suffix))
	registeredAt := molecule.GoU64ToMoleculeU64(subAccount.RegisteredAt)
	expiredAt := ConvertToExpiredAt(subAccount.ExpiredAt)
	status := molecule.GoU8ToMoleculeU8(subAccount.Status)
	records := ConvertToRecords(subAccount.Records)
	nonce := molecule.GoU64ToMoleculeU64(subAccount.Nonce)
	enableSubAccount := molecule.GoU8ToMoleculeU8(subAccount.EnableSubAccount)
	renewSubAccountPrice := molecule.GoU64ToMoleculeU64(subAccount.RenewSubAccountPrice)

	moleculeSubAccount := molecule.NewSubAccountBuilder().
		Lock(*lock).
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
	return &moleculeSubAccount
}

func (s *SubAccountBuilder) GenSubAccountBytes(p *SubAccountParam, subAccount *molecule.SubAccount) (bys []byte) {
	switch p.Version {
	case common.GoDataEntityVersion1:
		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Signature)))...)
		bys = append(bys, p.Signature...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.PrevRoot)))...)
		bys = append(bys, p.PrevRoot...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.CurrentRoot)))...)
		bys = append(bys, p.CurrentRoot...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Proof)))...)
		bys = append(bys, p.Proof...)

		versionBys := molecule.GoU32ToMoleculeU32(p.Version)
		bys = append(bys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
		bys = append(bys, versionBys.RawData()...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(subAccount.AsSlice())))...)
		bys = append(bys, subAccount.AsSlice()...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.EditKey)))...)
		bys = append(bys, p.EditKey...)

		bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.EditValue)))...)
		bys = append(bys, p.EditValue...)
	}
	return bys
}

func (s *SubAccountBuilder) GenSubAccountBuilder() *molecule.SubAccountBuilder {
	subAccountBuilder := s.MoleculeSubAccount.AsBuilder()
	switch string(s.EditKey) {
	case common.EditKeyOwner, common.EditKeyManager:
		return subAccountBuilder.Lock(*s.ConvertEditValueToLock())
	case common.EditKeyExpiredAt:
		return subAccountBuilder.ExpiredAt(*s.ConvertEditValueToExpiredAt())
	case common.EditKeyRecords:
		return subAccountBuilder.Records(*s.ConvertEditValueToRecords())
	}

	return &subAccountBuilder
}

func (s *SubAccountBuilder) GenWitness(p *SubAccountParam) ([]byte, error) {
	switch p.Action {
	case common.DasActionCreateSubAccount:
		witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.GenSubAccountBytes(p, ConvertToMoleculeSubAccount(p.SubAccount)))

		return witness, nil
	case common.DasActionEditSubAccount:
		// nonce increment on each transaction
		nonceUint64, _ := molecule.Bytes2GoU64(s.MoleculeSubAccount.Nonce().RawData())
		nonceUint64++
		nonce := molecule.GoU64ToMoleculeU64(nonceUint64)

		subAccount := s.GenSubAccountBuilder().Nonce(nonce).Build()
		switch string(p.EditKey) {
		case common.EditKeyOwner, common.EditKeyManager, common.EditKeyRecords:
			witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.GenSubAccountBytes(p, &subAccount))
			return witness, nil
		default:
			return nil, fmt.Errorf("not support edit key [%s]", string(s.EditKey))
		}
	case common.DasActionRenewSubAccount:
		// nonce increment on each transaction
		nonceUint64, _ := molecule.Bytes2GoU64(s.MoleculeSubAccount.Nonce().RawData())
		nonceUint64++
		nonce := molecule.GoU64ToMoleculeU64(nonceUint64)

		subAccount := s.GenSubAccountBuilder().Nonce(nonce).Build()
		switch string(p.EditKey) {
		case common.EditKeyExpiredAt:
			witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.GenSubAccountBytes(p, &subAccount))
			return witness, nil
		default:
			return nil, fmt.Errorf("not support edit key [%s]", string(s.EditKey))
		}
	case common.DasActionRecycleSubAccount:
		witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.GenSubAccountBytes(p, s.MoleculeSubAccount))
		return witness, nil
	}
	return nil, fmt.Errorf("not exist action [%s]", p.Action)
}
