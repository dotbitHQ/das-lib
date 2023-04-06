package witness

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nervosnetwork/ckb-sdk-go/types"
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

type ExpressionEntity struct {
	Type        ExpressionType     `json:"type"`
	Name        string             `json:"name,omitempty"`
	Symbol      SymbolType         `json:"symbol,omitempty"`
	Expressions ExpressionEntities `json:"expressions,omitempty"`
	Arguments   ExpressionEntities `json:"arguments,omitempty"`
	Value       interface{}        `json:"value,omitempty"`
	ValueType   ValueType          `json:"value_type,omitempty"`
}

type ExpressionEntities []ExpressionEntity

func NewSubAccountRuleSlice() *SubAccountRuleSlice {
	return &SubAccountRuleSlice{}
}

func (s *SubAccountRuleSlice) Parser(data []byte) (err error) {
	if err = json.Unmarshal(data, s); err != nil {
		return
	}
	for _, v := range *s {
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

func (s *SubAccountRuleSlice) Hit(account string) (hit bool, index int, err error) {
	account = strings.Split(account, ".")[0]
	for idx, v := range *s {
		hit, err = v.Ast.Process(true, account)
		if err != nil || hit {
			index = idx
			return
		}
	}
	return
}

func (s *SubAccountRuleSlice) ParseFromTx(tx *types.Transaction, action common.ActionDataType) error {
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		if actionDataType != action {
			return true, nil
		}
		ns := NewSubAccountRuleSlice()
		if err := ns.ParseFromWitnessData(dataBys); err != nil {
			return false, err
		}
		*s = append(*s, *ns...)
		return true, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *SubAccountRuleSlice) ParseFromWitnessData(data []byte) error {
	res := make(SubAccountRuleSlice, 0)
	if err := parseArrayWitnessData(data, func(data []byte) error {
		accountRule := &SubAccountRule{}
		if err := accountRule.ParseFromWitnessData(data); err != nil {
			return err
		}
		res = append(res, *accountRule)
		return nil
	}); err != nil {
		return err
	}
	*s = res
	return nil
}

func (s *SubAccountRuleSlice) GenWitnessData() []byte {
	subAccountSlice := *s
	count := len(subAccountSlice)
	itemDatas := make([][]byte, 0, count)

	for idx, v := range subAccountSlice {
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

		itemDatas = append(itemDatas, itemData)
	}
	res := genArrayWitnessData(itemDatas)
	return res
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
	itemDatas := make([][]byte, 0, len(expressions))
	for i := 0; i < len(expressions); i++ {
		itemDatas = append(itemDatas, expressions[i].GenWitnessData())
	}
	res := genArrayWitnessData(itemDatas)
	return res
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
		res = append(res, e.Arguments.GenWitnessData()...)
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
			buf.Write(gconv.Bytes(e.Value))
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case BinaryArray:
			res = append(res, 0x05)
			strArrays := gconv.Strings(e.Value)
			bsArray := make([][]byte, 0, len(strArrays))
			for _, v := range strArrays {
				bsArray = append(bsArray, common.Hex2Bytes(v))
			}
			res = append(res, genArrayWitnessData(bsArray)...)
		case String:
			res = append(res, 0x06)
			buf.Write(gconv.Bytes(e.Value))
			res = append(res, molecule.GoU32ToBytes(uint32(buf.Len()))...)
			res = append(res, buf.Bytes()...)
		case StringArray:
			res = append(res, 0x07)
			strArrays := gconv.Strings(e.Value)
			bsArray := make([][]byte, 0, len(strArrays))
			for _, v := range strArrays {
				bsArray = append(bsArray, []byte(v))
			}
			res = append(res, genArrayWitnessData(bsArray)...)
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
			hit, err = handleFunctionIncludeCharts(e, checkHit, account)
		case FunctionInList:
			hit, err = handleFunctionInList(e, checkHit, account)
		case FunctionOnlyIncludeCharset:
			hit, err = handleFunctionOnlyIncludeCharset(e, checkHit, account)
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
		if err := parseArrayWitnessData(data, func(data []byte) error {
			exp := &ExpressionEntity{}
			if err := exp.ParseFromWitnessData(data); err != nil {
				return err
			}
			expressions = append(expressions, *exp)
			return nil
		}); err != nil {
			return err
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
		if err := parseArrayWitnessData(data, func(data []byte) error {
			exp := &ExpressionEntity{}
			if err := exp.ParseFromWitnessData(data); err != nil {
				return err
			}
			expressions = append(expressions, *exp)
			return nil
		}); err != nil {
			return err
		}
		e.Arguments = expressions

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
			if err := parseArrayWitnessData(data, func(data []byte) error {
				strArrays = append(strArrays, common.Bytes2Hex(data))
				return nil
			}); err != nil {
				return err
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
			if err := parseArrayWitnessData(data, func(data []byte) error {
				strArrays = append(strArrays, string(data))
				return nil
			}); err != nil {
				return err
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

func handleFunctionIncludeCharts(exp *ExpressionEntity, checkHit bool, account string) (hit bool, err error) {
	if len(exp.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", exp.Name)
		return
	}
	accCharts := exp.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := exp.Arguments[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || value.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", exp.Name)
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

func handleFunctionInList(exp *ExpressionEntity, checkHit bool, account string) (hit bool, err error) {
	if len(exp.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", exp.Name)
		return
	}
	acc := exp.Arguments[0]
	if acc.Type != Variable || VariableName(acc.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return

	}
	value := exp.Arguments[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || (value.ValueType != BinaryArray && value.ValueType != StringArray) {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", exp.Name)
		return
	}

	if !checkHit {
		return
	}

	accBinary := common.Bytes2Hex(common.Blake2b([]byte(account))[:20])
	for _, v := range strArray {
		switch value.ValueType {
		case StringArray:
			if v == account {
				hit = true
				return
			}
		case BinaryArray:
			if v == accBinary {
				hit = true
				return
			}
		}
	}
	return
}

func handleFunctionOnlyIncludeCharset(exp *ExpressionEntity, checkHit bool, account string) (hit bool, err error) {
	if len(exp.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", exp.Name)
		return
	}
	accCharts := exp.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}
	value := exp.Arguments[1]
	val := gconv.String(value.Value)
	if val == "" ||
		value.Type != Value ||
		value.ValueType != Charset ||
		!CharsetTypes.Include(CharsetType(val)) {
		err = fmt.Errorf("function %s args[1] value must be one of: %#v", exp.Name, CharsetTypes)
		return
	}
	if !checkHit {
		return
	}
	inputAccCharts := []rune(account)

	switch CharsetType(val) {
	case Emoji:
		for _, v := range inputAccCharts {
			if !isEmoji(v) {
				return
			}
		}
	case Digit:
		if !govalidator.IsNumeric(account) {
			return
		}
	case En:
		if !govalidator.IsAlpha(account) {
			return
		}
	case ZhHans:
		for _, v := range inputAccCharts {
			if !isSimplifiedChar(v) {
				return
			}
		}
	case ZhHant:
		for _, v := range inputAccCharts {
			if !isTraditionalChar(v) {
				return
			}
		}
	case Ja:
		for _, v := range inputAccCharts {
			if !isJapaneseChar(v) {
				return
			}
		}
	case Ko:
		for _, v := range inputAccCharts {
			if !isKoreanChar(v) {
				return
			}
		}
	case Ru:
		for _, v := range inputAccCharts {
			if !isRussianChar(v) {
				return
			}
		}
	case Tr:
		for _, v := range inputAccCharts {
			if !isTurkishChar(v) {
				return
			}
		}
	case Th:
		for _, v := range inputAccCharts {
			if !isThaiChar(v) {
				return
			}
		}
	case Vi:
		for _, v := range inputAccCharts {
			if !isVietnameseChar(v) {
				return
			}
		}
	}
	hit = true
	return
}

func genArrayWitnessData(bsArray [][]byte) []byte {
	res := make([]byte, 0)

	count := len(bsArray)
	if count == 0 {
		length := molecule.GoU32ToMoleculeU32(molecule.HeaderSizeUint)
		res = append(res, length.RawData()...)
		return res
	}

	totalSize := molecule.HeaderSizeUint * uint32(count+1)
	offsets := make([]uint32, 0, count)
	offsets = append(offsets, totalSize)

	for i := 1; i < count; i++ {
		totalSize += uint32(len(bsArray[i-1]))
		offsets = append(offsets, offsets[i-1]+uint32(len(bsArray[i-1])))
	}
	totalSize += uint32(len(bsArray[count-1]))

	totalSizeBs := molecule.GoU32ToMoleculeU32(totalSize)
	res = append(res, totalSizeBs.RawData()...)

	for i := 0; i < count; i++ {
		offset := molecule.GoU32ToMoleculeU32(offsets[i])
		res = append(res, offset.RawData()...)
	}
	for i := 0; i < count; i++ {
		res = append(res, bsArray[i]...)
	}
	return res
}

func parseArrayWitnessData(data []byte, fn func(data []byte) error) error {
	index, indexLen, dataLen := uint32(0), molecule.HeaderSizeUint, uint32(0)

	dataLen, err := molecule.Bytes2GoU32(data[index : index+indexLen])
	if err != nil {
		return err
	}
	index += indexLen

	var number int
	if dataLen > molecule.HeaderSizeUint {
		countBytes, err := molecule.Bytes2GoU32(data[index : index+indexLen])
		if err != nil {
			return err
		}
		number = int(countBytes)/4 - 1
		index += indexLen
	}

	for i := 0; i < number; i++ {
		startIndex := indexLen * uint32(i+1)
		start, err := molecule.Bytes2GoU32(data[startIndex : startIndex+indexLen])
		if err != nil {
			return err
		}

		if i == number-1 {
			if err := fn(data[start:]); err != nil {
				return err
			}
		} else {
			endIndex := startIndex + indexLen
			end, err := molecule.Bytes2GoU32(data[endIndex : endIndex+indexLen])
			if err != nil {
				return err
			}
			if err := fn(data[start:end]); err != nil {
				return err
			}
		}
	}
	return nil
}

// Emoji unicode range
var emojiRanges = []struct {
	first, last rune
}{
	{0x1F600, 0x1F64F}, // Emoticons
	{0x1F300, 0x1F5FF}, // Misc Symbols and Pictographs
	{0x1F680, 0x1F6FF}, // Transport and Map
	{0x1F1E6, 0x1F1FF}, // Regional country flags
	{0x2600, 0x26FF},   // Misc symbols
	{0x2700, 0x27BF},   // Dingbats
	{0xFE00, 0xFE0F},   // Variation Selectors
	{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
	{0x1F018, 0x1F270}, // Various asian characters
	{0x238C, 0x2454},   // Misc items
	{0x20D0, 0x20FF},   // Combining Diacritical Marks for Symbols
}

func isEmoji(r rune) bool {
	for _, er := range emojiRanges {
		if r >= er.first && r <= er.last {
			return true
		}
	}
	return false
}

var traditionalRanges = []struct {
	first, last rune
}{
	{0x3400, 0x4DB5},   // CJK Unified Ideographs Extension A
	{0x4E00, 0x9FEF},   // CJK Unified Ideographs
	{0xF900, 0xFAFF},   // CJK Compatibility Ideographs
	{0x20000, 0x2A6D6}, // CJK Unified Ideographs Extension B
	{0x2A700, 0x2B734}, // CJK Unified Ideographs Extension C
	{0x2B740, 0x2B81D}, // CJK Unified Ideographs Extension D
	{0x2B820, 0x2CEA1}, // CJK Unified Ideographs Extension E
	{0x2CEB0, 0x2EBE0}, // CJK Unified Ideographs Extension F
	{0x2F800, 0x2FA1F}, // CJK Compatibility Ideographs Supplement
}

var simplifiedRanges = []struct {
	first, last rune
}{
	{0x4E00, 0x9FFF}, // CJK Unified Ideographs
}

func isSimplifiedChar(r rune) bool {
	for _, sr := range simplifiedRanges {
		if sr.first <= r && r <= sr.last {
			return true
		}
	}
	return false
}

func isTraditionalChar(r rune) bool {
	for _, v := range traditionalRanges {
		if v.first <= r && r <= v.last {
			return true
		}
	}
	return false
}

var japaneseRanges = []struct {
	first, last rune
}{
	{0x3040, 0x309F},   // Hiragana
	{0x30A0, 0x30FF},   // Katakana
	{0x4E00, 0x9FFF},   // Kanji (Common and Uncommon)
	{0x3400, 0x4DBF},   // Kanji (Rare)
	{0x1F200, 0x1F200}, // Katakana letter archaic E
	{0x1F210, 0x1F23B}, // Squared Katakana words
}

func isJapaneseChar(r rune) bool {
	for _, jr := range japaneseRanges {
		if jr.first <= r && r <= jr.last {
			return true
		}
	}
	return false
}

var koreanRanges = []struct {
	first, last rune
}{
	{0xAC00, 0xD7A3},   // Hangul Syllables
	{0x1100, 0x11FF},   // Hangul Jamo
	{0x3130, 0x318F},   // Hangul Compatibility Jamo
	{0xA960, 0xA97F},   // Hangul Jamo Extended-A
	{0xD7B0, 0xD7FF},   // Hangul Jamo Extended-B
	{0x3200, 0x321F},   // Enclosed CJK Letters and Months (Parenthesized Hangul)
	{0x3260, 0x327F},   // Enclosed CJK Letters and Months (Circled Hangul)
	{0x1F200, 0x1F200}, // Enclosed Ideographic Supplement (Squared Hangul)
}

func isKoreanChar(r rune) bool {
	for _, kr := range koreanRanges {
		if kr.first <= r && r <= kr.last {
			return true
		}
	}
	return false
}

func isRussianChar(r rune) bool {
	if r >= 0x0400 && r <= 0x04FF {
		return true
	}
	return false
}

var turkishRanges = []struct {
	first, last rune
}{
	{0x0060, 0x007A}, // Latin lowercase characters
	{0x0041, 0x005A}, // Latin uppercase characters
	{0x011E, 0x011F}, // G, g with breve
	{0x0130, 0x0131}, // I with dot above, i without dot
	{0x015E, 0x015F}, // S, s with cedilla
	{0x00C7, 0x00E7}, // C, c with cedilla
}

func isTurkishChar(r rune) bool {
	for _, tr := range turkishRanges {
		if tr.first <= r && r <= tr.last {
			return true
		}
	}
	return false
}

func isThaiChar(r rune) bool {
	if 0x0E00 <= r && r <= 0x0E7F {
		return true
	}
	return false
}

var vietnameseRanges = []struct {
	first, last rune
}{
	{0x0041, 0x005A}, // Latin uppercase characters
	{0x0061, 0x007A}, // Latin lowercase characters
	{0x00C0, 0x00C3}, // A with grave, A with tilde
	{0x00C8, 0x00CA}, // E with grave, E with circumflex
	{0x00CC, 0x00CD}, // I with grave, I with acute
	{0x00D2, 0x00D5}, // O with grave, O with tilde
	{0x00D9, 0x00DA}, // U with grave, U with acute
	{0x00DD, 0x00DD}, // Y with acute
	{0x0110, 0x0111}, // D with stroke
	{0x0128, 0x0129}, // I with tilde
	{0x0168, 0x0169}, // U with tilde
	{0x01A0, 0x01A1}, // O with horn
	{0x01AF, 0x01B0}, // U with horn
	{0x1EA0, 0x1EF9}, // Vietnamese additional characters
}

func isVietnameseChar(r rune) bool {
	for _, vn := range vietnameseRanges {
		if vn.first <= r && r <= vn.last {
			return true
		}
	}
	return false
}
