package common

type ActionDataType = string

const (
	ActionDataTypeActionData                    ActionDataType = "0x00000000" // action data
	ActionDataTypeAccountCell                   ActionDataType = "0x01000000" // account cell
	ActionDataTypeAccountSaleCell               ActionDataType = "0x02000000" // account sale cell
	ActionDataTypeAccountAuctionCell            ActionDataType = "0x03000000" // account auction cell
	ActionDataTypeProposalCell                  ActionDataType = "0x04000000" // proposal cell
	ActionDataTypePreAccountCell                ActionDataType = "0x05000000" // pre account cell
	ActionDataTypeIncomeCell                    ActionDataType = "0x06000000" // income cell
	ActionDataTypeOfferCell                     ActionDataType = "0x07000000" // offer cell
	ActionDataTypeSubAccount                    ActionDataType = "0x08000000" // sub account
	ActionDataTypeSubAccountMintSign            ActionDataType = "0x09000000"
	ActionDataTypeReverseSmt                    ActionDataType = "0x0a000000" // reverse smt
	ActionDataTypeSubAccountPriceRules          ActionDataType = "0x0b000000"
	ActionDataTypeSubAccountPreservedRules      ActionDataType = "0x0c000000"
	ActionDataTypeKeyListCfgCell                ActionDataType = "0x0d000000" // keylist config cell
	ActionDataTypeSubAccountRenewSign           ActionDataType = "0x0e000000" // sub_account renew sign
	ActionDataTypeKeyListCfgCellData            ActionDataType = "0x0f000000" //
	ActionDataTypeSubAccountCreateApprovalSign  ActionDataType = "0x10000000"
	ActionDataTypeSubAccountDelayApprovalSign   ActionDataType = "0x11000000"
	ActionDataTypeSubAccountRevokeApprovalSign  ActionDataType = "0x12000000"
	ActionDataTypeSubAccountFulfillApprovalSign ActionDataType = "0x13000000"
)

const (
	WitnessDas                  = "das"
	WitnessDasCharLen           = 3
	WitnessDasTableTypeEndIndex = 7
)

type DataType = int

const (
	DataTypeNew          DataType = 0
	DataTypeOld          DataType = 1
	DataTypeDep          DataType = 2
	GoDataEntityVersion1 uint32   = 1
	GoDataEntityVersion2 uint32   = 2
	GoDataEntityVersion3 uint32   = 3
	GoDataEntityVersion4 uint32   = 4
)

type EditKey = string

const (
	EditKeyOwner        EditKey = "owner"
	EditKeyManager      EditKey = "manager"
	EditKeyRecords      EditKey = "records"
	EditKeyExpiredAt    EditKey = "expired_at"
	EditKeyManual       EditKey = "manual"
	EditKeyCustomRule   EditKey = "custom_rule"
	EditKeyCustomScript EditKey = "custom_script"
	EditKeyApproval     EditKey = "approval"
)

type SubAction = string

const (
	SubActionCreate           SubAction = "create"
	SubActionEdit             SubAction = "edit"
	SubActionRenew            SubAction = "renew"
	SubActionRecycle          SubAction = "recycle"
	SubActionCreateApproval   SubAction = "create_approval"
	SubActionDelayApproval    SubAction = "delay_approval"
	SubActionRevokeApproval   SubAction = "revoke_approval"
	SubActionFullFillApproval SubAction = "fulfill_approval"
)


const (
	WitnessDataSizeLimit = 32 * 1e3
)

type WebAuchonKeyOperate = string

const (
	AddWebAuthnKey    WebAuchonKeyOperate = "add"
	DeleteWebAuthnKey WebAuchonKeyOperate = "delete"
)
