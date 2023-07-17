package witness

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"reflect"
	"strings"
)

var (
	log                   = mylog.NewLogger("witness", mylog.LevelDebug)
	ErrDataEntityOptIsNil = errors.New("DataEntityOpt is nil")
	ErrNotExistWitness    = errors.New("the witness wanted not exist")

	DataEntityVersion1 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion1)
	DataEntityVersion2 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion2)
	DataEntityVersion3 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion3)
	DataEntityVersion4 = molecule.GoU32ToMoleculeU32(common.GoDataEntityVersion4)
)

const DasWitnessTagName = "witness"

func GetWitnessDataFromTx(tx *types.Transaction, handle FuncParseWitness) error {
	inputsSize := len(tx.Inputs)
	witnessesSize := len(tx.Witnesses)
	for i := inputsSize; i < witnessesSize; i++ {
		dataBys := tx.Witnesses[i]
		if len(dataBys) <= common.WitnessDasTableTypeEndIndex+1 {
			continue
		} else if string(dataBys[0:common.WitnessDasCharLen]) != common.WitnessDas {
			continue
		} else {
			actionDataType := common.Bytes2Hex(dataBys[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])
			if goON, err := handle(actionDataType, dataBys[common.WitnessDasTableTypeEndIndex:], i); err != nil {
				return err
			} else if !goON {
				return nil
			}
		}
	}
	return nil
}

type FuncParseWitness func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error)

func getDataEntityOpt(dataBys []byte, dataType common.DataType) (*molecule.DataEntityOpt, *molecule.DataEntity, error) {
	data, err := molecule.DataFromSlice(dataBys, true)
	if err != nil {
		return nil, nil, fmt.Errorf("DataFromSlice err: %s", err.Error())
	}
	var dataEntityOpt *molecule.DataEntityOpt
	switch dataType {
	case common.DataTypeNew:
		dataEntityOpt = data.New()
	case common.DataTypeOld:
		dataEntityOpt = data.Old()
	case common.DataTypeDep:
		dataEntityOpt = data.Dep()
	}
	if dataEntityOpt == nil || dataEntityOpt.IsNone() {
		return nil, nil, ErrDataEntityOptIsNil
	}
	dataEntity, err := molecule.DataEntityFromSlice(dataEntityOpt.AsSlice(), true)
	if err != nil {
		return nil, nil, fmt.Errorf("DataEntityFromSlice err: %s", err.Error())
	}

	return dataEntityOpt, dataEntity, nil
}

func GenDasDataWitness(action common.ActionDataType, data *molecule.Data) []byte {
	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(action)...)
	tmp = append(tmp, data.AsSlice()...)
	return tmp
}

func GenDasDataWitnessWithByte(action common.ActionDataType, data []byte) []byte {
	tmp := append([]byte(common.WitnessDas), common.Hex2Bytes(action)...)
	tmp = append(tmp, data...)
	return tmp
}

type DasWitness interface {
	Gen() ([]byte, error)
	Parse([]byte) (DasWitness, error)
}

var TypeOfDasWitness = reflect.TypeOf((*DasWitness)(nil)).Elem()

func GenDasDataWitnessWithStruct(action common.ActionDataType, obj interface{}) ([]byte, error) {
	res := append([]byte(common.WitnessDas), common.Hex2Bytes(action)...)
	data, err := GenWitnessData(obj)
	if err != nil {
		return nil, err
	}
	res = append(res, data...)
	return res, nil
}

func GenWitnessData(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nil, errors.New("obj can't be nil")
	}
	v := reflect.ValueOf(obj)
	if v.IsNil() {
		return nil, errors.New("obj can't be nil")
	}
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Type().Kind() != reflect.Struct {
		return nil, errors.New("obj must struct pointer or struct")
	}

	res := make([]byte, 0)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		tag := v.Type().Field(i).Tag

		if !f.CanInterface() {
			return nil, fmt.Errorf("field: %s can't Interface()", f)
		}

		if f.Type().Implements(TypeOfDasWitness) {
			data, err := f.Convert(TypeOfDasWitness).Interface().(DasWitness).Gen()
			if err != nil {
				return nil, err
			}
			res = append(res, molecule.GoU32ToBytes(uint32(len(data)))...)
			res = append(res, data...)
			continue
		}

		switch f.Type().Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if hasOmitempty(tag) && f.Uint() == 0 {
				res = append(res, molecule.GoU32ToBytes(0)...)
				continue
			}
			byteBuf := bytes.NewBuffer([]byte{})
			var value interface{}
			switch f.Type().Kind() {
			case reflect.Uint8:
				value = uint8(f.Uint())
			case reflect.Uint16:
				value = uint16(f.Uint())
			case reflect.Uint32:
				value = uint32(f.Uint())
			case reflect.Uint64:
				value = f.Uint()
			}
			if err := binary.Write(byteBuf, binary.LittleEndian, value); err != nil {
				return nil, err
			}
			res = append(res, molecule.GoU32ToBytes(uint32(byteBuf.Len()))...)
			res = append(res, byteBuf.Bytes()...)
		case reflect.Slice:
			if f.Type().Elem().Kind() != reflect.Uint8 {
				return nil, fmt.Errorf("kind: [%s]{%s} no support now", reflect.Slice, f.Type().Elem().Kind())
			}
			res = append(res, molecule.GoU32ToBytes(uint32(f.Len()))...)
			res = append(res, f.Bytes()...)
		case reflect.String:
			res = append(res, molecule.GoU32ToBytes(uint32(f.Len()))...)
			res = append(res, []byte(f.String())...)
		}
	}
	return res, nil
}

func hasOmitempty(tag reflect.StructTag) bool {
	jsonTag := strings.TrimSpace(tag.Get(DasWitnessTagName))
	if jsonTag == "" {
		return false
	}
	parts := strings.Split(jsonTag, ",")
	if len(parts) > 1 && parts[1] == "omitempty" {
		return true
	}
	return false
}
