package wits

import (
	"github.com/shura1014/wits/render"
	"net/http"
)

// 渲染

// Redirect 重定向
// url 跳转地址
func (c *Context) Redirect(status int, url string) error {
	return c.Render(status, &render.RedirectRender{Url: url, Request: c.R})
}

// JSON json格式渲染
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

// ReturnNow 立刻返回，直接执行Flush
func (c *Context) ReturnNow(status int) {
	c.W.WriteStatus(status)
	c.W.Flush()
}

func (c *Context) Render(status int, r render.Render) error {
	err := r.Render(c.W, status)
	return err
}

func (c *Context) Fail(msg string, data any) {
	if RestResultEnable() {
		// 如果是RestResult模式，那么状态码应该是ok，根据code去判断
		_ = c.Back(http.StatusOK, restErrorCode, msg, data)
		return
	}
	_ = c.JSON(http.StatusInternalServerError, data)
}

func (c *Context) Success(msg string, data any) {
	if msg == "" {
		msg = restSuccessMsg
	}
	if RestResultEnable() {
		err := c.Back(http.StatusOK, restSuccessCode, msg, data)
		if err != nil {
			_ = c.Back(http.StatusOK, restErrorCode, msg, data)
		}
		return
	}

	err := c.JSON(http.StatusOK, data)
	if err != nil {
		_ = c.JSON(http.StatusInternalServerError, data)
	}
}

// Back Rest返回
func (c *Context) Back(httpStatus int, code int, msg string, data any) error {
	result := &RestResult{Code: code, Msg: msg, Data: data}
	return c.JSON(httpStatus, result)

}
