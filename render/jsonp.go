package render

import (
	"encoding/json"
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/wits/response"
	"html/template"
)

type JsonpRender struct {
	// 回调
	Callback string
	Data     any
	AbstractRender
}

func (r JsonpRender) Render(w response.Response, status int) error {
	r.contentTypeAndStatus(w, JAVASCRIPT, status)
	ret, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}

	if r.Callback == "" {
		_, err = w.Write(ret)
		return err
	}

	callback := template.JSEscapeString(r.Callback)
	if _, err = w.Write(stringutil.StringToBytes(callback)); err != nil {
		return err
	}

	if _, err = w.Write(stringutil.StringToBytes("(")); err != nil {
		return err
	}

	if _, err = w.Write(ret); err != nil {
		return err
	}

	if _, err = w.Write(stringutil.StringToBytes(");")); err != nil {
		return err
	}

	return nil
}
