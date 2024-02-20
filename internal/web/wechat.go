package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/auth2/wechat"
	ijwt "basic-go/webook/internal/web/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	key             []byte
	stateCookieName string
}

func NewWechatHandler(svc wechat.Service, hdl ijwt.Handler, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("qTzTTNzQcpXofciQynLVq1WbwRFeQrFn"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.OAuth2URL)

	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) OAuth2URL(context *gin.Context) {
	state := uuid.New()
	val, err := h.svc.AuthURL(context, state)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Msg:  "构造跳转 URL 失败",
			Code: 5,
		})
		return
	}
	err = h.setStateCookie(context, state)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Msg:  "服务器异常",
			Code: 5,
		})
	}
	context.JSON(http.StatusOK, Result{
		data: val,
	})
}

func (h *OAuth2WechatHandler) Callback(context *gin.Context) {
	err := h.verifyState(context)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Msg:  "state 被篡改",
			Code: 4,
		})
	}
	code := context.Query("code")
	//state := context.Query("state")

	wechatInfo, err := h.svc.VerifyCode(context, code)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	u, err := h.userSvc.FindOrCreateByWechat(context, wechatInfo)
	if err != nil {
		context.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 4,
		})
		return
	}

	err = h.SetLoginToken(context, u.ID)
	if err != nil {
		context.String(http.StatusOK, "系统错误")
		return
	}

	context.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
	return
}

func (h *OAuth2WechatHandler) setStateCookie(c *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.key)
	if err != nil {

		return err
	}
	c.SetCookie("jwt-state", tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(h.stateCookieName)
	if err != nil {
		return err
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.key, nil
	})
	if err != nil {
		return err
	}
	if state != sc.State {
		return errors.New("state被篡改了")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
