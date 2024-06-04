package common

type DidCellAction = string

const (
	DidCellActionDefault     DidCellAction = ""
	DidCellActionEditOwner   DidCellAction = "did_edit_owner"
	DidCellActionEditRecords DidCellAction = "did_edit_records"
	DidCellActionRenew       DidCellAction = "did_renew"
	DidCellActionRecycle     DidCellAction = "did_recycle"

	DidCellActionUpgrade  DidCellAction = "did_upgrade"
	DidCellActionRegister DidCellAction = "did_register"
	DidCellActionAuction  DidCellAction = "did_auction"
)

type AnyLockCodeHash = string

const (
	AnyLockCodeHashOfMainnetOmniLock  AnyLockCodeHash = "0x9b819793a64463aed77c615d6cb226eea5487ccfc0783043a587254cda2b6f26"
	AnyLockCodeHashOfTestnetOmniLock  AnyLockCodeHash = "0xf329effd1c475a2978453c8600e1eaf0bc2087ee093c3ee64cc96ec6847752cb"
	AnyLockCodeHashOfMainnetJoyIDLock AnyLockCodeHash = "0xd00c84f0ec8fd441c38bc3f87a371f547190f2fcff88e642bc5bf54b9e318323"
	AnyLockCodeHashOfTestnetJoyIDLock AnyLockCodeHash = "0xd23761b364210735c19c60561d213fb3beae2fd6172743719eff6920e020baac"
)
