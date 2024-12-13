package inputport

import (
	"crou-api/messages"
	"github.com/gofiber/fiber/v2"
)

type AuthInputPort interface {
	LoginUser(c *fiber.Ctx, req *messages.LoginRequest) (*messages.LoginResponse, error)
	RegisterUser(c *fiber.Ctx, req *messages.RegisterUserRequest) (*messages.RegisterUserResponse, error)
}
