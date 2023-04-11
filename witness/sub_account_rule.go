package witness

import (
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

type SubAccountRuleEntity struct {
	ParentAccount string              `json:"-"`
	Version       uint32              `json:"version"`
	Rules         SubAccountRuleSlice `json:"rules"`
}

type SubAccountRuleSlice []SubAccountRule

type SubAccountRule struct {
	Index uint32        `json:"index"`
	Name  string        `json:"name"`
	Note  string        `json:"note"`
	Price uint64        `json:"price,omitempty"`
	Ast   AstExpression `json:"ast"`
}

type AstExpression struct {
	Type       ExpressionType   `json:"type"`
	Expression ExpressionEntity `json:"expression"`
}

type ExpressionEntity struct {
	subAccountRuleEntity *SubAccountRuleEntity
	Name                 string         `json:"name,omitempty"`
	Symbol               SymbolType     `json:"symbol,omitempty"`
	Value                interface{}    `json:"value,omitempty"`
	ValueType            ValueType      `json:"value_type,omitempty"`
	Expressions          AstExpressions `json:"expressions,omitempty"`
}

type AstExpressions []AstExpression

func NewSubAccountRuleEntity(parentAccount string) *SubAccountRuleEntity {
	return &SubAccountRuleEntity{
		ParentAccount: parentAccount,
		Rules:         make(SubAccountRuleSlice, 0),
	}
}

func (s *SubAccountRuleEntity) ParseFromJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, s); err != nil {
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
		if _, err = v.Ast.Check(false, ""); err != nil {
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) Hit(account string) (hit bool, index int, err error) {
	account = strings.Split(account, ".")[0]
	for idx, v := range s.Rules {
		v.Ast.Expression.subAccountRuleEntity = s
		hit, err = v.Ast.Check(true, account)
		if err != nil || hit {
			index = idx
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) ParseFromTx(tx *types.Transaction, action common.ActionDataType) error {
	data := make([][]byte, 0)
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		if actionDataType != action {
			return true, nil
		}
		data = append(data, dataBys)
		return true, nil
	})
	if err != nil {
		return err
	}
	if err := s.ParseFromDasActionWitnessData(data); err != nil {
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
		index, indexLen, dataLen := uint32(0), molecule.HeaderSizeUint, uint32(0)
		if int(indexLen) > len(v) {
			return fmt.Errorf("data length error: %d", len(v))
		}

		dataLen, err := molecule.Bytes2GoU32(v[index : index+indexLen])
		if err != nil {
			return err
		}
		index += indexLen

		version, err := molecule.Bytes2GoU32(v[index : index+dataLen])
		if err != nil {
			return err
		}
		index += dataLen

		if s.Version > 0 && s.Version != version {
			return errors.New("version aberrant")
		}
		s.Version = version

		dataLen, err = molecule.Bytes2GoU32(v[index : index+indexLen])
		if err != nil {
			return err
		}
		index += indexLen

		v = v[index : index+dataLen]

		rules, err := molecule.SubAccountRulesFromSlice(v, true)
		if err != nil {
			return err
		}
		for i := uint(0); i < rules.ItemCount(); i++ {
			r := rules.Get(i)
			index, err := molecule.Bytes2GoU32(r.Index().RawData())
			if err != nil {
				return err
			}

			nameBytes, err := molecule.BytesFromSlice(r.Name().AsSlice(), true)
			if err != nil {
				return err
			}
			name := string(nameBytes.RawData())

			noteBytes, err := molecule.BytesFromSlice(r.Note().AsSlice(), true)
			if err != nil {
				return err
			}
			note := string(noteBytes.RawData())

			price, err := molecule.Bytes2GoU64(r.Price().RawData())
			if err != nil {
				return err
			}
			rule := NewSubAccountRule()
			rule.Index = index
			rule.Name = name
			rule.Note = note
			rule.Price = price

			exp, err := s.ParseFromMolecule(r.Ast())
			if err != nil {
				return err
			}
			rule.Ast = *exp

			s.Rules = append(s.Rules, *rule)
		}
	}
	sort.Slice(s.Rules, func(i, j int) bool {
		return s.Rules[i].Index < s.Rules[j].Index
	})
	return nil
}

func (s *SubAccountRuleEntity) ParseFromMolecule(astExp *molecule.ASTExpression) (*AstExpression, error) {
	expType, err := molecule.Bytes2GoU8(astExp.ExpressionType().AsSlice())
	if err != nil {
		return nil, err
	}

	ast := &AstExpression{}

	switch expType {
	case 0x00:
		ast.Type = Operator

		exp, err := molecule.ASTOperatorFromSlice(astExp.Expression().RawData(), true)
		if err != nil {
			return nil, err
		}
		symbol, err := molecule.Bytes2GoU8(exp.Symbol().AsSlice())
		if err != nil {
			return nil, err
		}

		switch symbol {
		case 0x00:
			ast.Expression.Symbol = Not
		case 0x01:
			ast.Expression.Symbol = And
		case 0x02:
			ast.Expression.Symbol = Or
		case 0x03:
			ast.Expression.Symbol = Gt
		case 0x04:
			ast.Expression.Symbol = Gte
		case 0x05:
			ast.Expression.Symbol = Lt
		case 0x06:
			ast.Expression.Symbol = Lte
		case 0x07:
			ast.Expression.Symbol = Equ
		}

		for i := uint(0); i < exp.Expressions().ItemCount(); i++ {
			r := exp.Expressions().Get(i)
			astExp, err := s.ParseFromMolecule(r)
			if err != nil {
				return nil, err
			}
			ast.Expression.Expressions = append(ast.Expression.Expressions, *astExp)
		}

	case 0x01:
		ast.Type = Function

		exp, err := molecule.ASTFunctionFromSlice(astExp.Expression().RawData(), true)
		if err != nil {
			return nil, err
		}

		name, err := molecule.Bytes2GoU8(exp.Name().AsSlice())
		if err != nil {
			return nil, err
		}

		switch name {
		case 0x00:
			ast.Expression.Name = string(FunctionIncludeCharts)
		case 0x01:
			ast.Expression.Name = string(FunctionIncludeCharts)
		case 0x02:
			ast.Expression.Name = string(FunctionInList)
		}

		for i := uint(0); i < exp.Arguments().ItemCount(); i++ {
			r := exp.Arguments().Get(i)
			astExp, err := s.ParseFromMolecule(r)
			if err != nil {
				return nil, err
			}
			ast.Expression.Expressions = append(ast.Expression.Expressions, *astExp)
		}
	case 0x02:
		ast.Type = Variable

		exp, err := molecule.ASTVariableFromSlice(astExp.Expression().RawData(), true)
		if err != nil {
			return nil, err
		}
		varName, err := molecule.Bytes2GoU8(exp.Name().AsSlice())
		if err != nil {
			return nil, err
		}

		switch varName {
		case 0x00:
			ast.Expression.Name = string(Account)
		case 0x01:
			ast.Expression.Name = string(AccountChars)
		case 0x02:
			ast.Expression.Name = string(AccountLength)
		}
	case 0x03:
		ast.Type = Value

		exp, err := molecule.ASTValueFromSlice(astExp.Expression().RawData(), true)
		if err != nil {
			return nil, err
		}
		valueType, err := molecule.Bytes2GoU8(exp.ValueType().AsSlice())
		if err != nil {
			return nil, err
		}

		switch valueType {
		case 0x00:
			ast.Expression.ValueType = Bool
			u8, err := molecule.Bytes2GoU8(exp.Value().RawData())
			if err != nil {
				return nil, err
			}
			ast.Expression.Value = gconv.Bool(u8)
		case 0x01:
			ast.Expression.ValueType = Uint8
			u8, err := molecule.Bytes2GoU8(exp.Value().RawData())
			if err != nil {
				return nil, err
			}
			ast.Expression.Value = u8
		case 0x02:
			ast.Expression.ValueType = Uint32
			u32, err := molecule.Bytes2GoU32(exp.Value().RawData())
			if err != nil {
				return nil, err
			}
			ast.Expression.Value = u32
		case 0x03:
			ast.Expression.ValueType = Uint64
			u64, err := molecule.Bytes2GoU64(exp.Value().RawData())
			if err != nil {
				return nil, err
			}
			ast.Expression.Value = u64
		case 0x04:
			ast.Expression.ValueType = Binary
			ast.Expression.Value = common.Bytes2Hex(exp.Value().RawData())
		case 0x05:
			ast.Expression.ValueType = BinaryArray
			strArrays := make([]string, 0)
			bytesVec, err := molecule.BytesVecFromSlice(exp.Value().RawData(), true)
			if err != nil {
				return nil, err
			}
			for i := uint(0); i < bytesVec.ItemCount(); i++ {
				getBytes := bytesVec.Get(i)
				strArrays = append(strArrays, common.Bytes2Hex(getBytes.RawData()))
			}
			ast.Expression.Value = strArrays
		case 0x06:
			ast.Expression.ValueType = String
			ast.Expression.Value = string(exp.Value().RawData())
		case 0x07:
			ast.Expression.ValueType = StringArray
			strArrays := make([]string, 0)
			bytesVec, err := molecule.BytesVecFromSlice(exp.Value().RawData(), true)
			if err != nil {
				return nil, err
			}
			for i := uint(0); i < bytesVec.ItemCount(); i++ {
				getBytes := bytesVec.Get(i)
				strArrays = append(strArrays, string(getBytes.RawData()))
			}
			ast.Expression.Value = strArrays
		case 0x08:
			ast.Expression.ValueType = Charset
			charset, err := molecule.Bytes2GoU32(exp.Value().RawData())
			if err != nil {
				return nil, err
			}
			ast.Expression.Value = charset
		}
	}
	return ast, nil
}

func (s *SubAccountRuleEntity) GenWitnessData(action common.ActionDataType) ([][]byte, error) {
	for _, v := range s.Rules {
		if string(v.Name) == "" {
			return nil, errors.New("name can't be empty")
		}
		if v.Price < 0 {
			return nil, errors.New("price can't be negative number")
		}
		if _, err := v.Ast.Check(false, ""); err != nil {
			return nil, err
		}
	}

	res := make([]molecule.SubAccountRules, 0)
	rulesBuilder := molecule.NewSubAccountRulesBuilder()

	for idx, v := range s.Rules {
		ruleBuilder := molecule.NewSubAccountRuleBuilder()
		ruleBuilder.Index(molecule.GoU32ToMoleculeU32(uint32(idx)))
		ruleBuilder.Name(molecule.GoString2MoleculeBytes(v.Name))
		ruleBuilder.Note(molecule.GoString2MoleculeBytes(v.Note))
		ruleBuilder.Price(molecule.GoU64ToMoleculeU64(v.Price))

		astExp, err := v.Ast.GenMoleculeASTExpression()
		if err != nil {
			return nil, err
		}
		ruleBuilder.Ast(*astExp)
		rule := ruleBuilder.Build()
		rules := rulesBuilder.Build()

		if rules.TotalSize()+rule.TotalSize()+4+12+common.WitnessDasTableTypeEndIndex > 32*1024 {
			res = append(res, rulesBuilder.Build())
			rulesBuilder = molecule.NewSubAccountRulesBuilder()
			if idx == len(s.Rules)-1 {
				rulesBuilder.Push(rule)
				res = append(res, rulesBuilder.Build())
			} else {
				rulesBuilder.Push(rule)
			}
		} else {
			rulesBuilder.Push(rule)
			if idx == len(s.Rules)-1 {
				res = append(res, rulesBuilder.Build())
			}
		}
	}

	resultBs := make([][]byte, 0)
	for _, v := range res {
		data := make([]byte, 0)

		versionBys := molecule.GoU32ToMoleculeU32(s.Version)
		data = append(data, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
		data = append(data, versionBys.RawData()...)

		data = append(data, molecule.GoU32ToBytes(uint32(v.TotalSize()))...)
		data = append(data, v.AsSlice()...)
		resultBs = append(resultBs, GenDasDataWitnessWithByte(action, data))
	}
	return resultBs, nil
}

func (e *AstExpression) GenMoleculeASTExpression() (*molecule.ASTExpression, error) {
	astExpBuilder := molecule.NewASTExpressionBuilder()
	switch e.Type {
	case Operator:
		astExpBuilder.ExpressionType(molecule.NewByte(0x00))
		expBuilder := molecule.NewASTOperatorBuilder()
		switch e.Expression.Symbol {
		case Not:
			expBuilder.Symbol(molecule.NewByte(0x00))
		case And:
			expBuilder.Symbol(molecule.NewByte(0x01))
		case Or:
			expBuilder.Symbol(molecule.NewByte(0x02))
		case Gt:
			expBuilder.Symbol(molecule.NewByte(0x03))
		case Gte:
			expBuilder.Symbol(molecule.NewByte(0x04))
		case Lt:
			expBuilder.Symbol(molecule.NewByte(0x05))
		case Lte:
			expBuilder.Symbol(molecule.NewByte(0x06))
		case Equ:
			expBuilder.Symbol(molecule.NewByte(0x07))
		}

		expsBuilder := molecule.NewASTExpressionsBuilder()
		for _, v := range e.Expression.Expressions {
			astExp, err := v.GenMoleculeASTExpression()
			if err != nil {
				return nil, err
			}
			expsBuilder.Push(*astExp)
		}
		expBuilder.Expressions(expsBuilder.Build())
		astOperator := expBuilder.Build()
		astExpBuilder.Expression(molecule.GoBytes2MoleculeBytes(astOperator.AsSlice()))

	case Function:
		astExpBuilder.ExpressionType(molecule.NewByte(0x01))
		expBuilder := molecule.NewASTFunctionBuilder()
		switch FunctionType(e.Expression.Name) {
		case FunctionIncludeCharts:
			expBuilder.Name(molecule.NewByte(0x00))
		case FunctionOnlyIncludeCharset:
			expBuilder.Name(molecule.NewByte(0x01))
		case FunctionInList:
			expBuilder.Name(molecule.NewByte(0x02))
		}

		expsBuilder := molecule.NewASTExpressionsBuilder()
		for _, v := range e.Expression.Expressions {
			astExp, err := v.GenMoleculeASTExpression()
			if err != nil {
				return nil, err
			}
			expsBuilder.Push(*astExp)
		}
		expBuilder.Arguments(expsBuilder.Build())
		astOperator := expBuilder.Build()
		astExpBuilder.Expression(molecule.GoBytes2MoleculeBytes(astOperator.AsSlice()))

	case Variable:
		astExpBuilder.ExpressionType(molecule.NewByte(0x02))
		expBuilder := molecule.NewASTVariableBuilder()
		switch VariableName(e.Expression.Name) {
		case Account:
			expBuilder.Name(molecule.NewByte(0x00))
		case AccountChars:
			expBuilder.Name(molecule.NewByte(0x01))
		case AccountLength:
			expBuilder.Name(molecule.NewByte(0x02))
		}
		exp := expBuilder.Build()
		astExpBuilder.Expression(molecule.GoBytes2MoleculeBytes(exp.AsSlice()))

	case Value:
		astExpBuilder.ExpressionType(molecule.NewByte(0x03))
		expBuilder := molecule.NewASTValueBuilder()
		switch e.Expression.ValueType {
		case Bool:
			expBuilder.ValueType(molecule.NewByte(0x00))
			u8 := molecule.GoU8ToMoleculeU8(gconv.Uint8(e.Expression.Value))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(u8.AsSlice()))
		case Uint8:
			expBuilder.ValueType(molecule.NewByte(0x01))
			u8 := molecule.GoU8ToMoleculeU8(gconv.Uint8(e.Expression.Value))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(u8.AsSlice()))
		case Uint32:
			expBuilder.ValueType(molecule.NewByte(0x02))
			u32 := molecule.GoU32ToMoleculeU32(gconv.Uint32(e.Expression.Value))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(u32.AsSlice()))
		case Uint64:
			expBuilder.ValueType(molecule.NewByte(0x03))
			u64 := molecule.GoU64ToMoleculeU64(gconv.Uint64(e.Expression.Value))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(u64.AsSlice()))
		case Binary:
			expBuilder.ValueType(molecule.NewByte(0x04))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(common.Hex2Bytes(gconv.String(e.Expression.Value))))
		case BinaryArray:
			expBuilder.ValueType(molecule.NewByte(0x05))
			bsVecBuilder := molecule.NewBytesVecBuilder()
			for _, v := range gconv.Strings(e.Expression.Value) {
				bsVecBuilder.Push(molecule.GoBytes2MoleculeBytes(common.Hex2Bytes(v)))
			}
			bsVec := bsVecBuilder.Build()
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(bsVec.AsSlice()))
		case String:
			expBuilder.ValueType(molecule.NewByte(0x06))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(gconv.Bytes(e.Expression.Value)))
		case StringArray:
			expBuilder.ValueType(molecule.NewByte(0x07))
			bsVecBuilder := molecule.NewBytesVecBuilder()
			for _, v := range gconv.Strings(e.Expression.Value) {
				bsVecBuilder.Push(molecule.GoBytes2MoleculeBytes([]byte(v)))
			}
			bsVec := bsVecBuilder.Build()
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(bsVec.AsSlice()))
		case Charset:
			expBuilder.ValueType(molecule.NewByte(0x08))
			u32 := molecule.GoU32ToMoleculeU32(gconv.Uint32(e.Expression.Value))
			expBuilder.Value(molecule.GoBytes2MoleculeBytes(u32.AsSlice()))
		}
		exp := expBuilder.Build()
		astExpBuilder.Expression(molecule.GoBytes2MoleculeBytes(exp.AsSlice()))
	}
	astExp := astExpBuilder.Build()
	return &astExp, nil
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
	_, err = s.Ast.Check(false, "")
	return
}

func (s *SubAccountRule) Hit(account string) (hit bool, err error) {
	account = strings.Split(account, ".")[0]
	return s.Ast.Check(true, account)
}

func (e *AstExpression) ReturnType() ReturnType {
	if e.Type == Operator || e.Type == Function || e.Type == Value && e.Expression.ValueType == Bool {
		return ReturnTypeBool
	}

	if e.Type == Value && (e.Expression.ValueType == Uint8 || e.Expression.ValueType == Uint32 || e.Expression.ValueType == Uint64) ||
		e.Type == Variable && VariableName(e.Expression.Name) == AccountLength {
		return ReturnTypeNumber
	}

	if e.Type == Value && e.Expression.ValueType == String ||
		e.Type == Variable && VariableName(e.Expression.Name) == Account ||
		e.Type == Value && e.Expression.ValueType == Binary {
		return ReturnTypeString
	}
	if e.Type == Variable && VariableName(e.Expression.Name) == AccountChars ||
		e.Type == Value && e.Expression.ValueType == StringArray ||
		e.Type == Value && e.Expression.ValueType == BinaryArray {
		return ReturnTypeStringArray
	}
	return ReturnTypeUnknown
}

func (e *AstExpression) Check(checkHit bool, account string) (hit bool, err error) {
	switch e.Type {
	case Function:
		funcName := FunctionType(e.Expression.Name)
		switch funcName {
		case FunctionIncludeCharts:
			hit, err = e.handleFunctionIncludeCharts(checkHit, account)
		case FunctionInList:
			hit, err = e.handleFunctionInList(checkHit, account)
		case FunctionOnlyIncludeCharset:
			hit, err = e.handleFunctionOnlyIncludeCharset(checkHit, account)
		default:
			err = fmt.Errorf("function %s can't be support", e.Expression.Name)
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

func (e *AstExpression) GetNumberValue(account string) float64 {
	if e.Type == Variable && VariableName(e.Expression.Name) == AccountLength {
		return float64(len([]rune(account)))
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeNumber {
		return gconv.Float64(e.Expression.Value)
	}
	return 0
}

func (e *AstExpression) GetStringValue(account string) string {
	if e.Type == Variable && VariableName(e.Expression.Name) == Account {
		return account
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeString {
		return gconv.String(e.Expression.Value)
	}
	return ""
}

func (e *AstExpression) ProcessOperator(checkHit bool, account string) (hit bool, err error) {
	switch e.Expression.Symbol {
	case And:
		for _, exp := range e.Expression.Expressions {
			exp.Expression.subAccountRuleEntity = e.Expression.subAccountRuleEntity
			rtType := exp.ReturnType()
			if rtType != ReturnTypeBool {
				return false, errors.New("operator 'and' every expression must be bool return")
			}
			hit, err := exp.Check(checkHit, account)
			if err != nil {
				return false, err
			}
			if checkHit && !hit {
				return false, nil
			}
		}
		return true, nil
	case Or:
		for _, exp := range e.Expression.Expressions {
			exp.Expression.subAccountRuleEntity = e.Expression.subAccountRuleEntity
			rtType := exp.ReturnType()
			if rtType != ReturnTypeBool {
				return false, errors.New("operator 'and' every expression must be bool return")
			}
			hit, err := exp.Check(checkHit, account)
			if err != nil {
				return false, err
			}
			if checkHit && hit {
				return true, nil
			}
		}
		return true, nil
	case Not:
		if len(e.Expression.Expressions) != 1 {
			return false, errors.New("operator not must have one expression")
		}
		exp := e.Expression.Expressions[0]
		exp.Expression.subAccountRuleEntity = e.Expression.subAccountRuleEntity

		rtType := exp.ReturnType()
		if rtType != ReturnTypeBool {
			return false, errors.New("operator 'not' expression must be bool return")
		}
		hit, err := exp.Check(checkHit, account)
		if err != nil {
			return false, err
		}
		if !hit {
			return true, nil
		}
	case Gt, Gte, Lt, Lte, Equ:
		if len(e.Expression.Expressions) != 2 {
			return false, errors.New("operator not must have two expression")
		}
		left := e.Expression.Expressions[0]
		right := e.Expression.Expressions[1]
		if !IsSameReturnType(left, right) {
			return false, errors.New("the comparison type operation must have same types on both sides")
		}
		left.Expression.subAccountRuleEntity = e.Expression.subAccountRuleEntity
		right.Expression.subAccountRuleEntity = e.Expression.subAccountRuleEntity

		switch left.ReturnType() {
		case ReturnTypeNumber:
			leftVal := left.GetNumberValue(account)
			rightVal := right.GetNumberValue(account)
			if e.Expression.Symbol == Gt {
				return leftVal > rightVal, nil
			}
			if e.Expression.Symbol == Gte {
				return leftVal >= rightVal, nil
			}
			if e.Expression.Symbol == Lt {
				return leftVal < rightVal, nil
			}
			if e.Expression.Symbol == Lte {
				return leftVal <= rightVal, nil
			}
			if e.Expression.Symbol == Equ {
				return leftVal == rightVal, nil
			}
		case ReturnTypeString:
			leftVal := left.GetStringValue(account)
			rightVal := right.GetStringValue(account)
			if e.Expression.Symbol == Gt {
				return leftVal > rightVal, nil
			}
			if e.Expression.Symbol == Gte {
				return leftVal >= rightVal, nil
			}
			if e.Expression.Symbol == Lt {
				return leftVal < rightVal, nil
			}
			if e.Expression.Symbol == Lte {
				return leftVal <= rightVal, nil
			}
			if e.Expression.Symbol == Equ {
				return leftVal == rightVal, nil
			}
		default:
			return false, fmt.Errorf("type %s is not currently supported", left.ReturnType())
		}
	default:
		err = fmt.Errorf("symbol %s can't be support", e.Expression.Symbol)
	}
	return
}

func (e *AstExpression) handleFunctionIncludeCharts(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expression.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Expression.Name)
		return
	}
	accCharts := e.Expression.Expressions[0]
	if accCharts.Type != Variable || VariableName(accCharts.Expression.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := e.Expression.Expressions[1]
	strArray := gconv.Strings(value.Expression.Value)
	if len(strArray) == 0 || value.Type != Value || value.Expression.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Expression.Name)
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

func (e *AstExpression) handleFunctionInList(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expression.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Expression.Name)
		return
	}
	value := e.Expression.Expressions[1]
	strArray := gconv.Strings(value.Expression.Value)
	if len(strArray) == 0 || value.Type != Value || (value.Expression.ValueType != BinaryArray && value.Expression.ValueType != StringArray) {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Expression.Name)
		return
	}

	if !checkHit {
		return
	}

	subAccount := fmt.Sprintf("%s.%s", account, e.Expression.subAccountRuleEntity.ParentAccount)
	subAccountId := common.Bytes2Hex(common.GetAccountIdByAccount(subAccount))
	for _, v := range strArray {
		switch value.Expression.ValueType {
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

func (e *AstExpression) handleFunctionOnlyIncludeCharset(checkHit bool, account string) (hit bool, err error) {
	if len(e.Expression.Expressions) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Expression.Name)
		return
	}
	accCharts := e.Expression.Expressions[0]
	if accCharts.Type != Variable || VariableName(accCharts.Expression.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := e.Expression.Expressions[1]
	val := common.AccountCharType(gconv.Uint32(value.Expression.Value))
	if _, ok := common.AccountCharTypeMap[val]; !ok ||
		value.Type != Value ||
		value.Expression.ValueType != Charset {
		err = fmt.Errorf("function %s args[1] charset %d no support ", e.Expression.Name, val)
		return
	}
	if !checkHit {
		return
	}

	chatSet := common.AccountCharTypeMap[val]
	for _, v := range []rune(account) {
		if _, ok := chatSet[string(v)]; !ok {
			return
		}
	}
	hit = true
	return
}

func IsSameReturnType(i, j AstExpression) bool {
	ir := i.ReturnType()
	jr := j.ReturnType()
	return ir == jr && ir != ReturnTypeUnknown
}
