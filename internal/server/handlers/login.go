package handlers

import (
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(adminTokenHash, jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {

		var req domain.LoginRequest
		if err := c.Bind().Body(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(domain.DefaultErrorResponse{
				Message: "Invalid request",
				Error:   err.Error(),
			})
		}
		if err := bcrypt.CompareHashAndPassword([]byte(adminTokenHash), []byte(req.AdminToken)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.DefaultErrorResponse{
				Message: "Invalid token",
				Error:   err.Error(),
			})
		}
		claims := jwt.MapClaims{
			"admin": true,
			"exp":   jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.DefaultErrorResponse{
				Message: "Could not sign token",
				Error:   err.Error(),
			})
		}
		return c.JSON(domain.LoginResponse{
			Token: tokenStr,
		})
	}
}
