package wits

import (
	"errors"
	"github.com/shura1014/acl"
	"github.com/shura1014/common/container/tree"
	"github.com/shura1014/common/utils/fileutil"
	"github.com/shura1014/logger"
	"github.com/shura1014/wits/render"
	"html/template"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// RouterHandler 路由的真正处理逻辑，不同的逻辑可以实现不同的功能，例如代理转发
type RouterHandler func(ctx *Context)

type Engine struct {
	mu sync.RWMutex
	*router

	/**********模版**************/
	render  render.HtmlTemplateRender
	delims  render.Delims
	FuncMap template.FuncMap
	/**********模版**************/

	// context使用频率高，需要资源复用
	pool sync.Pool

	// 应用防火墙
	acl acl.Acl

	// form表单上传文件的读取是是否全部加载到内存的阈值
	MaxMultipartMemory int64

	// log日志记录对象
	logger *logger.Logger

	// 过滤器，中间件只是对方法业务前后的处理，而过滤器是对整个请求的前后处理
	filterChain FilterChain

	// 全局异常错误处理逻辑
	handlerBizError HandlerBizError

	// 应用名称 如：console
	appName string

	// 需要在应用启动的时候执行
	PreRunFunc []func()

	// New出一个Engine的时候执行
	PostNewFunc []func()

	// 应用启动完后的预热操作
	WarmupFunc []func()

	// httpServer
	srv *http.Server

	// server创建后执行
	serverPostProcessor []ServerPostProcessor

	handler RouterHandler

	// 优雅停机钩子
	shutdownHook []func()
}

// New 得到一个全新Engine 执行 engine.Run() 即可启动一个web服务
func New() *Engine {
	e := &Engine{
		router:             &router{},
		FuncMap:            template.FuncMap{},
		delims:             render.Delims{Left: "{{", Right: "}}"},
		MaxMultipartMemory: 32 << 32,
	}
	e.RegisterPreRunFunc(e.defaultPreRun)
	e.RegisterPostNewFunc(e.defaultPostNew)
	e.pool.New = func() any {
		return e.allocateContext()
	}
	e.handler = e.handleRouter
	printBanner()
	e.postNew()
	return e
}

// Default 一般情况下使用默认的便可满足需求
func Default() *Engine {
	engine := New()
	engine.router.Use(Recovery)
	// 白名单
	engine.acl = acl.Default()
	engine.logger = applog
	engine.AddFilter(RequestLog)
	return engine
}

// Use 全局路由
func (e *Engine) Use(middlewares ...MiddlewareFunc) {
	e.router.Use(middlewares...)
}

// SetRouterHandler 注册路由处理函数
// 只有当需要做一些其他中间件的时候可能用到
func (e *Engine) SetRouterHandler(handler RouterHandler) {
	e.handler = handler
}

// SetAcl 设置防火墙
func (e *Engine) SetAcl(acl acl.Acl) {
	e.acl = acl
}

func (e *Engine) GetAcl() acl.Acl {
	return e.acl
}

// SetLogger 替换日志
func (e *Engine) SetLogger(logger *logger.Logger) {
	e.logger = logger
	applog = logger
}

func (e *Engine) Logger() *logger.Logger {
	return e.logger
}

// RegisterPreRunFunc 添加启动函数 Run 之前执行
func (e *Engine) RegisterPreRunFunc(preRun func()) {
	e.PreRunFunc = append(e.PreRunFunc, preRun)
}

// RegisterPostNewFunc New 之后执行
func (e *Engine) RegisterPostNewFunc(postNew func()) {
	e.PostNewFunc = append(e.PostNewFunc, postNew)
}

// RegisterWarmupFunc 添加预热执行函数
func (e *Engine) RegisterWarmupFunc(warmup func()) {
	e.PreRunFunc = append(e.WarmupFunc, warmup)
}

// AddFilter 添加过滤器
func (e *Engine) AddFilter(check *Filter) {
	includePath := check.Include
	if includePath != "" {
		check.includeTrie = &tree.Node{
			Name:     "/",
			Children: make([]*tree.Node, 0),
		}
		paths := strings.Split(includePath, ",")
		for _, path := range paths {
			check.includeTrie.Put(path)
		}
	}
	excludePath := check.Exclude

	if excludePath != "" {
		check.excludeTrie = &tree.Node{
			Name:     "/",
			Children: make([]*tree.Node, 0),
		}
		paths := strings.Split(excludePath, ",")
		for _, path := range paths {
			check.excludeTrie.Put(path)
		}
	}
	e.filterChain = append(e.filterChain, check)
	// 排序
	sort.Sort(e.filterChain)
}

// 分配context，只有当pool里面的不够用了才会分配
func (e *Engine) allocateContext() *Context {
	return &Context{e: e, keys: make(tree.Keys), logger: e.logger, W: Wrap(nil)}
}

// InitAclNodeList 初始化白名单列表
// 例如：fileName 文件名 ./file/acl.conf
func (e *Engine) InitAclNodeList(fileName string) {
	lines := fileutil.Read(fileName)
	if len(lines) > 0 {
		// 开启
		e.acl.Enable()
		Info("----------------------- Application Firewall -------------------------------")
		Info("acl path %s", fileName)
		for _, line := range lines {
			e.acl.ParseAclNode(line)
			Info(line)
		}
		Info("----------------------------------------------------------------------------")
	}
}

// SetHandlerBizError 设置异常处理函数
func (e *Engine) SetHandlerBizError(handler HandlerBizError) {
	e.handlerBizError = handler
}

/****************************************************模版**************************************************/

// SetHTMLTemplate 初始化模版
func (e *Engine) SetHTMLTemplate(template *template.Template) {
	e.render = render.HtmlTemplateRender{Template: template.Funcs(e.FuncMap)}
}

// SetFuncMap FuncMap
/*
func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

engine.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
})
*/
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.FuncMap = funcMap
}

// Delims 模版规则一般默认是{{xx}}
// engine.Delims("[{{", "}}]")
func (e *Engine) Delims(left, right string) *Engine {
	e.delims = render.Delims{Left: left, Right: right}
	return e
}

// LoadHTMLGlob 根据pattern规则一次性初始化多个模版
// 例如 file/**
func (e *Engine) LoadHTMLGlob(pattern string) {
	left := e.delims.Left
	right := e.delims.Right
	mustTemplate := template.Must(template.New("").Delims(left, right).Funcs(e.FuncMap).ParseGlob(pattern))

	e.SetHTMLTemplate(mustTemplate)
}

func (e *Engine) LoadHTMLFiles(files ...string) {
	mustTemplate := template.Must(template.New("").Delims(e.delims.Left, e.delims.Right).Funcs(e.FuncMap).ParseFiles(files...))
	e.SetHTMLTemplate(mustTemplate)
}

/***************************************************模版*************************************************/

/*
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
// ServeHTTP 实现Handler接口，是的Engine本身就是一个请求执行函数
// 目的是为了实现统一的入口 /** -> Engine
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.pool.Get().(*Context)
	ctx.Reset(w, r)
	// 每个请求将进入此处
	e.filterChain.doFilter(ctx, e.handler)
	e.pool.Put(ctx)
}

// 路由逻辑
func (e *Engine) handleRouter(ctx *Context) {
	//path := ctx.R.RequestURI
	path := ctx.R.URL.Path

	routerGroup := e.getRouterGroup(path)
	if routerGroup != nil {

		treeNode := routerGroup.treeNode.Get(routerGroup.Trim(path), ctx.keys)
		if treeNode != nil && treeNode.IsEnd {
			// 判断是否支持任意风格请求
			routerName := treeNode.RouterName
			handler, ok := routerGroup.handleMap[routerName][ANY]

			// 获取请求的协议 POST,GET
			method := ctx.R.Method
			if ok {
				routerGroup.methodHandle(routerName, ANY, handler, ctx)
				return
			}

			handler, ok = routerGroup.handleMap[routerName][method]
			if ok {
				routerGroup.methodHandle(routerName, method, handler, ctx)
				return
			}
			_ = ctx.String(http.StatusMethodNotAllowed, "%s not allowed request %s", path, method)

			return
		}
	}

	_ = ctx.String(http.StatusNotFound, "404 not found ")
	return
}

// Run 启动一个web服务
// ⚠️ 没有预热操作
func (e *Engine) Run(address ...string) {
	e.preRun()
	http.Handle("/", e)
	err := http.ListenAndServe(resolveAddress(address), nil)
	if err != nil {
		Fatal(err)
	}
}

// RunTLS 支持https启动
// ⚠️ 没有预热操作
func (e *Engine) RunTLS(certFile, keyFile string, address ...string) (err error) {
	e.preRun()
	addr := resolveAddress(address)
	Info("Listening and serving HTTPS on %s\n", addr)
	defer func() { Error(err) }()
	err = http.ListenAndServeTLS(addr, certFile, keyFile, e)
	return
}

// RunServer 建议都是用该函数启动
// 初始化 预热 优雅停机
func (e *Engine) RunServer(address ...string) {
	e.preRun()
	srv := &http.Server{
		//Addr:    resolveAddress(address),
		Handler: e,
	}
	e.srv = srv
	e.serverHandler()
	ln, err := net.Listen("tcp", resolveAddress(address))
	if err != nil {
		panic(err)
	}

	go func() {
		// 预热
		e.Warmup()
		// 服务连接
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			Error("listen: %s\n", err)
		}
	}()

	e.handlerSignal()
}

// AppName web服务名称
func (e *Engine) AppName() string {
	return e.appName
}

// Warmup 预热操作
func (e *Engine) Warmup() {
	for _, warmup := range e.WarmupFunc {
		warmup()
	}
}

func (e *Engine) reload() {

}

func (e *Engine) RegisterOnShutdown(hook func()) {
	e.mu.Lock()
	e.shutdownHook = append(e.shutdownHook, hook)
	e.mu.Unlock()
}
