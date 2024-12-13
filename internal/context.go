package app

import (
	"crou-api/internal/application/inputport"
	"crou-api/internal/application/usecase/auth"
	"crou-api/internal/application/usecase/user"
	"go.uber.org/fx"
)

var Ctx = fx.Module("usecase",
	fx.Provide(
		auth.NewOAuth2Service,
		auth.NewAuthUseCase,
		user.NewUserUseCase,
		NewInputPortProvider,
	),
)

type InputPortProvider struct {
	AuthUseCase   inputport.AuthInputPort
	OAuth2UseCase *auth.OAuth2Service
	UserService   *user.UserUseCase
}

func NewInputPortProvider(
	oauth2UseCase *auth.OAuth2Service,
	authUseCase inputport.AuthInputPort,
	userService *user.UserUseCase,
) *InputPortProvider {

	return &InputPortProvider{
		OAuth2UseCase: oauth2UseCase,
		AuthUseCase:   authUseCase,
		UserService:   userService,
	}
}
