package messages

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Nickname string `json:"nickname" validate:"required,max=20"`
	Birth    int32  `json:"birth" validate:"required,gte=19000101,lte=21000101"`
	Gender   string `json:"gender" validate:"required,oneof=M F"`
}

type RegisterUserResponse struct {
	UserID   uint   `json:"userId"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Birth    int32  `json:"birth"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginResponse struct {
	*AccessToken
}

type AccessToken struct {
	Token                 string `json:"token"`
	ExpiresIn             int64  `json:"expiresIn"`
	RefreshToken          string `json:"refreshToken"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
