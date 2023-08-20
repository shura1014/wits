package wits

import (
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/wits/session"
	"time"
)

var (
	IdName         string
	sessionManager *session.Manager
	enableSession  bool
)

// EnableSession 应用手动开启
func EnableSession() {
	initSession()
}

// NewSessionManager 如果是手动开启session EnableSession
// 那么也应该设置session的相应参数 ttl 超时时间，store 存储方式
func NewSessionManager(ttl int, store string) {
	session.New(time.Duration(ttl)*time.Minute, store)
}

func init() {
	if GetBool(AppSessionEnable) {
		initSession()
	}
}

func initSession() {
	enableSession = true
	IdName = defaultSessionIdName
	v := GetString(AppSessionIdName)
	if v != "" {
		IdName = v
	}
	sessionManager = session.New(GetTime(AppSessionTimeout)*time.Minute, GetString(AppSessionStore))
}
func getSessionFilter() *Filter {
	return &Filter{
		Order:   5,
		Include: GetString(AppSessionFilterInclude, "/**"),
		Exclude: GetString(AppSessionFilterExclude, "/login,/favicon.ico"),
		BeforeFunc: func(ctx *Context) bool {
			if !enableSession {
				return true
			}
			id, _ := ctx.Cookie(IdName)
			if id == "" {
				id = ctx.RequestHeader(IdName)
			}
			if id == "" {
				ctx.HandlerError(goerr.WithCode(sessionIdNotFound))
				return false
			}

			appSession := sessionManager.GetStore().GetSession(id)
			if appSession == nil {
				ctx.HandlerError(goerr.WithCode(sessionTimeout))
				return false
			}
			ctx.session = appSession
			return true
		},
		AfterFunc: func(ctx *Context) {
			if ctx.session.IsNew() || ctx.session.IsDirty() {
				ctx.session.Save()
				ctx.DEBUG("session %s save", ctx.GetSessionId())
			}
		},
	}

}
