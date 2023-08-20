package bind

import (
	"net/http"
)

type Bind interface {
	Name() string
	Bind(*http.Request, any) error
}

var JSON = JsonBind{}
var XML = XmlBind{}

func validate(obj any) error {
	if Validator == nil {
		return nil
	}
	return Validator.Validate(obj)
}

func DisableBindValidation() {
	Validator = nil
}
