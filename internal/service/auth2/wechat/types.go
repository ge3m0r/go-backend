package wechat

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=SCOPE&state=%s#wechat_redirect `

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://wangsan.com/oauth2/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
	l         logger.Logger
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`, s.appID,
		s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, errors.New("调用微信接口失败")
	}
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {

	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

func NewService(appID string, appSecret string, l logger.Logger) *service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
		l:         l,
	}
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`

	//错误返回

	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
