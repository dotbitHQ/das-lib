package smt

import (
	"bytes"
	"github.com/DeAccountSystems/das-lib/common"
)

type H256 []byte

func H256Zero() H256 {
	zero := make(H256, 32)
	return zero
}

func (h *H256) String() string {
	return common.Bytes2Hex(*h)
}

func (h *H256) IsZero() bool {
	zero := H256Zero()
	if bytes.Compare(*h, zero) == 0 {
		return true
	}
	return false
}

func (h *H256) SetBit(height byte) {
	bytePos := height / 8
	bitPos := height % 8
	(*h)[bytePos] |= 1 << bitPos
}

func (h *H256) CopyBits(height byte) *H256 {
	target := H256Zero()
	startByte := height / ByteSize
	copy(target[startByte:], (*h)[startByte:])
	remain := height % ByteSize
	if remain > 0 {
		target[startByte] &= 0b11111111 << remain
	}
	return &target
}

func (h *H256) ClearBit(height byte) {
	bytePos := height / 8
	bitPos := height % 8
	(*h)[bytePos] &= ^(1 << bitPos)
}

func (h *H256) ParentPath(height byte) *H256 {
	if height == 255 {
		tmp := H256Zero()
		return &tmp
	} else {
		return h.CopyBits(height + 1)
	}
}

func (h *H256) IsRight(height byte) bool {
	return h.GetBit(height)
}

func (h *H256) GetBit(height byte) bool {
	bytePos := height / 8
	bitPos := height % 8
	return (((*h)[bytePos] >> bitPos) & 1) != 0
}

func (h *H256) ForkHeight(key *H256) byte {
	for i := 0; i <= MaxU8; i++ {
		height := byte(i)
		if h.GetBit(height) != key.GetBit(height) {
			return height
		}
	}
	return 0
}

// qsort
type SortH256 []H256

func (s SortH256) Len() int { return len(s) }
func (s SortH256) Less(i, j int) bool {
	if res := bytes.Compare(s[i], s[j]); res > 0 {
		return false
	} else {
		return true
	}
}
func (s SortH256) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
