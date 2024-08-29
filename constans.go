package wits

const (
	AuthKey      = "Authorization"
	DefaultRealm = "Require"
	BasicPrefix  = "Basic "
	UserKey      = "UserKey"
	slash        = "/"
	AuthErrorMsg = "AuthErrorMsg"
	JWTToken     = "JWTToken"
	Bearer       = "Bearer "
	RS256        = "RS256"
	RS512        = "RS512"
	RS384        = "RS384"
	HS256        = "HS256"
	HS384        = "HS384"
	HS512        = "HS512"

	WWWAuthenticate          = "WWW-Authenticate"
	XRealIP                  = "X-Real-Ip"
	XForwardedFor            = "X-Forwarded-For"
	XRequestID               = "X-Request-Id"
	HeaderReferer            = "Referer"
	HeaderUserAgent          = "User-Agent"
	HeaderContentType        = "Content-Type"
	HeaderContentDisposition = "Content-Disposition"

	DefaultLogPrefix    = "WITS"
	AppLogPrefix        = "app.log.prefix"
	AppServerCorsEnable = "app.server.cors.enable"
	AppName             = "app.name"
	AppPort             = "app.server.port"
	AppAddress          = "app.server.address"

	appServerGracefulTime    = "app.server.graceful.time"
	defaultSessionIdName     = "app-session-id"
	AppSessionIdName         = "app.session.id"
	AppSessionEnable         = "app.session.enable"
	AppSessionTimeout        = "app.session.timeout"
	AppSessionStore          = "app.session.store"
	AppSessionFilterInclude  = "app.session.filter.Include"
	AppSessionFilterExclude  = "app.session.filter.Exclude"
	AppSessionCookiePath     = "app.session.cookie.path"
	AppSessionCookieMaxAge   = "app.session.cookie.maxAge"
	AppSessionCookieDomain   = "app.session.cookie.domain"
	AppSessionCookieSecure   = "app.session.cookie.secure"
	AppSessionCookieHttpOnly = "app.session.cookie.httponly"
	AppRestResultEnable      = "app.rest.result.enable"
	AppRecordCostEnable      = "app.server.cost.enable"
	AppGlobalTimeoutEnable   = "app.server.global.timeout.enable"
	AppGlobalTimeout         = "app.server.global.timeout"
	AppBannerEnable          = "app.banner.enable"
)
