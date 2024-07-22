package common

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type DidCellAction = string

const (
	DidCellActionDefault     DidCellAction = ""
	DidCellActionEditOwner   DidCellAction = "did_edit_owner"
	DidCellActionEditRecords DidCellAction = "did_edit_records"
	DidCellActionRenew       DidCellAction = "did_renew"

	DidCellActionUpgrade DidCellAction = "did_upgrade"
	DidCellActionUpdate  DidCellAction = "did_update"
	DidCellActionRecycle DidCellAction = "did_recycle"
)

type AnyLockCodeHash = string

const (
	AnyLockCodeHashOfMainnetOmniLock  AnyLockCodeHash = "0x9b819793a64463aed77c615d6cb226eea5487ccfc0783043a587254cda2b6f26"
	AnyLockCodeHashOfTestnetOmniLock  AnyLockCodeHash = "0xf329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb"
	AnyLockCodeHashOfMainnetJoyIDLock AnyLockCodeHash = "0xd00c84f0ec8fd441c38bc3f87a371f547190f2fcff88e642bc5bf54b9e318323"
	AnyLockCodeHashOfTestnetJoyIDLock AnyLockCodeHash = "0xd23761b364210735c19c60561d213fb3beae2fd6172743719eff6920e020baac"
	AnyLockCodeHashOfTestnetNoStrLock AnyLockCodeHash = "0x6ae5ee0cb887b2df5a9a18137315b9bdc55be8d52637b2de0624092d5f0c91d5"
	AnyLockCodeHashOfMainnetNoStrLock AnyLockCodeHash = "0x641a89ad2f77721b803cd50d01351c1f308444072d5fa20088567196c0574c68"
)

func GetDidCellTypeArgs(input *types.CellInput, outpointIndex uint64) ([]byte, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}
	bys, err := input.Serialize()
	if err != nil {
		return nil, fmt.Errorf("input.Serialize err: %s", err.Error())
	}
	bys2 := molecule.GoU64ToBytes(outpointIndex)
	bys = append(bys, bys2...)
	res, err := blake2b.Blake256(bys)
	if err != nil {
		return nil, fmt.Errorf("blake2b.Blake256: %s", err.Error())
	}
	return res, nil
}
