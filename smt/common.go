package smt

import "crypto/sha256"

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
