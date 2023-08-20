package bind

import (
	"encoding/xml"
	"net/http"
)

type XmlBind struct {
}

func (bind XmlBind) Name() string {
	return "xml"
}

func (bind XmlBind) Bind(req *http.Request, data any) error {
	decoder := xml.NewDecoder(req.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return validate(data)
}
