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
	ApiCodeSuccess        ApiCode = 0
	ApiCodeError500       ApiCode = 500
	ApiCodeParamsInvalid  ApiCode = 10000
	ApiCodeMethodNotExist ApiCode = 10001
	ApiCodeDbError        ApiCode = 10002
	ApiCodeCacheError     ApiCode = 10003
)

// unipay
const (
	ApiCodeOrderNotExist   ApiCode = 600000
	ApiCodeOrderUnPaid     ApiCode = 600001
	ApiCodePaymentNotExist ApiCode = 600002
	ApiCodeAmountIsTooLow  ApiCode = 600003
)
