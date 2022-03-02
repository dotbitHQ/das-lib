package smt

import (
	"errors"
	"github.com/DeAccountSystems/das-lib/common"
)

var (
	StoreErrorNotExist = errors.New("not exist")
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
	Left  MergeValue `json:"l" bson:"l"`
	Right MergeValue `json:"r" bson:"r"`
}

type Store interface {
	GetBranch(key BranchKey) (*BranchNode, error)
	InsertBranch(key BranchKey, node BranchNode) error
	RemoveBranch(key BranchKey) error
	UpdateRoot(root H256) error
	Root() (H256, error)
}

type DefaultStore struct {
	root        H256
	branchesMap map[string]*BranchNode
}

func newDefaultStore() *DefaultStore {
	return &DefaultStore{
		root:        H256Zero(),
		branchesMap: make(map[string]*BranchNode),
	}
}

func (d *DefaultStore) UpdateRoot(root H256) error {
	copy(d.root, root)
	return nil
}

func (d *DefaultStore) Root() (H256, error) {
	return d.root, nil
}

func (d *DefaultStore) GetBranch(key BranchKey) (*BranchNode, error) {
	keyHash := key.GetHash()
	if item, ok := d.branchesMap[keyHash]; ok {
		return item, nil
	}
	return nil, StoreErrorNotExist
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
