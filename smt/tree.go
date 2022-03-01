package smt

import (
	"bytes"
	"fmt"
)

type SparseMerkleTree struct {
	Store Store
	Root  H256
}

func NewSparseMerkleTree(store Store) *SparseMerkleTree {
	if store == nil {
		store = newDefaultStore()
	}
	return &SparseMerkleTree{
		Store: store,
		Root:  make([]byte, 32),
	}
}

func (s *SparseMerkleTree) Update(key, value H256) error {
	currentKey := key
	currentNode := MergeValueFromH256(value)
	for i := 0; i <= MaxU8; i++ {
		height := byte(i)

		parentKey := currentKey.ParentPath(height)
		parentBranchKey := BranchKey{
			Height:  height,
			NodeKey: *parentKey,
		}
		var left, right MergeValue

		parentBranch, err := s.Store.GetBranch(parentBranchKey)
		if err == nil {
			if currentKey.IsRight(height) {
				left, right = parentBranch.Left, currentNode
			} else {
				left, right = currentNode, parentBranch.Right
			}
		} else if currentKey.IsRight(height) {
			left, right = MergeValueFromZero(), currentNode
		} else {
			left, right = currentNode, MergeValueFromZero()
		}

		if !left.IsZero() || !right.IsZero() {
			if err := s.Store.InsertBranch(parentBranchKey, BranchNode{
				Left:  left,
				Right: right,
			}); err != nil {
				return fmt.Errorf("InsertBranch err: %s", err.Error())
			}
		} else {
			if err := s.Store.RemoveBranch(parentBranchKey); err != nil {
				return fmt.Errorf("RemoveBranch err: %s", err.Error())
			}
		}

		currentKey = *parentKey
		currentNode = Merge(height, *parentKey, left, right)

	}
	s.Root = currentNode.Hash()
	return nil
}

func (s *SparseMerkleTree) merkleProof(keys []H256) (*MerkleProof, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("EmptyKeys")
	} else if len(keys) > 1 {
		return nil, fmt.Errorf("len(keys) > 1")
	}
	// sort keys keys.sort_unstable();
	// sort.Sort(SortH256(keys))

	// Collect leaf bitmaps
	var leavesBitMap []H256
	for _, currentKey := range keys {
		bitmap := H256Zero()
		for i := 0; i <= MaxU8; i++ {
			height := byte(i)
			parentKey := currentKey.ParentPath(height)
			parentBranchKey := BranchKey{
				Height:  height,
				NodeKey: *parentKey,
			}
			parentBranch, err := s.Store.GetBranch(parentBranchKey)
			if err == nil {
				var sibling MergeValue
				if currentKey.IsRight(height) {
					sibling = parentBranch.Left
				} else {
					sibling = parentBranch.Right
				}
				if !sibling.IsZero() {
					bitmap.SetBit(height)
				}
			} else {
				// The key is not in the tree (support non-inclusion proof)
			}
		}
		leavesBitMap = append(leavesBitMap, bitmap)
	}
	var proof []MergeValue
	stackForkHeight := make([]byte, MaxStackSize)
	stackTop, leafIndex := 0, 0
	for leafIndex < len(keys) {
		leafKey := keys[leafIndex]
		forkHeight := byte(MaxU8)
		if leafIndex+1 < len(keys) {
			forkHeight = leafKey.ForkHeight(&keys[leafIndex+1])
		}
		for i := 0; i <= int(forkHeight); i++ {
			height := byte(i)
			if height == forkHeight && leafIndex+1 < len(keys) {
				// If it's not final round, we don't need to merge to root (height=255)
				break
			}
			parentKey := leafKey.ParentPath(height)
			isRight := leafKey.IsRight(height)

			// has non-zero sibling
			if stackTop > 0 && stackForkHeight[stackTop-1] == height {
				stackTop -= 1
			} else if leavesBitMap[leafIndex].GetBit(height) {
				parentBranchKey := BranchKey{
					Height:  height,
					NodeKey: *parentKey,
				}
				parentBranch, err := s.Store.GetBranch(parentBranchKey)
				if err == nil {
					var sibling MergeValue
					if isRight {
						sibling = parentBranch.Left
					} else {
						sibling = parentBranch.Right
					}
					if !sibling.IsZero() {
						proof = append(proof, sibling)
					} else {
						// unreachable!();
						return nil, fmt.Errorf("unreachable")
					}
				} else {
					// The key is not in the tree (support non-inclusion proof)
				}
			}
		}
		stackForkHeight[stackTop] = forkHeight
		stackTop++
		leafIndex++
	}
	return &MerkleProof{
		LeavesBitmap: leavesBitMap,
		MerklePath:   proof,
	}, nil
}

func (s *SparseMerkleTree) MerkleProof(keys, values []H256) (*CompiledMerkleProof, error) {
	merkleProof, err := s.merkleProof(keys)
	if err != nil {
		return nil, fmt.Errorf("merkleProof: %s", err.Error())
	}
	return merkleProof.compile(keys, values)
}

func Verify(root H256, proof *CompiledMerkleProof, keys, values []H256) (bool, error) {
	if proof == nil {
		return false, fmt.Errorf("proof is nil")
	}
	calculatedRoot, err := proof.computeRoot(keys, values)
	if err != nil {
		return false, fmt.Errorf("ComputeRoot err: %s", err.Error())
	}
	return bytes.Compare(root, *calculatedRoot) == 0, nil
}
