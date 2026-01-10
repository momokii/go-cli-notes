package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/momokii/go-cli-notes/internal/model"
	"github.com/momokii/go-cli-notes/internal/service"
)

// LinkHandler handles link HTTP requests
type LinkHandler struct {
	noteService any // NoteService interface
}

// GetOutgoingLinks handles GET /api/v1/notes/:id/links
func (h *LinkHandler) GetOutgoingLinks(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	// Get note service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	// Get outgoing links (links from this note to other notes)
	links, err := svc.GetOutgoingLinks(c.Context(), userID, noteID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get outgoing links")
	}

	// Build response with note details
	result := make([]*model.LinkDetail, 0, len(links))
	for _, link := range links {
		detail := &model.LinkDetail{
			ID:          link.ID,
			SourceID:    link.SourceNoteID,
			TargetID:    link.TargetNoteID,
			LinkContext: link.LinkContext,
			CreatedAt:   link.CreatedAt,
		}
		if link.TargetNote != nil {
			detail.TargetNote = link.TargetNote
		}
		result = append(result, detail)
	}

	return sendJSON(c, fiber.StatusOK, result)
}

// GetBacklinks handles GET /api/v1/notes/:id/backlinks
func (h *LinkHandler) GetBacklinks(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	noteID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid note ID")
	}

	// Get note service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	// Get backlinks (links from other notes to this note)
	links, err := svc.GetBacklinks(c.Context(), userID, noteID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get backlinks")
	}

	// Build response with note details
	result := make([]*model.LinkDetail, 0, len(links))
	for _, link := range links {
		detail := &model.LinkDetail{
			ID:          link.ID,
			SourceID:    link.SourceNoteID,
			TargetID:    link.TargetNoteID,
			LinkContext: link.LinkContext,
			CreatedAt:   link.CreatedAt,
		}
		if link.SourceNote != nil {
			detail.SourceNote = link.SourceNote
		}
		result = append(result, detail)
	}

	return sendJSON(c, fiber.StatusOK, result)
}

// GetLinkGraph handles GET /api/v1/notes/graph
func (h *LinkHandler) GetLinkGraph(c *fiber.Ctx) error {
	userIDStr, ok := getUserID(c)
	if !ok {
		return sendError(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return sendError(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get note service
	svc, ok := h.noteService.(*service.NoteService)
	if !ok {
		return sendError(c, fiber.StatusInternalServerError, "Service error")
	}

	// Get the full knowledge graph
	graph, err := svc.GetLinkGraph(c.Context(), userID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "Failed to get link graph")
	}

	return sendJSON(c, fiber.StatusOK, graph)
}
