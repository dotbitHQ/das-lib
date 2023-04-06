package witness

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAccountId(t *testing.T) {
	accounts := []string{"test", "reverse"}
	outs := make([]string, 0)
	for _, v := range accounts {
		out := common.Bytes2Hex(common.Blake2b([]byte(v))[:20])
		outs = append(outs, out)
	}
	t.Log(outs)
}

func TestRuleSpecialCharacters(t *testing.T) {
	rule := NewSubAccountRuleSlice()

	price := 100000000

	err := rule.Parser([]byte(fmt.Sprintf(`
[
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
                    "name": "account_chars"
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
        }
    }
]
`, price)))
	if err != nil {
		t.Fatal(err)
	}

	hit, _, err := rule.Hit("jerry.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, idx, err := rule.Hit("jerry⚠️.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, (*rule)[idx].Price, price)

	hit, idx, err = rule.Hit("jerry❌.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, (*rule)[idx].Price, price)

	hit, idx, err = rule.Hit("jerry✅.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, (*rule)[idx].Price, price)

	hit, _, err = rule.Hit("jerry💚.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	res := rule.GenWitnessData()
	t.Log(common.Bytes2Hex(res))

	parseRules := NewSubAccountRuleSlice()
	err = parseRules.ParseFromWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(*parseRules), 1)

	assert.EqualValues(t, (*parseRules)[0].Name, "特殊字符账户")
	assert.EqualValues(t, (*parseRules)[0].Price, price)
	assert.EqualValues(t, (*parseRules)[0].Ast.Type, Function)
	assert.EqualValues(t, (*parseRules)[0].Ast.Name, FunctionIncludeCharts)
	assert.EqualValues(t, len((*parseRules)[0].Ast.Arguments), 2)
	assert.EqualValues(t, (*parseRules)[0].Ast.Arguments[0].Type, Variable)
	assert.EqualValues(t, (*parseRules)[0].Ast.Arguments[0].Name, AccountChars)
	assert.EqualValues(t, (*parseRules)[0].Ast.Arguments[1].Type, Value)
	assert.EqualValues(t, (*parseRules)[0].Ast.Arguments[1].ValueType, StringArray)
	assert.EqualValues(t, (*parseRules)[0].Ast.Arguments[1].Value, []string{"⚠️", "❌", "✅"})
}

func TestAccountLengthPrice(t *testing.T) {
	rule := NewSubAccountRuleSlice()

	price100 := uint64(100000000)
	price10 := uint64(10000000)
	price1 := uint64(100000)

	err := rule.Parser([]byte(fmt.Sprintf(`
[
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
        }
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
        }
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
        }
    }
]
`, price100, price10, price1)))
	if err != nil {
		t.Fatal(err)
	}

	hit, idx, err := rule.Hit("1.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 0)
	assert.EqualValues(t, (*rule)[idx].Price, price100)

	hit, idx, err = rule.Hit("22.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 1)
	assert.EqualValues(t, (*rule)[idx].Price, price10)

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
	assert.EqualValues(t, (*rule)[idx].Price, price1)

	hit, idx, err = rule.Hit("999999999.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
	assert.Equal(t, idx, 2)
	assert.EqualValues(t, (*rule)[idx].Price, price1)

	res := rule.GenWitnessData()
	t.Log(common.Bytes2Hex(res))

	parseRules := NewSubAccountRuleSlice()
	err = parseRules.ParseFromWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(*parseRules), 3)

	assert.EqualValues(t, (*parseRules)[0].Name, "1 位账户")
	assert.EqualValues(t, (*parseRules)[0].Price, price100)
	assert.EqualValues(t, (*parseRules)[0].Ast.Type, Operator)
	assert.EqualValues(t, (*parseRules)[0].Ast.Symbol, Equ)
	assert.EqualValues(t, len((*parseRules)[0].Ast.Expressions), 2)
	assert.EqualValues(t, (*parseRules)[0].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, (*parseRules)[0].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, (*parseRules)[0].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, (*parseRules)[0].Ast.Expressions[1].ValueType, Uint8)
	assert.EqualValues(t, (*parseRules)[0].Ast.Expressions[1].Value, 1)

	assert.EqualValues(t, (*parseRules)[1].Price, price10)
	assert.EqualValues(t, (*parseRules)[1].Ast.Type, Operator)
	assert.EqualValues(t, (*parseRules)[1].Ast.Symbol, Equ)
	assert.EqualValues(t, len((*parseRules)[1].Ast.Expressions), 2)
	assert.EqualValues(t, (*parseRules)[1].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, (*parseRules)[1].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, (*parseRules)[1].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, (*parseRules)[1].Ast.Expressions[1].ValueType, Uint8)
	assert.EqualValues(t, (*parseRules)[1].Ast.Expressions[1].Value, 2)

	assert.EqualValues(t, (*parseRules)[2].Price, price1)
	assert.EqualValues(t, (*parseRules)[2].Ast.Type, Operator)
	assert.EqualValues(t, (*parseRules)[2].Ast.Symbol, Gte)
	assert.EqualValues(t, len((*parseRules)[2].Ast.Expressions), 2)
	assert.EqualValues(t, (*parseRules)[2].Ast.Expressions[0].Type, Variable)
	assert.EqualValues(t, (*parseRules)[2].Ast.Expressions[0].Name, AccountLength)
	assert.EqualValues(t, (*parseRules)[2].Ast.Expressions[1].Type, Value)
	assert.EqualValues(t, (*parseRules)[2].Ast.Expressions[1].ValueType, Uint8)
	assert.EqualValues(t, (*parseRules)[2].Ast.Expressions[1].Value, 8)

}

func TestRuleWhitelist(t *testing.T) {
	rule := NewSubAccountRuleSlice()

	price := 100000000

	err := rule.Parser([]byte(fmt.Sprintf(`
[
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
                        "0xc8988da7aa74e876576c44b1c0ac814457b3a461",
                        "0xcf45bb5b316a3d009fdc782dee25bf941a9daf0e"
                    ]
                }
            ]
        }
    }
]
`, price)))
	if err != nil {
		t.Fatal(err)
	}

	hit, _, err := rule.Hit("jerry.bit")
	assert.NoError(t, err)
	assert.False(t, hit)

	hit, _, err = rule.Hit("test.bit")
	assert.NoError(t, err)
	assert.True(t, hit)

	hit, _, err = rule.Hit("reverse.bit")
	assert.NoError(t, err)
	assert.True(t, hit)
}

func TestSubAccountRuleSlice_GenWitnessData(t *testing.T) {
	rule := NewSubAccountRuleSlice()

	price := 100000000

	err := rule.Parser([]byte(fmt.Sprintf(`
[
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
                        "0xb28072bd0201e6feeb4cd96a6879d6422f2218cd",
                        "0x75bc2d3192ec310b6ac2f826d3e19a5cfe9f080a"
                    ]
                }
            ]
        }
    }
]
`, price)))
	if err != nil {
		t.Fatal(err)
	}

	res := rule.GenWitnessData()
	t.Log(common.Bytes2Hex(res))

	parseRules := NewSubAccountRuleSlice()
	err = parseRules.ParseFromWitnessData(res)
	assert.NoError(t, err)
	assert.EqualValues(t, len(*parseRules), 1)

	parseRule := (*parseRules)[0]
	assert.EqualValues(t, parseRule.Name, "特殊账户")
	assert.EqualValues(t, parseRule.Note, "")
	assert.EqualValues(t, parseRule.Price, price)
	assert.EqualValues(t, parseRule.Ast.Type, "function")
	assert.EqualValues(t, parseRule.Ast.Name, "in_list")
	assert.EqualValues(t, len(parseRule.Ast.Arguments), 2)
	assert.EqualValues(t, parseRule.Ast.Arguments[0].Type, "variable")
	assert.EqualValues(t, parseRule.Ast.Arguments[0].Name, "account")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Type, "value")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].ValueType, "binary[]")
	assert.EqualValues(t, len(parseRule.Ast.Arguments[1].Value.([]string)), 2)
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Value.([]string)[0], "0xb28072bd0201e6feeb4cd96a6879d6422f2218cd")
	assert.EqualValues(t, parseRule.Ast.Arguments[1].Value.([]string)[1], "0x75bc2d3192ec310b6ac2f826d3e19a5cfe9f080a")
}
