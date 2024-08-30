package wits

import (
	"fmt"
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/wits/render"
	"github.com/shura1014/wits/response"
	"io"
	"net/http"
)

var (
	// 600表示即使使用了RestResult 也不处理
	// 比如文件系统，或者说结果本省就是RestResult
	restCode = 1000
	objSuf   = []byte{'}'}
	quot     = []byte{'"'}
	data     = []byte{'d', 'a', 't', 'a'}
)

// WrapPre {"code":200,"msg":"ok","data":""}
func WrapPre(response response.Response) {
	var (
		preString   = "{\"code\":%s,\"msg\":\"%s\",\"data\":"
		write       = response.Unwrap()
		state       = response.Status()
		contentType = response.ContentType()
		Code        = "200"
		Msg         = http.StatusText(state)
	)

	if render.JAVASCRIPT == contentType || render.HTML == contentType {
		return
	}

	if state == restCode {
		return
	}

	if state > 400 {
		Code = "10000"
		Msg = http.StatusText(state)
	}
	preString = fmt.Sprintf(preString, Code, Msg)
	_, _ = io.WriteString(write, preString)
	//_, _ = write.Write(objPre)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(code)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(colon)
	////--------code------------
	//_, _ = write.Write([]byte(Code))
	//
	//_, _ = write.Write(comma)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(msg)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(colon)
	//_, _ = write.Write(quot)
	////-----------msg------
	//_, _ = write.Write([]byte(Msg))
	//_, _ = write.Write(quot)
	//
	//_, _ = write.Write(comma)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(data)
	//_, _ = write.Write(quot)
	//_, _ = write.Write(colon)
	// ---------data-------
	switch contentType {
	case render.TEXT:
		_, _ = write.Write(quot)
	case render.JSON:

	default:
		_, _ = write.Write(quot)
	}
}

func WrapPost(response response.Response) {
	var (
		write       = response.Unwrap()
		contentType = response.ContentType()
		state       = response.Status()
	)

	if render.JAVASCRIPT == contentType || render.HTML == contentType {
		return
	}

	if state == restCode {
		return
	}
	switch contentType {
	case render.TEXT:
		_, _ = write.Write(quot)
	case render.JSON:
	default:
		_, _ = write.Write(quot)
	}
	_, _ = write.Write(objSuf)
}

// RestResult 统一返回结果
// {"code":200,"msg":"ok","data":""}
// 一般来说code 200 表示正常，其他表示异常
type RestResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
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
	_ = ctx.JSON(restCode, r)
	//ctx.W.Flush()
}

//func RestResultAdvice(next HandlerFunc) HandlerFunc {
//	return func(ctx *Context) {
//		ctx.SetHandlerBizError(commonHandler)
//		next(ctx)
//	}
//}

func RestResultEnable() bool {
	return GetBool(AppRestResultEnable)
}
