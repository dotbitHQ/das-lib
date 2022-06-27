package smt

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/minio/blake2b-simd"
)

type MergeValue struct {
	Value     H256 `json:"v" bson:"v"`
	BaseNode  H256 `json:"n" bson:"n"`
	ZeroBits  H256 `json:"b" bson:"b"`
	ZeroCount byte `json:"c" bson:"c"`
}

func (m *MergeValue) Hash() H256 {
	if m.Value != nil {
		return m.Value
	} else {
		var tmp []byte
		tmp = append(tmp, MergeZeros)
		tmp = append(tmp, m.BaseNode...)
		tmp = append(tmp, m.ZeroBits...)
		tmp = append(tmp, m.ZeroCount)
		res, _ := smtBlake256(tmp)
		return res
	}
}

func (m *MergeValue) IsZero() bool {
	if m.Value != nil {
		return m.Value.IsZero()
	} else {
		return false
	}
}

func (m *MergeValue) String() string {
	if m.Value != nil {
		return common.Bytes2Hex(m.Value)
	} else {
		return fmt.Sprintf("b:%s,n:%s,c:%d", common.Bytes2Hex(m.ZeroBits), common.Bytes2Hex(m.BaseNode), m.ZeroCount)
	}
}

func MergeValueFromH256(value H256) MergeValue {
	return MergeValue{
		Value:     value,
		BaseNode:  nil,
		ZeroBits:  nil,
		ZeroCount: 0,
	}
}

func MergeValueFromZero() MergeValue {
	return MergeValueFromH256(H256Zero())
}

func HashBaseNode(height byte, baseKey, baseValue H256) H256 {
	var tmp []byte
	tmp = append(tmp, height)
	tmp = append(tmp, baseKey...)
	tmp = append(tmp, baseValue...)
	res, _ := smtBlake256(tmp)
	return res
}
func Merge(height byte, nodeKey H256, lhs, rhs MergeValue) MergeValue {
	lhsZero, rhsZero := lhs.IsZero(), rhs.IsZero()
	if lhsZero && rhsZero {
		return MergeValueFromZero()
	}
	if lhsZero {
		return MergeWithZero(height, nodeKey, rhs, true)
	}
	if rhsZero {
		return MergeWithZero(height, nodeKey, lhs, false)
	}

	var data []byte
	data = append(data, MergeNormal)
	data = append(data, height)
	data = append(data, nodeKey...)
	data = append(data, lhs.Hash()...)
	data = append(data, rhs.Hash()...)
	res, _ := smtBlake256(data)
	return MergeValueFromH256(res)
}

func smtBlake256(data []byte) ([]byte, error) {
	config := &blake2b.Config{
		Size:   32,
		Person: []byte(PersonSparseMerkleTree),
	}
	hash, err := blake2b.New(config)
	if err != nil {
		return nil, err
	}
	hash.Write(data)
	return hash.Sum(nil), nil
}

func MergeWithZero(height byte, nodeKey H256, value MergeValue, setBit bool) MergeValue {
	if value.Value != nil {
		zeroBits := H256Zero()
		if setBit {
			zeroBits.SetBit(height)
		}
		baseNode := HashBaseNode(height, nodeKey, value.Value)
		return MergeValue{
			Value:     nil,
			BaseNode:  baseNode,
			ZeroBits:  zeroBits,
			ZeroCount: 1,
		}
	} else {
		tmp := MergeValue{
			Value:     nil,
			BaseNode:  H256Zero(),
			ZeroBits:  H256Zero(),
			ZeroCount: value.ZeroCount,
		}
		tmp.ZeroCount++
		copy(tmp.ZeroBits, value.ZeroBits)
		copy(tmp.BaseNode, value.BaseNode)
		if setBit {
			tmp.ZeroBits.SetBit(height)
		}

		return tmp
	}
}
