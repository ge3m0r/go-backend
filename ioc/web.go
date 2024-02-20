package ioc

import (
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/internal/web/middlewares"
	"basic-go/webook/pkg/ginx/middleware/ratelimit"
	"basic-go/webook/pkg/limiter"
	"basic-go/webook/pkg/logger"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()

	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable, hdl ijwt.Handler, l logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token", "x-fresh-token"},
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 12 * time.Hour,
		}), func(context *gin.Context) {
			println("这是跨域middleware")
		}, (&middlewares.LoginJWTMiddleWareBuilder{}).CheckLogin(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		middlewares.NewLogMiddlewareBuilder(func(ctx context.Context, al middlewares.AcessLog) {
			l.Debug("", logger.Field{Key: "req", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),
		middlewares.NewLoginJWTMiddleWareBuilder(hdl).CheckLogin(),
	}
}
