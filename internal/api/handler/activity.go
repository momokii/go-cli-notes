package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/repository"
)

// ActivityHandler handles activity HTTP requests
type ActivityHandler struct {
	activityRepo any // ActivityRepository interface
	noteService    any // NoteService interface (for trending/forgotten with note details)
}

// GetRecentActivity handles GET /api/v1/activity/recent
func (h *ActivityHandler) GetRecentActivity(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse limit
	limit := c.QueryInt("limit", 20)
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get activity repository
	repo, ok := h.activityRepo.(repository.ActivityRepository)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	activities, err := repo.GetRecent(c.Context(), userID, limit)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get recent activity")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{
		"activities": activities,
	})
}

// GetUserStats handles GET /api/v1/stats
func (h *ActivityHandler) GetUserStats(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get activity repository
	repo, ok := h.activityRepo.(repository.ActivityRepository)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	stats, err := repo.GetUserStats(c.Context(), userID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get user stats")
	}

	return sendJSON(c, fiber.StatusOK, stats)
}

// GetTrendingNotes handles GET /api/v1/notes/trending
func (h *ActivityHandler) GetTrendingNotes(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse limit
	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Get activity repository
	repo, ok := h.activityRepo.(repository.ActivityRepository)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	trending, err := repo.GetTrendingNotes(c.Context(), userID, limit)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get trending notes")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{
		"trending": trending,
	})
}

// GetForgottenNotes handles GET /api/v1/notes/forgotten
func (h *ActivityHandler) GetForgottenNotes(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse parameters
	days := c.QueryInt("days", 30)
	limit := c.QueryInt("limit", 10)

	if days < 1 {
		days = 30
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Get activity repository
	repo, ok := h.activityRepo.(repository.ActivityRepository)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	forgotten, err := repo.GetForgottenNotes(c.Context(), userID, days, limit)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get forgotten notes")
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{
		"forgotten": forgotten,
		"days":       days,
	})
}
