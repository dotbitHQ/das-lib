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

type AnyLockCodeHash string

const ()
