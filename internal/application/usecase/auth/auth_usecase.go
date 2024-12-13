package auth

import (
	"crou-api/config"
	"crou-api/config/database"
	"crou-api/errorcode"
	"crou-api/internal/application/inputport"
	"crou-api/internal/domains"
	"crou-api/messages"
	"crou-api/utils"
	"github.com/gofiber/fiber/v2"
)

type AuthUseCase struct {
	database    database.Persistent
	jwtProvider *utils.JwtProvider
}

func NewAuthUseCase(cnf *config.Config, database database.Persistent) inputport.AuthInputPort {
	return &AuthUseCase{
		database:    database,
		jwtProvider: utils.NewJwtService(cnf),
	}
}

// @Summary		로그인 API
// @Description	이메일과 비밀번호를 사용하여 로그인합니다.
// @Tags			인증
// @Accept			json
// @Produce		json
// @Success		200	{object}	messages.LoginResponse
// @Failure		401	{object}	server.Error
// @Failure		409	{object}	server.Error
// @Router			/v1/auth/login [post]
func (svc *AuthUseCase) LoginUser(c *fiber.Ctx, req *messages.LoginRequest) (*messages.LoginResponse, error) {
	user, err := svc.getUserByEmail(req.Email)
	if err != nil {
		return nil, errorcode.ErrInvalidEmailOrPassword
	}

	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, errorcode.ErrInvalidEmailOrPassword
	}

	token, err := svc.jwtProvider.GenerateJwt(nil, user.Email, user.Nickname)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return &messages.LoginResponse{
		AccessToken: token,
	}, nil
}

func (svc *AuthUseCase) RegisterUser(c *fiber.Ctx, req *messages.RegisterUserRequest) (*messages.RegisterUserResponse, error) {
	_, err := svc.getUserByEmail(req.Email)
	if err == nil {
		return nil, errorcode.ErrAlreadyUser
	}

	hash, err := utils.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	newUser := &domains.User{
		Email:    req.Email,
		Password: hash,
		Nickname: req.Nickname,
		Gender:   req.Gender,
		Birth:    req.Birth,
	}
	_, err = svc.createUser(newUser)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	}

	return &messages.RegisterUserResponse{
		UserID:   newUser.ID,
		Nickname: newUser.Nickname,
		Gender:   newUser.Gender,
		Birth:    newUser.Birth,
		Email:    newUser.Email,
	}, nil
}

func (svc *AuthUseCase) getUserByEmail(userEmail string) (*domains.User, error) {
	sql := svc.database.DB()
	user := domains.User{}
	result := sql.First(&user, "email = ? ", userEmail)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (svc *AuthUseCase) createUser(newUser *domains.User) (*domains.User, error) {
	sql := svc.database.DB()
	result := sql.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}
