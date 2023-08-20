package render

import (
	"fmt"
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/wits/response"
)

type StringRender struct {
	// 格式
	Format string

	Data []any

	AbstractRender
}

func (render *StringRender) Render(w response.Response, status int) error {
	render.contentTypeAndStatus(w, TEXT, status)
	if len(render.Data) > 0 {
		_, err := fmt.Fprintf(w, render.Format, render.Data...)
		return err
	}
	// 使用底层方法 地址强转
	// 直接强转会发生内存拷贝 (*[]byte) render.Format
	data := stringutil.StringToBytes(render.Format)
	_, err := w.Write(data)
	return err
}
