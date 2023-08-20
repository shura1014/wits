package render

import (
	"github.com/shura1014/wits/response"
	"html/template"
)

// Delims 自定义解析格式
type Delims struct {
	Left  string
	Right string
}

type HtmlTemplateRender struct {
	Template *template.Template
	Name     string
	Data     any
	AbstractRender
}

func (render HtmlTemplateRender) Render(w response.Response, status int) error {
	render.contentTypeAndStatus(w, HTML, status)

	if render.Name == "" {
		return render.Template.Execute(w, render.Data)
	}
	return render.Template.ExecuteTemplate(w, render.Name, render.Data)
}

// Instance /*获取render对象
func (render HtmlTemplateRender) Instance(name string, data any) Render {
	return &HtmlTemplateRender{
		Template: render.Template,
		Name:     name,
		Data:     data,
	}
}
