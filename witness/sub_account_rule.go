package witness

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clipperhouse/uax29/graphemes"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type (
	ReturnType            string
	ExpressionType        string
	ExpressionsType       []ExpressionType
	SymbolType            string
	SymbolsType           []SymbolType
	FunctionType          string
	FunctionsType         []FunctionType
	VariableName          string
	VariablesName         []VariableName
	ValueType             string
	ValuesType            []ValueType
	CharsetType           string
	CharsetsType          []CharsetType
	SubAccountRuleVersion uint32
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

	Not SymbolType = "not"
	And SymbolType = "and"
	Or  SymbolType = "or"
	Gt  SymbolType = ">"
	Gte SymbolType = ">="
	Lt  SymbolType = "<"
	Lte SymbolType = "<="
	Equ SymbolType = "=="

	FunctionIncludeCharts      FunctionType = "include_chars"
	FunctionIncludeWords       FunctionType = "include_words"
	FunctionOnlyIncludeCharset FunctionType = "only_include_charset"
	FunctionInList             FunctionType = "in_list"
	FunctionIncludeCharset     FunctionType = "include_charset"
	FunctionStartsWith         FunctionType = "starts_with"
	FunctionEndsWith           FunctionType = "ends_with"

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

	SubAccountRuleVersionV1 SubAccountRuleVersion = 1
)

var (
	Functions   = FunctionsType{FunctionIncludeCharts, FunctionIncludeWords, FunctionOnlyIncludeCharset, FunctionInList, FunctionIncludeCharset, FunctionStartsWith, FunctionEndsWith}
	Values      = ValuesType{Bool, Uint8, Uint32, Uint64, Binary, BinaryArray, String, StringArray, Charset}
	Variables   = VariablesName{Account, AccountChars, AccountLength}
	Operators   = SymbolsType{Not, And, Or, Gt, Gte, Lt, Lte, Equ}
	Expressions = ExpressionsType{Operator, Function, Variable, Value}

	ParentAccountError = errors.New("parent account can't be empty, please init from NewSubAccountRuleEntity func")
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

func (v *ValueType) Parse(data []byte) (interface{}, error) {
	switch *v {
	case Bool:
		val, err := molecule.Bytes2GoU8(data)
		if err != nil {
			return nil, err
		}
		return gconv.Bool(val), nil
	case Uint8:
		val, err := molecule.Bytes2GoU8(data)
		if err != nil {
			return nil, err
		}
		return val, nil
	case Uint32:
		val, err := molecule.Bytes2GoU32(data)
		if err != nil {
			return nil, err
		}
		return val, nil
	case Uint64:
		val, err := molecule.Bytes2GoU64(data)
		if err != nil {
			return nil, err
		}
		return val, nil
	case Binary:
		return common.Bytes2Hex(data), nil
	case BinaryArray:
		val := make([]string, 0)
		bytesVec, err := molecule.BytesVecFromSlice(data, true)
		if err != nil {
			return nil, err
		}
		for i := uint(0); i < bytesVec.ItemCount(); i++ {
			getBytes := bytesVec.Get(i)
			val = append(val, common.Bytes2Hex(getBytes.RawData()))
		}
		return val, nil
	case String:
		return string(data), nil
	case StringArray:
		val := make([]string, 0)
		bytesVec, err := molecule.BytesVecFromSlice(data, true)
		if err != nil {
			return nil, err
		}
		for i := uint(0); i < bytesVec.ItemCount(); i++ {
			getBytes := bytesVec.Get(i)
			val = append(val, string(getBytes.RawData()))
		}
		return val, nil
	case Charset:
		val, err := molecule.Bytes2GoU32(data)
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		return nil, fmt.Errorf("unknown value type: %s", *v)
	}
}

func (v *ValueType) Gen(data interface{}, preExp *AstExpression) (molecule.Bytes, error) {
	var res molecule.Bytes
	switch *v {
	case Bool:
		u8 := molecule.GoU8ToMoleculeU8(gconv.Uint8(data))
		res = molecule.GoBytes2MoleculeBytes(u8.AsSlice())
	case Uint8:
		u8 := molecule.GoU8ToMoleculeU8(gconv.Uint8(data))
		res = molecule.GoBytes2MoleculeBytes(u8.AsSlice())
	case Uint32:
		u32 := molecule.GoU32ToMoleculeU32(gconv.Uint32(data))
		res = molecule.GoBytes2MoleculeBytes(u32.AsSlice())
	case Uint64:
		u64 := molecule.GoU64ToMoleculeU64(gconv.Uint64(data))
		res = molecule.GoBytes2MoleculeBytes(u64.AsSlice())
	case Binary:
		res = molecule.GoBytes2MoleculeBytes(common.Hex2Bytes(gconv.String(data)))
	case BinaryArray:
		bsVecBuilder := molecule.NewBytesVecBuilder()
		strs := gconv.Strings(data)
		log.Infof("BinaryArray: %v, parentAccount: %s", strs, preExp.subAccountRuleEntity.ParentAccount)
		for _, v := range strs {
			if preExp != nil && preExp.Type == Variable && preExp.Name == string(Account) {
				account := strings.Split(v, ".")[0] + "." + preExp.subAccountRuleEntity.ParentAccount
				accountId := common.GetAccountIdByAccount(account)
				bsVecBuilder.Push(molecule.GoBytes2MoleculeBytes(accountId))
			} else {
				bsVecBuilder.Push(molecule.GoBytes2MoleculeBytes(common.Hex2Bytes(v)))
			}
		}
		bsVec := bsVecBuilder.Build()
		res = molecule.GoBytes2MoleculeBytes(bsVec.AsSlice())
	case String:
		res = molecule.GoBytes2MoleculeBytes(gconv.Bytes(data))
	case StringArray:
		bsVecBuilder := molecule.NewBytesVecBuilder()
		for _, v := range gconv.Strings(data) {
			bsVecBuilder.Push(molecule.GoBytes2MoleculeBytes([]byte(v)))
		}
		bsVec := bsVecBuilder.Build()
		res = molecule.GoBytes2MoleculeBytes(bsVec.AsSlice())
	case Charset:
		if _, ok := common.AccountCharTypeMap[common.AccountCharType(gconv.Uint32(data))]; !ok {
			return res, fmt.Errorf("invalid charset: %d", gconv.Uint32(data))
		}
		u32 := molecule.GoU32ToMoleculeU32(gconv.Uint32(data))
		res = molecule.GoBytes2MoleculeBytes(u32.AsSlice())
	}
	return res, nil
}

type SubAccountRuleEntity struct {
	ParentAccount string                `json:"-"`
	Version       SubAccountRuleVersion `json:"version"`
	Rules         SubAccountRuleSlice   `json:"rules"`
}

type SubAccountRuleSlice []*SubAccountRule

type SubAccountRule struct {
	Index  uint32        `json:"index"`
	Name   string        `json:"name"`
	Note   string        `json:"note"`
	Price  float64       `json:"price,omitempty"`
	Ast    AstExpression `json:"ast"`
	Status uint8         `json:"status"`
}

type AstExpression struct {
	subAccountRuleEntity *SubAccountRuleEntity

	Type        ExpressionType `json:"type"`
	Name        string         `json:"name,omitempty"`
	Symbol      SymbolType     `json:"symbol,omitempty"`
	Value       interface{}    `json:"value,omitempty"`
	ValueType   ValueType      `json:"value_type,omitempty"`
	Arguments   AstExpressions `json:"arguments,omitempty"`
	Expressions AstExpressions `json:"expressions,omitempty"`
}

type AstExpressions []*AstExpression

func NewSubAccountRuleEntity(parentAccount string, versions ...SubAccountRuleVersion) *SubAccountRuleEntity {
	entity := &SubAccountRuleEntity{
		ParentAccount: parentAccount,
		Version:       SubAccountRuleVersionV1,
		Rules:         make(SubAccountRuleSlice, 0),
	}
	if len(versions) > 0 {
		entity.Version = versions[0]
	}
	return entity
}

func (s *SubAccountRuleEntity) ParseFromJSON(data []byte) (err error) {
	if err = json.Unmarshal(data, s); err != nil {
		return
	}
	return s.Check()
}

func (s *SubAccountRuleEntity) Check() (err error) {
	if s.ParentAccount == "" {
		return ParentAccountError
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
		v.Ast.subAccountRuleEntity = s
		if _, err = v.Ast.Check(false, ""); err != nil {
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) Hit(account string) (hit bool, index int, err error) {
	if s.ParentAccount == "" {
		return false, -1, ParentAccountError
	}
	account = strings.Split(account, ".")[0]
	for idx, v := range s.Rules {
		v.Ast.subAccountRuleEntity = s
		if v.Status == 0 {
			continue
		}
		hit, err = v.Ast.Check(true, account)
		if err != nil || hit {
			index = idx
			return
		}
	}
	return
}

func (s *SubAccountRuleEntity) ParseFromTx(tx *types.Transaction, action common.ActionDataType) error {
	if s.ParentAccount == "" {
		return ParentAccountError
	}
	data := make([][]byte, 0)
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte, index int) (bool, error) {
		if actionDataType != action {
			return true, nil
		}
		data = append(data, dataBys)
		return true, nil
	})
	if err != nil {
		return err
	}
	return s.ParseFromWitnessData(data)
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

		if s.Version > 0 && uint32(s.Version) != version {
			return errors.New("version aberrant")
		}
		s.Version = SubAccountRuleVersion(version)

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

			exp, err := s.ParseFromMolecule(r.Ast())
			if err != nil {
				return err
			}

			status, err := molecule.Bytes2GoU8(r.Status().RawData())
			if err != nil {
				return err
			}

			rule := NewSubAccountRule()
			rule.Index = index
			rule.Name = name
			rule.Note = note
			rule.Price = float64(price)
			rule.Ast = *exp
			rule.Status = status

			s.Rules = append(s.Rules, rule)
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
		if int(symbol) >= len(Operators) {
			return nil, fmt.Errorf("symbol error: %d no support", symbol)
		}
		ast.Symbol = Operators[int(symbol)]

		for i := uint(0); i < exp.Expressions().ItemCount(); i++ {
			r := exp.Expressions().Get(i)
			astExp, err := s.ParseFromMolecule(r)
			if err != nil {
				return nil, err
			}
			if ast.Expressions == nil {
				ast.Expressions = make([]*AstExpression, 0, exp.Expressions().ItemCount())
			}
			ast.Expressions = append(ast.Expressions, astExp)
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
		if int(name) >= len(Functions) {
			return nil, fmt.Errorf("function error: %d no support", name)
		}
		ast.Name = string(Functions[int(name)])

		for i := uint(0); i < exp.Arguments().ItemCount(); i++ {
			r := exp.Arguments().Get(i)
			astExp, err := s.ParseFromMolecule(r)
			if err != nil {
				return nil, err
			}
			if ast.Arguments == nil {
				ast.Arguments = make([]*AstExpression, 0, exp.Arguments().ItemCount())
			}
			ast.Arguments = append(ast.Arguments, astExp)
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

		if int(varName) >= len(Variables) {
			return nil, fmt.Errorf("variable error: %d no support", varName)
		}
		ast.Name = string(Variables[int(varName)])

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
		if int(valueType) >= len(Values) {
			return nil, fmt.Errorf("value error: %d no support", valueType)
		}
		ast.ValueType = Values[int(valueType)]
		val, err := ast.ValueType.Parse(exp.Value().RawData())
		if err != nil {
			return nil, err
		}
		ast.Value = val
	}
	return ast, nil
}

func (s *SubAccountRuleEntity) GenWitnessData(action common.ActionDataType) ([][]byte, error) {
	res, err := s.GenData()
	if err != nil {
		return nil, err
	}
	return s.GenDasData(action, res)
}

func (s *SubAccountRuleEntity) GenWitnessDataWithRuleData(ruleData [][]byte) ([][]byte, error) {
	resultBs := make([][]byte, 0)
	for _, v := range ruleData {
		data := make([]byte, 0)

		versionBys := molecule.GoU32ToMoleculeU32(uint32(s.Version))
		data = append(data, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
		data = append(data, versionBys.RawData()...)

		data = append(data, molecule.GoU32ToBytes(uint32(len(v)))...)
		data = append(data, v...)
		resultBs = append(resultBs, data)
	}
	return resultBs, nil
}

func (s *SubAccountRuleEntity) GenDasData(action common.ActionDataType, ruleData [][]byte) ([][]byte, error) {
	res, err := s.GenWitnessDataWithRuleData(ruleData)
	if err != nil {
		return nil, err
	}
	resultBs := make([][]byte, 0)
	for _, v := range res {
		resultBs = append(resultBs, GenDasDataWitnessWithByte(action, v))
	}
	return resultBs, nil
}

const splitWitnessPreSize = uint(4 + 12 + common.WitnessDasTableTypeEndIndex)

func (s *SubAccountRuleEntity) GenData() ([][]byte, error) {
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
	tmpRules := make([]molecule.SubAccountRule, 0)

	for idx, v := range s.Rules {
		ruleBuilder := molecule.NewSubAccountRuleBuilder()
		ruleBuilder.Index(molecule.GoU32ToMoleculeU32(uint32(idx)))
		ruleBuilder.Name(molecule.GoString2MoleculeBytes(v.Name))
		ruleBuilder.Note(molecule.GoString2MoleculeBytes(v.Note))
		ruleBuilder.Price(molecule.GoU64ToMoleculeU64(uint64(v.Price)))
		ruleBuilder.Status(molecule.GoU8ToMoleculeU8(v.Status))

		astExp, err := v.Ast.GenMoleculeASTExpression(nil)
		if err != nil {
			return nil, err
		}
		ruleBuilder.Ast(*astExp)

		tmpRules = append(tmpRules, ruleBuilder.Build())
		rules := rulesBuilder.Set(tmpRules).Build()

		if rules.TotalSize()+splitWitnessPreSize < common.WitnessDataSizeLimit {
			if idx == len(s.Rules)-1 {
				res = append(res, rules)
				break
			}
			continue
		}
		if len(tmpRules) == 1 {
			return nil, fmt.Errorf("rule index: %d , size is too large", idx)
		}

		res = append(res, rulesBuilder.Set(tmpRules[:len(tmpRules)-1]).Build())
		if idx < len(s.Rules)-1 {
			tmpRules = []molecule.SubAccountRule{ruleBuilder.Build()}
			continue
		}

		rulesBuilder.Set([]molecule.SubAccountRule{ruleBuilder.Build()})
		lastRules := rulesBuilder.Build()
		if lastRules.TotalSize()+splitWitnessPreSize >= common.WitnessDataSizeLimit {
			return nil, fmt.Errorf("rule index: %d , size is too large", idx)
		}
		res = append(res, lastRules)
	}

	resultBs := make([][]byte, 0)
	for _, v := range res {
		resultBs = append(resultBs, v.AsSlice())
	}
	return resultBs, nil
}

func (e *AstExpression) Check(checkHit bool, account string) (hit bool, err error) {
	switch e.Type {
	case Function:
		funcName := FunctionType(e.Name)
		switch funcName {
		case FunctionIncludeCharts:
			hit, err = e.handleFunctionIncludeCharts(checkHit, account)
		case FunctionIncludeWords:
			hit, err = e.handleFunctionIncludeWords(checkHit, account)
		case FunctionOnlyIncludeCharset:
			hit, err = e.handleFunctionOnlyIncludeCharset(checkHit, account)
		case FunctionInList:
			hit, err = e.handleFunctionInList(checkHit, account)
		case FunctionIncludeCharset:
			hit, err = e.handleFunctionIncludeCharset(checkHit, account)
		case FunctionStartsWith:
			hit, err = e.handleFunctionStartsWith(checkHit, account)
		case FunctionEndsWith:
			hit, err = e.handleFunctionEndsWith(checkHit, account)
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

func (e *AstExpression) GenMoleculeASTExpression(preExp *AstExpression) (*molecule.ASTExpression, error) {
	astExpBuilder := molecule.NewASTExpressionBuilder()
	switch e.Type {
	case Operator:
		astExpBuilder.ExpressionType(molecule.NewByte(0x00))
		expBuilder := molecule.NewASTOperatorBuilder()

		for idx, v := range Operators {
			if v == e.Symbol {
				expBuilder.Symbol(molecule.NewByte(byte(idx)))
				break
			}
		}

		expsBuilder := molecule.NewASTExpressionsBuilder()
		for idx, v := range e.Expressions {
			v.subAccountRuleEntity = e.subAccountRuleEntity
			var astExp *molecule.ASTExpression
			var err error
			if idx == 0 {
				astExp, err = v.GenMoleculeASTExpression(e)
			} else {
				astExp, err = v.GenMoleculeASTExpression(e.Expressions[idx-1])
			}
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

		funcExist := false
		for idx, v := range Functions {
			if string(v) == e.Name {
				funcExist = true
				expBuilder.Name(molecule.NewByte(byte(idx)))
				break
			}
		}
		if !funcExist {
			return nil, fmt.Errorf("function %s can't be support", e.Name)
		}

		expsBuilder := molecule.NewASTExpressionsBuilder()
		for idx, v := range e.Arguments {
			v.subAccountRuleEntity = e.subAccountRuleEntity
			var astExp *molecule.ASTExpression
			var err error
			if idx == 0 {
				astExp, err = v.GenMoleculeASTExpression(e)
			} else {
				astExp, err = v.GenMoleculeASTExpression(e.Arguments[idx-1])
			}
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

		for idx, v := range Variables {
			if string(v) == e.Name {
				expBuilder.Name(molecule.NewByte(byte(idx)))
				break
			}
		}

		exp := expBuilder.Build()
		astExpBuilder.Expression(molecule.GoBytes2MoleculeBytes(exp.AsSlice()))

	case Value:
		astExpBuilder.ExpressionType(molecule.NewByte(0x03))
		expBuilder := molecule.NewASTValueBuilder()

		if preExp.Type == Variable &&
			preExp.Name == string(AccountLength) &&
			e.ReturnType() == ReturnTypeNumber {
			e.ValueType = Uint32
		}

		for idx, v := range Values {
			if v == e.ValueType {
				expBuilder.ValueType(molecule.NewByte(byte(idx)))
				break
			}
		}
		value, err := e.ValueType.Gen(e.Value, preExp)
		if err != nil {
			return nil, err
		}
		expBuilder.Value(value)

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

func (e *AstExpression) GetNumberValue(account string) float64 {
	if e.Type == Variable && VariableName(e.Name) == AccountLength {
		return float64(len([]rune(account)))
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeNumber {
		return gconv.Float64(e.Value)
	}
	return 0
}

func (e *AstExpression) GetStringValue(account string) string {
	if e.Type == Variable && VariableName(e.Name) == Account {
		return account
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeString {
		return gconv.String(e.Value)
	}
	return ""
}

func (e *AstExpression) GetStringArray(account string) []string {
	if e.Type == Variable && VariableName(e.Name) == AccountChars {
		return gconv.Strings([]rune(account))
	}
	if e.Type == Value && e.ReturnType() == ReturnTypeStringArray {
		return gconv.Strings(e.Value)
	}
	return nil
}

func (e *AstExpression) ProcessOperator(checkHit bool, account string) (hit bool, err error) {
	switch e.Symbol {
	case And:
		for _, exp := range e.Expressions {
			exp.subAccountRuleEntity = e.subAccountRuleEntity
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
		for _, exp := range e.Expressions {
			exp.subAccountRuleEntity = e.subAccountRuleEntity
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
		if len(e.Expressions) != 1 {
			return false, errors.New("operator not must have one expression")
		}
		exp := e.Expressions[0]
		exp.subAccountRuleEntity = e.subAccountRuleEntity

		rtType := exp.ReturnType()
		if rtType != ReturnTypeBool {
			return false, errors.New("operator 'not' expression must be bool return")
		}
		hit, err := exp.Check(checkHit, account)
		if err != nil {
			return false, err
		}
		return !hit, nil
	case Gt, Gte, Lt, Lte, Equ:
		if len(e.Expressions) != 2 {
			return false, errors.New("operator not must have two expression")
		}
		left := e.Expressions[0]
		right := e.Expressions[1]

		if left.ReturnType() == ReturnTypeUnknown {
			return false, fmt.Errorf("unknown type: %s %s", left.ValueType, left.Name)
		}
		if right.ReturnType() == ReturnTypeUnknown {
			return false, fmt.Errorf("unknown type: %s %s", right.ValueType, right.Name)
		}
		if left.ReturnType() != right.ReturnType() {
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
		case ReturnTypeStringArray:
			if e.Symbol != Equ {
				return false, fmt.Errorf("%s or %s only support '==' operator", StringArray, BinaryArray)
			}
			leftVal := left.GetStringArray(account)
			rightVal := right.GetStringArray(account)
			return reflect.DeepEqual(leftVal, rightVal), nil
		default:
			return false, fmt.Errorf("type %s is not currently supported", left.ReturnType())
		}
	default:
		err = fmt.Errorf("symbol %s can't be support", e.Symbol)
	}
	return
}

func (e *AstExpression) handleFunctionIncludeCharts(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}

	value := e.Arguments[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || value.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Name)
		return
	}
	if !checkHit {
		return
	}

	for _, v := range strArray {
		l := 0
		segments := graphemes.NewSegmenter([]byte(v))
		for segments.Next() {
			l++
		}
		if err := segments.Err(); err != nil {
			err = fmt.Errorf("segments.Err: %s", err.Error())
			return
		}
		if l > 1 {
			err = fmt.Errorf("function %s args[1] value must be single character", e.Name)
		}
		if strings.Contains(account, v) {
			hit = true
			return
		}
	}
	return
}

func (e *AstExpression) handleFunctionInList(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accountVar := e.Arguments[0]
	if accountVar.Type != Variable || VariableName(accountVar.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}

	value := e.Arguments[1]
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

func (e *AstExpression) handleFunctionOnlyIncludeCharset(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := e.Arguments[1]
	val := common.AccountCharType(gconv.Uint32(value.Value))
	if _, ok := common.AccountCharTypeMap[val]; !ok ||
		value.Type != Value ||
		value.ValueType != Charset {
		err = fmt.Errorf("function %s args[1] charset %d no support ", e.Name, val)
		return
	}
	if !checkHit {
		return
	}

	dasCore := core.NewDasCore(context.Background(), &sync.WaitGroup{})
	charsets, err := dasCore.GetAccountCharSetList(account + "." + e.subAccountRuleEntity.ParentAccount)
	if err != nil {
		return false, err
	}
	for _, v := range charsets {
		if v.CharSetName != val {
			return
		}
	}
	hit = true
	return
}

func (e *AstExpression) handleFunctionIncludeWords(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}

	value := e.Arguments[1]
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

func (e *AstExpression) handleFunctionIncludeCharset(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accCharts := e.Arguments[0]
	if accCharts.Type != Variable || VariableName(accCharts.Name) != AccountChars {
		err = fmt.Errorf("first args type must variable and name is %s", AccountChars)
		return
	}

	value := e.Arguments[1]
	val := common.AccountCharType(gconv.Uint32(value.Value))
	if _, ok := common.AccountCharTypeMap[val]; !ok ||
		value.Type != Value ||
		value.ValueType != Charset {
		err = fmt.Errorf("function %s args[1] charset %d no support ", e.Name, val)
		return
	}
	if !checkHit {
		return
	}

	dasCore := core.NewDasCore(context.Background(), &sync.WaitGroup{})
	charsets, err := dasCore.GetAccountCharSetList(account + "." + e.subAccountRuleEntity.ParentAccount)
	if err != nil {
		return false, err
	}
	for _, v := range charsets {
		if v.CharSetName == val {
			hit = true
			return
		}
	}
	return
}

func (e *AstExpression) handleFunctionStartsWith(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accountVar := e.Arguments[0]
	if accountVar.Type != Variable || VariableName(accountVar.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}

	value := e.Arguments[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || value.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Name)
		return
	}

	if !checkHit {
		return
	}

	for _, v := range strArray {
		if strings.HasPrefix(account, v) {
			hit = true
			return
		}
	}
	return
}

func (e *AstExpression) handleFunctionEndsWith(checkHit bool, account string) (hit bool, err error) {
	if len(e.Arguments) != 2 {
		err = fmt.Errorf("%s function args length must two", e.Name)
		return
	}
	accountVar := e.Arguments[0]
	if accountVar.Type != Variable || VariableName(accountVar.Name) != Account {
		err = fmt.Errorf("first args type must variable and name is %s", Account)
		return
	}

	value := e.Arguments[1]
	strArray := gconv.Strings(value.Value)
	if len(strArray) == 0 || value.Type != Value || value.ValueType != StringArray {
		err = fmt.Errorf("function %s args[1] value must be []string and length must > 0", e.Name)
		return
	}

	if !checkHit {
		return
	}

	for _, v := range strArray {
		if strings.HasSuffix(account, v) {
			hit = true
			return
		}
	}
	return
}
