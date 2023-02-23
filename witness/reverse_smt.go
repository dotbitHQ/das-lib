package witness

import (
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"reflect"
	"strings"
)

type ReverseSmtBuilder struct{}

func NewReverseSmtBuilder() *ReverseSmtBuilder {
	return &ReverseSmtBuilder{}
}

func (b *ReverseSmtBuilder) FromBytes(bs []byte) (*ReverseSmtRecord, error) {
	var err error
	var res ReverseSmtRecord
	index, indexLen, dataLen := uint32(0), uint32(4), uint32(0)

	if int(indexLen) > len(bs) {
		return nil, fmt.Errorf("data length error: %d", len(bs))
	}

	v := reflect.ValueOf(res)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}

		dataLen, err = molecule.Bytes2GoU32(bs[index : index+indexLen])
		if err != nil {
			return nil, err
		}
		if dataLen == 0 {
			index = index + indexLen
			continue
		}

		dataBs := bs[index+indexLen : index+indexLen+dataLen]

		switch f.Type().Kind() {
		case reflect.Uint8:
			u8, err := molecule.Bytes2GoU8(dataBs)
			if err != nil {
				return nil, err
			}
			f.Set(reflect.ValueOf(u8))
		case reflect.Uint32:
			u32, err := molecule.Bytes2GoU32(dataBs)
			if err != nil {
				return nil, err
			}
			f.Set(reflect.ValueOf(u32))
		case reflect.Uint64:
			u64, err := molecule.Bytes2GoU64(dataBs)
			if err != nil {
				return nil, err
			}
			f.Set(reflect.ValueOf(u64))
		case reflect.Slice:
			if f.Type().Elem().Kind() == reflect.Uint8 {
				f.Set(reflect.ValueOf(dataBs))
			}
		case reflect.String:
			f.Set(reflect.ValueOf(string(dataBs)))
		}
		index = index + indexLen + dataLen
	}
	return &res, nil
}

func (b *ReverseSmtBuilder) FromTx(tx *types.Transaction) ([]*ReverseSmtRecord, error) {
	resp := make([]*ReverseSmtRecord, 0)
	m := make(map[string]int)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeReverseSmt:
			reverseSmt, err := b.FromBytes(dataBys)
			if err != nil {
				return false, err
			}
			idx, ok := m[reverseSmt.Address]
			if ok {
				resp[idx] = reverseSmt
			} else {
				resp = append(resp, reverseSmt)
				m[reverseSmt.Address] = len(resp) - 1
			}
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("not exist reverse record")
	}
	return resp, nil
}

func (b *ReverseSmtBuilder) GenWitness(record *ReverseSmtRecord) ([]byte, error) {
	return record.GenWitness()
}

type ReverseSmtRecord struct {
	Version     ReverseSmtRecordVersion
	Action      ReverseSmtRecordAction
	Signature   string
	SignType    uint8
	Address     string
	Proof       string
	PrevNonce   uint32 `json:",omitempty"`
	PrevAccount string
	NextRoot    []byte
	NextAccount string
}

type ReverseSmtRecordVersion uint32

const (
	ReverseSmtRecordVersion1 ReverseSmtRecordVersion = 1
)

type ReverseSmtRecordAction string

const (
	ReverseSmtRecordActionUpdate ReverseSmtRecordAction = "update"
	ReverseSmtRecordActionRemove ReverseSmtRecordAction = "remove"
)

func (r *ReverseSmtRecord) GenBytes() ([]byte, error) {
	if r.Action != ReverseSmtRecordActionUpdate &&
		r.Action != ReverseSmtRecordActionRemove {
		return nil, errors.New("action must be update or remove")
	}

	v := reflect.ValueOf(r)
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	res := make([]byte, 0)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}
		tag := v.Type().Field(i).Tag.Get("json")

		switch f.Type().Kind() {
		case reflect.Uint8:
			value := uint8(f.Uint())
			bs := molecule.GoU8ToMoleculeU8(value)
			if tag != "" {
				parts := strings.Split(tag, ",")
				if len(parts) > 1 && parts[1] == "omitempty" && value == 0 {
					res = append(res, molecule.GoU32ToBytes(0)...)
					continue
				}
			}
			res = append(res, molecule.GoU32ToBytes(uint32(len(bs.RawData())))...)
			res = append(res, bs.RawData()...)
		case reflect.Uint32:
			value := uint32(f.Uint())
			bs := molecule.GoU32ToMoleculeU32(value)
			if tag != "" {
				parts := strings.Split(tag, ",")
				if len(parts) > 1 && parts[1] == "omitempty" && value == 0 {
					res = append(res, molecule.GoU32ToBytes(0)...)
					continue
				}
			}
			res = append(res, molecule.GoU32ToBytes(uint32(len(bs.RawData())))...)
			res = append(res, bs.RawData()...)
		case reflect.Uint64:
			value := f.Uint()
			bs := molecule.GoU64ToMoleculeU64(value)
			if tag != "" {
				parts := strings.Split(tag, ",")
				if len(parts) > 1 && parts[1] == "omitempty" && value == 0 {
					res = append(res, molecule.GoU32ToBytes(0)...)
					continue
				}
			}
			res = append(res, molecule.GoU32ToBytes(uint32(len(bs.RawData())))...)
			res = append(res, bs.RawData()...)
		case reflect.Slice:
			if f.Type().Elem().Kind() == reflect.Uint8 {
				res = append(res, molecule.GoU32ToBytes(uint32(f.Len()))...)
				res = append(res, f.Bytes()...)
			}
		case reflect.String:
			res = append(res, molecule.GoU32ToBytes(uint32(len([]byte(f.String()))))...)
			res = append(res, []byte(f.String())...)
		}
	}
	return res, nil
}

func (r *ReverseSmtRecord) GenWitness() ([]byte, error) {
	dataBys, err := r.GenBytes()
	if err != nil {
		return nil, fmt.Errorf("GenSubAccountNewBytes err: %s", err.Error())
	}
	witness := GenDasDataWitnessWithByte(common.ActionDataTypeReverseSmt, dataBys)
	return witness, nil
}
