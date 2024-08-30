package render

import (
	"github.com/shura1014/wits/response"
)

const (
	ContentType = "Content-Type"
	JSON        = "application/json; charset=utf-8"
	XML         = "application/xml; charset=utf-8"
	TEXT        = "text/plain; charset=utf-8"
	HTML        = "text/html; charset=utf-8"
	JAVASCRIPT  = "application/javascript; charset=utf-8"
)

type Render interface {

	// Render 响应逻辑
	//	w 响应流
	//	status 需要响应的状态
	Render(w response.Response, status int) error

	contentTypeAndStatus(w response.Response, contentType string, status int)
}

type AbstractRender struct {
}

func (render *AbstractRender) contentTypeAndStatus(w response.Response, contentType string, status int) {
	if contentType != "" {
		w.SetHeader(ContentType, contentType)
	}
	if status != 0 {
		w.WriteStatus(status)
	}
}
