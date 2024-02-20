package main

import (
	"basic-go/webook/internal/web/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	initViperV1()
	initLogger()
	server := InitWebServer()

	server.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello,启动了")
	})
	server.Run(":8080")
}

/*func initUserHdl(db *gorm.DB, redisClient redis.Cmdable, codeSvc service.codeService, server *gin.Engine) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	c := web.NewHandler(us, codeSvc)
	c.RegisterRoutes(server)
}*/

/*func initCodeSvc(redisClient redis.Cmdable) *service.codeService {
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	return service.NewCodeService(crepo, nil)
}*/
/*
func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(context *gin.Context) {
		println("这是跨域middleware")
	})
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//go get github.com/ulule/limiter/v3 另一个限流中间件
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	useJWT(server)
	return server
}*/

func useSession(server *gin.Engine) {
	login := &middlewares.LoginMiddleWareBuilder{}
	store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", []byte("qTzTTMzQcpXofciQynLVq1WbwRFeQrFn"), []byte("ik6K0pTEgJ8aqboo011NePKWmX837gxa"))
	//if err != nil {
	//panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}

func useJWT(server *gin.Engine) {
	loginJWT := &middlewares.LoginJWTMiddleWareBuilder{}
	server.Use(loginJWT.CheckLogin())
}

func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println(viper.Get("test.key"))

}

func initViperV1() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println(viper.Get("test.key"))
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3", "http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.ReadRemoteConfig()
}

func initLogger() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}
