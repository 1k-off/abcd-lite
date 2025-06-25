package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"golang.org/x/crypto/bcrypt"
)

func Auth(adminTokenHash string) fiber.Handler {
	return keyauth.New(keyauth.Config{
		Validator: func(c fiber.Ctx, key string) (bool, error) {
			if err := bcrypt.CompareHashAndPassword([]byte(adminTokenHash), []byte(key)); err != nil {
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}
			return true, nil
		},
	})
}
