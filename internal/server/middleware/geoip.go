package middleware

import (
	"net"

	"github.com/gofiber/fiber/v3"
	"github.com/oschwald/geoip2-golang"
)

// Middleware to block requests from denied countries
func CountryBlockMiddleware(db *geoip2.Reader, deniedCountries map[string]bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := net.ParseIP(c.IP())
		if ip == nil {
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}
		record, err := db.Country(ip)
		if err != nil {
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}
		if deniedCountries[record.Country.IsoCode] {
			return c.Status(fiber.StatusForbidden).SendString("Access denied in your country")
		}
		return c.Next()
	}
}
