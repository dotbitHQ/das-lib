package smt

import (
	"bytes"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
)

type MerkleProof struct {
	LeavesBitmap []H256
	MerklePath   []MergeValue
}

func (m *MerkleProof) compile(keys, values []H256) (*CompiledMerkleProof, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("keys is nil")
	} else if len(keys) != len(values) {
		return nil, fmt.Errorf("len(keys) != len(values)")
	} else if len(keys) != len(m.LeavesBitmap) {
		return nil, fmt.Errorf("len(keys) != LeavesCount")
	}
	// sort keys leaves.sort_unstable_by_key(|(k, _v)| *k);
	// sort.Sort(SortH256(keys))

	var proof []byte
	stackForkHeight := make([]byte, MaxStackSize)
	stackTop := 0
	leafIndex := 0
	merklePathIndex := 0
	for leafIndex < len(keys) {
		leafKey, _ := keys[leafIndex], values[leafIndex]
		forkHeight := byte(MaxU8)
		if leafIndex+1 < len(keys) {
			forkHeight = leafKey.ForkHeight(&keys[leafIndex+1])
		}
		proof = append(proof, common.Hex2Bytes("0x4C")...)
		zeroCount := uint16(0)
		for i := 0; i <= int(forkHeight); i++ {
			height := byte(i)
			if height == forkHeight && leafIndex+1 < len(keys) {
				// If it's not final round, we don't need to merge to root (height=255)
				break
			}
			var opCodeOpt string
			var siblingDataOpt []byte
			if stackTop > 0 && stackForkHeight[stackTop-1] == height {
				stackTop -= 1
				opCodeOpt = "0x48"
			} else if m.LeavesBitmap[leafIndex].GetBit(height) {
				if merklePathIndex >= len(m.MerklePath) {
					return nil, fmt.Errorf("CorruptedProof")
				}
				node := m.MerklePath[merklePathIndex]
				merklePathIndex += 1
				if node.Value != nil {
					opCodeOpt = "0x50"
					siblingDataOpt = node.Hash()
				} else {
					var buffer []byte
					buffer = append(buffer, node.ZeroCount)
					buffer = append(buffer, node.BaseNode...)
					buffer = append(buffer, node.ZeroBits...)
					opCodeOpt = "0x51"
					siblingDataOpt = buffer
				}
			} else {
				zeroCount += 1
				if zeroCount > 256 {
					return nil, fmt.Errorf("CorruptedProof")
				}
			}

			if opCode := opCodeOpt; opCode != "" {
				if zeroCount > 0 {
					n := byte(0)
					if zeroCount != 256 {
						n = byte(zeroCount)
					}
					proof = append(proof, common.Hex2Bytes("0x4F")...)
					proof = append(proof, n)
					zeroCount = 0
				}
				// note: opCode to []byte
				proof = append(proof, common.Hex2Bytes(opCode)...)
			}
			data := siblingDataOpt
			proof = append(proof, data...)
		}
		if zeroCount > 0 {
			n := byte(0)
			if zeroCount != 256 {
				n = byte(zeroCount)
			}
			proof = append(proof, common.Hex2Bytes("0x4F")...)
			proof = append(proof, n)
		}
		stackForkHeight[stackTop] = forkHeight
		stackTop += 1
		leafIndex += 1
	}
	if stackTop != 1 {
		return nil, fmt.Errorf("CorruptedProof")
	}
	if leafIndex != len(keys) {
		return nil, fmt.Errorf("CorruptedProof")
	}
	if merklePathIndex != len(m.MerklePath) {
		return nil, fmt.Errorf("CorruptedProof")
	}
	res := CompiledMerkleProof(proof)

	return &res, nil
}

type CompiledMerkleProof []byte

func (c *CompiledMerkleProof) computeRoot(keys, values []H256) (*H256, error) {
	// leaves.sort_unstable_by_key(|(k, _v)| *k);
	// sort.Sort(SortH256(keys))

	var stackHeights = make([]uint16, MaxStackSize)
	var stackKeys = make([]H256, MaxStackSize)
	var stackValues = make([]MergeValue, MaxStackSize)

	var proofIndex, proofLen, leaveIndex, stackTop = 0, len(*c), 0, 0
	for proofIndex < proofLen {
		code := (*c)[proofIndex]
		proofIndex++
		switch code {
		case 0x4C: // L : push leaf value
			if stackTop >= MaxStackSize {
				return nil, fmt.Errorf("CorruptedStack")
			}
			if leaveIndex >= len(keys) {
				return nil, fmt.Errorf("CorruptedStack")
			}
			k, v := keys[leaveIndex], values[leaveIndex]
			stackHeights[stackTop] = 0
			stackKeys[stackTop] = k
			stackValues[stackTop] = MergeValueFromH256(v)
			leaveIndex++
			stackTop++
		case 0x50: // P : hash stack top item with sibling node in proof
			if stackTop == 0 {
				return nil, fmt.Errorf("CorruptedStack")
			}
			if proofIndex+32 > proofLen {
				return nil, fmt.Errorf("CorruptedProof")
			}
			proofTmp := (*c)[proofIndex : proofIndex+32]
			siblingNode := MergeValueFromH256(H256(proofTmp))
			proofIndex += 32

			height, key, value := stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1]

			if height > 255 {
				return nil, fmt.Errorf("CorruptedProof")
			}
			heightPtr := byte(height) // u16 to u8
			parentKey := key.ParentPath(heightPtr)
			var parent MergeValue
			if key.GetBit(heightPtr) {
				parent = Merge(heightPtr, *parentKey, siblingNode, value)
			} else {
				parent = Merge(heightPtr, *parentKey, value, siblingNode)
			}

			stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1] = height+1, *parentKey, parent
		case 0x51:
			// Q : hash stack top item with sibling node in proof,
			// this is similar to P except that proof comes in using
			// MergeWithZero format.

			if stackTop == 0 {
				return nil, fmt.Errorf("CorruptedStack")
			}
			if proofIndex+65 > proofLen {
				return nil, fmt.Errorf("CorruptedProof")
			}

			zeroCount := (*c)[proofIndex]
			baseNode := H256((*c)[proofIndex+1 : proofIndex+33])
			zeroBits := H256((*c)[proofIndex+33 : proofIndex+65])
			proofIndex += 65
			siblingNode := MergeValue{
				Value:     nil,
				BaseNode:  baseNode,
				ZeroBits:  zeroBits,
				ZeroCount: zeroCount,
			}

			heightU16, key, value := stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1]
			if heightU16 > 255 {
				return nil, fmt.Errorf("CorruptedProof")
			}
			height := byte(heightU16) // u16 to u8
			parentKey := key.ParentPath(height)
			var parent MergeValue
			if key.GetBit(height) {
				parent = Merge(height, *parentKey, siblingNode, value)
			} else {
				parent = Merge(height, *parentKey, value, siblingNode)
			}

			stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1] = heightU16+1, *parentKey, parent
		case 0x48:
			// H : pop 2 items in stack hash them then push the result

			if stackTop < 2 {
				return nil, fmt.Errorf("CorruptedStack")
			}
			heightA, keyA, valueA := stackHeights[stackTop-2], stackKeys[stackTop-2], stackValues[stackTop-2]
			heightB, keyB, valueB := stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1]
			stackTop--

			if heightA != heightB {
				return nil, fmt.Errorf("CorruptedProof")
			}
			if heightA > 255 {
				return nil, fmt.Errorf("CorruptedProof")
			}

			heightPtr := byte(heightA)
			parentKeyA := keyA.ParentPath(heightPtr)
			parentKeyB := keyB.ParentPath(heightPtr)

			if bytes.Compare(*parentKeyA, *parentKeyB) != 0 {
				return nil, fmt.Errorf("CorruptedProof")
			}

			var parent MergeValue
			if keyA.GetBit(heightPtr) {
				parent = Merge(heightPtr, *parentKeyA, valueB, valueA)
			} else {
				parent = Merge(heightPtr, *parentKeyA, valueA, valueB)
			}

			stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1] = heightA+1, *parentKeyA, parent
		case 0x4F:
			// O : hash stack top item with n zero values
			if stackTop < 1 {
				return nil, fmt.Errorf("CorruptedStack")
			}
			if proofIndex >= proofLen {
				return nil, fmt.Errorf("CorruptedProof")
			}
			n := (*c)[proofIndex]
			proofIndex++

			zeroCount := uint16(256)
			if n != 0 {
				zeroCount = uint16(n)
			}
			baseHeight, key, value := stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1]
			if baseHeight > 255 {
				return nil, fmt.Errorf("CorruptedProof")
			}

			var parentKey *H256
			heightU16 := baseHeight
			for idx := uint16(0); idx < zeroCount; idx++ {
				if baseHeight+idx > 255 {
					return nil, fmt.Errorf("CorruptedProof")
				}
				heightU16 = baseHeight + idx
				height := byte(heightU16)
				parentKey = key.ParentPath(height)

				if key.GetBit(height) {
					value = Merge(height, *parentKey, MergeValueFromZero(), value)
				} else {
					value = Merge(height, *parentKey, value, MergeValueFromZero())
				}
			}
			stackHeights[stackTop-1], stackKeys[stackTop-1], stackValues[stackTop-1] = heightU16+1, *parentKey, value
		default:
			return nil, fmt.Errorf("CorruptedProof")
		}
	}
	if stackTop != 1 {
		return nil, fmt.Errorf("CorruptedStack")
	}
	if stackHeights[0] != 256 {
		return nil, fmt.Errorf("CorruptedProof")
	}
	if leaveIndex != len(keys) {
		return nil, fmt.Errorf("CorruptedProof")
	}
	root := stackValues[0].Hash()
	return &root, nil
}

func (c *CompiledMerkleProof) String() string {
	return common.Bytes2Hex(*c)
}
