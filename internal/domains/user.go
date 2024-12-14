package domains

import (
	"crou-api/common"
	"github.com/google/uuid"
	"time"
)

type User struct {
	common.UUIDModel
	Email        string `gorm:"unique;size:100"`
	Nickname     string `gorm:"size:20"`
	Password     string `gorm:"size:100"`
	Gender       string `gorm:"size:1"` // 성별 (M: 남성, F: 여성)
	Birth        int32  // 생년월일 (YYYYMMDD)
	ProfileImage *string
	OauthType    OauthType `gorm:"uniqueIndex:oauth_type_sub_idx;size:10"`
	OauthSub     *string   `gorm:"uniqueIndex:oauth_type_sub_idx;size:100"`
	OauthEmail   *string   `gorm:"unique;size:100"`
}

type UserDetail struct {
	UserId    uuid.UUID `gorm:"type:uuid;primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserNotificationSetting

	User User `gorm:"foreignKey:UserId;references:ID"`
}

type UserNotificationSetting struct {
	// 오전 알림
	MorningAppPush bool `gorm:"default:true"`
	// 저녁 알림
	EveningAppPush bool `gorm:"default:true"`
	// 주말 알림
	WeekendAppPush bool `gorm:"default:true"`
	// 마케팅, 혜택 알림
	MarketingAppPush bool `gorm:"default:false"`
	// 마케팅, 혜택 앱푸시 수정 날짜
	MarketingAppPushUpdatedAt *time.Time
}
