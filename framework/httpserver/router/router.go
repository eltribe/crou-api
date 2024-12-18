package router

import (
	"crou-api/config"
	"crou-api/framework/httpserver/middleware"
	app "crou-api/internal"
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
	noauth := api.Group(V1)
	noAuth(conf, noauth, stx)
	swaggerRoute(api)

	// ================== Auth Group ==================

	auth := api.Group(V1, middleware.JwtMiddleware(conf))
	auth.Get("/user/profile", func(c *fiber.Ctx) error {
		result, err := stx.UserService.GetUser(c)
		return response(c, result, err)
	})

	routine(conf, api, stx)
}
