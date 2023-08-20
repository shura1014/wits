package bind

import (
	"encoding/json"
	"errors"
	"net/http"
)

// UseNumber
// 对于需要转换为 map[string]any
// 解码器会将数字反编为float64
// UseNumber 后将会是一个json.Number对象
//
// 如果输入的数字比较大，这个表示会有损精度。所以可以 UseNumber() 启用 json.Number 来用字符串表示数字
var UseNumber = false

// StrictMatching
// 严格匹配每一个字段，如果json数据里有绑定对象里没有的字段会报错，表示严格匹配
// 如果存在报错 json: unknown field xxx
// 打开json的DisallowUnknownFields
var StrictMatching = false

func EnableDecoderUseNumber() {
	UseNumber = true
}

func DisableDecoderUseNumber() {
	UseNumber = false
}

func EnableStrictMatching() {
	StrictMatching = true
}

func DisableStrictMatching() {
	StrictMatching = false
}

type JsonBind struct {
}

func (bind JsonBind) Name() string {
	return "json"
}

func (bind JsonBind) Bind(req *http.Request, data any) error {
	body := req.Body
	if req == nil || body == nil {
		return errors.New("invalid request")
	}
	decoder := json.NewDecoder(body)

	if UseNumber {
		decoder.UseNumber()
	}

	if StrictMatching {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(data); err != nil {
		return err
	}
	return validate(data)
}
