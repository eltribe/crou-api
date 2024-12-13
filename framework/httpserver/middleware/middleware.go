package middleware

import (
	"crou-api/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"strings"
)

func JwtMiddleware(cnf *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT token from the header
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cnf.Auth.JWT.Secret), nil
		})

		log.Info(token, err)

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Validate the token and return the custom claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// do something with claims
			c.Set("email", claims["email"].(string))
			c.Set("type", claims["type"].(string))
			c.Set("sub", claims["sub"].(string))
			c.Locals("claims", claims)
			_ = claims
		} else {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}

		// Continue stack, if not valid it will return above
		return c.Next()
	}
}

func OptionalJwtMiddleware(cnf *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT token from the header
		authHeader := c.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cnf.Auth.JWT.Secret), nil
		})

		if err != nil {
			return c.Next()
		}

		// Validate the token and return the custom claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// do something with claims
			c.Set("email", claims["email"].(string))
			c.Set("type", claims["type"].(string))
			c.Set("sub", claims["sub"].(string))
			c.Locals("claims", claims)
			_ = claims
		}

		// Continue stack, if not valid it will return above
		return c.Next()
	}
}
