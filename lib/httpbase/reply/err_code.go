package reply

import "fmt"

type ErrCode int

//const (
//	ErrCodeOk              = 0
//	ErrCodeParamError      = 400
//	ErrCodePermissionError = 403
//	ErrCodeNotFound        = 404
//	ErrCodeServiceError    = 500
//)

const (
	ErrcodeOk           ErrCode = 1
	ErrcodeServiceError ErrCode = -100
	ErrcodeRpcError     ErrCode = -101
	ErrcodeParamError   ErrCode = -400
	ErrcodeMongoError   ErrCode = -401

	ErrCodeInvalidAuthToken = -1000
	ErrCodeExpiredAuthToken = -1001
)

var serviceErrCodeName = map[ErrCode]string{
	ErrcodeOk:           "ERRCODE_OK",
	ErrcodeParamError:   "请求参数错误",
	ErrcodeServiceError: "服务内部错误，请稍后重试！",
	ErrcodeRpcError:     "内部服务调用错误",
	ErrcodeMongoError:   "操作mongo错误",

	ErrCodeInvalidAuthToken: "无效的token",
	ErrCodeExpiredAuthToken: "token过期了",
}

func (e ErrCode) ServiceErrCodeName(args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf("%v", args)
	}
	if s, ok := serviceErrCodeName[e]; ok {
		return s
	}
	return fmt.Sprintf("unknown error code: %d", int(e))
}
