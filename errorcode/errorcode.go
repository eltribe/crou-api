package errorcode

import (
	"net/http"
)

type UseCaseError struct {
	Message   string `json:"message"`
	Code      int    `json:"code"`
	ErrorCode int    `json:"errorCode"`
}

func (error UseCaseError) Error() string {
	return error.Message
}

func NewUseCaseError(errorCode int, message string) *UseCaseError {
	return &UseCaseError{
		Message:   message,
		Code:      http.StatusConflict,
		ErrorCode: errorCode,
	}
}

var ErrAlreadyUser = NewUseCaseError(-100, "이미 존재하는 사용자입니다")
var ErrInvalidEmailOrPassword = NewUseCaseError(-101, "이메일 또는 비밀번호가 올바르지 않습니다")
