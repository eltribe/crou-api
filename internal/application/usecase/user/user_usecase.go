package user

import (
	"crou-api/config/database"
	"crou-api/internal/domains"
	"crou-api/messages"
	"crou-api/utils"
	"github.com/gofiber/fiber/v2"
)

type UserUseCase struct {
	database   database.Persistent
	jwtService *utils.JwtProvider
}

func NewUserUseCase(database database.Persistent) *UserUseCase {
	return &UserUseCase{
		database: database,
	}
}

// GetUser godoc
//
//	@Summary		유저 정보 조회  API
//	@Description	JWT 기반 유저정보를 조회합니다.
//	@Tags			유저정보
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	messages.User
//	@Failure		401	{object}	server.Error
//	@Failure		409	{object}	server.Error
//	@Router			/v1/user/profile [get]
func (svc *UserUseCase) GetUser(c *fiber.Ctx) (*messages.User, error) {
	claims := svc.jwtService.GetClaims(c)
	user, err := svc.GetUserByOauthInfo(domains.OauthType(claims.Type), claims.Sub)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	}

	return &messages.User{
		UserID:     user.ID,
		Nickname:   user.Nickname,
		OauthType:  user.OauthType,
		OauthSub:   *user.OauthSub,
		OauthEmail: *user.OauthEmail,
	}, nil
}

func (svc *UserUseCase) GetUserByOauthInfo(oauthType domains.OauthType, oauthSub string) (*domains.User, error) {
	sql := svc.database.DB()
	user := domains.User{}
	result := sql.First(&user, "oauth_type = ? and oauth_sub = ?", oauthType, oauthSub)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (svc *UserUseCase) GetUserById(userId uint) (*domains.User, error) {
	sql := svc.database.DB()
	user := domains.User{}
	result := sql.First(&user, userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (svc *UserUseCase) createUser(newUser *domains.User) (*domains.User, error) {
	sql := svc.database.DB()
	result := sql.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return newUser, nil
}

func (svc *UserUseCase) getUserByEmail(userEmail string) (*domains.User, error) {
	sql := svc.database.DB()
	user := domains.User{}
	result := sql.First(&user, "email = ? ", userEmail)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (svc *UserUseCase) getUserDetailByEmail(userEmail string) (*domains.UserDetail, error) {
	sql := svc.database.DB()
	userDetail := domains.UserDetail{}
	result := sql.First(&userDetail, "user_email = ? ", userEmail)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userDetail, nil
}

func (svc *UserUseCase) updateUser(user *domains.User) error {
	sql := svc.database.DB()
	result := sql.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
