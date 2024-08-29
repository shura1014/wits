package wits

import (
	"net/url"
	"strconv"
)

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.R != nil {
			c.queryCache = c.R.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

// GetQueryArray
// url.Values 是一个 map[string][]string
// 所以同一个key可以有多个值
// 提供一个获取数组的方法
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

// QueryArray 业务不需要ok bool参数
func (c *Context) QueryArray(key string) (values []string) {
	values, _ = c.GetQueryArray(key)
	return
}

// GetQuery 一般来说一个key对应一个参数，是大部分业务所需求的
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

// GetQueryInt 获取int类型的参数
func (c *Context) GetQueryInt(key string) (value int) {

	if val, ok := c.GetQuery(key); ok {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal
		}
	}
	return 0
}

// GetQueryInt64 获取int64类型参数
func (c *Context) GetQueryInt64(key string) (value int64) {

	if val, ok := c.GetQuery(key); ok {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return intVal
		}
	}
	return 0
}

// GetQueryBool 获取bool类型参数
func (c *Context) GetQueryBool(key string) (value bool) {

	if val, ok := c.GetQuery(key); ok {
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			return boolVal
		}
	}
	return false
}

// GetQueryFloat64 获取float类型参数
func (c *Context) GetQueryFloat64(key string) (value float64) {

	if val, ok := c.GetQuery(key); ok {
		floatVal, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return floatVal
		}
	}
	return 0
}

// Query
// 绝大部分业务不想处理bool，希望框架返回一个值就行啦
// curl  http://127.0.0.1:8888/user/add\?id=1001\&person%5Bname%5D=wendell\&person%5Bage%5D=21
func (c *Context) Query(key string) (value string) {
	value, _ = c.GetQuery(key)
	return
}

// QueryOrDefault 带有默认值的返回，不需要业务再去判空
func (c *Context) QueryOrDefault(key, defaultValue string) (value string) {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// QueryMap
/*
请求 curl 127.0.0.1:8888/user/person[name]=wendell&person[age]=21
需要返回
{
	"name":"wendell"
	"age": 21
}
*/
func (c *Context) QueryMap(key string) (dict map[string]string) {
	dict, _ = c.GetQueryMap(key)
	return
}

func (c *Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCache()
	return c.get(c.queryCache, key)
}
