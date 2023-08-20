package wits

import (
	"crypto/subtle"
	"encoding/base64"
	"github.com/shura1014/common/utils/stringutil"
	"net/http"
)

// AuthCheckHandler 定义一个没有认证的回调
// DefaultAuthCheckHandler 默认回调
// Accounts
// 逻辑完全由业务自己定义
// 业务将user加载到Accounts对象
// 业务设置没有检查到user的处理逻辑
type AuthCheckHandler func(ctx *Context)

var DefaultAuthCheckHandler = func(ctx *Context) {
	msg := ctx.GetString(AuthErrorMsg)
	if msg == "" {
		msg = "Authorization failure"
	}
	_ = ctx.String(http.StatusUnauthorized, msg)

}

type Accounts struct {
	UnAuthHandler AuthCheckHandler
	Users         map[string]string
}

type authPair struct {
	base64Value string
	user        string
}

type authPairs []authPair

// BasicAuth 过滤器
// 返回一个过滤器
// 将过滤器添加到engine即可使用
func (a *Accounts) BasicAuth() *Filter {
	Debug("Enable basic auth")
	// 加载账号
	pairs := processAccounts(*a)
	return &Filter{
		Order: 10,
		BeforeFunc: func(ctx *Context) bool {
			user, found := pairs.searchCredential(ctx.RequestHeader(AuthKey))
			if !found {
				ctx.Header(WWWAuthenticate, DefaultRealm)
				a.unAuthHandler(ctx)
				return false
			}
			// 存放username 可供后续业务使用
			ctx.Set(UserKey, user)
			return true
		},
	}
}

func (a *Accounts) unAuthHandler(ctx *Context) {
	if a.UnAuthHandler != nil {
		a.UnAuthHandler(ctx)
	} else {
		DefaultAuthCheckHandler(ctx)
	}
}

// 加载账号
// base64处理
func processAccounts(accounts Accounts) authPairs {
	length := len(accounts.Users)
	assert(length > 0, "Empty list of authorized credentials")
	pairs := make(authPairs, 0, length)
	for user, password := range accounts.Users {
		assert(user != "", "User can not be empty")
		value := BasicAuth(user, password)
		pairs = append(pairs, authPair{
			base64Value: value,
			user:        user,
		})
	}
	return pairs
}

// BasicAuth
// base64处理
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// 寻找可认证账号 遍历
func (a authPairs) searchCredential(authValue string) (string, bool) {
	authValue = authValue[len(BasicPrefix):]
	if authValue == "" {
		return "", false
	}
	for _, pair := range a {
		if subtle.ConstantTimeCompare(stringutil.StringToBytes(pair.base64Value), stringutil.StringToBytes(authValue)) == 1 {
			return pair.user, true
		}
	}
	return "", false
}
