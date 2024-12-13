package domains

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email        string `gormadapter:"unique;size:100"`
	Nickname     string `gormadapter:"size:20"`
	Password     string `gormadapter:"size:100"`
	Gender       string `gormadapter:"size:1"` // 성별 (M: 남성, F: 여성)
	Birth        int32  // 생년월일 (YYYYMMDD)
	ProfileImage *string
	OauthType    OauthType `gormadapter:"uniqueIndex:oauth_type_sub_idx;size:10"`
	OauthSub     *string   `gormadapter:"uniqueIndex:oauth_type_sub_idx;size:100"`
	OauthEmail   *string   `gormadapter:"unique;size:100"`
}

type UserDetail struct {
	UserId    uint `gormadapter:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserNotificationSetting

	User User `gormadapter:"foreignKey:UserId;references:ID"`
}

type UserNotificationSetting struct {
	// 오전 알림
	MorningAppPush bool `gormadapter:"default:true"`
	// 저녁 알림
	EveningAppPush bool `gormadapter:"default:true"`
	// 주말 알림
	WeekendAppPush bool `gormadapter:"default:true"`
	// 마케팅, 혜택 알림
	MarketingAppPush bool `gormadapter:"default:false"`
	// 마케팅, 혜택 앱푸시 수정 날짜
	MarketingAppPushUpdatedAt *time.Time
}
