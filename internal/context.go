package app

import (
	"crou-api/internal/application/inputport"
	"crou-api/internal/application/usecase/auth"
	"crou-api/internal/application/usecase/routine"
	"crou-api/internal/application/usecase/user"
	"go.uber.org/fx"
)

var Ctx = fx.Module("usecase",
	fx.Provide(
		auth.NewOAuth2Service,
		auth.NewAuthUseCase,
		user.NewUserUseCase,
		routine.NewRoutineUseCase,
		NewInputPortProvider,
	),
)

type InputPortProvider struct {
	AuthInputPort    inputport.AuthInputPort
	OAuth2UseCase    *auth.OAuth2Service
	UserService      *user.UserUseCase
	RoutineInputPort inputport.RoutineInputPort
}

func NewInputPortProvider(
	authInputPort inputport.AuthInputPort,
	oauth2UseCase *auth.OAuth2Service,
	userService *user.UserUseCase,
	routineInputPort inputport.RoutineInputPort,
) *InputPortProvider {

	return &InputPortProvider{
		AuthInputPort:    authInputPort,
		OAuth2UseCase:    oauth2UseCase,
		UserService:      userService,
		RoutineInputPort: routineInputPort,
	}
}
