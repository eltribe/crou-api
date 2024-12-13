package outputport

import "crou-api/internal/domains"

type UserDataOutputPort interface {
	GetUserByEmail(userEmail string) (*domains.User, error)
	CreateUser(newUser *domains.User) (*domains.User, error)
}
