package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/service"
)

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
	}

	// Call service (type assertion)
	svc, ok := h.authService.(*service.AuthService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	user, err := svc.Register(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusCreated, user)
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Call service
	svc, ok := h.authService.(*service.AuthService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	resp, err := svc.Login(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, resp)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Call service
	svc, ok := h.authService.(*service.AuthService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	resp, err := svc.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, resp)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Call service
	svc, ok := h.authService.(*service.AuthService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	// TODO: Get user ID from context and revoke token
	_ = svc

	return sendJSON(c, fiber.StatusOK, fiber.Map{"message": "Logged out"})
}

// handleError maps service errors to HTTP status codes
func handleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for specific error types by inspecting the error message
	switch {
	case errMsg == "resource not found" || errMsg == "find user: resource not found":
		return sendError(c, fiber.StatusNotFound, "Resource not found")
	case errMsg == "unauthorized access":
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	case errMsg == "invalid credentials":
		return sendError(c, fiber.StatusUnauthorized, "Invalid email or password")
	case errMsg == "check email exists: resource not found" || errMsg == "check username exists: resource not found":
		// These are actually success cases (user doesn't exist yet)
		return sendError(c, fiber.StatusInternalServerError, "Internal error")
	default:
		// Log the actual error for debugging
		return sendError(c, fiber.StatusInternalServerError, "Internal server error: "+errMsg)
	}
}
