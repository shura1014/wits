package render

import (
	"encoding/json"
	"github.com/shura1014/wits/response"
)

type JsonRender struct {
	Data any
	// 带有html格式的json是否不被编码
	Pure bool
	// 是否展开json
	Expand bool
	AbstractRender
}

func (render JsonRender) Render(w response.Response, status int) error {
	render.contentTypeAndStatus(w, JSON, status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(!render.Pure)
	if render.Expand {
		encoder.SetIndent("", DefaultIndent)
	}
	err := encoder.Encode(render.Data)
	//jsonData, err := json.Marshal(render.Data)
	//if err != nil {
	//	return err
	//}
	//_, err = w.Write(jsonData)
	return err
}
