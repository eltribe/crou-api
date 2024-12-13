package utils

import (
	"crou-api/config"
	"crou-api/internal/domains"
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtProvider struct {
	config *config.JWTConfig
}

func NewJwtService(config *config.Config) *JwtProvider {
	return &JwtProvider{config: &config.Auth.JWT}
}

func GetAccessUser(svc *JwtProvider, c *fiber.Ctx) *messages.JwtClaims {
	claims := svc.GetClaims(c)
	return claims
}

func (srv *JwtProvider) GenerateTemporaryJwt(oauthType domains.OauthType, sub, email string) (*messages.AccessToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expires := time.Now().Add(time.Hour * time.Duration(srv.config.ExpiresHours)).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["type"] = oauthType
	claims["sub"] = sub
	claims["email"] = email
	claims["approve"] = false
	claims["exp"] = expires

	tokenString, err := token.SignedString([]byte(srv.config.Secret))
	if err != nil {
		return nil, err
	}
	return &messages.AccessToken{
		Token:     tokenString,
		ExpiresIn: expires,
	}, nil
}

func (srv *JwtProvider) GenerateJwt(oauthType *domains.OauthType, email, nickname string) (*messages.AccessToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expires := time.Now().Add(time.Hour * time.Duration(srv.config.ExpiresHours)).Unix()

	claims := token.Claims.(jwt.MapClaims)

	if oauthType == nil {
		claims["oauth"] = "X"
	} else {
		claims["oauth"] = *oauthType
	}
	//claims["sub"] = sub
	claims["email"] = email
	claims["nickname"] = nickname
	claims["exp"] = expires

	tokenString, err := token.SignedString([]byte(srv.config.Secret))
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshTokenExpires := time.Now().Add(time.Hour * time.Duration(srv.config.RefreshTokenExpiresHours)).Unix()
	refreshTokenString, err := Hash(srv.config.Secret)
	if err != nil {
		return nil, err
	}

	return &messages.AccessToken{
		Token:                 tokenString,
		ExpiresIn:             expires,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresIn: refreshTokenExpires,
	}, nil
}

func (srv *JwtProvider) GenerateOauthJwt(oauthType domains.OauthType, sub, email string, nickname string) (*messages.AccessToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expires := time.Now().Add(time.Hour * time.Duration(srv.config.ExpiresHours)).Unix()

	claims := token.Claims.(jwt.MapClaims)

	claims["type"] = oauthType
	claims["sub"] = sub
	claims["email"] = email
	claims["nickname"] = nickname
	claims["exp"] = expires

	tokenString, err := token.SignedString([]byte(srv.config.Secret))
	if err != nil {
		return nil, err
	}

	// refresh token
	refreshTokenExpires := time.Now().Add(time.Hour * time.Duration(srv.config.RefreshTokenExpiresHours)).Unix()
	refreshTokenString, err := Hash(srv.config.Secret)
	if err != nil {
		return nil, err
	}

	return &messages.AccessToken{
		Token:                 tokenString,
		ExpiresIn:             expires,
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresIn: refreshTokenExpires,
	}, nil
}

func (srv *JwtProvider) GetClaims(c *fiber.Ctx) *messages.JwtClaims {
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return nil
	}

	return &messages.JwtClaims{
		Sub:      claims["sub"].(string),
		Type:     claims["type"].(string),
		Email:    claims["email"].(string),
		Nickname: claims["nickname"].(string),
	}
}
