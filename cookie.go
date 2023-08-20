package wits

import (
	"net/http"
	"net/url"
	"time"
)

var defaultCookieOptions *CookieOptions

func init() {
	defaultCookieOptions = GetCookieOptionsByEnv()
}

// CookieOptions 可选项
type CookieOptions struct {
	MaxAge           time.Duration
	Path, Domain     string
	Secure, HttpOnly bool
}

func GetCookieOptionsByEnv() *CookieOptions {
	return &CookieOptions{
		Path:     GetString(AppSessionCookiePath, "/"),
		Domain:   GetString(AppSessionCookieDomain, ""),
		MaxAge:   time.Hour * GetTime(AppSessionCookieMaxAge, 24*7),
		Secure:   GetBool(AppSessionCookieSecure),
		HttpOnly: GetBool(AppSessionCookieHttpOnly),
	}
}

func SetCookie(name, value string, ctx *Context, options ...*CookieOptions) {
	option := defaultCookieOptions
	if len(options) > 0 {
		option = options[0]
	}
	sameSite := ctx.sameSite
	if IsOpenCors() {
		// 跨站使用
		sameSite = http.SameSiteNoneMode
		option.Secure = true
	}
	http.SetCookie(ctx.W, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   int(option.MaxAge),
		Expires:  time.Now().Add(option.MaxAge),
		Path:     option.Path,
		Domain:   option.Domain,
		SameSite: sameSite,
		Secure:   option.Secure,
		HttpOnly: option.HttpOnly,
	})
}
