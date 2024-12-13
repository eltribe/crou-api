package app

import (
	"crou-api/internal/application/auth"
	"crou-api/internal/application/user"
	"go.uber.org/fx"
)

var Ctx = fx.Module("ctx",
	fx.Provide(
		user.NewUserService,
		auth.NewAuthService,
		NewServiceContext,
	),
)

type ServiceContext struct {
	AuthService *auth.AuthService
	UserService *user.UserService
}

func NewServiceContext(
	authService *auth.AuthService,
	userService *user.UserService,
) *ServiceContext {

	return &ServiceContext{
		AuthService: authService,
		UserService: userService,
	}
}
