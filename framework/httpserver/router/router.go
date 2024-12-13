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

func Route(
	conf *config.Config,
	api *fiber.App,
	stx *app.ServiceContext,
) {

	noauthGroup := api.Group(V1)

	noauthGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("👋")
	})
	noauthGroup.Get("/auth/google", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}

		result, err := stx.AuthService.OauthGoogleLogin(c, &req)
		return response(c, result, err)
	})
	noauthGroup.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")
		result, err := stx.AuthService.OauthGoogleCallback(c, code, state)
		return response(c, result, err)
	})

	noauthGroup.Get("/auth/naver", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}

		result, err := stx.AuthService.OauthNaverLogin(c, &req)
		return response(c, result, err)
	})

	noauthGroup.Get("/auth/naver/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")
		result, err := stx.AuthService.OauthNaverCallback(c, code, state)
		return response(c, result, err)
	})

	noauthGroup.Post("/auth/refresh", func(c *fiber.Ctx) error {
		req := messages.RefreshTokenRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthService.Refresh(c, &req)
		return response(c, result, err)
	})

	authGroup := api.Group(V1, middleware.JwtMiddleware(conf))
	authGroup.Post("/join", func(c *fiber.Ctx) error {
		req := messages.JoinRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthService.Join(c, &req)
		return response(c, result, err)
	})

	authGroup.Get("/user/profile", func(c *fiber.Ctx) error {
		result, err := stx.UserService.GetUser(c)
		return response(c, result, err)
	})
}
