package messages

type JwtClaims struct {
	Sub      string `json:"sub"`
	Type     string `json:"type"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type OauthLoginRequest struct {
	CallbackUrl string `json:"callbackUrl" validate:"required"`
}

type OauthLoginUrl struct {
	Url string `json:"url"`
}

type GoogleOauthUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Locale     string `json:"locale"`
}

type NaverOauthUserInfo struct {
	Resultcode string `json:"resultcode"`
	Message    string `json:"message"`
	Response   struct {
		Email        string `json:"email"`
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
		Age          string `json:"age"`
		Gender       string `json:"gender"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Birthday     string `json:"birthday"`
		Birthyear    string `json:"birthyear"`
		Mobile       string `json:"mobile"`
	} `json:"response"`
}
