package wits

import (
	"context"
	"github.com/shura1014/common/container/tree"
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/logger"
	"github.com/shura1014/wits/render"
	"github.com/shura1014/wits/response"
	"github.com/shura1014/wits/session"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Context struct {
	W response.Response

	R *http.Request
	e *Engine

	// 保存参数 c.Request.URL.Query
	queryCache url.Values

	// 表单内容缓存
	formCache url.Values

	bodyMap map[string]any

	// 上下文参数（可以业务自定义），以及路径参数 /user/:id
	keys tree.Keys
	mu   sync.RWMutex

	logger *logger.Logger

	handlerBizError HandlerBizError

	// 是否允许跨站请求携带cookie Lax Strict None
	sameSite http.SameSite

	session *session.Session
}

// Copy 拷贝上下文为一个新对象
// 由于context是一个资源复用的对象
// 如果业务使用携程需要拷贝一份copy
func (c *Context) Copy() *Context {
	cp := Context{
		R: c.R,
		W: c.W,
		//params:          c.params,
		e:               c.e,
		logger:          c.logger,
		handlerBizError: c.handlerBizError,
	}
	cp.keys = tree.Keys{}
	for k, v := range c.keys {
		cp.keys[k] = v
	}
	//paramCopy := make(Params, len(*cp.params))
	//copy(paramCopy, *cp.params)
	//cp.params = &paramCopy
	return &cp
}

func (c *Context) Ctx() context.Context {
	return c.R.Context()
}

func (c *Context) GetUser() string {
	return c.GetString(UserKey)
}

// File 文件操作
func (c *Context) File(filepath string) {
	http.ServeFile(c.W, c.R, filepath)
}

// c.FileFromFS("./gin.go", Dir(".", false))

func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.SetHeader(HeaderContentDisposition, `attachment; filename="`+escapeQuotes(filename)+`"`)
	} else {
		c.W.SetHeader(HeaderContentDisposition, `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

/*******************************参数处理start********************************************/

func (c *Context) get(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts, exist
}

// SaveUploadedFile 拿到文件后可以上传到指定的目录
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

/*******************************参数处理end**********************************************/

func (c *Context) JSONP(code int, obj any) {
	callback := c.QueryOrDefault("callback", "")
	if callback == "" {
		c.Render(code, &render.JsonRender{Data: obj})
		return
	}
	c.Render(code, &render.JsonpRender{Callback: callback, Data: obj})
}

// DEBUG INFO ERROR 日志打印
/***************************提供框架使用框架日志打印方法start******************************/
func (c *Context) DEBUG(msg any, data ...any) {
	c.logger.DebugSkip(c.Ctx(), msg, 1, data...)

}

func (c *Context) INFO(msg any, data ...any) {
	c.logger.InfoSkip(c.Ctx(), msg, 1, data...)

}

func (c *Context) WARN(msg any, data ...any) {
	c.logger.WarnSkip(c.Ctx(), msg, 1, data...)

}

func (c *Context) ERROR(msg any, data ...any) {
	c.logger.ErrorSkip(c.Ctx(), msg, 1, data...)
}

/***************************提供框架使用框架日志打印方法end********************************/

// HandlerError SetHandlerBizError GetHandlerBizError Panic
/******************************业务异常处理start****************************************/
func (c *Context) HandlerError(err *goerr.BizError) {
	bizHandler := c.GetHandlerBizError()
	if bizHandler != nil {
		bizHandler.HandlerError(c, err)
		return
	}
	c.ERROR(err.DetailMsg())
	c.Fail("业务执行异常", nil)
}

func (c *Context) SetHandlerBizError(handler HandlerBizError) {
	c.handlerBizError = handler
}

func (c *Context) GetHandlerBizError() HandlerBizError {
	return c.handlerBizError
}

func (c *Context) Panic(msg string) {
	panic(goerr.TextSkip(1, msg))
}

/******************************业务异常处理end****************************************/

// SetSameSite with cookie
// SameSiteDefaultMode SameSite = iota + 1
//	SameSiteLaxMode
//	SameSiteStrictMode
//	SameSiteNoneMode
/******************************cookie start****************************************/
func (c *Context) SetSameSite(samesite http.SameSite) {
	c.sameSite = samesite
}

// SetCookie 添加cookie到响应头
func (c *Context) SetCookie(name, value string, options ...*CookieOptions) {
	SetCookie(name, value, c, options...)
}

// Cookie 查找cookie
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.R.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

/******************************cookie end******************************************/

/*****************************session start****************************************/

// NewSession 创建一个session
// 一般只给登录请求使用，并且需要开启session支持 app.session.enable
func (c *Context) NewSession() *session.Session {
	if enableSession {
		return sessionManager.NewSession()
	}
	return nil
}

// GetSessionId 需要开启session支持
func (c *Context) GetSessionId() string {
	if enableSession && c.session != nil {
		return c.session.GetId()
	}
	return ""
}

func (c *Context) Session() *session.Session {
	return c.session
}

/*****************************session end******************************************/

// Reset 由于使用了对象池资源复用，复用的时候需要清理上一次缓存的数据
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.R = r
	c.W.Reset(w)
	c.queryCache = nil
	c.formCache = nil
	c.bodyMap = nil
	c.handlerBizError = nil
	c.sameSite = http.SameSiteDefaultMode
	c.keys = make(tree.Keys)
	c.session = nil
	c.handlerBizError = c.e.handlerBizError

	// todo Max body size limit

}

func (c *Context) Shutdown() {
	c.e.Shutdown()
}
