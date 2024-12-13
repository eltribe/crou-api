package gormadapter

import (
	"crou-api/config/database"
	"crou-api/internal/application/outputport"
	"crou-api/internal/domains"
)

type UserGorm struct {
	database database.Persistent
}

func NewUserGorm(database database.Persistent) outputport.UserDataOutputPort {
	return &UserGorm{database: database}
}

func (gorm *UserGorm) GetUserByEmail(userEmail string) (*domains.User, error) {
	sql := gorm.database.DB()
	user := domains.User{}
	result := sql.First(&user, "email = ? ", userEmail)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (gorm *UserGorm) CreateUser(newUser *domains.User) (*domains.User, error) {
	sql := gorm.database.DB()
	result := sql.Create(newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	sql.Create(&domains.UserDetail{
		UserId: newUser.ID,
	})

	return newUser, nil
}
