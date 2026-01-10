package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/util"
)

// AuthConfig is the configuration for Auth middleware
type AuthConfig struct {
	JWTManager *util.JWTManager
	SkipPaths  map[string]bool
}

// Auth middleware validates JWT tokens
func Auth(jwtManager any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip auth for certain paths
		path := c.Path()
		if isPublicPath(path) {
			return c.Next()
		}

		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return unauthorized(c, "Missing authorization header")
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return unauthorized(c, "Invalid authorization format")
		}

		token := parts[1]

		// Validate token using JWT manager
		// Type assertion for JWT manager
		jwtMgr, ok := jwtManager.(*util.JWTManager)
		if !ok {
			return unauthorized(c, "Invalid JWT manager")
		}

		claims, err := jwtMgr.ValidateAccessToken(token)
		if err != nil {
			return unauthorized(c, "Invalid token")
		}

		// Parse user ID from claims
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return unauthorized(c, "Invalid user ID in token")
		}

		// Store user info in context
		c.Locals("user_id", userID.String())
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// isPublicPath checks if a path should skip authentication
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/api/v1/auth/register",
		"/api/v1/auth/login",
		"/health",
		"/docs",
		"/openapi",
	}

	for _, p := range publicPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// unauthorized returns an unauthorized error response
func unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "unauthorized",
		"message": message,
	})
}
