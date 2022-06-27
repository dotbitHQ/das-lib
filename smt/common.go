package smt

import (
	"crypto/sha256"
	"github.com/dotbitHQ/das-lib/common"
)

const (
	ByteSize     = 8
	MaxU8        = 255
	MaxStackSize = 257
)

const (
	MergeNormal            byte = 1
	MergeZeros             byte = 2
	PersonSparseMerkleTree      = "ckb-default-hash"
)

func Sha256(src string) []byte {
	m := sha256.New()
	m.Write([]byte(src))
	return m.Sum(nil)
}

func AccountIdToSmtH256(accountId string) H256 {
	bys := common.Hex2Bytes(accountId)
	key := H256Zero()
	copy(key, bys)
	return key
}
