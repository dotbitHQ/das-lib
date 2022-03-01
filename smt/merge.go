package smt

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/minio/blake2b-simd"
)

type MergeValue interface {
	Hash() H256
	IsZero() bool
	String() string
}

type MergeValueH256 H256

func (m *MergeValueH256) Hash() H256 {
	return H256(*m)
}

func (m *MergeValueH256) IsZero() bool {
	tmp := H256(*m)
	return tmp.IsZero()
}

func (m *MergeValueH256) String() string {
	return common.Bytes2Hex(*m)
}

type MergeValueZero struct {
	BaseNode  H256
	ZeroBits  H256
	ZeroCount byte
}

func (m *MergeValueZero) Hash() H256 {
	var tmp []byte
	tmp = append(tmp, MergeZeros)
	tmp = append(tmp, m.BaseNode...)
	tmp = append(tmp, m.ZeroBits...)
	tmp = append(tmp, m.ZeroCount)
	res, _ := smtBlake256(tmp)
	return res
}

func (m *MergeValueZero) IsZero() bool {
	return false
}

func (m *MergeValueZero) String() string {
	return fmt.Sprintf("ZeroBits:%s,BaseNode:%s,ZeroCount:%d", common.Bytes2Hex(m.ZeroBits), common.Bytes2Hex(m.BaseNode), m.ZeroCount)
}

func MergeValueFromH256(value H256) MergeValue {
	tmp := MergeValueH256(value)
	return &tmp
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
	if v, ok := (value).(*MergeValueH256); ok {
		zeroBits := H256Zero()
		if setBit {
			zeroBits.SetBit(height)
		}
		tmp := H256(*v)
		baseNode := HashBaseNode(height, nodeKey, tmp)

		return &MergeValueZero{
			BaseNode:  baseNode,
			ZeroBits:  zeroBits,
			ZeroCount: 1,
		}
	} else if z, ok := (value).(*MergeValueZero); ok {
		tmp := MergeValueZero{
			BaseNode:  H256Zero(),
			ZeroBits:  H256Zero(),
			ZeroCount: z.ZeroCount,
		}
		copy(tmp.ZeroBits, z.ZeroBits)
		copy(tmp.BaseNode, z.BaseNode)
		if setBit {
			tmp.ZeroBits.SetBit(height)
		}
		tmp.ZeroCount++
		return &tmp
	}
	return nil
}
