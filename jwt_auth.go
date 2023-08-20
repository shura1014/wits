package wits

import (
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type JwtAuth struct {
	UnAuthHandler AuthCheckHandler
	Include       string
	Exclude       string
	//key
	Key []byte
}

func (j *JwtAuth) Auth() *Filter {
	return &Filter{
		Order:   10,
		Include: j.Include,
		Exclude: j.Exclude,
		BeforeFunc: func(ctx *Context) bool {
			token := ctx.R.Header.Get(AuthKey)
			token = strings.TrimPrefix(token, Bearer)
			if token == "" {
				j.unAuthHandler(ctx)
				return false
			}
			t := j.ParseToken(token, ctx)
			if t == nil {
				j.unAuthHandler(ctx)
				return false
			}
			return true
		},
	}
}

func (j *JwtAuth) unAuthHandler(ctx *Context) {
	if j.UnAuthHandler != nil {
		j.UnAuthHandler(ctx)
	} else {
		DefaultAuthCheckHandler(ctx)
	}
}

func (j *JwtAuth) ParseToken(token string, ctx *Context) *jwt.Token {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return j.Key, nil
	})
	if err != nil {
		ctx.Set(AuthErrorMsg, "jwt解析异常:"+err.Error())
		return nil
	}
	if !t.Valid {
		ctx.Set(AuthErrorMsg, "invalid token")
		return nil
	}
	ctx.Set("claims", t.Claims.(jwt.MapClaims))
	return t
}

type JwtToken struct {
	Key     []byte
	TimeOut time.Duration
	Alg     string
	GetData func(ctx *Context) (map[string]any, error)
}

// GetToken
// 获取token，是一个前置中间件
// 应该使用在某一个路由中，比如/login
// 其中提供一个获取数据的函数，以获得数据生成token
func (j *JwtToken) GetToken() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			data, err := j.GetData(ctx)
			if err != nil {
				_ = ctx.Fail("Token: 获取数据异常" + err.Error())
				return
			}
			if data == nil {
				_ = ctx.Fail("Token: 为查询到数据")
				return
			}

			// 获取token串
			if j.Alg == "" {
				j.Alg = HS256
			}
			signingMethod := jwt.GetSigningMethod(j.Alg)
			token := jwt.New(signingMethod)
			claims := token.Claims.(jwt.MapClaims)
			if data != nil {
				for key, value := range data {
					claims[key] = value
				}
			}
			//now := time.Now()

			now := func() time.Time {
				return time.Now()
			}
			expire := now().Add(j.TimeOut)
			// jwt过期时间
			claims["exp"] = expire.Unix()
			// jwt签发时间
			claims["iat"] = now().Unix()

			tokenString, errToken := token.SignedString(j.Key)
			if errToken != nil {
				_ = ctx.Fail("Token: jwt创建异常: " + errToken.Error())
				return
			}

			ctx.Set(JWTToken, tokenString)
			next(ctx)
		}
	}
}

// RefreshToken 刷新令牌
// token 令牌 date 续期时间
func RefreshToken(token string, date time.Duration, key []byte) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return key, nil
	})
	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(date).Unix()
	if err != nil {
		return "", nil
	}
	return t.SignedString(key)
}
