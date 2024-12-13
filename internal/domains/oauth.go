package domains

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

type RefreshToken struct {
	OauthType  OauthType `gormadapter:"uniqueIndex:oauth_type_sub_idx"`
	OauthSub   string    `gormadapter:"uniqueIndex:oauth_type_sub_idx"`
	OauthEmail string    `gormadapter:"unique"`
}
