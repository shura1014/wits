package wits

import (
	"bytes"
	"fmt"
	"github.com/shura1014/wits/bind"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXmlBind(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := CreateTestContext(w)

	ctx.R, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`<?xml version="1.0" encoding="UTF-8"?>
		<root>
			<foo>FOO</foo>
		   	<bar>BAR</bar>
		</root>`))

	var obj struct {
		Foo string `xml:"foo"`
		Bar string `xml:"bar"`
		T   string `xml:"t" bind:"required"`
	}

	ctx.R.Header.Add("Content-Type", "application/xml")
	err := ctx.BindXml(&obj)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(obj)
}

func TestJsonBind(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := CreateTestContext(w)

	ctx.R, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"name\":\"wendell\"}"))
	//ctx.R, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"name\":\"wendell\",\"age\":17}"))
	//ctx.R, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"name\":\"wendell\",\"age\":101}"))

	var User struct {
		Name string `json:"name"`
		Age  int    `json:"age" bind:"required,max=100,min=18"`
	}

	ctx.R.Header.Add("Content-Type", "application/json")
	err := ctx.BindJSON(&User)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(User)
}

func TestEnableDecoderUseNumber(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := CreateTestContext(w)

	ctx.R, _ = http.NewRequest("POST", "/", bytes.NewBufferString("{\"name\":\"wendell\",\"age\":17,\"weight\":\"140\"}"))

	var User struct {
		Name string
		Age  int
	}

	ctx.R.Header.Add("Content-Type", "application/json")

	ctx.EnableStrictMatching(func() {
		fmt.Println(bind.StrictMatching)
		err := ctx.BindJSON(&User)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println(User)
	})

	fmt.Println(bind.StrictMatching)

}
