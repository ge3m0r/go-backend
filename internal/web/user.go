// 与http打交道
package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	ijwt "basic-go/webook/internal/web/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

// 所有用户有关的路由
type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewHandler(svc service.UserService, hdl ijwt.Handler, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		codeSvc:        codeSvc,
		Handler:        hdl,
	}

}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	{
		ug.POST("/login", h.LoginJWT)
		//server.PUT("/user/signup", h.Signup)
		ug.POST("/signup", h.SignUp)
		ug.POST("/loginout", h.LogoutJWT)
		ug.GET("/profile", h.Profile)
		ug.POST("/edit", h.Edit)
		ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
		ug.POST("/login_sms", h.LoginSMS)
		ug.GET("/refresh_token", h.RefreshToken)
	}

}

func (h *UserHandler) SendSMSLoginCode(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	if req.Phone == "" {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号",
		})
	}
	err := h.codeSvc.Send(c, bizLogin, req.Phone)
	switch err {
	case nil:

		c.JSON(http.StatusOK, Result{

			Msg: "验证码发送成功",
		})
	case service.ErrCodeSendTooMany:
		zap.L().Warn("频繁发送验证码", zap.Error(err))
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

	}
}

func (h *UserHandler) LoginSMS(c *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	ok, err := h.codeSvc.Verify(c, bizLogin, req.Phone, req.Code)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("手机验证码验证失败", zap.Error(err))
		return
	}
	if !ok {
		c.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不正确,请重新输入",
		})
		return
	}
	u, err := h.svc.FindOrCreate(c, req.Phone)
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.SetLoginToken(c, u.ID)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
	}

	c.JSON(http.StatusOK, Result{

		Msg: "登录成功",
	})
}

func (h *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		return
	}
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "非法邮箱")
		return
	}
	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次输入密码不对")
		return
	}
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		c.String(http.StatusOK, "密码必须包含数字，字母特殊字符，并且不少于8位")
		return
	}
	err = h.svc.Signup(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	switch err {
	case nil:

		c.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateUser:
		c.String(http.StatusOK, "邮箱错误，请换一个")
	default:
		c.String(http.StatusOK, "系统错误")

	}

}

func (h *UserHandler) Profile(c *gin.Context) {

	type ProfileReq struct {
		Id int64 `json:"id"`
	}
	var req ProfileReq
	if err := c.Bind(&req); err != nil {
		return
	}
	sess := sessions.Default(c)
	uc := sess.Get("userId")
	u, err := h.svc.Profile(c, domain.User{
		ID: uc.(int64),
	})
	println(err)
	switch err {
	case nil:
		c.JSON(http.StatusOK, &domain.User{
			NickName: u.NickName,
			Birthday: u.Birthday,
			AboutMe:  u.AboutMe,
		})
	default:
		c.String(http.StatusOK, "系统错误")

	}
}

func (h *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req EditReq
	if err := c.Bind(&req); err != nil {
		return
	}
	sess := sessions.Default(c)
	uc := sess.Get("userId")
	_, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		c.String(http.StatusOK, "生日格式不对")
		return
	}
	err = h.svc.Edit(c, domain.User{
		ID:       uc.(int64),
		NickName: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	switch err {
	case nil:
		c.String(http.StatusOK, "修改成功")
	default:
		c.String(http.StatusOK, "系统错误")

	}
}

func (h *UserHandler) LoginJWT(c *gin.Context) {
	type Req struct {
		Email    string
		Password string
	}
	var req Req
	if err := c.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(c, req.Email, req.Password)

	switch err {
	case nil:
		err = h.SetLoginToken(c, u.ID)
		if err != nil {
			c.String(http.StatusOK, "系统错误")
		}
		c.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		c.String(http.StatusOK, "用户名或者密码不对")

	default:
		c.String(http.StatusOK, "系统错误")

	}
	if err != nil {

		return
	}

}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	//约定， 前端  Authorization 里边带上相关 token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaimes
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录成功",
	})
}
