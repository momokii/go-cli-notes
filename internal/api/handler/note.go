package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/service"
)

// Create handles note creation
func (h *NoteHandler) Create(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req model.CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	note, err := svc.Create(c.Context(), userID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusCreated, note)
}

// List handles note listing
func (h *NoteHandler) List(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	noteType := c.Query("type")
	search := c.Query("search")
	tagID := c.Query("tag")

	filter := model.NoteFilter{
		Page:   page,
		Limit:  limit,
		Search: search,
	}

	if noteType != "" {
		nt := model.NoteType(noteType)
		filter.NoteType = &nt
	}

	if tagID != "" {
		filter.TagID = &tagID
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	notes, total, err := svc.List(c.Context(), userID, filter)
	if err != nil {
		return handleError(c, err)
	}

	// Calculate pagination
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{
		"notes": notes,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetByID handles getting a single note
func (h *NoteHandler) GetByID(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteIDStr := c.Params("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	note, err := svc.GetByID(c.Context(), userID, noteID)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, note)
}

// Update handles note update
func (h *NoteHandler) Update(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteIDStr := c.Params("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	var req model.UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	note, err := svc.Update(c.Context(), userID, noteID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, note)
}

// Delete handles note deletion
func (h *NoteHandler) Delete(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteIDStr := c.Params("id")
	noteID, err := uuid.Parse(noteIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	if err := svc.Delete(c.Context(), userID, noteID); err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{"message": "Note deleted"})
}

// GetOrCreateDailyNote handles GET /api/v1/notes/daily/:date
func (h *NoteHandler) GetOrCreateDailyNote(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get date from path parameter (YYYY-MM-DD format)
	dateStr := c.Params("date")
	if dateStr == "" {
		return sendError(c, fiber.StatusBadRequest, "Date parameter is required")
	}

	// Call service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	note, isCreated, err := svc.GetOrCreateDailyNote(c.Context(), userID, dateStr)
	if err != nil {
		return handleError(c, err)
	}

	return sendJSON(c, fiber.StatusOK, fiber.Map{
		"note":       note,
		"is_created": isCreated,
		"date":       dateStr,
	})
}
