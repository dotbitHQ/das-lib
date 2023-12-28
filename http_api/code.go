package http_api

type ApiCode = int
type ApiResp struct {
	ErrNo  ApiCode     `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func ApiRespOK(data interface{}) ApiResp {
	return ApiResp{
		ErrNo:  ApiCodeSuccess,
		ErrMsg: "",
		Data:   data,
	}
}

func ApiRespErr(errNo ApiCode, errMsg string) ApiResp {
	return ApiResp{
		ErrNo:  errNo,
		ErrMsg: errMsg,
		Data:   nil,
	}
}

func (a *ApiResp) ApiRespErr(errNo ApiCode, errMsg string) {
	a.ErrNo = errNo
	a.ErrMsg = errMsg
}

func (a *ApiResp) ApiRespOK(data interface{}) {
	a.ErrNo = ApiCodeSuccess
	a.Data = data
}

const (
	TextSystemUpgrade = "The service is under maintenance, please try again later."
)

// common
const (
	ApiCodeSuccess        ApiCode = 0
	ApiCodeError500       ApiCode = 500
	ApiCodeParamsInvalid  ApiCode = 10000
	ApiCodeMethodNotExist ApiCode = 10001
	ApiCodeDbError        ApiCode = 10002
	ApiCodeCacheError     ApiCode = 10003

	ApiCodeTransactionNotExist ApiCode = 11001
	//ApiCodePermissionDenied    ApiCode = 11002
	ApiCodeNotSupportAddress   ApiCode = 11005
	ApiCodeInsufficientBalance ApiCode = 11007
	ApiCodeTxExpired           ApiCode = 11008
	ApiCodeAmountInvalid       ApiCode = 11010
	ApiCodeRejectedOutPoint    ApiCode = 11011
	ApiCodeSyncBlockNumber     ApiCode = 11012
	ApiCodeOperationFrequent   ApiCode = 11013
	ApiCodeNotEnoughChange     ApiCode = 11014
)

// reverse
const (
	ApiCodeReverseAlreadyExist ApiCode = 12001
	ApiCodeReverseNotExist     ApiCode = 12002
)

// account-indexer
const (
	ApiCodeAccountFormatInvalid   ApiCode = 20006
	ApiCodeIndexerAccountNotExist ApiCode = 20007
	ApiCodeAccountOnLock          ApiCode = 20008
)

// register
const (
	ApiCodeNotOpenForRegistration       ApiCode = 30001
	ApiCodeAccountNotExist              ApiCode = 30003
	ApiCodeAccountAlreadyRegister       ApiCode = 30004
	ApiCodeAccountLenInvalid            ApiCode = 30014
	ApiCodeOrderNotExist                ApiCode = 30006
	ApiCodeAccountIsExpired             ApiCode = 30010
	ApiCodePermissionDenied             ApiCode = 30011
	ApiCodeAccountContainsInvalidChar   ApiCode = 30015
	ApiCodeReservedAccount              ApiCode = 30017
	ApiCodeInviterAccountNotExist       ApiCode = 30018
	ApiCodeSystemUpgrade                ApiCode = 30019
	ApiCodeRecordInvalid                ApiCode = 30020
	ApiCodeRecordsTotalLengthExceeded   ApiCode = 30021
	ApiCodeSameLock                     ApiCode = 30023
	ApiCodeChannelAccountNotExist       ApiCode = 30026
	ApiCodeOrderPaid                    ApiCode = 30027
	ApiCodeUnAvailableAccount           ApiCode = 30029
	ApiCodeAccountStatusNotNormal       ApiCode = 30031 //repeat
	ApiCodeAccountStatusOnSaleOrAuction ApiCode = 30031
	ApiCodePayTypeInvalid               ApiCode = 30032
	ApiCodeSameOrderInfo                ApiCode = 30033
	ApiCodeSigErr                       ApiCode = 30034 // contracte -31
	ApiCodeOnCross                      ApiCode = 30035
	ApiCodeSubAccountNotEnabled         ApiCode = 30036
	ApiCodeParentAccountExpired         ApiCode = 30036
	ApiCodeAfterGracePeriod             ApiCode = 30037
	ApiCodeCouponInvalid                ApiCode = 30038
	ApiCodeCouponUsed                   ApiCode = 30039
	ApiCodeCouponUnopen                 ApiCode = 30040
	ApiCodeReverseSmtPending            ApiCode = 30040
	ApiCodeAccountStatusOnCross         ApiCode = 30041
	ApiCodeNoAccountPermissions         ApiCode = 30042

	ApiCodeAuctionAccountNotFound ApiCode = 30404
	ApiCodeAuctionAccountBided    ApiCode = 30405
	ApiCodeAuctionOrderNotFound   ApiCode = 30406
)

// sub_account
const (
	ApiCodeEnableSubAccountIsOn               ApiCode = 40000
	ApiCodeNotExistEditKey                    ApiCode = 40001
	ApiCodeNotExistConfirmAction              ApiCode = 40002
	ApiCodeSignError                          ApiCode = 40003
	ApiCodeNotExistSignType                   ApiCode = 40004
	ApiCodeNotSubAccount                      ApiCode = 40005
	ApiCodeEnableSubAccountIsOff              ApiCode = 40006
	ApiCodeCreateListCheckFail                ApiCode = 40007
	ApiCodeTaskInProgress                     ApiCode = 40008
	ApiCodeDistributedLockPreemption          ApiCode = 40009
	ApiCodeRecordDoing                        ApiCode = 40010
	ApiCodeUnableInit                         ApiCode = 40011
	ApiCodeNotHaveManagementPermission        ApiCode = 40012
	ApiCodeSmtDiff                            ApiCode = 40013
	ApiCodeSuspendOperation                   ApiCode = 40014
	ApiCodeTaskNotExist                       ApiCode = 40015
	ApiCodeSameCustomScript                   ApiCode = 40016
	ApiCodeNotExistCustomScriptConfigPrice    ApiCode = 40017
	ApiCodeCustomScriptSet                    ApiCode = 40018
	ApiCodeProfitNotEnough                    ApiCode = 40019
	ApiCodeNoSupportPaymentToken              ApiCode = 40020
	ApiCodeSubAccOrderNotExist                ApiCode = 40021 //remove
	ApiCodeRuleDataErr                        ApiCode = 40022
	ApiCodeParentAccountNotExist              ApiCode = 40023
	ApiCodeSubAccountMinting                  ApiCode = 40024
	ApiCodeSubAccountMinted                   ApiCode = 40025
	ApiCodeBeyondMaxYears                     ApiCode = 40026
	ApiCodeHitBlacklist                       ApiCode = 40027
	ApiCodeNoTSetRules                        ApiCode = 40028
	ApiCodeTokenIdNotSupported                ApiCode = 40029
	ApiCodeNoSubAccountDistributionPermission ApiCode = 40030
	ApiCodeSubAccountNoEnable                 ApiCode = 40031
	ApiCodeAutoDistributionClosed             ApiCode = 40032
	ApiCodeAccountCanNotBeEmpty               ApiCode = 40033
	ApiCodePriceRulePriceNotBeLessThanMin     ApiCode = 40034
	ApiCodePriceMostReserveTwoDecimal         ApiCode = 40035
	ApiCodeConfigSubAccountPending            ApiCode = 40036
	ApiCodeAccountRepeat                      ApiCode = 40037
	ApiCodeInListMostBeLessThan1000           ApiCode = 40038
	ApiCodePreservedRulesMostBeOne            ApiCode = 40039
	ApiCodeRuleSizeExceedsLimit               ApiCode = 40040
	ApiCodeRuleFormatErr                      ApiCode = 40041
	ApiCodeExceededMaxLength                  ApiCode = 40042
	ApiCodeInvalidCharset                     ApiCode = 40043
	ApiCodeAccountNameErr                     ApiCode = 40044
	ApiCodeAccountLengthMostBeLessThan42      ApiCode = 40045
	ApiCodeAccountCharsetNotSupport           ApiCode = 40046
	ApiCodeAccountExpiringSoon                ApiCode = 40047
	ApiCodeUSDPricingTooLow                   ApiCode = 40048
	ApiCodeUSDPricingBelowMin                 ApiCode = 40049
	ApiCodeAccountRenewNoSupportCustomScript  ApiCode = 40050
	ApiCodeSubAccountRenewing                 ApiCode = 40051
	ApiCodeApprovalAlreadyExist               ApiCode = 40052
	ApiCodeAccountApprovalNotExist            ApiCode = 40053
	ApiCodeAccountApprovalProtected           ApiCode = 40054
	ApiCodeCouponCidNotExist                  ApiCode = 40055
	ApiCodeCouponPaid                         ApiCode = 40056
	ApiCodeCouponUnpaid                       ApiCode = 40057
	ApiCodeUnauthorized                       ApiCode = 40058
	ApiCodeCouponOpenTimeNotArrived           ApiCode = 40059
	ApiCodeCouponExpired                      ApiCode = 40060
	ApiCodeCouponErrAccount                   ApiCode = 40061
	ApiCodeOrderClosed                        ApiCode = 40062
)

// multi_device
const (
	ApiCodeHasNoAccessToCreate  ApiCode = 60000
	ApiCodeCreateConfigCellFail ApiCode = 60001
	ApiCodeHasNoAccessToRemove  ApiCode = 60002
)

// unipay - 600XXX
const (
	ApiCodeUnipayOrderNotExist  ApiCode = 600000 //remove
	ApiCodeOrderUnPaid          ApiCode = 600001
	ApiCodePaymentNotExist      ApiCode = 600002
	ApiCodeAmountIsTooLow       ApiCode = 600003
	ApiCodePaymentMethodDisable ApiCode = 600004
)

// remote_sign - 601XXX
const (
	ApiCodeServiceNotActivated    ApiCode = 601000
	ApiCodeAddressStatusNotNormal ApiCode = 601001
	ApiCodeUnsupportedAddrChain   ApiCode = 601002
	ApiCodeUnsupportedSignType    ApiCode = 601003
	ApiCodeIpBlockingAccess       ApiCode = 601004
	ApiCodeKeyDiff                ApiCode = 601005
	ApiCodeWalletAddrNotExist     ApiCode = 601006
)

//padge
const (
	ApiCodeUserNotExist               ApiCode = 70001
	ApiCodeGroupNotExist              ApiCode = 70002
	ApiCodePadgeNotExist              ApiCode = 70003
	ApiCodeReceiveNotExist            ApiCode = 70004
	ApiCodeInsufficientIssuance       ApiCode = 70005
	ApiCodeDistributeNotExist         ApiCode = 70006
	ApiCodeDeviceNotExist             ApiCode = 70007
	ApiCodeAlreadyBoundUser           ApiCode = 70008
	ApiCodeNotTheManagerOfDevice      ApiCode = 70009
	ApiCodeDeviceAlreadyUnboundUser   ApiCode = 70010
	ApiCodeAlreadyMinted              ApiCode = 70011
	ApiCodeInsufficientNumOfAI        ApiCode = 70012
	ApiCodeFailedToVerifySignature    ApiCode = 70013
	ApiCodeIssuerAlreadyExist         ApiCode = 70014
	ApiCodeNotTheManagerOfDid         ApiCode = 70015
	ApiCodeIssuerNotExist             ApiCode = 70016
	ApiCodeDidNotExist                ApiCode = 70017
	ApiCodeNotTheManagerOfPadge       ApiCode = 70017
	ApiCodePadgeAlreadyBoundIssuer    ApiCode = 70018
	ApiCodeDeviceAlreadyUnboundPadge  ApiCode = 70019
	ApiCodeDeviceAlreadyBoundPadge    ApiCode = 70020
	ApiCodeIssuanceHasBeenSet         ApiCode = 70021
	ApiCodeDistributeTypeAlreadyExist ApiCode = 70022
	ApiCodeAIFailedToDeduceAINum      ApiCode = 70023
	ApiCodeNumRemainingZero           ApiCode = 70024
	ApiCodeDistributionClosed         ApiCode = 70025
	ApiCodeInsufficientCredit         ApiCode = 70026
	ApiCodeAlreadyReceived            ApiCode = 70027
	ApiCodeNotInReceiveTime           ApiCode = 70028
)
