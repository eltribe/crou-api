package auth

import (
	"context"
	"crou-api/config"
	"crou-api/config/database"
	"crou-api/internal/application/usecase/user"
	"crou-api/internal/domains"
	messages2 "crou-api/messages"
	"crou-api/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"io"
	"time"
)

var store = session.New(session.Config{
	Expiration: time.Second * 180,
})

type OAuth2Service struct {
	database    database.Persistent
	config      *config.Config
	auth        *config.OAuth
	userService *user.UserService
	jwtProvider *utils.JwtProvider
}

func NewOAuth2Service(database database.Persistent, auth *config.OAuth, config *config.Config, userService *user.UserService) *OAuth2Service {
	return &OAuth2Service{
		database:    database,
		auth:        auth,
		config:      config,
		userService: userService,
		jwtProvider: utils.NewJwtService(config),
	}
}

// OauthGoogleLogin godoc
//
//	@Summary		Oauth2.0 Google 인증 API
//	@Description	Oauth2.0 Google 인증 URL을 반환합니다
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			messages.OauthLoginRequest	query		messages.OauthLoginRequest	true	"OauthLoginRequest"
//	@Success		200							{object}	messages.OauthLoginUrl
//	@Failure		409							{object}	server.Error
//	@Router			/v1/auth/google [get]
func (srv *OAuth2Service) OauthGoogleLogin(c *fiber.Ctx, req *messages2.OauthLoginRequest) (*messages2.OauthLoginUrl, error) {
	state := utils.RandToken(5)
	//sess, err := store.Get(c)
	//if err != nil {
	//	return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	//}
	//
	//sess.Set(config.AUTH_STATE_KEY, state)
	//// Save session
	//if err := sess.Save(); err != nil {
	//	return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	//}

	srv.auth.GoogleOauthConfig.RedirectURL = req.CallbackUrl
	url := srv.auth.GoogleOauthConfig.AuthCodeURL(state)
	return &messages2.OauthLoginUrl{
		Url: url,
	}, nil
}

// OauthNaverLogin godoc
//
//	@Summary		Oauth2.0 Naver 인증 API
//	@Description	Oauth2.0 Naver 인증 URL을 반환합니다
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			messages.OauthLoginRequest	query		messages.OauthLoginRequest	true	"OauthLoginRequest"
//	@Success		200							{object}	messages.OauthLoginUrl
//	@Failure		409							{object}	server.Error
//	@Router			/v1/auth/naver [get]
func (srv *OAuth2Service) OauthNaverLogin(c *fiber.Ctx, req *messages2.OauthLoginRequest) (*messages2.OauthLoginUrl, error) {
	state := utils.RandToken(5)
	//sess, err := store.Get(c)
	//if err != nil {
	//	return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	//}
	//
	//sess.Set(config.AUTH_STATE_KEY, state)
	//// Save session
	//if err := sess.Save(); err != nil {
	//	return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	//}

	srv.auth.NaverOauthConfig.RedirectURL = req.CallbackUrl
	url := srv.auth.NaverOauthConfig.AuthCodeURL(state)
	return &messages2.OauthLoginUrl{
		Url: url,
	}, nil
}

// OauthGoogleCallback godoc
//
//	@Summary		Oauth2.0 Google 인증 완료 처리 API
//	@Description	Oauth2.0 Google 인증 완료 처리 후 accessToken을 발급합니다.
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"code"
//	@Param			state	query		string	false	"state"
//	@Success		200		{object}	messages.AccessToken
//	@Failure		409		{object}	server.Error
//	@Router			/v1/auth/google/callback [get]
func (srv *OAuth2Service) OauthGoogleCallback(c *fiber.Ctx, code string, state string) (*messages2.AccessToken, error) {

	//sess, _ := store.Get(c)
	//storedState := sess.Get(config.AUTH_STATE_KEY)
	//if storedState == nil {
	//	return nil, fiber.NewError(fiber.StatusUnauthorized, "state 가 만료되었습니다. 로그인을 다시 시도해주세요")
	//}
	//if storedState != state {
	//	return nil, fiber.NewError(fiber.StatusUnauthorized, "state 값이 유효하지 않습니다")
	//}

	// Use the authorization code that is pushed to the redirect URL
	token, err := srv.auth.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if !token.Valid() {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid oauth token")
	}

	// Getting the userinfo
	client := srv.auth.GoogleOauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get(config.GoogleUserInfoAPIEndpoint)

	defer userInfoResp.Body.Close()
	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	var googleUser messages2.GoogleOauthUserInfo
	_ = json.Unmarshal(userInfo, &googleUser)

	// 이메일, 고유 ID 기반 유저 체크
	userInfo2, err := srv.userService.GetUserByOauthInfo(domains.GOOGLE, googleUser.Sub)
	if err != nil {
		jwtToken, err := srv.jwtProvider.GenerateTemporaryJwt(domains.GOOGLE, googleUser.Sub, googleUser.Email)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusPaymentRequired, jwtToken.Token)
	}
	// 성공시 JWT 재할당
	jwtToken, err := srv.jwtProvider.GenerateOauthJwt(domains.GOOGLE, googleUser.Sub, googleUser.Email, userInfo2.Nickname)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	srv.saveRefreshToken(jwtToken.RefreshToken, userInfo2.ID)

	return jwtToken, nil
}

// OauthNaverCallback godoc
//
//	@Summary		Oauth2.0 Naver 인증 완료 처리 API
//	@Description	Oauth2.0 Naver 인증 완료 처리 후 accessToken을 발급합니다.
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"code"
//	@Param			state	query		string	false	"state"
//	@Success		200		{object}	messages.AccessToken
//	@Failure		409		{object}	server.Error
//	@Router			/v1/auth/naver/callback [get]
func (srv *OAuth2Service) OauthNaverCallback(c *fiber.Ctx, code string, state string) (*messages2.AccessToken, error) {

	//sess, _ := store.Get(c)
	//storedState := sess.Get(config.AUTH_STATE_KEY)
	//if storedState == nil {
	//	log.Info("state 가 만료되었습니다")
	//	return nil, fiber.NewError(fiber.StatusUnauthorized, "state 가 만료되었습니다. 로그인을 다시 시도해주세요")
	//}
	//if storedState != state {
	//	log.Info("state 값이 유효하지 않습니다")
	//	return nil, fiber.NewError(fiber.StatusUnauthorized, "state 값이 유효하지 않습니다")
	//}

	// Use the authorization code that is pushed to the redirect URL
	token, err := srv.auth.NaverOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if !token.Valid() {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid oauth token")
	}

	// Getting the userinfo
	client := srv.auth.NaverOauthConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get(config.NaverUserInfoAPIEndpoint)

	defer userInfoResp.Body.Close()
	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	var naverUser messages2.NaverOauthUserInfo
	_ = json.Unmarshal(userInfo, &naverUser)

	// 이메일, 고유 ID 기반 유저 체크
	userInfo2, err := srv.userService.GetUserByOauthInfo(domains.NAVER, naverUser.Response.ID)
	if err != nil {
		jwtToken, err := srv.jwtProvider.GenerateTemporaryJwt(domains.NAVER, naverUser.Response.ID, naverUser.Response.Email)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return nil, fiber.NewError(fiber.StatusPaymentRequired, jwtToken.Token)
	}

	jwtToken, err := srv.jwtProvider.GenerateOauthJwt(domains.NAVER, naverUser.Response.ID, naverUser.Response.Email, userInfo2.Nickname)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	srv.saveRefreshToken(jwtToken.RefreshToken, userInfo2.ID)
	return jwtToken, nil
}

// Join godoc
//
//	@Summary		회원가입 API
//	@Description	Oauth 이후 회원가입을 요청합니다
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			messages.JoinRequest	body		messages.JoinRequest	true	"JoinRequest"
//	@Success		200						{object}	messages.AccessToken
//	@Failure		401						{object}	server.Error
//	@Failure		409						{object}	server.Error
//	@Router			/v1/join [post]
//func (srv *OAuth2Service) Join(c *fiber.Ctx, input *messages2.JoinRequest) (*messages2.AccessToken, error) {
//	jwtClaims := srv.jwtProvider.GetClaims(c)
//	newUser := &domains.User{
//		Nickname:   input.Nickname,
//		OauthType:  domains.ParseOauthType(jwtClaims.Type),
//		OauthSub:   &jwtClaims.Sub,
//		OauthEmail: &jwtClaims.Email,
//	}
//	_, err := srv.userService.CreateUser(newUser)
//	if err != nil {
//		return nil, fiber.NewError(fiber.StatusConflict, err.Error())
//	}
//	// JWT 생성
//	jwtToken, err := srv.jwtProvider.GenerateOauthJwt(domains.GOOGLE, *newUser.OauthSub, *newUser.OauthEmail, newUser.Nickname)
//	srv.saveRefreshToken(jwtToken.RefreshToken, newUser.ID)
//	if err != nil {
//		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
//	}
//	return jwtToken, nil
//}

// Refresh godoc
//
//	@Summary		리프레쉬 토큰 요청  API
//	@Description	리프레쉬 토큰을 요청합니다.
//	@Tags			인증
//	@Accept			json
//	@Produce		json
//	@Param			messages.RefreshTokenRequest	body		messages.RefreshTokenRequest	true	"RefreshTokenRequest"
//	@Success		200						{object}	messages.AccessToken
//	@Failure		401						{object}	server.Error
//	@Failure		409						{object}	server.Error
//	@Router			/v1/auth/refresh [post]
func (srv *OAuth2Service) Refresh(c *fiber.Ctx, input *messages2.RefreshTokenRequest) (*messages2.AccessToken, error) {

	if err := utils.VerifyPassword(input.RefreshToken, srv.config.Auth.JWT.Secret); err != nil {
		log.Info(`올바르지 않는 Refresh Token입니다.`, err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, `올바르지 않는 Refresh Token입니다.`)
	}
	storedUserId, err := srv.getRefreshToken(input.RefreshToken)
	if err != nil || storedUserId == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Refresh Token이 존재하지 않습니다. 다시 인증해주세요")
	}

	userInfo, err := srv.userService.GetUserById(utils.StringToUint32(storedUserId))
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "존재하지 않는 유저입니다.")
	}

	// JWT 생성
	srv.removeRefreshToken(input.RefreshToken)
	jwtToken, err := srv.jwtProvider.GenerateOauthJwt(userInfo.OauthType, *userInfo.OauthSub, *userInfo.OauthEmail, userInfo.Nickname)
	srv.saveRefreshToken(jwtToken.RefreshToken, userInfo.ID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return jwtToken, nil
}

func (srv *OAuth2Service) saveRefreshToken(refreshToken string, userId uint) {
	srv.database.REDIS().Set(context.Background(), "refreshtoken:"+refreshToken, utils.UintToString(uint64(userId)), time.Hour*time.Duration(srv.config.Auth.JWT.RefreshTokenExpiresHours))
}

func (srv *OAuth2Service) getRefreshToken(refreshToken string) (string, error) {
	return srv.database.REDIS().Get(context.Background(), "refreshtoken:"+refreshToken).Result()
}

func (srv *OAuth2Service) removeRefreshToken(refreshToken string) {
	srv.database.REDIS().Del(context.Background(), "refreshtoken:"+refreshToken)
}
