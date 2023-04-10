package witness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"sort"
	"strings"
)

type (
	ReturnType      string
	ExpressionType  string
	ExpressionsType []ExpressionType
	SymbolType      string
	SymbolsType     []SymbolType
	FunctionType    string
	FunctionsType   []FunctionType
	VariableName    string
	VariablesName   []VariableName
	ValueType       string
	ValuesType      []ValueType
	CharsetType     string
	CharsetsType    []CharsetType
)

const (
	ReturnTypeBool        ReturnType = "bool"
	ReturnTypeNumber      ReturnType = "number"
	ReturnTypeString      ReturnType = "string"
	ReturnTypeStringArray ReturnType = "string[]"
	ReturnTypeUnknown     ReturnType = "unknown"

	Operator ExpressionType = "operator"
	Function ExpressionType = "function"
	Variable ExpressionType = "variable"
	Value    ExpressionType = "value"

	And SymbolType = "and"
	Or  SymbolType = "or"
	Not SymbolType = "not"
	Gt  SymbolType = ">"
	Gte SymbolType = ">="
	Lt  SymbolType = "<"
	Lte SymbolType = "<="
	Equ SymbolType = "=="

	FunctionIncludeCharts      FunctionType = "include_chars"
	FunctionOnlyIncludeCharset FunctionType = "only_include_charset"
	FunctionInList             FunctionType = "in_list"

	Account       VariableName = "account"
	AccountChars  VariableName = "account_chars"
	AccountLength VariableName = "account_length"

	Bool        ValueType = "bool"
	Uint8       ValueType = "uint8"
	Uint32      ValueType = "uint32"
	Uint64      ValueType = "uint64"
	Binary      ValueType = "binary"
	BinaryArray ValueType = "binary[]"
	String      ValueType = "string"
	StringArray ValueType = "string[]"
	Charset     ValueType = "charset_type"

	Emoji  CharsetType = "Emoji"
	Digit  CharsetType = "Digit"
	En     CharsetType = "En"
	ZhHans CharsetType = "ZhHans"
	ZhHant CharsetType = "ZhHant"
	Ja     CharsetType = "Ja"
	Ko     CharsetType = "Ko"
	Ru     CharsetType = "Ru"
	Tr     CharsetType = "Tr"
	Th     CharsetType = "Th"
	Vi     CharsetType = "Vi"
)

var (
	CharsetTypes = CharsetsType{
		Emoji,
		Digit,
		En,
		ZhHans,
		ZhHant,
		Ja,
		Ko,
		Ru,
		Tr,
		Th,
		Vi,
	}
)

func (fs FunctionsType) Include(functionType FunctionType) bool {
	for _, v := range fs {
		if v == functionType {
			return true
		}
	}
	return false
}

func (cs CharsetsType) Include(charset CharsetType) bool {
	for _, v := range cs {
		if v == charset {
			return true
		}
	}
	return false
}

type SubAccountRule struct {
	Index uint32           `json:"index"`
	Name  string           `json:"name"`
	Note  string           `json:"note"`
	Price uint64           `json:"price,omitempty"`
	Ast   ExpressionEntity `json:"ast"`
}

type SubAccountRuleSlice []SubAccountRule

type SubAccountRuleEntity struct {
	ParentAccount string
	Rules         SubAccountRuleSlice
}

type ExpressionEntity struct {
	subAccountRuleEntity *SubAccountRuleEntity

	Type        ExpressionType     `json:"type"`
	Name        string             `json:"name,omitempty"`
	Symbol      SymbolType         `json:"symbol,omitempty"`
	Expressions ExpressionEntities `json:"expressions,omitempty"`
	Value       interface{}        `json:"value,omitempty"`
	ValueType   ValueType          `json:"value_type,omitempty"`
}

type ExpressionEntities []ExpressionEntity

func NewSubAccountRuleEntity(parentAccount string) *SubAccountRuleEntity {
	return &SubAccountRuleEntity{
		ParentAccount: parentAccount,
		Rules:         make(SubAccountRuleSlice, 0),
	}
}

func (s *SubAccountRuleEntity) ParseFromJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, &s.Rules); err != nil {
		return
	}
	for _, v := range s.Rules {
		if string(v.Name) == "" {
			err = errors.New("name can't be empty")
			return
		}
		if v.Price < 0 {
			err = errors.New("price can't be negative number")
			return
		}
		if _, err = v.Ast.Process(false, ""); err != nil {
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) Hit(account string) (hit bool, index int, err error) {
	account = strings.Split(account, ".")[0]
	for idx, v := range s.Rules {
		v.Ast.subAccountRuleEntity = s
		hit, err = v.Ast.Process(true, account)
		if err != nil || hit {
			index = idx
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) ParseFromTx(tx *types.Transaction, action common.ActionDataType) error {
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		if actionDataType != action {
			return true, nil
		}
		accountRule := &SubAccountRule{}
		if err := accountRule.ParseFromWitnessData(dataBys); err != nil {
			return false, err
		}
		s.Rules = append(s.Rules, *accountRule)
		return true, nil
	})
	if err != nil {
		return err
	}
	sort.Slice(s.Rules, func(i, j int) bool {
		return s.Rules[i].Index < s.Rules[j].Index
	})
	return nil
}

func (s *SubAccountRuleEntity) ParseFromDasActionWitnessData(data [][]byte) error {
	resData := make([][]byte, 0, len(data))
	for _, v := range data {
		action := ParserWitnessAction(v)
		if action != common.ActionDataTypeSubAccountPriceRules &&
			action != common.ActionDataTypeSubAccountPreservedRules {
			return fmt.Errorf("no support action: %s", action)
		}
		if len(v) <= common.WitnessDasTableTypeEndIndex {
			return errors.New("data length error")
		}
		v = v[common.WitnessDasTableTypeEndIndex:]
		resData = append(resData, v)
	}
	return s.ParseFromWitnessData(resData)
}

func (s *SubAccountRuleEntity) ParseFromWitnessData(data [][]byte) error {
	for _, v := range data {
		accountRule := &SubAccountRule{}
		if err := accountRule.ParseFromWitnessData(v); err != nil {
			return err
		}
		s.Rules = append(s.Rules, *accountRule)
	}
	sort.Slice(s.Rules, func(i, j int) bool {
		return s.Rules[i].Index < s.Rules[j].Index
	})
	return nil
}

func (s *SubAccountRuleEntity) GenWitnessData(action common.ActionDataType) ([][]byte, error) {
	for _, v := range s.Rules {
		if string(v.Name) == "" {
			return nil, errors.New("name can't be empty")
		}
		if v.Price < 0 {
			return nil, errors.New("price can't be negative number")
		}
		if _, err := v.Ast.Process(false, ""); err != nil {
			return nil, err
		}
	}

	res := make([][]byte, 0)
	for idx, v := range s.Rules {
		itemData := make([]byte, 0)

		index := molecule.GoU32ToMoleculeU32(uint32(idx))
		itemData = append(itemData, molecule.GoU32ToBytes(uint32(len(index.RawData())))...)
		itemData = append(itemData, index.RawData()...)

		itemData = append(itemData, molecule.GoU32ToBytes(uint32(len(v.Name)))...)
		itemData = append(itemData, []byte(v.Name)...)

		itemData = append(itemData, molecule.GoU32ToBytes(uint32(len(v.Note)))...)
		itemData = append(itemData, []byte(v.Note)...)

		price := molecule.GoU64ToMoleculeU64(v.Price)
		itemData = append(itemData, molecule.GoU32ToBytes(uint32(len(price.RawData())))...)
		itemData = append(itemData, price.RawData()...)

		astData := v.Ast.GenWitnessData()
		itemData = append(itemData, molecule.GoU32ToBytes(uint32(len(astData)))...)
		itemData = append(itemData, astData...)

		res = append(res, GenDasDataWitnessWithByte(action, itemData))
	}
	return res, nil
}

func NewSubAccountRule() *SubAccountRule {
	return &SubAccountRule{}
}

func (s *SubAccountRule) Parser(data []byte) (err error) {
	if err = json.Unmarshal(data, s); err != nil {
		return
	}
	if string(s.Name) == "" {
		err = errors.New("name can't be empty")
		return
	}
	if s.Price < 0 {
		err = errors.New("price can't be negative number")
		return
	}
	_, err = s.Ast.Process(false, "")
	return
}

func (s *SubAccountRule) ParseFromWitnessData(data []byte) error {
	index, indexLen, dataLen := uint32(0), molecule.HeaderSizeUint, uint32(0)
	if int(indexLen) > len(data) {
		return fmt.Errorf("data length error: %d", len(data))
	}

	dataLen, err := molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	index += indexLen

	ruleIndex, err := molecule.Bytes2GoU32(data[index : index+dataLen])
	if err != nil {
		return err
	}
	s.Index = ruleIndex
	index += dataLen

	dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	s.Name = string(data[index+indexLen : index+indexLen+dataLen])
	index += indexLen + dataLen

	dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	s.Note = string(data[index+indexLen : index+indexLen+dataLen])
	index += indexLen + dataLen

	dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	s.Price, err = molecule.Bytes2GoU64(data[index+indexLen : index+indexLen+dataLen])
	if err != nil {
		return err
	}
	index += indexLen + dataLen

	dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	index += indexLen

	ast := &ExpressionEntity{}
	if err := ast.ParseFromWitnessData(data[index : index+dataLen]); err != nil {
		return err
	}
	index += dataLen
	s.Ast = *ast
	return nil
}

func (s *SubAccountRule) Hit(account string) (hit bool, err error) {
	account = strings.Split(account, ".")[0]
	return s.Ast.Process(true, account)
}

func (e *ExpressionEntities) GenWitnessData() []byte {
	expressions := *e
	builder := molecule.NewBytesVecBuilder()
	for i := 0; i < len(expressions); i++ {
		builder.Push(*molecule.BytesFromSliceUnchecked(expressions[i].GenWitnessData()))
	}
	bytesVec := builder.Build()
	return bytesVec.AsSlice()
}

func (e *ExpressionEntity) GenWitnessData() []byte {
	res := make([]byte, 0)
	switch e.Type {
	case Operator:
		res = append(res, 0x00)
		switch e.Symbol {
		case Not:
			res = append(res, 0x00)
		case And:
			res = append(res, 0x01)
		case Or:
			res = append(res, 0x02)
		case Gt:
			res = append(res, 0x03)
		case Gte:
			res = append(res, 0x04)
		case Lt:
			res = append(res, 0x05)
		case Lte:
			res = append(res, 0x06)
		case Equ:
			res = append(res, 0x07)
		}
		res = append(res, e.Expressions.GenWitnessData()...)
	case Function:
		res = append(res, 0x01)
		switch FunctionType(e.Name) {
		case FunctionIncludeCharts:
			res = append(res, 0x00)
		case FunctionOnlyIncludeCharset:
			res = append(res, 0x01)
		case FunctionInList:
			res = append(res, 0x02)
		}
		res = append(res, e.Expressions.GenWitnessData()...)
	case Variable:
		res = append(res, 0x02)
		switch VariableName(e.Name) {
		case Account:
			res = append(res, 0x00)
		case AccountChars:
			res = append(res, 0x01)
		case AccountLength:
			res = append(res, 0x02)
		}
	case Value:
		buf := bytes.NewBuffer([]byte{})
		res = append(res, 0x03)
		switch e.ValueType {
		case Bool:
			res = append(res, 0x00)
			v := molecule.GoU8ToMoleculeU8(gconv.Uint8(e.Value))
			buf.Write(v.AsSlice())
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case Uint8:
			res = append(res, 0x01)
			v := molecule.GoU8ToMoleculeU8(gconv.Uint8(e.Value))
			buf.Write(v.AsSlice())
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case Uint32:
			res = append(res, 0x02)
			v := molecule.GoU32ToMoleculeU32(gconv.Uint32(e.Value))
			buf.Write(v.AsSlice())
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case Uint64:
			res = append(res, 0x03)
			v := molecule.GoU64ToMoleculeU64(gconv.Uint64(e.Value))
			buf.Write(v.AsSlice())
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case Binary:
			res = append(res, 0x04)
			buf.Write(common.Hex2Bytes(gconv.String(e.Value)))
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case BinaryArray:
			res = append(res, 0x05)
			strArrays := gconv.Strings(e.Value)
			builder := molecule.NewBytesVecBuilder()
			for _, v := range strArrays {
				builder.Push(*molecule.BytesFromSliceUnchecked(common.Hex2Bytes(v)))
			}
			bytesVec := builder.Build()
			res = append(res, bytesVec.AsSlice()...)
		case String:
			res = append(res, 0x06)
			buf.Write(gconv.Bytes(e.Value))
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case StringArray:
			res = append(res, 0x07)
			strArrays := gconv.Strings(e.Value)
			builder := molecule.NewBytesVecBuilder()
			for _, v := range strArrays {
				builder.Push(*molecule.BytesFromSliceUnchecked([]byte(v)))
			}
			bytesVec := builder.Build()
			res = append(res, bytesVec.AsSlice()...)
		case Charset:
			res = append(res, 0x08)
			buf.Write(gconv.Bytes(e.Value))
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		}
	}
	return res
}

func (e *ExpressionEntity) ReturnType() ReturnType {
	if e.Type == Operator || e.Type == Function || e.Type == Value && e.ValueType == Bool {
		return ReturnTypeBool
	}

	if e.Type == Value && (e.ValueType == Uint8 || e.ValueType == Uint32 || e.ValueType == Uint64) ||
		e.Type == Variable && VariableName(e.Name) == AccountLength {
		return ReturnTypeNumber
	}

	if e.Type == Value && e.ValueType == String ||
		e.Type == Variable && VariableName(e.Name) == Account ||
		e.Type == Value && e.ValueType == Binary {
		return ReturnTypeString
	}
	if e.Type == Variable && VariableName(e.Name) == AccountChars ||
		e.Type == Value && e.ValueType == StringArray ||
		e.Type == Value && e.ValueType == BinaryArray {
		return ReturnTypeStringArray
	}
	return ReturnTypeUnknown
}

func IsSameReturnType(i, j ExpressionEntity) bool {
	ir := i.ReturnType()
	jr := j.ReturnType()
	return ir == jr && ir != ReturnTypeUnknown
}

func (e *ExpressionEntity) GetNumberValue(account string) float64 {
	if e.Type == Variable && VariableName(e.Name) == AccountLength {
		return float64(len([]rune(account)))
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeNumber {
		return gconv.Float64(e.Value)
	}
	return 0
}

func (e *ExpressionEntity) GetStringValue(account string) string {
	if e.Type == Variable && VariableName(e.Name) == Account {
		return account
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeString {
		return gconv.String(e.Value)
	}
	return ""
}

func (e *ExpressionEntity) ProcessOperator(checkHit bool, account string) (hit bool, err error) {
	switch e.Symbol {
	case And:
		for _, exp := range e.Expressions {
			exp.subAccountRuleEntity = e.subAccountRuleEntity
			rtType := exp.ReturnType()
			if rtType != ReturnTypeBool {
				return false, errors.New("operator 'and' every expression must be bool return")
			}
			hit, err := exp.Process(checkHit, account)
			if err != nil {
				return false, err
			}
			if checkHit && !hit {
				return false, nil
			}
		}
		return true, nil
	case Or:
		for _, exp := range e.Expressions {
			exp.subAccountRuleEntity = e.subAccountRuleEntity
			rtType := exp.ReturnType()
			if rtType != ReturnTypeBool {
				return false, errors.New("operator 'and' every expression must be bool return")
			}
			hit, err := exp.Process(checkHit, account)
			if err != nil {
				return false, err
			}
			if checkHit && hit {
				return true, nil
			}
		}
		return true, nil
	case Not:
		if len(e.Expressions) != 1 {
			return false, errors.New("operator not must have one expression")
		}
		exp := e.Expressions[0]
		exp.subAccountRuleEntity = e.subAccountRuleEntity

		rtType := exp.ReturnType()
		if rtType != ReturnTypeBool {
			return false, errors.New("operator 'not' expression must be bool return")
		}
		hit, err := exp.Process(checkHit, account)
		if err != nil {
			return false, err
		}
		if !hit {
			return true, nil
		}
	case Gt, Gte, Lt, Lte, Equ:
		if len(e.Expressions) != 2 {
			return false, errors.New("operator not must have two expression")
		}
		left := e.Expressions[0]
		right := e.Expressions[1]
		if !IsSameReturnType(left, right) {
			return false, errors.New("the comparison type operation must have same types on both sides")
		}
		left.subAccountRuleEntity = e.subAccountRuleEntity
		right.subAccountRuleEntity = e.subAccountRuleEntity

		switch left.ReturnType() {
		case ReturnTypeNumber:
			leftVal := left.GetNumberValue(account)
			rightVal := right.GetNumberValue(account)
			if e.Symbol == Gt {
				return leftVal > rightVal, nil
			}
			if e.Symbol == Gte {
				return leftVal >= rightVal, nil
			}
			if e.Symbol == Lt {
				return leftVal < rightVal, nil
			}
			if e.Symbol == Lte {
				return leftVal <= rightVal, nil
			}
			if e.Symbol == Equ {
				return leftVal == rightVal, nil
			}
		case ReturnTypeString:
			leftVal := left.GetStringValue(account)
			rightVal := right.GetStringValue(account)
			if e.Symbol == Gt {
				return leftVal > rightVal, nil
			}
			if e.Symbol == Gte {
				return leftVal >= rightVal, nil
			}
			if e.Symbol == Lt {
				return leftVal < rightVal, nil
			}
			if e.Symbol == Lte {
				return leftVal <= rightVal, nil
			}
			if e.Symbol == Equ {
				return leftVal == rightVal, nil
			}
		default:
			return false, fmt.Errorf("type %s is not currently supported", left.ReturnType())
		}
	default:
		err = fmt.Errorf("symbol %s can't be support", e.Symbol)
	}
	return
}

func (e *ExpressionEntity) Process(checkHit bool, account string) (hit bool, err error) {
	switch e.Type {
	case Function:
		funcName := FunctionType(e.Name)
		switch funcName {
		case FunctionIncludeCharts:
			hit, err = e.handleFunctionIncludeCharts(checkHit, account)
		case FunctionInList:
			hit, err = e.handleFunctionInList(checkHit, account)
		case FunctionOnlyIncludeCharset:
			hit, err = e.handleFunctionOnlyIncludeCharset(checkHit, account)
		default:
			err = fmt.Errorf("function %s can't be support", e.Name)
			return
		}
		if hit && checkHit || err != nil {
			return
		}
	case Operator:
		hit, err = e.ProcessOperator(checkHit, account)
	case Value, Variable:
		err = fmt.Errorf("can't Process %s", e.Type)
	default:
		err = fmt.Errorf("expression %s can't be support", e.Type)
		return
	}
	return
}

func (e *ExpressionEntity) ParseFromWitnessData(data []byte) error {
	index, indexLen, dataLen := uint32(0), molecule.HeaderSizeUint, uint32(0)

	var err error
	expType := data[0]
	index += 1

	switch expType {
	case 0x00:
		e.Type = Operator
		symbol := data[index : index+1][0]
		switch symbol {
		case 0x00:
			e.Symbol = Not
		case 0x01:
			e.Symbol = And
		case 0x02:
			e.Symbol = Or
		case 0x03:
			e.Symbol = Gt
		case 0x04:
			e.Symbol = Gte
		case 0x05:
			e.Symbol = Lt
		case 0x06:
			e.Symbol = Lte
		case 0x07:
			e.Symbol = Equ
		}
		index += 1

		data = data[index:]

		expressions := make([]ExpressionEntity, 0)

		bytesVec := molecule.BytesVecFromSliceUnchecked(data)
		for i := uint(0); i < bytesVec.ItemCount(); i++ {
			getBytes := bytesVec.Get(i)
			exp := &ExpressionEntity{}
			if err := exp.ParseFromWitnessData(getBytes.AsSlice()); err != nil {
				return err
			}
			expressions = append(expressions, *exp)
		}
		e.Expressions = expressions
	case 0x01:
		e.Type = Function

		name := data[index : index+1][0]
		switch name {
		case 0x00:
			e.Name = string(FunctionIncludeCharts)
		case 0x01:
			e.Name = string(FunctionIncludeCharts)
		case 0x02:
			e.Name = string(FunctionInList)
		}
		index += 1

		data = data[index:]

		expressions := make([]ExpressionEntity, 0)
		bytesVec := molecule.BytesVecFromSliceUnchecked(data)
		for i := uint(0); i < bytesVec.ItemCount(); i++ {
			getBytes := bytesVec.Get(i)
			exp := &ExpressionEntity{}
			if err := exp.ParseFromWitnessData(getBytes.AsSlice()); err != nil {
				return err
			}
			expressions = append(expressions, *exp)
		}
		e.Expressions = expressions
	case 0x02:
		e.Type = Variable

		varName := data[index : index+1][0]
		switch varName {
		case 0x00:
			e.Name = string(Account)
		case 0x01:
			e.Name = string(AccountChars)
		case 0x02:
			e.Name = string(AccountLength)
		}
		index += 1
	case 0x03:
		e.Type = Value

		valueType := data[index : index+1][0]
		index += 1

		switch valueType {
		case 0x00:
			e.ValueType = Bool
			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			u8, err := molecule.Bytes2GoU8(data[index : index+dataLen])
			if err != nil {
				return err
			}
			e.Value = gconv.Bool(u8)
			index += dataLen

		case 0x01:
			e.ValueType = Uint8

			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			e.Value, err = molecule.Bytes2GoU8(data[index : index+dataLen])
			if err != nil {
				return err
			}
			index += dataLen

		case 0x02:
			e.ValueType = Uint32

			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			e.Value, err = molecule.Bytes2GoU32(data[index : index+dataLen])
			if err != nil {
				return err
			}
			index += dataLen

		case 0x03:
			e.ValueType = Uint64

			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			e.Value, err = molecule.Bytes2GoU64(data[index : index+dataLen])
			if err != nil {
				return err
			}
			index += dataLen

		case 0x04:
			e.ValueType = Binary

			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			e.Value = common.Bytes2Hex(data[index : index+dataLen])
			index += dataLen

		case 0x05:
			e.ValueType = BinaryArray
			data = data[index:]
			strArrays := make([]string, 0)
			bytesVec := molecule.BytesVecFromSliceUnchecked(data)
			for i := uint(0); i < bytesVec.ItemCount(); i++ {
				getBytes := bytesVec.Get(i)
				strArrays = append(strArrays, common.Bytes2Hex(getBytes.AsSlice()))
			}
			e.Value = strArrays
		case 0x06:
			e.ValueType = String
			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen
			e.Value = string(data[index : index+dataLen])
			index += dataLen

		case 0x07:
			e.ValueType = StringArray
			data = data[index:]
			strArrays := make([]string, 0)
			bytesVec := molecule.BytesVecFromSliceUnchecked(data)
			for i := uint(0); i < bytesVec.ItemCount(); i++ {
				getBytes := bytesVec.Get(i)
				strArrays = append(strArrays, string(getBytes.AsSlice()))
			}
			e.Value = strArrays

		case 0x08:
			e.ValueType = Charset

			dataLen, err = molecule.Bytes2GoU32(data[index : index+indexLen])
			if err != nil {
				return err
			}
			index += indexLen

			e.Value = string(data[index : index+dataLen])
			index += dataLen
		}
	}
	return nil
}

func (e *ExpressionEntity) handleFunctionIncludeCharts(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Expressions[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := e.Expressions[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || value.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Name)
		return
	}
	if !checkHit {
		return
	}

	for _, v := range strArray {
		if strings.Contains(account, v) {
			hit = true
			return
		}
	}
	return
}

func (e *ExpressionEntity) handleFunctionInList(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	acc := e.Expressions[0]
	if acc.Type != Variable || VariableName(acc.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}
	value := e.Expressions[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || (value.ValueType != BinaryArray && value.ValueType != StringArray) {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Name)
		return
	}

	if !checkHit {
		return
	}

	subAccount := fmt.Sprintf("%s.%s", account, e.subAccountRuleEntity.ParentAccount)
	subAccountId := common.Bytes2Hex(common.GetAccountIdByAccount(subAccount))
	for _, v := range strArray {
		switch value.ValueType {
		case StringArray:
			if v == account {
				hit = true
				return
			}
		case BinaryArray:
			if v == subAccountId {
				hit = true
				return
			}
		}
	}
	return
}

func (e *ExpressionEntity) handleFunctionOnlyIncludeCharset(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Expressions[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}
	value := e.Expressions[1]
	val := gconv.String(value.Value)
	if val == "" ||
		value.Type != Value ||
		value.ValueType != Charset ||
		!CharsetTypes.Include(CharsetType(val)) {
		err = fmt.Errorf("function %s args[1] value must be one of: %#v", e.Name, CharsetTypes)
		return
	}
	if !checkHit {
		return
	}
	inputAccCharts := []rune(account)

	switch CharsetType(val) {
	case Emoji:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeEmojiMap[string(v)]; !ok {
				return
			}
		}
	case Digit:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeDigitMap[string(v)]; !ok {
				return
			}
		}
	case En:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeEnMap[string(v)]; !ok {
				return
			}
		}
	case ZhHans:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeHanSMap[string(v)]; !ok {
				return
			}
		}
	case ZhHant:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeHanTMap[string(v)]; !ok {
				return
			}
		}
	case Ja:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeJaMap[string(v)]; !ok {
				return
			}
		}
	case Ko:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeKoMap[string(v)]; !ok {
				return
			}
		}
	case Ru:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeRuMap[string(v)]; !ok {
				return
			}
		}
	case Tr:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeTrMap[string(v)]; !ok {
				return
			}
		}
	case Th:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeThMap[string(v)]; !ok {
				return
			}
		}
	case Vi:
		for _, v := range inputAccCharts {
			if _, ok := common.CharSetTypeViMap[string(v)]; !ok {
				return
			}
		}
	}
	hit = true
	return
}
