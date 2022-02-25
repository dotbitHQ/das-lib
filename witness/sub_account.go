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
}

type SubAccountParam struct {
	Action      string
	SubAction   string
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
	Lock                 *types.Script
	AccountId            string
	Account              string
	Suffix               string
	RegisteredAt         uint64
	ExpiredAt            uint64
	Status               uint8
	Records              []*SubAccountRecord
	Nonce                uint64
	EnableSubAccount     uint8
	RenewSubAccountPrice uint64
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

			signatureLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Signature = dataBys[length:signatureLen]
			index = length + int(signatureLen)

			prevRootLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.PrevRoot = dataBys[index+length : prevRootLen]
			index = length + int(prevRootLen)

			currentRootLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.CurrentRoot = dataBys[index+length : currentRootLen]
			index = length + int(currentRootLen)

			proofLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Proof = dataBys[index+length : proofLen]
			index = length + int(proofLen)

			versionLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Version, _ = molecule.Bytes2GoU32(dataBys[index+length : versionLen])
			index = length + int(versionLen)

			subAccountLen := molecule.BytesToGoU32Big(dataBys[index:length])
			subAccountBys := dataBys[index+length : subAccountLen]
			index = length + int(subAccountLen)

			keyLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.EditKey = dataBys[index+length : keyLen]
			index = length + int(keyLen)

			valueLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.EditValue = dataBys[index+length : valueLen]

			switch resp.Version {
			case common.GoDataEntityVersion1:
				subAccount, err := molecule.SubAccountFromSlice(subAccountBys, false)
				if err != nil {
					return false, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
				}
				resp.MoleculeSubAccount = subAccount
				resp.SubAccount.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
				resp.SubAccount.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
				resp.SubAccount.Account = common.AccountCharsToAccount(subAccount.Account())
				resp.SubAccount.Suffix = string(subAccount.Suffix().RawData())
				resp.SubAccount.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
				resp.SubAccount.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
				resp.SubAccount.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
				resp.SubAccount.Records = ConvertToRecordList(subAccount.Records())
				resp.SubAccount.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
				resp.SubAccount.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
				resp.SubAccount.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())
				respMap[resp.SubAccount.Account] = &resp
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

func (s *SubAccountBuilder) ConvertToSubAccount(sub *SubAccount) {
	switch string(s.EditKey) {
	case "lock":
		sub.Lock = s.ConvertToLock()
	case "id":
		sub.AccountId = s.ConvertToAccountId()
	case "account":
		sub.Account = s.ConvertToAccount()
	case "suffix":
		sub.Suffix = s.ConvertToSuffix()
	case "registered_at":
		sub.RegisteredAt = s.ConvertToRegisteredAt()
	case "expired_at":
		sub.ExpiredAt = s.ConvertToExpiredAt()
	case "status":
		sub.Status = s.ConvertToStatus()
	case "records":
		sub.Records = s.ConvertToRecords()
	case "nonce":
		sub.Nonce = s.ConvertToNonce()
	case "enable_sub_account":
		sub.EnableSubAccount = s.ConvertToEnableSubAccount()
	case "renew_sub_account_price":
		sub.RenewSubAccountPrice = s.ConvertToRenewSubAccountPrice()
	}
}

func (s *SubAccountBuilder) ConvertToLock() *types.Script {
	lock, _ := molecule.ScriptFromSlice(s.EditValue, false)
	return molecule.MoleculeScript2CkbScript(lock)
}

func (s *SubAccountBuilder) ConvertToAccountId() string {
	return common.Bytes2Hex(s.EditValue)
}

func (s *SubAccountBuilder) ConvertToAccount() string {
	account, _ := molecule.AccountCharsFromSlice(s.EditValue, false)
	return common.AccountCharsToAccount(account)
}

func (s *SubAccountBuilder) ConvertToSuffix() string {
	return string(s.EditValue)
}

func (s *SubAccountBuilder) ConvertToRegisteredAt() uint64 {
	registeredAt, _ := molecule.Bytes2GoU64(s.EditValue)
	return registeredAt
}

func (s *SubAccountBuilder) ConvertToExpiredAt() uint64 {
	expiredAt, _ := molecule.Bytes2GoU64(s.EditValue)
	return expiredAt
}

func (s *SubAccountBuilder) ConvertToStatus() uint8 {
	status, _ := molecule.Bytes2GoU8(s.EditValue)
	return status
}

func (s *SubAccountBuilder) ConvertToRecords() []*SubAccountRecord {
	records, _ := molecule.RecordsFromSlice(s.EditValue, false)
	return ConvertToRecordList(records)
}

func (s *SubAccountBuilder) ConvertToNonce() uint64 {
	nonce, _ := molecule.Bytes2GoU64(s.EditValue)
	return nonce
}

func (s *SubAccountBuilder) ConvertToEnableSubAccount() uint8 {
	enableSubAccount, _ := molecule.Bytes2GoU8(s.EditValue)
	return enableSubAccount
}

func (s *SubAccountBuilder) ConvertToRenewSubAccountPrice() uint64 {
	renewSubAccountPrice, _ := molecule.Bytes2GoU64(s.EditValue)
	return renewSubAccountPrice
}

type SubAccountRecord struct {
	Key   string
	Type  string
	Label string
	Value string
	TTL   uint32
}

func ConvertToRecordList(records *molecule.Records) []*SubAccountRecord {
	var list []*SubAccountRecord
	for index, lenRecords := uint(0), records.Len(); index < lenRecords; index++ {
		record := records.Get(index)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		list = append(list, &SubAccountRecord{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   ttl,
		})
	}
	return list
}

func ConvertToRecordHash(records *molecule.Records) []byte {
	bys, _ := blake2b.Blake256(records.AsSlice())
	return bys
}

/****************************************** Parting Line ******************************************/

func (s *SubAccountBuilder) ConvertToMoleculeSubAccount(p *SubAccountParam) *molecule.SubAccount {
	return nil
}

func (s *SubAccountBuilder) genOldSubAccountBytes(p *SubAccountParam) (bys []byte) {
	switch s.Version {
	case common.GoDataEntityVersion1:
		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.Signature)))...)
		bys = append(bys, p.Signature...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.PrevRoot)))...)
		bys = append(bys, p.PrevRoot...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.CurrentRoot)))...)
		bys = append(bys, p.CurrentRoot...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(s.Proof)))...)
		bys = append(bys, s.Proof...)

		versionBys := molecule.GoU32ToMoleculeU32(s.Version)
		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(versionBys.RawData())))...)
		bys = append(bys, versionBys.RawData()...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(s.MoleculeSubAccount.AsSlice())))...)
		bys = append(bys, s.MoleculeSubAccount.AsSlice()...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.EditKey)))...)
		bys = append(bys, p.EditKey...)
	}
	return bys
}

func (s *SubAccountBuilder) genNewSubAccountBytes(p *SubAccountParam) (bys []byte) {
	switch p.Version {
	case common.GoDataEntityVersion1:
		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.Signature)))...)
		bys = append(bys, p.Signature...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.PrevRoot)))...)
		bys = append(bys, p.PrevRoot...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.CurrentRoot)))...)
		bys = append(bys, p.CurrentRoot...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.Proof)))...)
		bys = append(bys, p.Proof...)

		versionBys := molecule.GoU32ToMoleculeU32(p.Version)
		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(versionBys.RawData())))...)
		bys = append(bys, versionBys.RawData()...)

		moleculeSubAccount := s.ConvertToMoleculeSubAccount(p)
		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(moleculeSubAccount.AsSlice())))...)
		bys = append(bys, moleculeSubAccount.AsSlice()...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.EditKey)))...)
		bys = append(bys, p.EditKey...)

		bys = append(bys, molecule.GoU32ToBytesBig(uint32(len(p.EditValue)))...)
		bys = append(bys, p.EditValue...)
	}
	return bys
}

func (s *SubAccountBuilder) GenWitness(p *SubAccountParam) ([]byte, error) {
	switch p.Action {
	case common.DasActionRenewAccount:
		witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.genOldSubAccountBytes(p))
		bys := molecule.GoU64ToMoleculeU64(p.SubAccount.ExpiredAt)

		witness = append(witness, molecule.GoU32ToBytesBig(uint32(len(bys.RawData())))...)
		return append(witness, bys.RawData()...), nil
	case common.DasActionEditRecords:
		witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.genOldSubAccountBytes(p))
		records := molecule.RecordsDefault()

		if len(p.SubAccount.Records) > 0 {
			recordsBuilder := records.AsBuilder()
			for _, v := range p.SubAccount.Records {
				record := molecule.RecordDefault()
				recordBuilder := record.AsBuilder()
				recordBuilder.RecordKey(molecule.GoString2MoleculeBytes(v.Key)).
					RecordType(molecule.GoString2MoleculeBytes(v.Type)).
					RecordLabel(molecule.GoString2MoleculeBytes(v.Label)).
					RecordValue(molecule.GoString2MoleculeBytes(v.Value)).
					RecordTtl(molecule.GoU32ToMoleculeU32(v.TTL))
				recordsBuilder.Push(recordBuilder.Build())
			}
			records = recordsBuilder.Build()
		}

		witness = append(witness, molecule.GoU32ToBytesBig(uint32(len(records.AsSlice())))...)
		return append(witness, records.AsSlice()...), nil
	case common.DasActionEditManager, common.DasActionTransferAccount:
		witness := GenDasSubAccountWitness(common.ActionDataTypeSubAccount, s.genOldSubAccountBytes(p))
		lock := molecule.CkbScript2MoleculeScript(p.SubAccount.Lock)

		witness = append(witness, molecule.GoU32ToBytesBig(uint32(len(lock.AsSlice())))...)
		return append(witness, lock.AsSlice()...), nil

	case common.DasActionEnableSubAccount:
	case common.DasActionCreateSubAccount:
	case common.DasActionEditSubAccount:
	case common.DasActionRenewSubAccount:
	case common.DasActionRecycleSubAccount:

	}
	return nil, fmt.Errorf("not exist action [%s]", p.Action)
}
