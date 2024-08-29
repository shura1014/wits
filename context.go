package wits

import (
	"context"
	"github.com/shura1014/common/container/tree"
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/logger"
	"github.com/shura1014/wits/bind"
	"github.com/shura1014/wits/render"
	"github.com/shura1014/wits/response"
	"github.com/shura1014/wits/session"
	"io"
	"mime/multipart"
	"net"
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

func (c *Context) Redirect(status int, url string) error {
	return c.Render(status, &render.RedirectRender{Url: url, Request: c.R})
}

func (c *Context) JSON(status int, data any) error {
	return c.Render(status, &render.JsonRender{Data: data})
}

// PureJSON 带有html格式的json不被编码
func (c *Context) PureJSON(status int, data any) error {
	return c.Render(status, &render.JsonRender{Data: data, Pure: true})
}

// ExpandJSON 是否展开json
func (c *Context) ExpandJSON(status int, data any) error {
	return c.Render(status, &render.JsonRender{Data: data, Expand: true})
}

func (c *Context) ExpandAndPureJSON(status int, data any) error {
	return c.Render(status, &render.JsonRender{Data: data, Pure: true, Expand: true})
}

func (c *Context) XML(status int, data any) error {
	return c.Render(status, &render.XmlRender{Data: data})
}

func (c *Context) ExpandXML(status int, data any) error {
	return c.Render(status, &render.XmlRender{Data: data, Expand: true})
}

func (c *Context) String(status int, format string, values ...any) error {
	return c.Render(status, &render.StringRender{Format: format, Data: values})
}

func (c *Context) HTMLTemplate(code int, name string, obj any) error {
	instance := c.e.render.Instance(name, obj)
	return c.Render(code, instance)
}

func (c *Context) ByteToString(status int, data []byte) error {
	return c.Render(status, &render.StringRender{Format: string(data)})
}

func (c *Context) Back(code int, msg string) error {
	return c.String(code, msg)

}

func (c *Context) Success(msg string) error {
	return c.Back(http.StatusOK, msg)
}

func (c *Context) ReturnNow(status int) {
	c.W.WriteStatus(status)
	c.W.Flush()
}

func (c *Context) Fail(msg string) error {
	return c.Back(http.StatusInternalServerError, msg)
}

func (c *Context) Render(status int, r render.Render) error {
	err := r.Render(c.W, status)
	return err
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

// RequestHeader
/**************头部操作****************************/
func (c *Context) RequestHeader(key string) string {
	return c.R.Header.Get(key)
}

func (c *Context) ContentType() string {
	return c.R.Header.Get(HeaderContentType)
}

func (c *Context) GetMethod() string {
	return c.R.Method
}

func (c *Context) GetPath() string {
	return c.R.URL.Path
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.W.Header().Del(key)
		return
	}
	c.W.Header().Set(key, value)
}

func (c *Context) UserAgent() string {
	return c.RequestHeader(HeaderUserAgent)
}

func (c *Context) Referer() string {
	return c.RequestHeader(HeaderReferer)
}

// RemoteIP 获取客户端ip
func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.R.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

func (c *Context) ClientIP() string {
	var clientIP string

	// 检查该ip是不是本地ip 或者说它是一个代理ip，因此需要找出原ip
	proxyIps := c.RequestHeader(XForwardedFor)
	if proxyIps != "" {
		// XForwardedFor可能经过多重代理，取第一个就行
		i := strings.IndexAny(proxyIps, ",")
		if i > 0 {
			clientIP = strings.TrimSpace(proxyIps[:i])
		}
		clientIP = strings.TrimPrefix(clientIP, "[")
		clientIP = strings.TrimSuffix(clientIP, "]")
		return clientIP
	}

	clientIP = c.RequestHeader(XRealIP)
	if clientIP == "" {
		clientIP = c.RemoteIP()
	}

	clientIP = strings.TrimPrefix(clientIP, "[")
	clientIP = strings.TrimSuffix(clientIP, "]")
	return clientIP
}

/**************头部操作****************************/

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

/******************************参数绑定start*******************************************/

// BindJSON
// curl -X POST -d '{"name":"wendell","age":28}'  http://127.0.0.1:8888/user/bind/json
func (c *Context) BindJSON(obj any) error {
	return c.MustBindWith(obj, bind.JSON)
}

func (c *Context) EnableDecoderUseNumber(fun func()) {
	number := bind.UseNumber
	bind.EnableDecoderUseNumber()
	defer func() {
		bind.UseNumber = number
	}()
	fun()
}

func (c *Context) EnableStrictMatching(fun func()) {
	s := bind.StrictMatching
	bind.EnableStrictMatching()
	defer func() {
		bind.StrictMatching = s
	}()
	fun()
}

func (c *Context) DisableDecoderUseNumber(fun func()) {
	number := bind.UseNumber
	bind.DisableDecoderUseNumber()
	defer func() {
		bind.UseNumber = number
	}()
	fun()
}

func (c *Context) DisableStrictMatching(fun func()) {
	s := bind.StrictMatching
	bind.DisableStrictMatching()
	defer func() {
		bind.StrictMatching = s
	}()
	fun()
}

// BindXml
/*
@Example
type User struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
}

curl -X POST -d '<?xml version="1.0" encoding="UTF-8"?><root><age>25</age><name>juan</name></root>' -H 'Content-Type: application/xml'  http://127.0.0.1:8888/user/bind/xml
*/
func (c *Context) BindXml(obj any) error {
	return c.MustBindWith(obj, bind.XML)
}

func (c *Context) MustBindWith(obj any, b bind.Bind) error {
	if err := c.ShouldBindWith(obj, b); err != nil {
		//c.W.WriteStatus(http.StatusBadRequest)
		return err
	}
	return nil
}

func (c *Context) ShouldBindWith(obj any, b bind.Bind) error {
	return b.Bind(c.R, obj)
}

/******************************参数绑定end*********************************************/

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
	_ = c.Fail("业务执行异常")
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
