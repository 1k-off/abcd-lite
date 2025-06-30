package handlers

import (
	"time"

	"github.com/1k-off/abcd-lite/internal/server/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(adminTokenHash, jwtSecret, env string) fiber.Handler {
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
		secure := true
		if env != "production" {
			secure = false
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    tokenStr,
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			Secure:   secure,
			SameSite: fiber.CookieSameSiteStrictMode,
			Path:     "/",
		})

		if c.Accepts("json") == "json" {
			return c.Status(fiber.StatusOK).JSON(domain.LoginResponse{
				Token: tokenStr,
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func Logout(env string) fiber.Handler {
	return func(c fiber.Ctx) error {
		secure := true
		if env != "production" {
			secure = false
		}
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HTTPOnly: true,
			Secure:   secure,
			SameSite: fiber.CookieSameSiteStrictMode,
			Path:     "/",
		})
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out"})
	}
}

func AuthStatus(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var tokenStr string
		if cookie := c.Cookies("jwt"); cookie != "" {
			tokenStr = cookie
		} else if auth := c.Get("Authorization"); auth != "" {
			if len(auth) > 7 && auth[:7] == "Bearer " {
				tokenStr = auth[7:]
			}
		}
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
		}
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"authenticated": false})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"authenticated": true})
	}
}
