package wits

import "github.com/shura1014/wits/bind"

/******************************参数绑定start*******************************************/

// BindJSON
// curl -X POST -d '{"name":"wendell","age":28}'  http://127.0.0.1:8888/user/bind/json
func (c *Context) BindJSON(obj any) error {
	return c.MustBindWith(obj, bind.JSON)
}

func (c *Context) EnableDecoderUseNumber(fun func()) {
	number := bind.UseNumber
	bind.EnableDecoderUseNumber()
	defer func() {
		bind.UseNumber = number
	}()
	fun()
}

func (c *Context) EnableStrictMatching(fun func()) {
	s := bind.StrictMatching
	bind.EnableStrictMatching()
	defer func() {
		bind.StrictMatching = s
	}()
	fun()
}

func (c *Context) DisableDecoderUseNumber(fun func()) {
	number := bind.UseNumber
	bind.DisableDecoderUseNumber()
	defer func() {
		bind.UseNumber = number
	}()
	fun()
}

func (c *Context) DisableStrictMatching(fun func()) {
	s := bind.StrictMatching
	bind.DisableStrictMatching()
	defer func() {
		bind.StrictMatching = s
	}()
	fun()
}

// BindXml
/*
@Example
type User struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
}

curl -X POST -d '<?xml version="1.0" encoding="UTF-8"?><root><age>25</age><name>juan</name></root>' -H 'Content-Type: application/xml'  http://127.0.0.1:8888/user/bind/xml
*/
func (c *Context) BindXml(obj any) error {
	return c.MustBindWith(obj, bind.XML)
}

func (c *Context) MustBindWith(obj any, b bind.Bind) error {
	if err := c.ShouldBindWith(obj, b); err != nil {
		//c.W.WriteStatus(http.StatusBadRequest)
		return err
	}
	return nil
}

func (c *Context) ShouldBindWith(obj any, b bind.Bind) error {
	return b.Bind(c.R, obj)
}

/******************************参数绑定end*********************************************/
