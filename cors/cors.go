package cors

import (
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/wits"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Cors 跨域中间件，一般用于开发阶段
// https://www.w3.org/TR/cors/
type Cors struct {
	AllowAllOrigins  bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func (cors *Cors) AddAllowMethods(methods ...string) {
	cors.AllowMethods = append(cors.AllowMethods, methods...)
}

func (cors *Cors) AddAllowOrigins(origin ...string) {
	cors.AllowOrigins = append(cors.AllowOrigins, origin...)
}

func (cors *Cors) AddAllowHeaders(headers ...string) {
	cors.AllowHeaders = append(cors.AllowHeaders, headers...)
}

func (cors *Cors) AddExposeHeaders(headers ...string) {
	cors.ExposeHeaders = append(cors.ExposeHeaders, headers...)
}

func DefaultConfig() Cors {
	return Cors{
		AllowMethods:     []string{"POST, GET, OPTIONS, DELETE"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func Default() wits.MiddlewareFunc {
	return func(next wits.HandlerFunc) wits.HandlerFunc {
		return func(ctx *wits.Context) {
			origin := ctx.RequestHeader("Origin")
			if origin != "" {
				ctx.Header("Access-Control-Allow-Origin", ctx.RequestHeader("Origin"))
			} else {
				referer := ctx.Referer()
				if origin != "" {
					// 去掉 http://
					referer = referer[6:]
					if pos := strings.Index(referer, "/"); pos != -1 {
						referer = referer[:pos]
					}
				}
				ctx.Header("Access-Control-Allow-Origin", referer)

			}
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			ctx.Header("Access-Control-Max-Age", "3600")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Origin, Content-Length, Content-Type, Set-Cookie")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			if ctx.GetMethod() == "OPTIONS" {
				ctx.ReturnNow(http.StatusNoContent)
				return
			}
			next(ctx)
		}
	}
}

func New(cors Cors) wits.MiddlewareFunc {
	return func(next wits.HandlerFunc) wits.HandlerFunc {
		return func(ctx *wits.Context) {
			origin := ctx.RequestHeader("Origin")

			if !cors.validateOrigin(origin) {
				ctx.ReturnNow(http.StatusForbidden)
				return
			}

			if cors.AllowAllOrigins {
				ctx.Header("Access-Control-Allow-Origin", "*")
			} else if len(cors.AllowOrigins) > 0 {

			} else {
				ctx.Header("Access-Control-Allow-Origin", origin)
			}
			if len(cors.AllowMethods) > 0 {
				ctx.Header("Access-Control-Allow-Methods", strings.Join(cors.AllowMethods, ","))

			}
			if cors.MaxAge > time.Duration(0) {
				ctx.Header("Access-Control-Max-Age", strconv.FormatInt(int64(cors.MaxAge/time.Second), 10))
			}
			if len(cors.AllowHeaders) > 0 {
				ctx.Header("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders, ","))
			}
			if len(cors.ExposeHeaders) > 0 {
				ctx.Header("Access-Control-Expose-Headers", strings.Join(cors.ExposeHeaders, ","))
			}
			if cors.AllowCredentials {
				ctx.Header("Access-Control-Allow-Credentials", "true")
			}
			if ctx.GetMethod() == "OPTIONS" {
				ctx.ReturnNow(http.StatusNoContent)
				return
			}
			next(ctx)
		}
	}
}

func (cors Cors) validateOrigin(origin string) bool {
	if stringutil.IsArray(cors.AllowOrigins, "*") {
		return true
	}
	return stringutil.IsArray(cors.AllowOrigins, origin)
}
