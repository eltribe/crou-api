package router

import (
	"crou-api/config"
	"crou-api/framework/httpserver/middleware"
	app "crou-api/internal"
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	noauth := api.Group(V1)
	noauth.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	noauth.Get("/oauth2/google", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthGoogleLogin(c, &req)
		return response(c, result, err)
	})
	noauth.Get("/oauth2/google/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthGoogleCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	noauth.Get("/oauth2/naver", func(c *fiber.Ctx) error {
		req := messages.OauthLoginRequest{}
		if err := queryValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.OauthNaverLogin(c, &req)
		return response(c, result, err)
	})

	noauth.Get("/oauth2/naver/callback", func(c *fiber.Ctx) error {
		result, err := stx.OAuth2UseCase.OauthNaverCallback(c, c.Query("code"), c.Query("state"))
		return response(c, result, err)
	})

	noauth.Post("/oauth2/refresh", func(c *fiber.Ctx) error {
		req := messages.RefreshTokenRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.OAuth2UseCase.Refresh(c, &req)
		return response(c, result, err)
	})

	noauth.Post("/auth/join", func(c *fiber.Ctx) error {
		req := messages.RegisterUserRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthInputPort.RegisterUser(c, &req)
		return response(c, result, err)
	})

	noauth.Post("/auth/login", func(c *fiber.Ctx) error {
		req := messages.LoginRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.AuthInputPort.LoginUser(c, &req)
		return response(c, result, err)
	})

	// ================== Auth Group ==================

	auth := api.Group(V1, middleware.JwtMiddleware(conf))
	auth.Get("/user/profile", func(c *fiber.Ctx) error {
		result, err := stx.UserService.GetUser(c)
		return response(c, result, err)
	})

	auth.Post("/routine", func(c *fiber.Ctx) error {
		req := messages.CreateRoutineRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.RoutineInputPort.CreateRoutine(c, req)
		return response(c, result, err)
	})

	auth.Get("/routine", func(c *fiber.Ctx) error {
		result, err := stx.RoutineInputPort.GetRoutines(c)
		return response(c, result, err)
	})

	auth.Put("/routine", func(c *fiber.Ctx) error {
		req := messages.UpdateRoutineRequest{}
		if err := bodyValidator(c, &req); err != nil {
			return err
		}
		result, err := stx.RoutineInputPort.UpdateRoutine(c, req)
		return response(c, result, err)
	})

	auth.Delete("/routine/:routineId", func(c *fiber.Ctx) error {
		routineId := c.Params("routineId")
		err := stx.RoutineInputPort.DeleteRoutine(c, uuid.MustParse(routineId))
		return response(c, nil, err)
	})
}
