package common

type DidCellAction string

const (
	DidCellActionDefault     DidCellAction = ""
	DidCellActionEditOwner   DidCellAction = "edit_owner"
	DidCellActionEditRecords DidCellAction = "edit_records"
	DidCellActionRenew       DidCellAction = "renew"
	DidCellActionRecycle     DidCellAction = "recycle"

	DidCellActionUpgrade  DidCellAction = "upgrade"
	DidCellActionRegister DidCellAction = "register"
	DidCellActionAuction  DidCellAction = "auction"
)
