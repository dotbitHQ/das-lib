package witness

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sort"
	"strings"
)

//
//const (
//	SubAccountCurrentVersion = common.GoDataEntityVersion1
//)
//
//type SubAccountBuilder struct {
//	Signature         []byte
//	SignRole          []byte
//	PrevRoot          []byte
//	CurrentRoot       []byte
//	Proof             []byte
//	Version           uint32
//	SubAccount        *SubAccount
//	EditKey           []byte
//	EditValue         []byte
//	Account           string
//	CurrentSubAccount *SubAccount
//}
//
//type SubAccountParam struct {
//	Signature      []byte
//	SignRole       []byte
//	PrevRoot       []byte
//	CurrentRoot    []byte
//	Proof          []byte
//	SubAccount     *SubAccount
//	EditKey        string
//	EditLockArgs   []byte
//	EditRecords    []Record
//	RenewExpiredAt uint64
//}
//
//type SubAccount struct {
//	Lock                 *types.Script           `json:"lock"`
//	AccountId            string                  `json:"account_id"`
//	AccountCharSet       []common.AccountCharSet `json:"account_char_set"`
//	Suffix               string                  `json:"suffix"`
//	RegisteredAt         uint64                  `json:"registered_at"`
//	ExpiredAt            uint64                  `json:"expired_at"`
//	Status               uint8                   `json:"status"`
//	Records              []Record                `json:"records"`
//	Nonce                uint64                  `json:"nonce"`
//	EnableSubAccount     uint8                   `json:"enable_sub_account"`
//	RenewSubAccountPrice uint64                  `json:"renew_sub_account_price"`
//}
//
//type SubAccountEditValue struct {
//	LockArgs  string   `json:"lock_args"`
//	Records   []Record `json:"records"`
//	ExpiredAt uint64   `json:"expired_at"`
//}
//
//func SubAccountBuilderFromTx(tx *types.Transaction) (*SubAccountBuilder, error) {
//	respMap, err := SubAccountBuilderMapFromTx(tx)
//	if err != nil {
//		return nil, err
//	}
//	for k, _ := range respMap {
//		return respMap[k], nil
//	}
//	return nil, fmt.Errorf("not exist sub account")
//}
//
//func SubAccountBuilderMapFromTx(tx *types.Transaction) (map[string]*SubAccountBuilder, error) {
//	var respMap = make(map[string]*SubAccountBuilder)
//
//	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
//		switch actionDataType {
//		case common.ActionDataTypeSubAccount:
//			builder, err := SubAccountBuilderFromBytes(dataBys)
//			if err != nil {
//				return false, err
//			}
//
//			currentSubAccount := *builder.SubAccount
//			builder.CurrentSubAccount = &currentSubAccount
//
//			editKey := string(builder.EditKey)
//			if editKey != "" {
//				builder.CurrentSubAccount.Nonce++
//			}
//			switch editKey {
//			case common.EditKeyOwner:
//				builder.CurrentSubAccount.Lock = &types.Script{
//					CodeHash: builder.SubAccount.Lock.CodeHash,
//					HashType: builder.SubAccount.Lock.HashType,
//					Args:     builder.EditValue,
//				}
//				builder.CurrentSubAccount.Records = nil
//			case common.EditKeyManager:
//				builder.CurrentSubAccount.Lock = &types.Script{
//					CodeHash: builder.SubAccount.Lock.CodeHash,
//					HashType: builder.SubAccount.Lock.HashType,
//					Args:     builder.EditValue,
//				}
//			case common.EditKeyRecords:
//				records := builder.ConvertEditValueToRecords()
//				builder.CurrentSubAccount.Records = ConvertToRecords(records)
//			case common.EditKeyExpiredAt:
//				expiredAt := builder.ConvertEditValueToExpiredAt()
//				builder.CurrentSubAccount.ExpiredAt, _ = molecule.Bytes2GoU64(expiredAt.RawData())
//			}
//
//			respMap[builder.SubAccount.AccountId] = builder
//		}
//		return true, nil
//	})
//
//	if err != nil {
//		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
//	}
//	if len(respMap) == 0 {
//		return nil, fmt.Errorf("not exist sub account")
//	}
//	return respMap, nil
//}
//
//func SubAccountBuilderFromBytes(dataBys []byte) (*SubAccountBuilder, error) {
//	var resp SubAccountBuilder
//	index, length := uint32(0), uint32(4)
//
//	signatureLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.Signature = dataBys[index : index+signatureLen]
//	index += signatureLen
//
//	signRoleLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.SignRole = dataBys[index : index+signRoleLen]
//	index += signRoleLen
//
//	prevRootLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.PrevRoot = dataBys[index : index+prevRootLen]
//	index += prevRootLen
//
//	currentRootLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.CurrentRoot = dataBys[index : index+currentRootLen]
//	index += currentRootLen
//
//	proofLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.Proof = dataBys[index : index+proofLen]
//	index += proofLen
//
//	versionLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.Version, _ = molecule.Bytes2GoU32(dataBys[index : index+versionLen])
//	index += versionLen
//
//	subAccountLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	subAccountBys := dataBys[index : index+subAccountLen]
//	index += subAccountLen
//
//	keyLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.EditKey = dataBys[index : index+keyLen]
//	index += keyLen
//
//	valueLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
//	index += length
//	resp.EditValue = dataBys[index : index+valueLen]
//	index += valueLen
//
//	switch resp.Version {
//	case common.GoDataEntityVersion1:
//		subAccount, err := ConvertToSubAccount(subAccountBys)
//		if err != nil {
//			return nil, fmt.Errorf("ConvertToSubAccount err: %s", err.Error())
//		}
//		resp.SubAccount = subAccount
//		resp.Account = subAccount.Account()
//		return &resp, nil
//	default:
//		subAccount, err := ConvertToSubAccount(subAccountBys)
//		if err != nil {
//			return nil, fmt.Errorf("ConvertToSubAccount err: %s", err.Error())
//		}
//		resp.SubAccount = subAccount
//		resp.Account = subAccount.Account()
//		return &resp, nil
//	}
//}
//
//func (s *SubAccountBuilder) ConvertToEditValue() (*SubAccountEditValue, error) {
//	var editValue SubAccountEditValue
//	editKey := string(s.EditKey)
//	switch editKey {
//	case common.EditKeyOwner, common.EditKeyManager:
//		editValue.LockArgs = common.Bytes2Hex(s.EditValue)
//	case common.EditKeyRecords:
//		records := s.ConvertEditValueToRecords()
//		editValue.Records = ConvertToRecords(records)
//	case common.EditKeyExpiredAt:
//		expiredAt := s.ConvertEditValueToExpiredAt()
//		editValue.ExpiredAt, _ = molecule.Bytes2GoU64(expiredAt.RawData())
//	default:
//		return nil, fmt.Errorf("not support edit key[%s]", editKey)
//	}
//	return &editValue, nil
//}
//
//func (s *SubAccountBuilder) ConvertEditValueToExpiredAt() *molecule.Uint64 {
//	expiredAt, _ := molecule.Uint64FromSlice(s.EditValue, true)
//	return expiredAt
//}
//
//func (s *SubAccountBuilder) ConvertEditValueToRecords() *molecule.Records {
//	records, _ := molecule.RecordsFromSlice(s.EditValue, true)
//	return records
//}

/****************************************** Parting Line ******************************************/

//func ConvertToSubAccount(slice []byte) (*SubAccount, error) {
//	subAccount, err := molecule.SubAccountFromSlice(slice, true)
//	if err != nil {
//		return nil, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
//	}
//	var tmp SubAccount
//	tmp.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
//	tmp.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
//	tmp.AccountCharSet = common.ConvertToAccountCharSets(subAccount.Account())
//	tmp.Suffix = string(subAccount.Suffix().RawData())
//	tmp.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
//	tmp.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
//	tmp.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
//	tmp.Records = ConvertToRecords(subAccount.Records())
//	tmp.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
//	tmp.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
//	tmp.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())
//
//	return &tmp, nil
//}
//
//func (s *SubAccount) ConvertToCellSubAccount() *molecule.SubAccount {
//	lock := molecule.CkbScript2MoleculeScript(s.Lock)
//	accountChars := common.ConvertToAccountChars(s.AccountCharSet)
//	accountId, _ := molecule.AccountIdFromSlice(common.Hex2Bytes(s.AccountId), true)
//	suffix := molecule.GoBytes2MoleculeBytes([]byte(s.Suffix))
//	registeredAt := molecule.GoU64ToMoleculeU64(s.RegisteredAt)
//	expiredAt := molecule.GoU64ToMoleculeU64(s.ExpiredAt)
//	status := molecule.GoU8ToMoleculeU8(s.Status)
//	records := ConvertToCellRecords(s.Records)
//	nonce := molecule.GoU64ToMoleculeU64(s.Nonce)
//	enableSubAccount := molecule.GoU8ToMoleculeU8(s.EnableSubAccount)
//	renewSubAccountPrice := molecule.GoU64ToMoleculeU64(s.RenewSubAccountPrice)
//
//	moleculeSubAccount := molecule.NewSubAccountBuilder().
//		Lock(lock).
//		Id(*accountId).
//		Account(*accountChars).
//		Suffix(suffix).
//		RegisteredAt(registeredAt).
//		ExpiredAt(expiredAt).
//		Status(status).
//		Records(*records).
//		Nonce(nonce).
//		EnableSubAccount(enableSubAccount).
//		RenewSubAccountPrice(renewSubAccountPrice).
//		Build()
//	return &moleculeSubAccount
//}
//
//func (s *SubAccount) Account() string {
//	var account string
//	for _, v := range s.AccountCharSet {
//		account += v.Char
//	}
//	return account + s.Suffix
//}
//
//func (s *SubAccount) ToH256() []byte {
//	moleculeSubAccount := s.ConvertToCellSubAccount()
//	bys, _ := blake2b.Blake256(moleculeSubAccount.AsSlice())
//	return bys
//}
//
//func (p *SubAccountParam) GenSubAccountBytes() (bys []byte) {
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Signature)))...)
//	bys = append(bys, p.Signature...)
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.SignRole)))...)
//	bys = append(bys, p.SignRole...)
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.PrevRoot)))...)
//	bys = append(bys, p.PrevRoot...)
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.CurrentRoot)))...)
//	bys = append(bys, p.CurrentRoot...)
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Proof)))...)
//	bys = append(bys, p.Proof...)
//
//	versionBys := molecule.GoU32ToMoleculeU32(SubAccountCurrentVersion)
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
//	bys = append(bys, versionBys.RawData()...)
//
//	subAccount := p.SubAccount.ConvertToCellSubAccount()
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(subAccount.AsSlice())))...)
//	bys = append(bys, subAccount.AsSlice()...)
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len([]byte(p.EditKey))))...)
//	bys = append(bys, p.EditKey...)
//
//	var editValue []byte
//	switch p.EditKey {
//	case common.EditKeyOwner, common.EditKeyManager:
//		editValue = p.EditLockArgs
//	case common.EditKeyRecords:
//		records := ConvertToCellRecords(p.EditRecords)
//		editValue = records.AsSlice()
//	case common.EditKeyExpiredAt:
//		expiredAt := molecule.GoU64ToMoleculeU64(p.RenewExpiredAt)
//		editValue = expiredAt.AsSlice()
//	}
//
//	bys = append(bys, molecule.GoU32ToBytes(uint32(len(editValue)))...)
//	bys = append(bys, editValue...)
//	return
//}
//
//func (p *SubAccountParam) NewSubAccountWitness() ([]byte, error) {
//	bys := p.GenSubAccountBytes()
//	witness := GenDasDataWitnessWithByte(common.ActionDataTypeSubAccount, bys)
//	return witness, nil
//}

// ===================== outputs data ====================

type FlagType uint8
type AutoDistribution uint8

const (
	FlagTypeDefault     FlagType = 0
	FlagTypeCustomPrice FlagType = 1
	FlagTypeCustomRule  FlagType = 255

	AutoDistributionDefault AutoDistribution = 0
	AutoDistributionEnable  AutoDistribution = 1
)

type SubAccountCellDataDetail struct {
	Action             common.DasAction
	SmtRoot            []byte           // 32
	DasProfit          uint64           // 8
	OwnerProfit        uint64           // 8
	Flag               FlagType         // 1
	CustomScriptArgs   []byte           // 32
	CustomScriptConfig []byte           // 10
	AutoDistribution   AutoDistribution // 1
	PriceRulesHash     []byte           // 10
	PreservedRulesHash []byte           // 10
}

func (s *SubAccountCellDataDetail) HasCustomScriptArgs() bool {
	defaultCustomScriptArgs := make([]byte, 32)
	if len(s.CustomScriptArgs) > 0 && bytes.Compare(defaultCustomScriptArgs, s.CustomScriptArgs) != 0 {
		return true
	}
	return false
}

func (s *SubAccountCellDataDetail) IsSameCustomScriptConfig(customScriptConfigHex string) bool {
	customScriptConfig := common.Hex2Bytes(customScriptConfigHex)
	if len(s.CustomScriptConfig) > 0 && len(customScriptConfig) > 0 {
		return bytes.Compare(customScriptConfig, s.CustomScriptConfig) == 0
	}
	return false
}

func (s *SubAccountCellDataDetail) ArgsAndConfigHash() string {
	data := append(s.CustomScriptArgs, s.CustomScriptConfig...)
	//fmt.Println(len(data), data)
	return common.Bytes2Hex(common.Blake2b(data)[10:])
}

func ConvertSubAccountCellOutputData(data []byte) (detail SubAccountCellDataDetail) {
	if len(data) >= 32 {
		detail.SmtRoot = data[:32]
	}
	if len(data) >= 40 {
		detail.DasProfit, _ = molecule.Bytes2GoU64(data[32:40])
	}
	if len(data) >= 48 {
		detail.OwnerProfit, _ = molecule.Bytes2GoU64(data[40:48])
	}
	if len(data) >= 49 {
		detail.Flag = FlagType(data[48:49][0])
	}

	switch detail.Flag {
	case FlagTypeDefault, FlagTypeCustomPrice:
		if len(data) >= 81 {
			detail.CustomScriptArgs = data[49:81]
		}
		if len(data) >= 91 {
			detail.CustomScriptConfig = data[81:91]
		}
	case FlagTypeCustomRule:
		if len(data) >= 50 {
			detail.AutoDistribution = AutoDistribution(data[49:50][0])
		}
		if len(data) >= 60 {
			detail.PriceRulesHash = data[50:60]
		}
		if len(data) >= 70 {
			detail.PreservedRulesHash = data[60:70]
		}
	}
	return
}

func BuildSubAccountCellOutputData(detail SubAccountCellDataDetail) []byte {
	dasProfit := molecule.GoU64ToMoleculeU64(detail.DasProfit)
	data := append(detail.SmtRoot, dasProfit.RawData()...)

	ownerProfit := molecule.GoU64ToMoleculeU64(detail.OwnerProfit)
	data = append(data, ownerProfit.RawData()...)

	data = append(data, uint8(detail.Flag))
	switch detail.Flag {
	case FlagTypeDefault, FlagTypeCustomPrice:
		if len(detail.CustomScriptArgs) > 0 {
			data = append(data, detail.CustomScriptArgs...)
		}
		if len(detail.CustomScriptConfig) > 0 {
			data = append(data, detail.CustomScriptConfig...)
		}
	case FlagTypeCustomRule:
		data = append(data, uint8(detail.AutoDistribution))
		if len(detail.PriceRulesHash) <= 0 {
			detail.PriceRulesHash = make([]byte, 10)
		}
		if len(detail.PriceRulesHash) <= 0 {
			detail.PriceRulesHash = make([]byte, 10)
		}
		data = append(data, detail.PriceRulesHash...)
		data = append(data, detail.PreservedRulesHash...)
	}
	return data
}

// ===================== custom script config ====================
type CustomScriptConfig struct {
	Header    string                      `json:"header"`  // 10
	Version   uint32                      `json:"version"` // 4
	Body      map[uint8]CustomScriptPrice `json:"body"`
	MaxLength uint8                       `json:"max_length"`
}

func (c *CustomScriptConfig) GetPrice(length uint8) (*CustomScriptPrice, error) {
	if length > c.MaxLength {
		length = c.MaxLength
	}
	price, ok := c.Body[length]
	if !ok {
		return &price, ErrCustomScriptPriceNotExist
	}
	return &price, nil
}

func (c *CustomScriptConfig) GetPriceBySubAccount(subAccount string) (*CustomScriptPrice, error) {
	index := strings.Index(subAccount, ".")
	if index == -1 {
		return nil, fmt.Errorf("sub-account is invalid")
	}
	accLen := common.GetAccountLength(subAccount[:index])
	return c.GetPrice(accLen)
}

type CustomScriptPrice struct {
	New   uint64 `json:"new"`
	Renew uint64 `json:"renew"`
}

const (
	Script001 = "script-001"
)

var ErrCustomScriptPriceNotExist = errors.New("CustomScriptPrice not exist")

func ConvertCustomScriptConfigByTx(tx *types.Transaction) ([]byte, *CustomScriptConfig, error) {
	for _, wit := range tx.Witnesses {
		tmp, err := ConvertCustomScriptConfig(wit)
		if err != nil {
			continue
		} else if tmp != nil {
			return wit, tmp, nil
		}
	}
	return nil, nil, ErrNotExistCustomScriptConfig
}

var ErrNotExistCustomScriptConfig = errors.New("not exist custom-script-config")

func ConvertCustomScriptConfig(wit []byte) (*CustomScriptConfig, error) {
	var res CustomScriptConfig
	res.Body = make(map[uint8]CustomScriptPrice)

	if len(wit) < 14 {
		return nil, fmt.Errorf("len is invalid")
	}
	header := wit[:10]
	script001 := []byte(Script001)
	if bytes.Compare(header, script001) != 0 {
		return nil, fmt.Errorf("header is invalid")
	}
	res.Header = string(header)
	res.Version, _ = molecule.Bytes2GoU32(wit[10:14])

	body := wit[14:]
	moleculePriceList, err := molecule.PriceConfigListFromSlice(body, true)
	if err != nil {
		return nil, fmt.Errorf("PriceConfigListFromSlice err: %s", err.Error())
	}
	for i, count := uint(0), moleculePriceList.Len(); i < count; i++ {
		price, err := molecule.PriceConfigFromSlice(moleculePriceList.Get(i).AsSlice(), true)
		if err != nil {
			return nil, fmt.Errorf("PriceConfigFromSlice err: %s", err.Error())
		}
		length, err := molecule.Bytes2GoU8(price.Length().RawData())
		if err != nil {
			return nil, fmt.Errorf("price.Length() err: %s", err.Error())
		}
		var tmp CustomScriptPrice
		tmp.New, _ = molecule.Bytes2GoU64(price.New().RawData())
		tmp.Renew, _ = molecule.Bytes2GoU64(price.Renew().RawData())

		res.Body[length] = tmp
		if res.MaxLength < length {
			res.MaxLength = length
		}
	}

	if res.Header == "" {
		return nil, ErrNotExistCustomScriptConfig
	}

	return &res, nil
}

func BuildCustomScriptConfig(csc CustomScriptConfig) (wit []byte, hash []byte) {
	wit = append(wit, []byte(csc.Header)...)
	wit = append(wit, molecule.GoU32ToBytes(csc.Version)...)

	var sortList []int
	for i, _ := range csc.Body {
		sortList = append(sortList, int(i))
	}
	sort.Ints(sortList)

	moleculePriceList := molecule.NewPriceConfigListBuilder()
	for _, v := range sortList {
		moleculePrice := molecule.NewPriceConfigBuilder()
		moleculePrice.New(molecule.GoU64ToMoleculeU64(csc.Body[uint8(v)].New))
		moleculePrice.Renew(molecule.GoU64ToMoleculeU64(csc.Body[uint8(v)].Renew))
		moleculePrice.Length(molecule.GoU8ToMoleculeU8(uint8(v)))
		moleculePriceList.Push(moleculePrice.Build())
	}
	res := moleculePriceList.Build()
	moleculePriceListBys := res.AsSlice()

	wit = append(wit, moleculePriceListBys...)
	return wit, common.Blake2b(moleculePriceListBys)[:10]
}
