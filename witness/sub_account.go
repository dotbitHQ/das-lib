package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type SubAccountBuilder struct {
	Signature      string
	PrevRoot       string
	CurrentRoot    string
	Proof          string
	Version        uint32
	SubAccount     *SubAccount
	Key            string
	Value          []byte
	SubAccountData *molecule.SubAccountData
}

type SubAccountParam struct {
	Action      string
	SubAction   string
	Signature   string
	PrevRoot    string
	CurrentRoot string
	Proof       string
	Version     uint32
	SubAccount  *SubAccount
	Key         string
	Value       []byte
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
	Nonce                uint32
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
			resp.Signature = common.Bytes2Hex(dataBys[length:signatureLen])
			index = length + int(signatureLen)

			prevRootLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.PrevRoot = common.Bytes2Hex(dataBys[index+length : prevRootLen])
			index = length + int(prevRootLen)

			currentRootLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.CurrentRoot = common.Bytes2Hex(dataBys[index+length : currentRootLen])
			index = length + int(currentRootLen)

			proofLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Proof = common.Bytes2Hex(dataBys[index+length : proofLen])
			index = length + int(proofLen)

			versionLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Version = molecule.BytesToGoU32Big(dataBys[index+length : versionLen])
			index = length + int(versionLen)

			subAccountLen := molecule.BytesToGoU32Big(dataBys[index:length])
			subAccountBys := dataBys[index+length : subAccountLen]
			index = length + int(subAccountLen)

			keyLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Key = common.Bytes2Hex(dataBys[index+length : keyLen])
			index = length + int(keyLen)

			valueLen := molecule.BytesToGoU32Big(dataBys[index:length])
			resp.Value = dataBys[index+length : valueLen]
			index = length + int(valueLen)

			switch resp.Version {
			case common.GoDataEntityVersion1:
				subAccountData, err := molecule.SubAccountDataFromSlice(subAccountBys, false)
				if err != nil {
					return false, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
				}
				resp.SubAccountData = subAccountData
				resp.SubAccount.Lock = molecule.MoleculeScript2CkbScript(subAccountData.Lock())
				resp.SubAccount.AccountId = common.Bytes2Hex(subAccountData.Id().RawData())
				resp.SubAccount.Account = common.AccountCharsToAccount(subAccountData.Account())
				resp.SubAccount.Suffix = string(subAccountData.Suffix().RawData())
				resp.SubAccount.RegisteredAt, _ = molecule.Bytes2GoU64(subAccountData.RegisteredAt().RawData())
				resp.SubAccount.ExpiredAt, _ = molecule.Bytes2GoU64(subAccountData.ExpiredAt().RawData())
				resp.SubAccount.Status, _ = molecule.Bytes2GoU8(subAccountData.Status().RawData())
				resp.SubAccount.Records = ConvertToRecordList(subAccountData.Records())
				resp.SubAccount.Nonce, _ = molecule.Bytes2GoU32(subAccountData.Nonce().RawData())
				resp.SubAccount.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccountData.EnableSubAccount().RawData())
				resp.SubAccount.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccountData.RenewSubAccountPrice().RawData())
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

// ConvertToSubAccount value convert use case
func (a *SubAccountBuilder) ConvertToSubAccount() (sub *SubAccount) {
	switch a.Key {
	case "lock":
		sub.Lock = a.ConvertToLock()
	case "id":
		sub.AccountId = a.ConvertToAccountId()
	case "account":
		sub.Account = a.ConvertToAccount()
	case "suffix":
		sub.Suffix = a.ConvertToSuffix()
	case "registered_at":
		sub.RegisteredAt = a.ConvertToRegisteredAt()
	case "expired_at":
		sub.ExpiredAt = a.ConvertToExpiredAt()
	case "status":
		sub.Status = a.ConvertToStatus()
	case "records":
		sub.Records = a.ConvertToRecords()
	case "nonce":
		sub.Nonce = a.ConvertToNonce()
	case "enable_sub_account":
		sub.EnableSubAccount = a.ConvertToEnableSubAccount()
	case "renew_sub_account_price":
		sub.RenewSubAccountPrice = a.ConvertToRenewSubAccountPrice()
	}
	return
}

func (a *SubAccountBuilder) ConvertToLock() *types.Script {
	lock, _ := molecule.ScriptFromSlice(a.Value, false)
	return molecule.MoleculeScript2CkbScript(lock)
}

func (a *SubAccountBuilder) ConvertToAccountId() string {
	return common.Bytes2Hex(a.Value)
}

func (a *SubAccountBuilder) ConvertToAccount() string {
	account, _ := molecule.AccountCharsFromSlice(a.Value, false)
	return common.AccountCharsToAccount(account)
}

func (a *SubAccountBuilder) ConvertToSuffix() string {
	return string(a.Value)
}

func (a *SubAccountBuilder) ConvertToRegisteredAt() uint64 {
	registeredAt, _ := molecule.Bytes2GoU64(a.Value)
	return registeredAt
}

func (a *SubAccountBuilder) ConvertToExpiredAt() uint64 {
	expiredAt, _ := molecule.Bytes2GoU64(a.Value)
	return expiredAt
}

func (a *SubAccountBuilder) ConvertToStatus() uint8 {
	status, _ := molecule.Bytes2GoU8(a.Value)
	return status
}

func (a *SubAccountBuilder) ConvertToRecords() []*SubAccountRecord {
	records, _ := molecule.RecordsFromSlice(a.Value, false)
	return ConvertToRecordList(records)
}

func (a *SubAccountBuilder) ConvertToNonce() uint32 {
	nonce, _ := molecule.Bytes2GoU32(a.Value)
	return nonce
}

func (a *SubAccountBuilder) ConvertToEnableSubAccount() uint8 {
	enableSubAccount, _ := molecule.Bytes2GoU8(a.Value)
	return enableSubAccount
}

func (a *SubAccountBuilder) ConvertToRenewSubAccountPrice() uint64 {
	renewSubAccountPrice, _ := molecule.Bytes2GoU64(a.Value)
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
