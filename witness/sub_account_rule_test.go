package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAccountId(t *testing.T) {
	accounts := []string{"test.test.bit", "reverse.test.bit"}
	outs := make([]string, 0)
	for _, v := range accounts {
		out := common.Bytes2Hex(common.GetAccountIdByAccount(v))
		outs = append(outs, out)
	}
	t.Log(outs)
}

func TestRuleSpecialCharacters(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")

	price := 100000000

	err := rule.ParseFromJSON([]byte(fmt.Sprintf(`
{
    "version": 1,
    "rules": [
        {
            "name": "特殊字符账户",
            "note": "",
            "price": %d,
            "ast": {
                "type": "function",
                "name": "include_chars",
                "arguments": [
                    {
                        "type": "variable",
                        "name": "account"
                    },
                    {
                        "type": "value",
                        "value_type": "string[]",
                        "value": [
                            "⚠️",
                            "❌",
                            "✅"
                        ]
                    }
                ]
            },
			"status": 1
        }
    ]
}
`, price)))
	if err != nil {
		t.Fatal(err)
	}

	witness, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)
	for _, v := range witness {
		t.Log(common.Bytes2Hex(v))
	}

	hit, _, err := rule.Hit("jerry.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, idx, err := rule.Hit("jerry⚠️.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, rule.Rules[idx].Price, price)

	hit, idx, err = rule.Hit("jerry❌.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, rule.Rules[idx].Price, price)

	hit, idx, err = rule.Hit("jerry✅.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, rule.Rules[idx].Price, price)

	hit, _, err = rule.Hit("jerry💚.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	res, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)

	parseRules := NewSubAccountRuleEntity("test.bit")
	err = parseRules.ParseFromDasActionWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(parseRules.Rules), 1)

	assert.EqualValues(t, parseRules.Rules[0].Name, "特殊字符账户")
	assert.EqualValues(t, parseRules.Rules[0].Price, price)
	assert.EqualValues(t, parseRules.Rules[0].Status, 1)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Type, Function)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Name, FunctionIncludeCharts)
	assert.EqualValues(t, len(parseRules.Rules[0].Ast.Arguments), 2)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Arguments[0].Type, Variable)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Arguments[0].Name, Account)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Arguments[1].Type, Value)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Arguments[1].ValueType, StringArray)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Arguments[1].Value, []string{"⚠️", "❌", "✅"})
}

func TestAccountLengthPrice(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")

	price100 := uint64(100000000)
	price10 := uint64(10000000)
	price1 := uint64(100000)

	err := rule.ParseFromJSON([]byte(fmt.Sprintf(`
{
    "version": 1,
    "rules": [
        {
            "name": "1 位账户",
            "note": "",
            "price": %d,
            "ast": {
                "type": "operator",
                "symbol": "==",
                "expressions": [
                    {
                        "type": "variable",
                        "name": "account_length"
                    },
                    {
                        "type": "value",
                        "value_type": "uint8",
                        "value": 1
                    }
                ]
            },
			"status": 1
        },
        {
            "name": "2 位账户",
            "note": "",
            "price": %d,
            "ast": {
                "type": "operator",
                "symbol": "==",
                "expressions": [
                    {
                        "type": "variable",
                        "name": "account_length"
                    },
                    {
                        "type": "value",
                        "value_type": "uint8",
                        "value": 2
                    }
                ]
            },
			"status": 1
        },
        {
            "name": "8 位及以上账户",
            "note": "",
            "price": %d,
            "ast": {
                "type": "operator",
                "symbol": ">=",
                "expressions": [
                    {
                        "type": "variable",
                        "name": "account_length"
                    },
                    {
                        "type": "value",
                        "value_type": "uint8",
                        "value": 8
                    }
                ]
            },
			"status": 1
        }
    ]
}
`, price100, price10, price1)))
	if err != nil {
		t.Fatal(err)
	}

	witness, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)
	for _, v := range witness {
		t.Log(common.Bytes2Hex(v))
	}

	hit, idx, err := rule.Hit("1.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, rule.Rules[idx].Price, price100)

	hit, idx, err = rule.Hit("22.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 1)
	assert.EqualValues(t, rule.Rules[idx].Price, price10)

	hit, _, err = rule.Hit("333.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("4444.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("55555.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("666666.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("7777777.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, idx, err = rule.Hit("88888888.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 2)
	assert.EqualValues(t, rule.Rules[idx].Price, price1)

	hit, idx, err = rule.Hit("999999999.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 2)
	assert.EqualValues(t, rule.Rules[idx].Price, price1)

	res, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)

	parseRules := NewSubAccountRuleEntity("test.bit")
	err = parseRules.ParseFromDasActionWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(parseRules.Rules), 3)

	assert.EqualValues(t, parseRules.Rules[0].Name, "1 位账户")
	assert.EqualValues(t, parseRules.Rules[0].Price, price100)
	assert.EqualValues(t, parseRules.Rules[0].Status, 1)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Type, Operator)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Symbol, Equ)
	assert.EqualValues(t, len(parseRules.Rules[0].Ast.Expressions), 2)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Expressions[1].ValueType, Uint32)
	assert.EqualValues(t, parseRules.Rules[0].Ast.Expressions[1].Value, 1)

	assert.EqualValues(t, parseRules.Rules[1].Price, price10)
	assert.EqualValues(t, parseRules.Rules[1].Status, 1)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Type, Operator)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Symbol, Equ)
	assert.EqualValues(t, len(parseRules.Rules[1].Ast.Expressions), 2)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Expressions[1].ValueType, Uint32)
	assert.EqualValues(t, parseRules.Rules[1].Ast.Expressions[1].Value, 2)

	assert.EqualValues(t, parseRules.Rules[2].Price, price1)
	assert.EqualValues(t, parseRules.Rules[2].Status, 1)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Type, Operator)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Symbol, Gte)
	assert.EqualValues(t, len(parseRules.Rules[2].Ast.Expressions), 2)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Expressions[1].ValueType, Uint32)
	assert.EqualValues(t, parseRules.Rules[2].Ast.Expressions[1].Value, 8)

}

func TestRuleWhitelist(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")

	price := 100000000

	err := rule.ParseFromJSON([]byte(fmt.Sprintf(`
{
    "version": 1,
    "rules": [
        {
            "name": "特殊账户",
            "note": "",
            "price": %d,
            "ast": {
                "type": "function",
                "name": "in_list",
                "arguments": [
                    {
                        "type": "variable",
                        "name": "account"
                    },
                    {
                        "type": "value",
                        "value_type": "binary[]",
                        "value": [
                            "test",
                            "reverse"
                        ]
                    }
                ]
            },
			"status": 1
        }
    ]
}
`, price)))
	if err != nil {
		t.Fatal(err)
	}

	witness, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)
	for _, v := range witness {
		t.Log(common.Bytes2Hex(v))
	}

	res, err := rule.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)

	parseRules := NewSubAccountRuleEntity("test.bit")
	err = parseRules.ParseFromDasActionWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(parseRules.Rules), 1)

	hit, _, err := parseRules.Hit("jerry")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = parseRules.Hit("test")
	assert.NoError(t, err)
	assert.True(t, hit)

	hit, _, err = parseRules.Hit("reverse")
	assert.NoError(t, err)
	assert.True(t, hit)

	parseRule := parseRules.Rules[0]
	assert.EqualValues(t, parseRule.Name, "特殊账户")
	assert.EqualValues(t, parseRule.Note, "")
	assert.EqualValues(t, parseRule.Price, price)
	assert.EqualValues(t, parseRule.Status, 1)
	assert.EqualValues(t, parseRule.Ast.Type, "function")
	assert.EqualValues(t, parseRule.Ast.Name, "in_list")
	assert.EqualValues(t, len(parseRule.Ast.Arguments), 2)
	assert.EqualValues(t, parseRule.Ast.Arguments[0].Type, "variable")
	assert.EqualValues(t, parseRule.Ast.Arguments[0].Name, "account")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Type, "value")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].ValueType, "binary[]")
	assert.EqualValues(t, len(parseRule.Ast.Arguments[1].Value.([]string)), 2)
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Value.([]string)[0], "0x6ade4c435b8f3c4cf52336c9dd9dac71ed98520d")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Value.([]string)[1], "0xa84c83477c8f43670e70cef260da053818d770a5")
}

func TestSubAccountRuleDataSize(t *testing.T) {
	rulesEntity := NewSubAccountRuleEntity("test.bit")
	for i := 0; i < 100; i++ {
		rulesEntity.Rules = append(rulesEntity.Rules, &SubAccountRule{
			Index:  uint32(i),
			Name:   fmt.Sprintf("test rule %d", i),
			Note:   fmt.Sprintf("this is test rule %d", i),
			Price:  1,
			Status: 1,
			Ast: AstExpression{
				Type:   Operator,
				Symbol: And,
				Expressions: AstExpressions{
					{
						Type:   Operator,
						Symbol: Equ,
						Expressions: AstExpressions{
							{
								Type: Variable,
								Name: string(AccountLength),
							},
							{
								Type:      Value,
								ValueType: Uint32,
								Value:     uint32(i + 1),
							},
						},
					},
					{
						Type: Function,
						Name: string(FunctionOnlyIncludeCharset),
						Arguments: []*AstExpression{
							{
								Type: Variable,
								Name: string(AccountChars),
							},
							{
								Type:      Value,
								ValueType: Charset,
								Value:     common.AccountCharTypeEn,
							},
						},
					},
					{
						Type: Function,
						Name: string(FunctionIncludeWords),
						Arguments: []*AstExpression{
							{
								Type: Variable,
								Name: string(Account),
							},
							{
								Type:      Value,
								ValueType: StringArray,
								Value:     []string{"t"}},
						},
					},
				},
			},
		})
	}
	witness, err := rulesEntity.GenWitnessData(common.ActionDataTypeSubAccountPriceRules)
	assert.NoError(t, err)

	entity := NewSubAccountRuleEntity("test.bit")
	err = entity.ParseFromDasActionWitnessData(witness)
	assert.NoError(t, err)
}

func TestSubAccountRule_IncludeCharset(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")
	err := rule.ParseFromJSON([]byte(`
{
    "version": 1,
    "rules": [
        {
            "name": "func_include_charset",
            "note": "",
            "price": 1,
            "ast": {
                "type": "function",
                "name": "include_charset",
                "arguments": [
                    {
                        "type": "variable",
                        "name": "account_chars"
                    },
                    {
                        "type": "value",
                        "value_type": "charset_type",
                        "value": 2
                    }
                ]
            },
			"status": 1
        }
    ]
}
`))
	if err != nil {
		t.Fatal(err)
	}
	hit, _, err := rule.Hit("test")
	assert.NoError(t, err)
	assert.True(t, hit)
}

func TestSubAccountRule_StartsWith(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")
	err := rule.ParseFromJSON([]byte(`
{
    "version": 1,
    "rules": [
        {
            "name": "func_starts_with",
            "note": "",
            "price": 1,
            "ast": {
                "type": "function",
                "name": "starts_with",
                "arguments": [
                    {
                        "type": "variable",
                        "name": "account"
                    },
                    {
                        "type": "value",
                        "value_type": "string[]",
                        "value": [
                          "a"
                        ]
                    }
                ]
            },
			"status": 1
        }
    ]
}
`))
	if err != nil {
		t.Fatal(err)
	}
	hit, _, err := rule.Hit("test")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("abc")
	assert.NoError(t, err)
	assert.True(t, hit)
}

func TestSubAccountRule_EndsWith(t *testing.T) {
	rule := NewSubAccountRuleEntity("test.bit")
	err := rule.ParseFromJSON([]byte(`
{
    "version": 1,
    "rules": [
        {
            "name": "func_ends_with",
            "note": "",
            "price": 1,
            "ast": {
                "type": "function",
                "name": "ends_with",
                "arguments": [
                    {
                        "type": "variable",
                        "name": "account"
                    },
                    {
                        "type": "value",
                        "value_type": "string[]",
                        "value": [
                          "test"
                        ]
                    }
                ]
            },
			"status": 1
        }
    ]
}
`))
	if err != nil {
		t.Fatal(err)
	}
	hit, _, err := rule.Hit("test")
	assert.NoError(t, err)
	assert.True(t, hit)

	hit, _, err = rule.Hit("xxxtest")
	assert.NoError(t, err)
	assert.True(t, hit)

	hit, _, err = rule.Hit("testxxx")
	assert.NoError(t, err)
	assert.False(t, hit)
}
