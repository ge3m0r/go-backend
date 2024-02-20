//go:build wireinject

package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis, ioc.InitLogger,

		dao.NewUserDAO,

		cache.NewCodeCache, cache.NewUserCache,

		repository.NewUserRepository, repository.NewCodeRepository,

		ioc.InitSMSService, ioc.InitWechatService, service.NewUserService, service.NewCodeService,

		web.NewHandler, ijwt.NewRedisJWTHandler, web.NewWechatHandler,

		ioc.InitGinMiddlewares,

		ioc.InitWebServer,
	)
	return gin.Default()

}
