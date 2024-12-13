package router

import (
	"crou-api/config"
	"crou-api/framework/httpserver/middleware"
	app "crou-api/internal"
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
)

const (
	V1 = "/v1"
)

func v1Url(path string) string {
	return V1 + path
}

func Route(
	conf *config.Config,
	api *fiber.App,
	stx *app.InputPortProvider,
) {

	noauthGroup := api.Group(V1)
	noauthGroup.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	noauthGroup.Get("/oauth2/google", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthGoogleLogin(c, &req)
		return response(c, result, err)
	})
	noauthGroup.Get("/oauth2/google/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthGoogleCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	noauthGroup.Get("/oauth2/naver", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthNaverLogin(c, &req)
		return response(c, result, err)
	})

	noauthGroup.Get("/oauth2/naver/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthNaverCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	noauthGroup.Post("/oauth2/refresh", func(c *fiber.Ctx) error {
		req := messages.RefreshTokenRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.Refresh(c, &req)
		return response(c, result, err)
	})

	noauthGroup.Post("/auth/join", func(c *fiber.Ctx) error {
		req := messages.RegisterUserRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthUseCase.RegisterUser(c, &req)
		return response(c, result, err)
	})

	noauthGroup.Post("/auth/login", func(c *fiber.Ctx) error {
		req := messages.LoginRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthUseCase.LoginUser(c, &req)
		return response(c, result, err)
	})

	// ================== Auth Group ==================

	authGroup := api.Group(V1, middleware.JwtMiddleware(conf))
	authGroup.Get("/user/profile", func(c *fiber.Ctx) error {
		result, err := stx.UserService.GetUser(c)
		return response(c, result, err)
	})
}
