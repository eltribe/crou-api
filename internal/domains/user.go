package domains

import "gorm.io/gorm"

type OauthType string

const (
	GOOGLE OauthType = "GOOGLE"
	NAVER  OauthType = "NAVER"
)

func (t OauthType) String() string {
	return string(t)
}

func ParseOauthType(t string) OauthType {
	return OauthType(t)
}

type User struct {
	gorm.Model
	Nickname   string
	OauthType  OauthType `gorm:"uniqueIndex:oauth_type_sub_idx"`
	OauthSub   string    `gorm:"uniqueIndex:oauth_type_sub_idx"`
	OauthEmail string    `gorm:"unique"`
	Taste      string
}

type UserDetail struct {
	gorm.Model
	UserEmail        string `gorm:"index"`
	LikeCount        uint   // 좋아요 개수
	ClipCount        uint   // 클립 개수
	PrecisionUpCount uint   // 정확도 Up
	User             User   `gorm:"foreignKey:UserEmail;references:OauthEmail"`
}

type RefreshToken struct {
	Nickname   string
	OauthType  OauthType `gorm:"uniqueIndex:oauth_type_sub_idx"`
	OauthSub   string    `gorm:"uniqueIndex:oauth_type_sub_idx"`
	OauthEmail string    `gorm:"unique"`
	Taste      string
}
