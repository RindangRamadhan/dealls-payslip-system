package middleware

import (
	"strings"

	"github.com/RindangRamadhan/dealls-payslip-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isPublicPath(c.Path()) {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.WriteError(c, fiber.StatusUnauthorized, "Authorization header required", nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.WriteError(c, fiber.StatusUnauthorized, "Invalid token", nil)
		}

		// Set ke context Fiber
		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdminRaw := c.Locals("is_admin")
		isAdmin, ok := isAdminRaw.(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		return c.Next()
	}
}

func isPublicPath(path string) bool {
	publicPaths := []string{
		"/health",
		"/metrics",
		"/swagger",
		"/v1/auth/login",
	}

	for _, prefix := range publicPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}
