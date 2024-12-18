package router

import (
	_ "crou-api/docs"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func swaggerRoute(
	api *fiber.App,
) {
	api.Get("/swagger/*", swagger.HandlerDefault) // default
}
