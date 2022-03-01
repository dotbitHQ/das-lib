package smt

import (
	"errors"
	"github.com/DeAccountSystems/das-lib/common"
)

var (
	StoreErrorNotExistBranch = errors.New("not exist branch")
)

type BranchKey struct {
	Height  byte
	NodeKey H256
}

func (b *BranchKey) GetHash() string {
	var tmp []byte
	tmp = append(tmp, b.Height)
	tmp = append(tmp, byte('-'))
	tmp = append(tmp, b.NodeKey...)
	return common.Bytes2Hex(tmp)
}

type BranchNode struct {
	Left  MergeValue
	Right MergeValue
}

type Store interface {
	GetBranch(key BranchKey) (*BranchNode, error)
	InsertBranch(key BranchKey, node BranchNode) error
	RemoveBranch(key BranchKey) error
}

type DefaultStore struct {
	branchesMap map[string]*BranchNode
}

func newDefaultStore() *DefaultStore {
	return &DefaultStore{
		branchesMap: make(map[string]*BranchNode),
	}
}

func (d *DefaultStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	if item, ok := d.branchesMap[keyHash]; ok {
		return item, nil
	}
	return nil, StoreErrorNotExistBranch
}

func (d *DefaultStore) InsertBranch(key BranchKey, node BranchNode) error {
	keyHash := key.GetHash()
	d.branchesMap[keyHash] = &BranchNode{
		Left:  node.Left,
		Right: node.Right,
	}
	return nil
}

func (d *DefaultStore) RemoveBranch(key BranchKey) error {
	keyHash := key.GetHash()
	delete(d.branchesMap, keyHash)
	return nil
}
