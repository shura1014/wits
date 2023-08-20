package wits

import (
	"github.com/shura1014/common/goerr"
)

type HandlerBizError interface {
	HandlerError(ctx *Context, err *goerr.BizError)
}

type NothingHandler struct {
}

func (biz NothingHandler) HandlerError(ctx *Context, err *goerr.BizError) {
}

var (
	sessionTimeout    = &goerr.ErrorCode{ErrCode: 10000, ErrMsg: "session time out"}
	sessionIdNotFound = &goerr.ErrorCode{ErrCode: 10001, ErrMsg: "not found session id"}
)
