package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth struct
type OAuth struct {
	GoogleOauthConfig *oauth2.Config
	NaverOauthConfig  *oauth2.Config
}

const (

	// 인증 후 유저 정보를 가져오기 위한 API
	GoogleUserInfoAPIEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	NaverUserInfoAPIEndpoint  = "https://openapi.naver.com/v1/nid/me"

	// 인증 권한 범위. 여기에서는 프로필 정보 권한만 사용
	ScopeEmail   = "https://www.googleapis.com/auth/userinfo.email"
	ScopeProfile = "https://www.googleapis.com/auth/userinfo.profile"
)

func NewOauth(conf *Config) *OAuth {
	return &OAuth{
		GoogleOauthConfig: &oauth2.Config{
			RedirectURL:  conf.Auth.Google.Redirect,
			ClientID:     conf.Auth.Google.ClientID,
			ClientSecret: conf.Auth.Google.ClientSecret,
			Scopes:       []string{ScopeEmail, ScopeProfile},
			Endpoint:     google.Endpoint,
		},
		NaverOauthConfig: &oauth2.Config{
			RedirectURL:  conf.Auth.Naver.Redirect,
			ClientID:     conf.Auth.Naver.ClientID,
			ClientSecret: conf.Auth.Naver.ClientSecret,
			Scopes:       []string{"email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
				TokenURL: "https://nid.naver.com/oauth2.0/token",
			},
		}}
}
