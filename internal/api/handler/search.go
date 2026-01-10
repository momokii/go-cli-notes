package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/service"
)

// SearchHandler handles search requests
type SearchHandler struct {
	noteService any // NoteService interface
}

// Search handles GET /api/v1/search
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse query parameters
	query := c.Query("q")
	if query == "" {
		return sendError(c, fiber.StatusBadRequest, "Query parameter 'q' is required")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	noteType := c.Query("type")
	tagID := c.Query("tag_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Build filter
	filter := model.NoteFilter{
		Page:   page,
		Limit:  limit,
		Search: query,
		SortBy: "created_at",
	}

	if noteType != "" {
		nt := model.NoteType(noteType)
		filter.NoteType = &nt
	}

	if tagID != "" {
		filter.TagID = &tagID
	}

	// Get note service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	// Search notes (uses List method internally with search filter)
	notes, total, err := svc.Search(c.Context(), userID, filter)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to search notes")
	}

	// Build search results
	results := make([]*model.SearchResult, 0, len(notes))
	for _, note := range notes {
		results = append(results, &model.SearchResult{
			Note:    note,
			Rank:    1.0, // Default rank (can be improved with ts_rank)
			Snippet: generateSnippet(note.Content, query),
		})
	}

	// Calculate pagination
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	response := &model.SearchResponse{
		Query:   query,
		Results: results,
		Pagination: &model.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return sendJSON(c, fiber.StatusOK, response)
}

// generateSnippet creates a highlighted snippet from content
func generateSnippet(content, query string) string {
	// Simple snippet generation - take first 200 chars
	// In production, use ts_headline from PostgreSQL
	maxLen := 200
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}
