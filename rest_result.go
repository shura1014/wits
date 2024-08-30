package wits

import (
	"github.com/shura1014/common/goerr"
)

var (
	// -1表示即使使用了RestResult 也不处理
	// 比如文件系统，或者说结果本省就是RestResult
	restErrorCode   = -1
	restSuccessCode = 1000
	restSuccessMsg  = "success"
)

// RestResult 统一返回结果
// {"code":200,"msg":"ok","data":""}
// 一般来说code 200 表示正常，其他表示异常
type RestResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func RestSuccess(date any) *RestResult {
	r := &RestResult{}
	r.Code = 200
	r.Msg = restSuccessMsg
	r.Data = date
	return r
}

func RestFailed(msg string) *RestResult {
	r := &RestResult{}
	r.Code = -1
	r.Msg = msg
	r.Data = nil
	return r
}

var commonHandler = &CommonHandler{}

type CommonHandler struct {
}

func (handler *CommonHandler) HandlerError(ctx *Context, err *goerr.BizError) {
	ctx.ERROR(err.DetailMsg())
	r := &RestResult{
		Code: err.Code.Code(),
		Msg:  err.Code.Message(),
		Data: err.Data,
	}
	_ = ctx.JSON(restErrorCode, r)
	//ctx.W.Flush()
}

func RestResultEnable() bool {
	return GetBool(AppRestResultEnable)
}
