package wits

func (c *Context) initBodyCache() {
	if c.bodyMap == nil {
		c.bodyMap = make(map[string]any)
		if c.R.ContentLength == 0 {
			return
		}
		if c.ContentType() == "application/json" {
			err := c.BindJSON(&c.bodyMap)
			if err != nil {
				c.Panic("Parse the body in json format failed")
			}
			return
		}
		if c.ContentType() == "application/xml" {
			err := c.BindXml(&c.bodyMap)
			if err != nil {
				c.Panic("Parse the body in xml format failed")
			}
			return
		}
	}
}

func (c *Context) BodyMap() map[string]any {
	c.initBodyCache()
	return c.bodyMap
}

func (c *Context) GetBody(key string) (v any, ok bool) {
	c.initBodyCache()
	v, ok = c.bodyMap[key]
	return
}

func (c *Context) Body(key string) (v any) {
	c.initBodyCache()
	v = c.bodyMap[key]
	return
}

// GetBodyString 获取string类型的参数
func (c *Context) GetBodyString(key string) (v string) {

	if body, ok := c.GetBody(key); ok {
		v = body.(string)
	}
	return
}

// GetBodyInt 获取int类型参数
func (c *Context) GetBodyInt(key string) (v int) {
	if body, ok := c.GetBody(key); ok {
		v = body.(int)
	}
	return
}

// GetBodyInt64 获取int64类型参数
func (c *Context) GetBodyInt64(key string) (v int64) {
	if body, ok := c.GetBody(key); ok {
		v = body.(int64)
	}
	return
}

// GetBodyBool 获取bool类型参数
func (c *Context) GetBodyBool(key string) (v bool) {
	if body, ok := c.GetBody(key); ok {
		v = body.(bool)
	}
	return
}

// GetBodyLong 获取long类型参数
func (c *Context) GetBodyLong(key string) (v int64) {
	if body, ok := c.GetBody(key); ok {
		v = body.(int64)
	}
	return
}

// GetBodyFloat 获取float类型参数
func (c *Context) GetBodyFloat(key string) (v float64) {
	if body, ok := c.GetBody(key); ok {
		v = body.(float64)
	}
	return
}
