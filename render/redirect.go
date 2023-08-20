package render

import (
	"errors"
	"fmt"
	"github.com/shura1014/wits/response"
	"net/http"
)

type RedirectRender struct {
	Url     string
	Request *http.Request
	AbstractRender
}

func (render *RedirectRender) Render(w response.Response, status int) error {
	//render.contentTypeAndStatus(w, "", status)
	if status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect && status != http.StatusCreated {
		return errors.New(fmt.Sprintf("cannot redirect with status code %d", status))
	}
	http.Redirect(w, render.Request, render.Url, status)
	return nil
}
