package router

import (
	"crou-api/config"
	app "crou-api/internal"
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
)

func noAuth(
	conf *config.Config,
	router fiber.Router,
	stx *app.InputPortProvider,
) {

	router.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	router.Get("/oauth2/google", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthGoogleLogin(c, &req)
		return response(c, result, err)
	})
	router.Get("/oauth2/google/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthGoogleCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	router.Get("/oauth2/naver", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthNaverLogin(c, &req)
		return response(c, result, err)
	})

	router.Get("/oauth2/naver/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthNaverCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	router.Post("/oauth2/refresh", func(c *fiber.Ctx) error {
		req := messages.RefreshTokenRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.Refresh(c, &req)
		return response(c, result, err)
	})

	router.Post("/auth/join", func(c *fiber.Ctx) error {
		req := messages.RegisterUserRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthInputPort.RegisterUser(c, &req)
		return response(c, result, err)
	})

	router.Post("/auth/login", func(c *fiber.Ctx) error {
		req := messages.LoginRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthInputPort.LoginUser(c, &req)
		return response(c, result, err)
	})
}
