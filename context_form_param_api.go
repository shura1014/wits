package wits

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

/*
form表单的结构体

	type Form struct {
		Value map[string][]string
		File  map[string][]*FileHeader
	}

如果是Value作为字符串存储在内存中
文件，存储在内存或磁盘上
如果文件过大，那么就会写入磁盘并刷新缓冲区
这个阈值可以设置 MaxMultipartMemory
*/
func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.R
		if err := req.ParseMultipartForm(c.e.MaxMultipartMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				log.Panicf("error on parse multipart form array: %v", err)
			}
		}
		c.formCache = req.PostForm
	}
}

// PostForm
// application/x-www-form-urlencoded
// curl -X POST -d "sex=女&person[weight]=148&person[hight]=1.75"  http://127.0.0.1:8888/user/add
func (c *Context) PostForm(key string) (value string) {
	value, _ = c.GetPostForm(key)
	return
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return
}

func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

func (c *Context) PostFormMap(key string) (dicts map[string]string) {
	dicts, _ = c.GetPostFormMap(key)
	return
}

func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.initFormCache()
	return c.get(c.formCache, key)
}

// PostFormInt 获取int类型参数
func (c *Context) PostFormInt(key string) int {
	if val, ok := c.GetPostForm(key); ok {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal
		}
	}
	return 0
}

func (c *Context) PostFormInt64(key string) int64 {
	if val, ok := c.GetPostForm(key); ok {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return intVal
		}
	}
	return 0
}

func (c *Context) PostFormFloat64(key string) float64 {
	if val, ok := c.GetPostForm(key); ok {
		floatVal, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return floatVal
		}
	}
	return 0
}

func (c *Context) PostFormBool(key string) bool {
	if val, ok := c.GetPostForm(key); ok {
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			return boolVal
		}
	}
	return false
}

// FormFile
// 返回第一个文件，一般业务就使用这一个
// 如果有多个文件可以使用  MultipartForm
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.R.MultipartForm == nil {
		if err := c.R.ParseMultipartForm(c.e.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	file, header, err := c.R.FormFile(name)
	if err != nil {
		return nil, err
	}
	file.Close()
	return header, err
}

// MultipartForm
// 通过 Form.file 拿到多个文件，是一个map对象
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.R.ParseMultipartForm(c.e.MaxMultipartMemory)
	return c.R.MultipartForm, err
}
