package messages

import (
	"crou-api/internal/domains"
	"github.com/google/uuid"
)

type JoinRequest struct {
	Nickname string `json:"nickname" validate:"required"`
}

type User struct {
	UserID            uuid.UUID         `json:"userID"`
	Nickname          string            `json:"nickname"`
	OauthType         domains.OauthType `json:"oauthType"`
	OauthSub          string            `json:"oauthSub"`
	OauthEmail        string            `json:"oauthEmail"`
	Taste             string            `json:"taste"`
	NotificationCount uint              `json:"notificationCount,omitempty"`
}
