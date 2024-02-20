package middlewares

import (
	ijwt "basic-go/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddleWareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddleWareBuilder(hdl ijwt.Handler) *LoginJWTMiddleWareBuilder {
	return &LoginJWTMiddleWareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddleWareBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms" ||
			path == "/users/login_sms/code/send" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			return
		}

		tokenStr := m.ExtractToken(context)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//if uc.UserAgent != context.GetHeader("User-Agent") {
		//	context.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//
		//expireTime := uc.ExpiresAt
		//
		//if expireTime.Sub(time.Now()) < time.Second*50 {
		//	uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
		//	tokenStr, err = token.SignedString(web.JWTKey)
		//	context.Header("x-jwt-token", tokenStr)
		//	if err != nil {
		//		log.Println(err)
		//	}
		//}

		err = m.CheckSession(context, uc.Ssid)
		if err != nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		context.Set("user", uc)
	}
}
