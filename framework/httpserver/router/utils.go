package router

import (
	"crou-api/errorcode"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-querystring/query"
)

// 생성한 구조체에 대해 Validate 함수 생성
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// requestValidator 함수는 req를 파싱하고 검증합니다.
// T 형식의 req를 인자로 받고 검증 에러가 없다면 nil을 반환합니다.
func bodyValidator[T any](c *fiber.Ctx, req *T) error {
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate struct
	if err := validate.Struct(req); err != nil {
		return err
	}

	return nil
}

func queryValidator[T any](c *fiber.Ctx, req *T) error {
	if err := c.QueryParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate struct
	if err := validate.Struct(req); err != nil {
		return err
	}

	return nil
}

// 공통 응답
func response(c *fiber.Ctx, response interface{}, err error) error {
	if err != nil {
		var useCaseError *errorcode.UseCaseError
		if errors.As(err, &useCaseError) {
			return c.Status(useCaseError.Code).JSON(useCaseError)
		}
		return err
	}
	return c.JSON(response)
}

func responseRedirect(c *fiber.Ctx, url string, response interface{}, err error) error {
	if err != nil {
		return err
	}
	v, _ := query.Values(response)
	fmt.Print(v.Encode())
	return c.Redirect(url+"?"+v.Encode(), 301)
}
