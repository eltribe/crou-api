package app

import (
	"crou-api/internal/application/inputport"
	"crou-api/internal/application/usecase/auth"
	"crou-api/internal/application/usecase/user"
	"go.uber.org/fx"
)

var Ctx = fx.Module("ctx",
	fx.Provide(
		auth.NewOAuth2Service,
		auth.NewAuthUseCase,
		user.NewUserService,
		NewInputPortProvider,
	),
)

type InputPortProvider struct {
	AuthUseCase   inputport.AuthInputPort
	OAuth2UseCase *auth.OAuth2Service
	UserService   *user.UserService
}

func NewInputPortProvider(
	oauth2UseCase *auth.OAuth2Service,
	authUseCase inputport.AuthInputPort,
	userService *user.UserService,
) *InputPortProvider {

	return &InputPortProvider{
		OAuth2UseCase: oauth2UseCase,
		AuthUseCase:   authUseCase,
		UserService:   userService,
	}
}
