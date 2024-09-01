package wits

import (
	"github.com/shura1014/common/container/tree"
	"net/http"
	"path/filepath"
	"strings"
)

// Static StaticFS
// 作为静态资源服务器，实际上注册的是GET HEAD请求
// 使用示例，必须以/**结尾
// group.Static("/static/**", "./file")
func (r *routerGroup) Static(handlerPath, root string) {
	r.StaticFS(handlerPath, http.Dir(root))
}

func (r *routerGroup) StaticFS(handlerPath string, fs http.FileSystem) {
	if !strings.HasSuffix(handlerPath, "/**") {
		panic("静态资源文件路由路径需要以/**结尾")
	}

	handler := r.createStaticHandler(handlerPath, fs)
	r.GET(handlerPath, handler)
	r.HEAD(handlerPath, handler)
}

func (r *routerGroup) createStaticHandler(handlerPath string, fs http.FileSystem) HandlerFunc {

	relativePath := strings.TrimSuffix(filepath.Join(r.groupName+handlerPath), "**")
	fileServer := http.StripPrefix(relativePath, http.FileServer(fs))
	return func(ctx *Context) {
		// 检查一下文件是否存在
		file := ctx.GetString(tree.FullMatchKey)
		f, err := fs.Open(file)
		if err != nil {
			ctx.W.WriteStatus(http.StatusNotFound)
			return
		}
		_ = f.Close()
		ctx.W.WriteStatus(restErrorCode)
		fileServer.ServeHTTP(ctx.W, ctx.R)

	}
}
