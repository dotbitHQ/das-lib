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
}

func TestAccountLengthPrice(t *testing.T) {
	rule := NewSubAccountRuleSlice()

	price100 := 100000000
	price10 := 10000000
	price1 := 100000

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
}
