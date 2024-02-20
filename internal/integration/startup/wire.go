//go:build wireinject

package startup

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, InitRedis,

		dao.NewUserDAO,

		cache.NewCodeCache, cache.NewUserCache,

		repository.NewUserRepository, repository.NewCodeRepository,

		ioc.InitSMSService, service.NewUserService, service.NewCodeService,

		web.NewHandler,

		ioc.InitGinMiddlewares,

		ioc.InitWebServer,
	)
	return gin.Default()

}
