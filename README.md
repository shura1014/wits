# wits

一个简易的web框架

主要功能包括

路由匹配(前缀树)，上下文，拦截器，过滤器 认证（basic，jwt） 数据绑定，响应（json，jsonp,xml，string，重定向，html）支持https

支持中间件扩展 全局中间件以及路由级别中间件，提供 跨域中间件 、访问日志中间件、异常恢复中间件、请求耗时中间件

生命周期扩展机制

统一的异常处理

统一响应设置

cookie session的支持

IP黑白名单设置

路径参数支持 /api/getId/:id

优雅停机

请求超时中间件

配置文件 支持 yaml toml properties

环境变量，支持方便获取参数变量，操作系统变量，配置文件变量

banner打印

# 快速使用

```go

func main() {
    engine := wits.Default()
    group := engine.Group("api")
    group.GET("probe", func (ctx *wits.Context) {
        _ = ctx.Success(timeutil.Now())
    })
    engine.RunServer()
}


```

```text
        .__  __          
__  _  _|__|/  |_  ______
\ \/ \/ /  \   __\/  ___/
 \     /|  ||  |  \___ \ 
  \/\_/ |__||__| /____  >
                      \/ 1.0.0
Hello, Welcome To Use Wits.

[WITS] 2023-08-20 17:05:48 /Users/wendell/GolandProjects/shura/wits/utils.go:40 WARN The port is not set and no environment variables app.server.address or app.server.port are configured. 
[WITS] 2023-08-20 17:05:48 /Users/wendell/GolandProjects/shura/wits/utils.go:41 DEBUG The App will Running Listen Address is :8888 Get by default. 
```

# web容器engine常用api

```text
RegisterWarmupFunc      预热
RegisterPostNewFunc     New之后执行
RegisterPreRunFunc      Run之前执行
SetHandlerBizError      统一异常处理
RegisterOnShutdown      优雅停机钩子
ServerPostProcessor     http.server创建后执行
SetLogger               自定义日志
Use                     注册全局中间件
AppName                 获取app名字
AddFilter               全局过滤器
InitAclNodeList // engine.InitAclNodeList("./file/acl.conf") 黑白名单
Group // engine.Group("user") 添加一个路由组
LoadHTMLGlob            加载模版

RunServer               启动
RunTLS                  https启动
```

# 上下文Context 常用api

```text
Shutdown() 
GetSessionId()
Cookie(name string)
// 参数绑定
BindXml(obj any)
BindJson(obj any)
// 响应
JSONP(code int, obj any)
Redirect(status int, url string)
JSON(status int, data any)
PureJSON(status int, data any)
ExpandJSON(status int, data any)
ExpandAndPureJSON(status int, data any)
XML(status int, data any)
ExpandXML(status int, data any)
String(status int, format string, values ...any)
HTMLTemplate(code int, name string, obj any) error
Back(code int, msg string)
Success(msg string)
Fail(msg string)
File(filepath string)
FileFromFS(filepath string, fs http.FileSystem)
FileAttachment(filepath, filename string)

Query(key string)
QueryOrDefault(key, defaultValue string)
QueryMap(key string)
PostForm(key string)
PostFormMap(key string)
PostFormArray(key string)
FormFile(name string)
MultipartForm()
BodyMap()
Body(key string)


DEBUG(msg any, data ...any)
INFO(msg any, data ...any)
WARM(msg any, data ...any)
ERROR(msg any, data ...any)

UserAgent()
Referer()
RemoteIP()
ClientIP()
ContentType()
GetMethod()
GetPath()


```

## 黑名名单配置
engine.InitAclNodeList("./file/acl.conf")
```text
::1
127.0.0.1
192.168.0.0/16
```

## 静态资源服务

group.Static("/static/**", "./file")


## 路径参数
```go
group.GET("/refresh/:cacheName", func(ctx *wits.Context) {
		cacheName := ctx.GetString("cacheName")
		ctx.String(http.StatusOK, "ok")
})
```

## 优雅停机
kill -15 pid

## 中间件
```go
全局中间件

userGroup.Use(mock.TestMiddleWare)

路由级别
group.GET("probe", func(ctx *wits.Context) {
    _ = ctx.Success(timeutil.Now())
}, func(handlerFunc wits.HandlerFunc) wits.HandlerFunc {
    return func(ctx *wits.Context) {
}
})
```

## 全局异常处理
engine.SetHandlerBizError(BizHandler{})

## 认证
engine.AddFilter(mock.Accounts.BasicAuth())

## 文件处理

```go
userGroup.GET("file", func(c *wits.Context) {
		c.File("./file/file.xlsx")
})

userGroup.GET("fileFromFS", func(c *wits.Context) {
    c.FileFromFS("file.xlsx", http.Dir("./file"))
})

userGroup.GET("fileAttachment", func(c *wits.Context) {
    c.FileAttachment("./file/下载.txt", "下载.txt")
})
```

## 模版
```go
templateMust := template.Must(template.New("t").Parse(`Hello {{.name}}`))
engine.SetHTMLTemplate(templateMust)
userGroup.GET("template1", func(c *wits.Context) {
    c.HTMLTemplate(http.StatusCreated, "t", wits.Map{"name": "wendell"})
})


// ------------
engine.Delims("[{{", "}}]")
engine.SetFuncMap(template.FuncMap{
    "formatAsDate": formatAsDate,
})
engine.LoadHTMLGlob("file/**")

userGroup.GET("template2", func(c *wits.Context) {
    c.HTMLTemplate(http.StatusOK, "hello.tmpl", wits.Map{"name": "world"})

})
```

## 参数绑定
```go
func(c *wits.Context) {
    u := &User{}
    err := c.BindJSON(u)
    if err != nil {
        c.ERROR(err.Error())
    }
    c.DEBUG(u)
    _ = c.Success("success")
}
```

## 超时设置
1. 配置方式
app.server.global.timeout.enable=true
app.server.global.timeout=5000
2. 手动开启方式
EnableGlobalTimeout()

## session
1. 配置
app.session.enable = true
2. 手动
EnableSession()


> 主要配置

```text
app.session.enable
app.session.id
app.session.store
app.session.filter.Include
app.session.filter.Exclude
app.session.cookie.path
app.session.cookie.maxAge
app.session.cookie.domain
app.session.cookie.secure
app.session.cookie.httponly
```

## RestResult返回
app.rest.result.enable
```text
curl 127.0.0.1:8888/api/probe
{"code":200,"msg":"OK","data":"2023-08-20 17:55:19"}
```