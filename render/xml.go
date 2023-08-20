package render

import (
	"encoding/xml"
	"github.com/shura1014/wits/response"
)

var DefaultIndent = "    "

type XmlRender struct {
	Data   any
	Expand bool
	AbstractRender
}

func (render XmlRender) Render(w response.Response, status int) error {
	render.contentTypeAndStatus(w, XML, status)
	encoder := xml.NewEncoder(w)
	if render.Expand {
		encoder.Indent("", DefaultIndent)
	}
	err := encoder.Encode(render.Data)
	return err
}
